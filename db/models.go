package db

import "time"

type ImagePair struct {
  ID          string      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
  Original    string      `gorm:"type:text;notNull"`
  Negative    string      `gorm:"type:text;notNull"`
  CreatedAt   time.Time   `gorm:"type:timestamptz;index;default:now();notNull"`
  Hash        string      `gorm:"type:varchar(32);index;notNull"`
  OrigMime    string      `gorm:"type:varchar(32);notNull"`
  NegMime     string      `gorm:"type:varchar(32);notNull"`
}
