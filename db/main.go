package db

import (
  "crypto/md5"
  "encoding/base64"
  "encoding/hex"
  "fmt"
  "log"
  "os"

  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

var DB *gorm.DB

func init() {
  dsn := fmt.Sprintf(
    "host=%s port=%s dbname=%s user=%s password=%s",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_DBNAME"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASS"),
  )
  var err error

  DB, err = gorm.Open(postgres.Open(dsn))
  if err != nil {
    log.Fatalln("Failed to connect to database", err)
  }
  log.Printf("Connected to database")

  DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
  DB.AutoMigrate(&ImagePair{})
}

func GetLastImages(n int) []ImagePair {
  pairs := make([]ImagePair, n)
  DB.Order("created_at desc").Limit(n).Find(&pairs)

  return pairs
}

func AddImage(original, negative []byte, origMime, negMime string) *ImagePair {
  hash := md5.Sum(original)

  pair := ImagePair{
    Original: base64.RawStdEncoding.EncodeToString(original[:]),
    Negative: base64.RawStdEncoding.EncodeToString(negative[:]),
    Hash: hex.EncodeToString(hash[:]),
    OrigMime: origMime,
    NegMime: negMime,
  }

  DB.Create(&pair)

  return &pair
}

func FindImageByHash(hash string) *ImagePair {
  var pairs []ImagePair
  var count int64

  DB.Where("hash = ?", hash).Limit(1).Find(&pairs).Count(&count)

  if count > 0 {
    return &pairs[0]
  }
  return nil
}
