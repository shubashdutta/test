package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/shubash/pipo/router"
)

func main() {
	
	r := router.Router()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	fmt.Println("server is ready .........")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
