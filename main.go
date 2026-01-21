package main

import (
	"database/sql"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {

	db := initDatabase()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	
	r := gin.Default()
	
	
	r.LoadHTMLGlob(filepath.Join("templates", "*.html"))

	r.GET("/", func(c *gin.Context) {
		stats, err := getStats(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "DB-fel: %v", err)
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Stats": stats, 
		})
	})	


	r.POST("/play", func(c *gin.Context) {

		playerChoice := c.PostForm("choice")

		if playerChoice != "sten" && playerChoice != "sax" && playerChoice != "påse" {
			c.String(http.StatusBadRequest, "Ogiltigt val")
			return
		}


		computerChoice := getComputerChoice()
		
		result := getResults(playerChoice, computerChoice)


		if err := insertRound(db, playerChoice, computerChoice,result); err != nil {
			c.String(http.StatusInternalServerError, "Kunde inte spara i DB: %v", err)
			return
		}	

		stats, err := getStats(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "DB-fel: %v", err)
			return
		}


		c.HTML(http.StatusOK, "result.html", gin.H{
			"PlayerChoice": playerChoice,
			"ComputerChoice": computerChoice,
			"Result": result,
			"Stats": stats,
		})
	})








	_ = r.Run("0.0.0.0:8080")

}