package util
import (
	"math/rand"
	"time"
)
const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
// GenerateRandomCode membuat string acak dengan panjang tertentu
func GenerateRandomCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}