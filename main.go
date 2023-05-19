package main

import (
    "log"
    "net/http"
	"govtech/handlers"
)


func main() {
    log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/products", handlers.HandleProducts)
	http.HandleFunc("/products/", handlers.HandleProduct)
	http.HandleFunc("/products/reviews/", handlers.HandleProductReviews)
    http.ListenAndServe(":8080", nil)
}