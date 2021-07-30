package image

import (
  "bytes"
  "image"
  "image/jpeg"
  "image/png"
)

type ImageCodec interface {
  Decode(*bytes.Reader, *image.Image) bool
  Encode(*bytes.Buffer, *image.NRGBA) error
  Mimetype() string
}

type PNGImageCodec struct {}
type JPEGImageCodec struct {}

func (d *PNGImageCodec) Decode(r *bytes.Reader, img *image.Image) bool {
  var err error

  *img, err = png.Decode(r)
  return err == nil
}

func (d *PNGImageCodec) Encode(buf *bytes.Buffer, img *image.NRGBA) error {
  return png.Encode(buf, img)
}

func (d *PNGImageCodec) Mimetype() string {
  return "image/png"
}

func (d *JPEGImageCodec) Decode(r *bytes.Reader, img *image.Image) bool {
  var err error

  *img, err = jpeg.Decode(r)
  return err == nil
}

func (d *JPEGImageCodec) Encode(buf *bytes.Buffer, img *image.NRGBA) error {
  return jpeg.Encode(buf, img, nil)
}

func (d *JPEGImageCodec) Mimetype() string {
  return "image/jpeg"
}

var imageCodecs = []ImageCodec{
  &PNGImageCodec{},
  &JPEGImageCodec{},
}
