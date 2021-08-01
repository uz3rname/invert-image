package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/uz3rname/invert-image/db"
	"github.com/uz3rname/invert-image/image"
	"github.com/uz3rname/invert-image/task"
)

const (
  UUIDRegex = "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"
)

func loadJson(path string, into interface{}) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    panic("Error reading file \"" + path + "\"")
  }
  err = json.Unmarshal(data, into)
  if err != nil {
    panic("Invalid JSON \"" + path + "\"")
  }
}

func Test_negativeImage(t *testing.T) {
  store := db.NewMemoryStore()
  taskStore := task.NewMemoryTaskStore()
  manager := task.NewTaskManager(taskStore)
  manager.RegisterHandler(
    "invert-image",
    &image.InvertImageTaskHandler{Store: store},
  )
  manager.StartWorkers(1)

  app := CreateApp(&Options{
    Store: store,
    TaskManager: manager,
  })

  var smallFile, largeFile, responseDto InvertImageResponse
  loadJson("../test/image-pair.json", &smallFile)
  loadJson("../test/large-image.json", &largeFile)

  t.Run("Post image", func (t *testing.T) {
    apitest.New().
      HandlerFunc(FiberToHandlerFunc(app)).
      Post("/negative_image").
      JSON(fiber.Map{
        "data": smallFile.Pair.Original.Base64,
      }).
      Expect(t).
      Status(http.StatusOK).
      Assert(jsonpath.Equal(
        "$.pair.original.base64",
        smallFile.Pair.Original.Base64,
      )).
      Assert(jsonpath.Equal(
        "$.pair.negative.base64",
        smallFile.Pair.Negative.Base64,
      )).
      Assert(jsonpath.Equal(
        "$.pair.original.mimeType",
        smallFile.Pair.Original.MimeType,
      )).
      Assert(jsonpath.Equal(
        "$.pair.negative.mimeType",
        smallFile.Pair.Negative.MimeType,
      )).
      End().
      JSON(&responseDto)
  })

  t.Run("Post same image", func (t *testing.T) {
    apitest.New().
      HandlerFunc(FiberToHandlerFunc(app)).
      Post("/negative_image").
      JSON(fiber.Map{
        "data": smallFile.Pair.Original.Base64,
      }).
      Expect(t).
      Status(http.StatusOK).
      Assert(jsonpath.Equal("$.pair.id", responseDto.Pair.ID)).
      End()
  })

  t.Run("Get last images", func (t *testing.T) {
    apitest.New().
      HandlerFunc(FiberToHandlerFunc(app)).
      Get("/get_last_images").
      Expect(t).
      Status(http.StatusOK).
      Assert(jsonpath.Equal("$.status", "ok")).
      Assert(jsonpath.Equal("$.items[0].id", responseDto.Pair.ID)).
      End()
  })

  t.Run("Upload large image", func (t *testing.T) {
    apitest.New().
      HandlerFunc(FiberToHandlerFunc(app)).
      Post("/negative_image").
      JSON(fiber.Map{
        "data": largeFile.Pair.Original.Base64,
      }).
      Expect(t).
      Status(http.StatusOK).
      Assert(jsonpath.Equal("$.status", "defered")).
      Assert(jsonpath.Present("$.taskId")).
      Assert(jsonpath.Matches("$.taskId", UUIDRegex)).
      End().
      JSON(&responseDto)
  })

  t.Run("Waiting for task to complete", func (t *testing.T) {
    for {
      switch status, _ := manager.GetStatus(responseDto.TaskID); status {
      case task.TaskDone:
        return
      case task.TaskNew, task.TaskRunning:
        time.Sleep(time.Second)
      default:
        t.Fail()
        return
      }
    }
  })

  t.Run("Getting large file", func (t *testing.T) {
    apitest.New().
      HandlerFunc(FiberToHandlerFunc(app)).
      Get("/get_last_images").
      Expect(t).
      Status(http.StatusOK).
      Assert(jsonpath.Equal(
        "$.items[0].original.base64",
        largeFile.Pair.Original.Base64,
      )).
      Assert(jsonpath.Equal(
        "$.items[0].negative.base64",
        largeFile.Pair.Negative.Base64,
      )).
      Assert(jsonpath.Equal(
        "$.items[0].original.mimeType",
        largeFile.Pair.Original.MimeType,
      )).
      Assert(jsonpath.Equal(
        "$.items[0].negative.mimeType",
        largeFile.Pair.Negative.MimeType,
      )).
      End()
  })
}

func FiberToHandlerFunc(app *fiber.App) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    resp, err := app.Test(r)
    if err != nil {
      panic(err)
    }

    for k, vv := range resp.Header {
      for _, v := range vv {
        w.Header().Add(k, v)
      }
    }
    w.WriteHeader(resp.StatusCode)

    if _, err := io.Copy(w, resp.Body); err != nil {
      panic(err)
    }
  }
}
