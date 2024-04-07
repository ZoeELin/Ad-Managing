package main

import (
	"flag"
	"fmt"

	"ad-proj/controllers"
	"ad-proj/database"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func main() {

	const (
		// Parameters that db connect to
		HOST = "localhost"
		// HOST     = "postgres_db"
		DATABASE = "postgres"
		USER     = "postgres"
		PASSWORD = "postgres"
		PORT     = 5400
		// PORT     = 5432
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		HOST, PORT, USER, DATABASE, PASSWORD)

	// Controller to initiallize the database
	initFlag := flag.Bool("init", false, "Initialize the database")
	flag.Parse()

	// if add -init, initiallize db
	if *initFlag {
		database.DbInit(dsn)
		database.DatasetInit(dsn)
		fmt.Println("Initialization completed.")
	}

	r := gin.Default()

	database.ConnectDatabase(dsn)

	r.POST("/api/v1/ad", controllers.CreateAd)
	r.GET("/api/v1/ad", controllers.ListAds)

	fmt.Println("Server is running on port 8088")
	r.Run(":5001")
}
