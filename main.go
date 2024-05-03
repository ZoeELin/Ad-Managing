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
		HOST = "dpg-cof4s30cmk4c7380b7ig-a"
		// HOST     = "dpg-cof4s30cmk4c7380b7ig-a.oregon-postgres.render.com"
		DATABASE = "ad_proj"
		USER     = "admin"
		PASSWORD = "AD8Mi6fAnpPvoMsXOHeuxOWhRy5Ghrti"
		PORT     = 5432
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=allow",
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
	r.GET("/", router.HelloGo)
	r.POST("/api/v1/ad", router.CreateAd)
	r.GET("/api/v1/ad", router.ListAds)

	r.Run(":8088")
}
