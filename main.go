package main

import (
	"context"
	"log"
	"os"

	// "github.com/NachoxMacho/opnlab/database"
	"github.com/NachoxMacho/opnlab/opnsense"
	"github.com/NachoxMacho/opnlab/portal"
	"github.com/NachoxMacho/opnlab/proxmox"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	godotenv.Load()
	// Run any migrations and by extension test database connection
	// if err := database.Initialize(); err != nil {
	// 	log.Fatal(err)
	// }

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	err := redisClient.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	val, err := redisClient.Get(context.Background(), "key").Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(val)

	err = proxmox.InitalizeConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = opnsense.InitalizeConfig()
	if err != nil {
		log.Fatal(err)
	}

	// This loads the views folder as html templates so they can be referred to by name in all routes
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New(logger.ConfigDefault))

	// We do this so the portal package can define it's own routes under this path
	err = portal.RegisterRoutes(app.Group("/portal"))
	if err != nil {
		log.Fatal(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the app, go to /portal for a start page")
	})

	app.Static("/css", "./css")

	log.Fatal(app.Listen(":42069"))
}
