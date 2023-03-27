package main

import (
	"log"

	"github.com/sivaramsajeev/terraform-provider-student/api/server"
)

func main() {
	store := server.NewRedisStore()
	studentSvc := server.NewService("localhost:8888", store)
	err := studentSvc.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
