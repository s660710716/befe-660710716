package main

import(
	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-contrib/cors"

	_ "week11-assignment/docs"
)

  type Book struct {
      ID            int       `json:"id"`
      Title         string    `json:"title"`
      Author        string    `json:"author"`
      ISBN          string    `json:"isbn"`
      Year          int       `json:"year"`
      Price         float64   `json:"price"`

      // ฟิลด์ใหม่
      Category      string    `json:"category"`
      OriginalPrice *float64  `json:"original_price,omitempty"`
      Discount      int       `json:"discount"`
      CoverImage    string    `json:"cover_image"`
      Rating        float64   `json:"rating"`
      ReviewsCount  int       `json:"reviews_count"`
      IsNew         bool      `json:"is_new"`
      Pages         *int      `json:"pages,omitempty"`
      Language      string    `json:"language"`
      Publisher     string    `json:"publisher"`
      Description   string    `json:"description"`

      CreatedAt     time.Time `json:"created_at"`
      UpdatedAt     time.Time `json:"updated_at"`
  }
type ErrorResponse struct {
    Message string `json:"message"`
}

func getEnv(key, defaultValue string) string{
	if value := os.Getenv(key); value != ""{
		return value
	}
	return defaultValue 
}

var db *sql.DB

func initDB(){
	var err error
	host := getEnv("DB_HOST", "")
	name := getEnv("DB_NAME", "")
	user := getEnv("DB_USER","")
	password := getEnv("DB_PASSWORD","")
	port := getEnv("DB_PORT","")

	conSt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	//fmt.Println(conSt)

	db, err = sql.Open("postgres", conSt)
	if err !=nil {
		log.Fatal("failed to open database")
	}
	
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	log.Print("successfully connect to database")

	// กำหนดจำนวน Connection สูงสุด
	db.SetMaxOpenConns(25)

	// กำหนดจำนวน Idle connection สูงสุด
	db.SetMaxIdleConns(20)

	// กำหนดอายุของ Connection
	db.SetConnMaxLifetime(5 * time.Minute)
}

// @Summary Get all book
// @Description Get details of a book by ID
// @Tags Books
// @Produce  json
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books [get]
func getAllBooks(c *gin.Context) {
    debug := " "
    category := c.Query("category")
    var rows *sql.Rows
    var err error
    // ลูกค้าถาม "มีหนังสืออะไรบ้าง"
    if category !=""{
     rows, err = db.Query("SELECT id, title, author, isbn, year, price, created_at, updated_at, category, rating FROM books WHERE category = $1", category)
     debug= "GetAllCate"
     log.Print("check: ",debug)
     }else{
    rows, err = db.Query("SELECT id, title, author, isbn, year, price, created_at, updated_at, category, rating FROM books")
  
    debug= "GetAll"
    log.Print("check: ",debug)
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close() // ต้องปิด rows เสมอ เพื่อคืน Connection กลับ pool
   
    var books []Book
    for rows.Next() {
        var book Book
            if category != "" {
        err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Category)
    } else {
        err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt, &book.Category, &book.Rating)
    }
    
    if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}
    log.Print("check: ",debug)
	c.JSON(http.StatusOK, books)
}

// @Summary Get book by ID
// @Description Get details of a book by ID
// @Tags Books
// @Produce  json
// @Param   id   path      int     true  "Book ID"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books/{category} [get]
func getBook(c *gin.Context) {
    id := c.Query("id")
    var book Book

    // QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์ 0 หรือ 1 แถว
    err := db.QueryRow("SELECT id, title, author, category FROM books WHERE id = $1", id).
        Scan(&book.ID, &book.Title, &book.Author, &book.Category)

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, book)
}

// @Summary Get book by Category
// @Description Get details of a book by Category
// @Tags Books
// @Produce  json
// @Param   category   path      int     true  "Book Category"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books/{category} [get]
func getBookByCategory(c *gin.Context) {
    category := c.Param("category")
  
    var rows *sql.Rows
    var err error
    // QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์ 0 หรือ 1 แถว
    rows, err = db.Query("SELECT id, title, author, category FROM books WHERE category = $1", category)

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var books []Book
    for rows.Next() {
        var book Book
        err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
        if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

func Search (c *gin.Context){
    q := c.Query("q")

    var rows *sql.Rows
    var err error
    
    //searchQuery := q
    searchQuery := "%"+q+"%"
    rows, err = db.Query(`SELECT id, title, author, category, isbn, year, price, created_at, updated_at FROM books 
    WHERE title LIKE $1
    OR author   LIKE $1
    OR category LIKE $1
    OR isbn     LIke $1`, searchQuery)
    log.Print("check: search->", searchQuery)
    log.Print("SELECT id, title, author, category, isbn, year, price, created_at, updated_at FROM books WHERE title LIKE $1", searchQuery)
    
    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var books []Book
    for rows.Next() {
        var book Book
        err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
        if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)

}
func getFeatureBooks(c *gin.Context){

    var rows *sql.Rows
    var err error
    // QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์ 0 หรือ 1 แถว
    rows, err = db.Query("SELECT id, title, author, category FROM books WHERE rateing > 4.5")

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var books []Book
    for rows.Next() {
        var book Book
        err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
        if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)

}

func getNewBooks(c *gin.Context){
   
    var rows *sql.Rows
    var err error
    // QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์ 0 หรือ 1 แถว
    rows, err = db.Query("SELECT id, title, author, category FROM books WHERE is_new =true")

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var books []Book
    for rows.Next() {
        var book Book
        err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
        if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)

}

func getDiscountedBooks(c *gin.Context){

    var rows *sql.Rows
    var err error
    // QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์ 0 หรือ 1 แถว
    rows, err = db.Query("SELECT id, title, author, category FROM books WHERE discount > 0")

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var books []Book
    for rows.Next() {
        var book Book
        err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
        if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)

}

// @Summary Post book
// @Description post new book
// @Tags Books
// @Produce  json
// @Param   book  body  Book  true  "Book JSON"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books/ [post]
func createBook(c *gin.Context) {
    var newBook Book
    
    if err := c.ShouldBindJSON(&newBook); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // ใช้ RETURNING เพื่อดึงค่าที่ database generate (id, timestamps)
    var id int
    var createdAt, updatedAt time.Time

    err := db.QueryRow(
        `INSERT INTO books (title, author, isbn, year, price)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, created_at, updated_at`,
        newBook.Title, newBook.Author, newBook.ISBN, newBook.Year, newBook.Price,
    ).Scan(&id, &createdAt, &updatedAt)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    newBook.ID = id
    newBook.CreatedAt = createdAt
    newBook.UpdatedAt = updatedAt

    c.JSON(http.StatusCreated, newBook) // ใช้ 201 Created
}

// @Summary update book by ID
// @Description update details of a book by ID
// @Tags Books
// @Produce  json
// @Param   id   path      int     true  "Book ID"
// @Param   book  body  Book  true  "Book JSON"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books/{id} [put]
func updateBook(c *gin.Context) {
    var ID int
    id := c.Param("id")
    var updateBook Book

    if err := c.ShouldBindJSON(&updateBook); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var updatedAt time.Time
    err := db.QueryRow(
        `UPDATE books
         SET title = $1, author = $2, isbn = $3, year = $4, price = $5
         WHERE id = $6
         RETURNING ID,updated_at`,
        updateBook.Title, updateBook.Author, updateBook.ISBN,
        updateBook.Year, updateBook.Price, id,
    ).Scan(&ID, &updatedAt)

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    updateBook.ID = ID
	updateBook.UpdatedAt = updatedAt
	c.JSON(http.StatusOK, updateBook)
}

// @Summary Delete book by ID
// @Description Delete details of a book by ID
// @Tags Books
// @Produce  json
// @Param   id   path      int     true  "Book ID"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books/{id} [delete]
func deleteBook(c *gin.Context) {
    id := c.Param("id")

    result, err := db.Exec("DELETE FROM books WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

// @title           Simple API Example
// @version         1.0
// @description     This is a simple example of using Gin with Swagger.
// @host localhost:8080
// @host 127.0.0.1:8080
// @BasePath        /api/v1
func main(){
	initDB()
	defer db.Close()

	r := gin.Default()	
    r.Use(cors.Default())

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		err := db.Ping()
		if err != nil{
			c.JSON(http.StatusServiceUnavailable, gin.H{"message":"unhealty", "error":err})
			return
		}
		c.JSON(200, gin.H{"message": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/books", getAllBooks)
	 	api.GET("/books/:id", getBook)
	 	api.POST("/books", createBook)
	 	api.PUT("/books/:id", updateBook)
	 	api.DELETE("/books/:id", deleteBook)

        api.GET("categories/:category", getBookByCategory)
        api.GET("books/search", Search)// - ค้นหา
        api.GET("books/featured", getFeatureBooks)// - หนังสือแนะนำ
        api.GET("books/new", getNewBooks)
        api.GET("books/discounted", getDiscountedBooks)// - หนังสือลดราคา

	 }

	r.Run(":8080")
}





