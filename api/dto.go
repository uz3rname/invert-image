package api

import (
  "fmt"
  "time"

  "github.com/uz3rname/invert-image/db"
)

type ErrorDTO struct {
  Status      string          `json:"status" enums:"error"`
  Message     string          `json:"message"`
}

type EncodedImage struct {
  Base64      string          `json:"base64"`
  MimeType    string          `json:"mimeType"`
}

type ImagePair struct {
  ID          string          `json:"id"`
  Original    EncodedImage    `json:"original"`
  Negative    EncodedImage    `json:"negative"`
  Hash        string          `json:"hash"`
  CreatedAt   time.Time       `json:"createdAt"`
}

type ImagePairListDTO struct {
  Status      string          `json:"status" enums:"ok"`
  Items       []*ImagePair    `json:"items"`
  Count       int             `json:"count"`
}

type UploadImageDTO struct {
  Data        string          `json:"data" validate:"required"`
}

type InvertSuccessDTO struct {
  Status      string          `json:"status" enums:"ok"`
  Pair        *ImagePair      `json:"pair"`
}

func makeError(msg string, err ...error) ErrorDTO {
  var strings []interface{}

  for _, e := range err {
    strings = append(strings, e.Error())
  }
  return ErrorDTO{"error", fmt.Sprintf(msg, strings...)}
}

func serializeImagePair(pair *db.ImagePair) *ImagePair {
  return &ImagePair{
    ID: pair.ID,
    Original: EncodedImage{
      Base64: pair.Original,
      MimeType: pair.OrigMime,
    },
    Negative: EncodedImage{
      Base64: pair.Negative,
      MimeType: pair.NegMime,
    },
    CreatedAt: pair.CreatedAt,
    Hash: pair.Hash,
  }
}
