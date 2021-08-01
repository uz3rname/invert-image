package image

import (
  "bytes"
  "context"

  "github.com/disintegration/imaging"
  "github.com/uz3rname/invert-image/db"
  "github.com/uz3rname/invert-image/task"
)

type InvertImageTaskHandler struct {
  Store db.Store
}

type InvertImageTaskInput struct {
  Data []byte
}

type InvertImageTaskResult struct {
  Pair *db.ImagePair
}

func (h *InvertImageTaskHandler) Run(
  ctx context.Context,
  input task.TaskInput,
) (task.TaskResult, error) {
  var result InvertImageTaskResult

  data := input.(*InvertImageTaskInput)
  img, codec, err := decodeImage(data.Data[:])
  if err != nil {
    return &result, err
  }

  negImg := imaging.Invert(*img)
  var buf bytes.Buffer
  err = codec.Encode(&buf, negImg)
  if err != nil {
    return &result, err
  }
  negData := buf.Bytes()

  mimetype := codec.Mimetype()
  result.Pair = h.Store.AddImage(
    Serialize(data.Data[:]),
    Serialize(negData[:]),
    mimetype,
    mimetype,
    Hash(data.Data[:]),
  )

  return &result, nil
}
