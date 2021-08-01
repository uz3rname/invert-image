package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/uz3rname/invert-image/api"
	"github.com/uz3rname/invert-image/db"
)

func main() {
  store := db.NewPostgresStore(fmt.Sprintf(
    "host=%s port=%s dbname=%s user=%s password=%s",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_DBNAME"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASS"),
  ))

  app := fiber.New()
  app.Use(logger.New())
  app.Mount("/api", api.CreateApp(&api.Options{
    Store: store,
  }))
  app.Static("/", "./public")

  err := app.Listen(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
  if err != nil {
    log.Fatalf(err.Error())
  }
}
