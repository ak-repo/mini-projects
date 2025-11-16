package main

import (
	"log"

	"github.com/ak-repo/code-board/services/user-auth-service/internals/app"
	router "github.com/ak-repo/code-board/services/user-auth-service/internals/routes"
)

func main() {
	app := app.NewApp()

	r := router.NewRouter(app)

	addr := "0.0.0.0:8000"
	log.Printf("🌐 Server running at %s", addr)
	r.Run(addr)

}
