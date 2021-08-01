package image

import (
  "crypto/md5"
  "encoding/base64"
  "encoding/hex"
)

func Hash(data []byte) string {
  hash := md5.Sum(data[:])
  return hex.EncodeToString(hash[:])
}

func Serialize(data []byte) string {
  return base64.RawStdEncoding.EncodeToString(data[:])
}

func Deserialize(s string) ([]byte, error) {
  data, err := base64.RawStdEncoding.DecodeString(s)
  return data, err
}
