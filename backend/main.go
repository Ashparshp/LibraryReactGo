package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

type Book struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	ISBN     string `json:"isbn"`
	Quantity int    `json:"quantity"`
}

func main() {
	dsn := "host=localhost user=library_user password=password dbname=library_db port=5432 sslmode=disable"
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(&Book{}).Error; err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

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
	c.JSON(http.StatusOK, gin.H{"message": "Book added successfully"})
}

func searchBook(c *gin.Context) {
	query := c.Query("query")
	var books []Book
	db.Where("title LIKE ? OR author LIKE ? OR isbn LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&books)
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
	c.JSON(http.StatusOK, gin.H{"message": "Book removed successfully"})
}

func getTotalBooks(c *gin.Context) {
	var totalBooks int
	if err := db.Model(&Book{}).Count(&totalBooks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total books"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"totalBooks": totalBooks})
}
