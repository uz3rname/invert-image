package db

import (
  "log"

  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
)

type PostgresStore struct {
  db *gorm.DB
}

func NewPostgresStore(dsn string) Store {
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Error),
  })
  if err != nil {
    log.Fatalln("Failed to connect to database", err)
  }

  log.Printf("Connected to database")

  db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
  db.AutoMigrate(&ImagePair{})

  return &PostgresStore{db}
}

func (s *PostgresStore) GetLastImages(n int) []*ImagePair {
  pairs := make([]*ImagePair, n)
  s.db.Order("created_at desc").Limit(n).Find(&pairs)

  return pairs
}

func (s *PostgresStore) AddImage(
  original, negative, origMime, negMime, hash string,
) *ImagePair {
  pair := ImagePair{
    Original: original,
    Negative: negative,
    Hash: hash,
    OrigMime: origMime,
    NegMime: negMime,
  }

  s.db.Create(&pair)

  return &pair
}

func (s *PostgresStore) FindImageByHash(hash string) (*ImagePair, bool) {
  var pairs []ImagePair
  var count int64

  s.db.Where("hash = ?", hash).Limit(1).Find(&pairs).Count(&count)

  if count > 0 {
    return &pairs[0], true
  }
  return nil, false
}
