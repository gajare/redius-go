package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gajare/redius-go/config"
	"github.com/gajare/redius-go/db"
	"github.com/gajare/redius-go/router"
)

func main() {
	config.LoadEnv()
	db.InitPostgres()
	db.InitRedis()

	r := router.Router()
	port := os.Getenv("PORT")
	fmt.Printf("Server started at :%s\n", port)
	http.ListenAndServe(":"+port, r)
}
