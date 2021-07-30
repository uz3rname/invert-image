package main

import (
  "log"
  "os"

  swagger "github.com/arsmn/fiber-swagger/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
  "github.com/gofiber/fiber/v2"
  _ "github.com/joho/godotenv/autoload"
  apiController "github.com/uz3rname/invert-image/api"
  _ "github.com/uz3rname/invert-image/docs"
)

func createApp() *fiber.App {
  app := fiber.New()

  app.Use(logger.New())

  app.Static("/", "./public")

  apiRouter := app.Group("/api")

  apiRouter.Post("/negative_image", apiController.NegativeImage)
  apiRouter.Get("/get_last_images", apiController.GetLastImages)

  apiRouter.Get("/docs/*", swagger.Handler)
  apiRouter.Get("/docs/*", swagger.New(swagger.Config{}))

  return app
}

// @title Backend test task
// @version 1.0
func main() {
  app := createApp()
  err := app.Listen(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
  if err != nil {
    log.Fatalf(err.Error())
  }
}
