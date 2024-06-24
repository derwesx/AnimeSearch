package anime

import (
	"AnimeSearch/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type DBMaintainer struct {
	animeCount int
	db         *gorm.DB
	rdb        *redis.Client
}

func generateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}

func generateRandomHash(length int) string {
	randomBytes, err := generateRandomBytes(length) // 32 bytes = 256 bits
	if err != nil {
		return ""
	}

	hasher := sha256.New()
	hasher.Write(randomBytes)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)[:length]
}

func hashToInt(s string) uint64 {
	h := sha256.New()
	h.Write([]byte(s))
	hashBytes := h.Sum(nil)
	return binary.BigEndian.Uint64(hashBytes)
}

func hashToID(hash string, count int) int {
	id := hashToInt(hash) % uint64(count)
	return int(id)
}

func (handler *DBMaintainer) update() {
	for {
		var counter int
		handler.db.Model(&models.Anime{}).Select("count(*)").Scan(&counter)
		handler.animeCount = counter
		fmt.Println("Updating...\nCurrent anime count: ", handler.animeCount)
		time.Sleep(60 * time.Second)
	}
}

func (handler *DBMaintainer) getAnimeCount() int {
	return handler.animeCount
}

func NewMaintainer(db *gorm.DB, rdb *redis.Client) *DBMaintainer {
	maintainer := &DBMaintainer{
		db:  db,
		rdb: rdb,
	}
	go maintainer.update()
	return maintainer
}
