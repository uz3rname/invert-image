package image

import (
  "bytes"
  "errors"
  "image"

  "github.com/disintegration/imaging"
)

func decodeImage(data []byte) (*image.Image, ImageCodec, error) {
  var img image.Image

  for _, codec := range imageCodecs {
    reader := bytes.NewReader(data)
    if codec.Decode(reader, &img) {
      return &img, codec, nil
    }
  }

  return nil, nil, errors.New("Unsupported image format")
}

func GetInvertedImage(data []byte, encodeWith ImageCodec) ([]byte, error, string) {
  img, codec, err := decodeImage(data[:])
  if err != nil {
    return nil, err, ""
  }

  neg := imaging.Invert(*img)
  var buf bytes.Buffer
  encodeWith.Encode(&buf, neg)

  return buf.Bytes(), nil, codec.Mimetype()
}
