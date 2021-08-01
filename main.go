package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/uz3rname/invert-image/api"
	"github.com/uz3rname/invert-image/db"
	"github.com/uz3rname/invert-image/image"
	"github.com/uz3rname/invert-image/task"
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

  taskStore := task.NewMemoryTaskStore()
  manager := task.NewTaskManager(taskStore)
  manager.RegisterHandler(
    "invert-image",
    &image.InvertImageTaskHandler{Store: store},
  )
  numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
  if err != nil {
    numWorkers = runtime.NumCPU()
  }
  manager.StartWorkers(numWorkers)

  app := fiber.New(fiber.Config{
    BodyLimit: 1024 * 1024 * 1024,
  })
  app.Use(logger.New())
  app.Mount("/api", api.CreateApp(&api.Options{
    Store: store,
    TaskManager: manager,
  }))
  app.Static("/", "./public")

  err = app.Listen(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
  if err != nil {
    log.Fatalf(err.Error())
  }
}
