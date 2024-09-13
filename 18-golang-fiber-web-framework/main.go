package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Prefork:      true,
	})

	app.Use("/api", func(c *fiber.Ctx) error {
		fmt.Println("i'm middleware before processing request")
		err := c.Next()
		fmt.Println("i'm middleware after processing request")
		return err
	})
	if fiber.IsChild() {
		fmt.Println("Child process")
	} else {
		fmt.Println("Main process")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	err := app.Listen("localhost:3000")
	if err != nil {
		panic(err)
	}

}
