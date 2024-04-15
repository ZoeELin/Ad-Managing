package main

import (
	"flag"
	"fmt"

	"ad-proj/database"
	"ad-proj/router"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func main() {

	const (
		// Parameters that db connect to
		HOST     = "localhost"
		DATABASE = "postgres"
		USER     = "postgres"
		PASSWORD = "postgres"
		PORT     = 5400
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

	database.ConnectDatabase(dsn)

	r := gin.Default()
	r.POST("/api/v1/ad", router.CreateAd)
	r.GET("/api/v1/ad", router.ListAds)

	r.Run(":8088")
}
