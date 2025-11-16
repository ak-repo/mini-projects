package main

import "github.com/gofiber/fiber/v2"



func setupRoutes(app *fiber.App){
	app.Get("/:url",routes.ResolveURL)
	

}


func main(){

}