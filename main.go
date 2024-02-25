package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/NachoxMacho/opnlab/database"
	"github.com/NachoxMacho/opnlab/portal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	godotenv.Load()
	// Run any migrations and by extension test database connection
	if err := database.Initialize(); err != nil {
		log.Fatal(err)
	}

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

	log.Println(getProxmoxData[ProxmoxNode]("https://proxmox.docker.homelab:8006/api2/json/nodes"))
	log.Println(getProxmoxData[ProxmoxVM]("https://proxmox.docker.homelab:8006/api2/json/nodes/prox/qemu"))

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

type ProxmoxNode struct {
	Node   string `json:"node"`
	Status string `json:"status"`
}

type ProxmoxVM struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ProxmoxResponses interface{ ProxmoxNode | ProxmoxVM }

func getProxmoxData[ResponseType ProxmoxResponses](url string) []ResponseType {

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	defaultTransport := http.DefaultTransport.(*http.Transport)

	// Create new Transport that ignores self-signed SSL
	customTransport := &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: customTransport}
	request.Header.Set("Authorization", os.Getenv("PVE_TOKEN"))
	request.Header.Set("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	type proxmoxData struct {
		Data []ResponseType
	}

	res := proxmoxData{}
	json.Unmarshal(bodyBytes, &res)
	return res.Data
}
