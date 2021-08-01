package api

import (
	"log"
	"os"
	"strconv"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	_ "github.com/uz3rname/invert-image/api/docs"
	"github.com/uz3rname/invert-image/db"
	"github.com/uz3rname/invert-image/image"
	"github.com/uz3rname/invert-image/task"
)

type appState struct {
  validator *validator.Validate
  logger *log.Logger
  store db.Store
  manager *task.TaskManager
}

const (
  MaxFileSize = 1024 * 1024
)

// NegativeImage
// @Description Create negative of image
// @Summary creates inversion of image
// @Tags images
// @Accept json
// @Produce json
// @Param image body UploadImageDTO true "Base64 encoded image"
// @Success 200 {object} InvertImageResponse
// @Success 400 {object} ErrorDTO
// @Router /negative_image [post]
func negativeImage(state *appState) fiber.Handler {
  return func(ctx *fiber.Ctx) error {
    var dto UploadImageDTO

    err := ctx.BodyParser(&dto)
    if err != nil {
      return ctx.Status(400).JSON(makeError("Invalid request: %s", err))
    }

    err = state.validator.Struct(&dto)
    if err != nil {
      return ctx.Status(400).JSON(makeError("Validation error: %s", err))
    }

    data, err := image.Deserialize(dto.Data)
    if err != nil {
      return ctx.Status(400).JSON(
        makeError("Failed to decode base64 image: %s", err),
      )
    }

    hash := image.Hash(data[:])

    if pair, ok := state.store.FindImageByHash(hash); ok {
      state.logger.Printf(
        "Found already processed image (MD5 sum: \"%s\"), skipping",
        pair.Hash,
      )
      return ctx.JSON(&InvertImageResponse{
        Status: "ok",
        Pair: serializeImagePair(pair),
      })
    }

    taskInput := &image.InvertImageTaskInput{
      Data: data[:],
    }

    if len(data) > MaxFileSize {
      taskId, err := state.manager.PushTask("invert-image", taskInput)
      if err != nil {
        return ctx.Status(400).JSON(makeError("Couldn't start task"))
      }

      return ctx.JSON(&InvertImageResponse{
        Status: "defered",
        TaskID: taskId,
      })
    } else {
      result, err := state.manager.RunTask("invert-image", taskInput)
      if err != nil {
        return ctx.Status(400).JSON(makeError("Image inversion error: %s", err))
      }

      pair := result.(*image.InvertImageTaskResult).Pair

      state.logger.Printf(
        "Added image \"%s\", MD5 sum: \"%s\"",
        pair.ID,
        pair.Hash,
      )

      return ctx.JSON(&InvertImageResponse{
        Status: "ok",
        Pair: serializeImagePair(pair),
      })
    }
  }
}

// GetLastImages
// @Description Get last images
// @Summary returns last images from db
// @Tags images
// @Accept json
// @Produce json
// @Param count query int false "Number of images to return"
// @Success 200 {object} ImagePairListDTO
// @Success 400 {object} ErrorDTO
// @Router /get_last_images [get]
func getLastImages(state *appState) fiber.Handler {
  return func(ctx *fiber.Ctx) error {
    count, err := strconv.Atoi(ctx.Query("count", "3"))
    if err != nil {
      return ctx.JSON(makeError("Invalid count parameter"))
    }

    dbPairs := state.store.GetLastImages(count)
    var pairs []*ImagePair

    for _, pair := range dbPairs {
      pairs = append(pairs, serializeImagePair(pair))
    }

    return ctx.JSON(ImagePairListDTO{
      Status: "ok",
      Items: pairs,
      Count: len(dbPairs),
    })
  }
}

type Options struct {
  Store db.Store
  TaskManager *task.TaskManager
}

// @title Backend test task
// @version 1.0
// @BasePath /api
func CreateApp(options *Options) *fiber.App {
  app := fiber.New(fiber.Config{
    BodyLimit: 1024 * 1024 * 1024,
  })
  state := &appState{
    validator: validator.New(),
    logger: log.New(os.Stdout, "| ", log.Ltime | log.Lmsgprefix),
    store: options.Store,
    manager: options.TaskManager,
  }

  app.Get("/docs", swagger.Handler)
  app.Get("/docs/*", swagger.Handler)
  app.Get("/docs/*", swagger.New(swagger.Config{}))

  app.Post("/negative_image", negativeImage(state))
  app.Get("/get_last_images", getLastImages(state))

  return app
}
