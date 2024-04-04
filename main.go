package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"ad-proj/controllers"
	"ad-proj/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get parameter from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}

	// Set the parameters to dsn
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := "postgres"
	dbname := "postgres"
	pass := os.Getenv("PASSWORD")
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)

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
	r.Run(":8088")
}
