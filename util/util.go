package util

import (
	"math/rand"
	log "github.com/sirupsen/logrus"
	"unsafe"
	"time"
)

var logError = log.Error
// String A-Z and 0-9.
const strAtoZ0to9 = "ABDEFGHJKLMNPQRSTVWXY123456789"

// IsNilInterface returns if an interface contains nil type or value.
// See https://dev.to/pauljlucas/go-tcha-when-nil--nil-hic for detailed information.
func IsNilInterface(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}

// GenerateRandString returns a random string with given length.
func GenerateRandString(l int) string {
	str := "ABDEFGHJKLMNPQRSTVWXY123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

