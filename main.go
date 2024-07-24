package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Clip represents the clip item structure
type Clip struct {
	ID          int        `db:"id" json:"id"`
	URL         string     `db:"url" json:"url"`
	Content     string     `db:"content" json:"content"`
	DeleteAfter *time.Time `db:"delete_after" json:"delete_after"`
}

var db *sqlx.DB

func main() {
	var err error
	db, err = sqlx.Connect("postgres", "user=postgres password=password dbname=golang_internet_clipboard sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// HTML Routes
	r.GET("/", showIndexPage)
	r.GET("/:url", showClipFormOrContent)
	r.POST("/:url", createClipForm)

	// API Routes
	api := r.Group("/api")
	{
		api.POST("/:url", createClip)
		api.GET("/:url", getClip)
	}

	r.Run(":8080")
}

// HTML Handlers

func showIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func showClipFormOrContent(c *gin.Context) {
	url := c.Param("url")

	var clip Clip
	err := db.Get(&clip, "SELECT * FROM clips WHERE url = $1", url)
	if err == sql.ErrNoRows {
		// No content, show form
		c.HTML(http.StatusOK, "clip.html", gin.H{
			"url": url,
		})
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	location, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to load time zone")
		return
	}
	now := time.Now().In(location)

	// fmt.Println(time.Now(), "------", clip.DeleteAfter, "-----", time.Now().After(*clip.DeleteAfter))

	if clip.DeleteAfter == nil || clip.DeleteAfter.IsZero() || now.After(*clip.DeleteAfter) {
		// Show content and delete if is view once or if time has exceeded
		db.Exec("DELETE FROM clips WHERE url = $1", url)
		c.HTML(http.StatusOK, "clip_content.html", gin.H{
			"clip": clip,
		})
	} else {
		// Content exists and not yet expired, show content
		c.HTML(http.StatusOK, "clip_content.html", gin.H{
			"clip": clip,
		})
	}
}

func createClipForm(c *gin.Context) {
	url := c.Param("url")
	content := c.PostForm("content")
	deleteAfterOption := c.PostForm("delete_after")

	var deleteAfter *time.Time // Use pointer to time.Time to handle nullable timestamps

	if deleteAfterOption != "" {
		var da time.Time
		switch deleteAfterOption {
		case "once_viewed":
			deleteAfter = nil
		case "1_minute":
			da = time.Now().Add(1 * time.Minute)
		case "1_hour":
			da = time.Now().Add(1 * time.Hour)
		case "1_day":
			da = time.Now().Add(24 * time.Hour)
		default:
			c.String(http.StatusBadRequest, "Invalid delete_after option")
			return
		}
		deleteAfter = &da
	}

	clip := Clip{
		URL:         url,
		Content:     content,
		DeleteAfter: deleteAfter, // Assign the pointer to DeleteAfter
	}

	_, err := db.NamedExec(`
		INSERT INTO clips (url, content, delete_after)
		VALUES (:url, :content, :delete_after)`, &clip)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Render the success page
	c.HTML(http.StatusOK, "success.html", gin.H{
		"url": url,
	})
}

// API Handlers

func createClip(c *gin.Context) {
	url := c.Param("url")
	var clip Clip
	if err := c.ShouldBindJSON(&clip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clip.URL = url
	_, err := db.NamedExec(`INSERT INTO clips (url, content) VALUES (:url, :content)`, &clip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, clip)
}

func getClip(c *gin.Context) {
	url := c.Param("url")

	var clip Clip
	err := db.Get(&clip, "SELECT * FROM clips WHERE url = $1", url)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"message": "No content"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clip)
}
