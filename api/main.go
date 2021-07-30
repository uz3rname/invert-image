package api

import (
  "crypto/md5"
  "encoding/base64"
  "encoding/hex"
  "log"
  "os"
  "strconv"

  "github.com/go-playground/validator/v10"
  "github.com/gofiber/fiber/v2"
  "github.com/uz3rname/invert-image/db"
  "github.com/uz3rname/invert-image/image"
)

var App *fiber.App
var validate *validator.Validate
var logger *log.Logger

func init() {
  validate = validator.New()
  logger = log.New(os.Stdout, "| ", log.Ltime | log.Lmsgprefix)
}

// NegativeImage
// @Description Create negative of image
// @Summary creates inversion of image
// @Tags images
// @Accept json
// @Produce json
// @Param image body UploadImageDTO true "Base64 encoded image"
// @Success 200 {object} InvertSuccessDTO
// @Success 400 {object} ErrorDTO
// @Router /api/negative_image [post]
func NegativeImage(ctx *fiber.Ctx) error {
  var dto UploadImageDTO

  err := ctx.BodyParser(&dto)
  if err != nil {
    return ctx.Status(400).JSON(makeError("Invalid request: %s", err))
  }

  err = validate.Struct(&dto)
  if err != nil {
    return ctx.Status(400).JSON(makeError("Validation error: %s", err))
  }

  data, err := base64.RawStdEncoding.DecodeString(dto.Data)
  if err != nil {
    return ctx.Status(400).JSON(
      makeError("Failed to decode base64 image: %s", err),
    )
  }

  hash := md5.Sum(data[:])

  if pair := db.FindImageByHash(hex.EncodeToString(hash[:])); pair != nil {
    logger.Printf(
      "Found already processed image (MD5 sum: \"%s\"), skipping",
      pair.Hash,
    )
    return ctx.JSON(&InvertSuccessDTO{
      Status: "ok",
      Pair: serializeImagePair(pair),
    })
  }

  encoder := image.PNGImageCodec{};
  neg, err, mimetype := image.GetInvertedImage(data[:], &encoder)
  if err != nil {
    return ctx.Status(400).JSON(makeError("Unsupported image format"))
  }

  pair := db.AddImage(data[:], neg[:], mimetype, encoder.Mimetype())
  logger.Printf("Added image \"%s\", MD5 sum: \"%s\"", pair.ID, pair.Hash)

  return ctx.JSON(&InvertSuccessDTO{
    Status: "ok",
    Pair: serializeImagePair(pair),
  })
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
// @Router /api/get_last_images [get]
func GetLastImages(ctx *fiber.Ctx) error {
  count, err := strconv.Atoi(ctx.Query("count", "3"))
  if err != nil {
    return ctx.JSON(makeError("Invalid count parameter"))
  }

  dbPairs := db.GetLastImages(count)
  var pairs []*ImagePair

  for _, pair := range dbPairs {
    pairs = append(pairs, serializeImagePair(&pair))
  }

  return ctx.JSON(ImagePairListDTO{
    Status: "ok",
    Items: pairs,
    Count: len(dbPairs),
  })
}
