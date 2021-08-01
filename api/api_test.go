package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/uz3rname/invert-image/db"
)

func Test_negativeImage(t *testing.T) {
  store := db.NewMemoryStore()
  app := CreateApp(&Options{
    Store: store,
  })

  data, err := ioutil.ReadFile("../test/image-pair.json")
  if err != nil {
    panic("Error reading file")
  }
  var dto, responseDto InvertSuccessDTO
  json.Unmarshal(data, &dto)

  t.Run("Post image", func (t *testing.T) {
    apitest.New().
      HandlerFunc(FiberToHandlerFunc(app)).
      Post("/negative_image").
      JSON(fiber.Map{
        "data": dto.Pair.Original.Base64,
      }).
      Expect(t).
      Status(http.StatusOK).
      Assert(func (res *http.Response, req *http.Request) error {
        data := make([]byte, res.ContentLength)
        res.Body.Read(data)
        json.Unmarshal(data, &responseDto)
        return nil
      }).
      Assert(jsonpath.Equal(`$.pair.original.base64`, dto.Pair.Original.Base64)).
      Assert(jsonpath.Equal(`$.pair.negative.base64`, dto.Pair.Negative.Base64)).
      End()
  })

  t.Run("Post same image", func (t *testing.T) {
    apitest.New().
      HandlerFunc(FiberToHandlerFunc(app)).
      Post("/negative_image").
      JSON(fiber.Map{
        "data": dto.Pair.Original.Base64,
      }).
      Expect(t).
      Status(http.StatusOK).
      Assert(jsonpath.Equal(`$.pair.id`, responseDto.Pair.ID)).
      End()
  })

  t.Run("Get last images", func (t *testing.T) {
    apitest.New().
      HandlerFunc(FiberToHandlerFunc(app)).
      Get("/get_last_images").
      Expect(t).
      Status(http.StatusOK).
      Assert(jsonpath.Equal(`$.status`, "ok")).
      Assert(jsonpath.Equal(`$.items[0].id`, responseDto.Pair.ID)).
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
