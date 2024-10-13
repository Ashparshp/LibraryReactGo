package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *gorm.DB
var mongoClient *mongo.Client

type Book struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	ISBN     string `json:"isbn"`
	Quantity int    `json:"quantity"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("POSTGRES_DSN")
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(&Book{}).Error; err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)
	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.TODO())

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/add_book", addBook)
	r.GET("/search_book", searchBook)
	r.PUT("/borrow_book/:isbn", borrowBook)
	r.PUT("/return_book/:isbn", returnBook)
	r.GET("/list_books", listBooks)
	r.DELETE("/remove_book/:isbn", removeBook)
	r.GET("/total_books", getTotalBooks)
	r.Run(":8080")
}

func addBook(c *gin.Context) {
	var newBook Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&newBook).Error; err != nil {
		log.Printf("Error creating book: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	collection := mongoClient.Database("library_db").Collection("books")
	_, err := collection.InsertOne(context.TODO(), newBook)
	if err != nil {
		log.Printf("Error inserting book into MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book added successfully"})
}

func searchBook(c *gin.Context) {
	query := c.Query("query")
	var books []Book

	db.Where("title LIKE ? OR author LIKE ? OR isbn LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&books)

	collection := mongoClient.Database("library_db").Collection("books")
	filter := bson.M{"$or": []bson.M{
		{"title": bson.M{"$regex": query, "$options": "i"}},
		{"author": bson.M{"$regex": query, "$options": "i"}},
		{"isbn": bson.M{"$regex": query, "$options": "i"}},
	}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error searching books in MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var book Book
		cursor.Decode(&book)
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

func borrowBook(c *gin.Context) {
	isbn := c.Param("isbn")
	var book Book
	if err := db.Where("isbn = ?", isbn).First(&book).Error; err != nil {
		log.Printf("Error finding book: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	if book.Quantity > 0 {
		book.Quantity--
		db.Save(&book)

		collection := mongoClient.Database("library_db").Collection("books")
		_, err := collection.UpdateOne(context.TODO(), bson.M{"isbn": isbn}, bson.M{"$inc": bson.M{"quantity": -1}})
		if err != nil {
			log.Printf("Error updating book in MongoDB: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Book borrowed successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
	}
}

func returnBook(c *gin.Context) {
	isbn := c.Param("isbn")
	var book Book
	if err := db.Where("isbn = ?", isbn).First(&book).Error; err != nil {
		log.Printf("Error finding book: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	book.Quantity++
	db.Save(&book)

	collection := mongoClient.Database("library_db").Collection("books")
	_, err := collection.UpdateOne(context.TODO(), bson.M{"isbn": isbn}, bson.M{"$inc": bson.M{"quantity": 1}})
	if err != nil {
		log.Printf("Error updating book in MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

func listBooks(c *gin.Context) {
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		log.Printf("Error listing books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

func removeBook(c *gin.Context) {
	isbn := c.Param("isbn")
	if err := db.Where("isbn = ?", isbn).Delete(&Book{}).Error; err != nil {
		log.Printf("Error deleting book: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	collection := mongoClient.Database("library_db").Collection("books")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"isbn": isbn})
	if err != nil {
		log.Printf("Error deleting book from MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book removed successfully"})
}

func getTotalBooks(c *gin.Context) {
	var totalBooks int
	if err := db.Model(&Book{}).Count(&totalBooks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total books"})
		return
	}

	collection := mongoClient.Database("library_db").Collection("books")
	count, err := collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total books from MongoDB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"totalBooks": (totalBooks + int(count)) / 2})
}
