package services

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func GenerateULID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
