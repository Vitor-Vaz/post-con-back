package main

import (
	"fmt"
	"log"

	"post-con-back/extension/database"
	"post-con-back/internal/app"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	r := app.NewRouter(db)

	fmt.Println("post-con-back listening on :8080")
	log.Fatal(r.Run(":8080"))
}
