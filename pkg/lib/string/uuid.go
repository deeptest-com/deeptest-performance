package _stringUtils

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func Uuid() string {
	uid := uuid.NewV4().String()
	return strings.Replace(uid, "-", "", -1)
}

func UuidWithSep() string {
	uid := uuid.NewV4().String()
	return uid
}

func Ulid() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	rand, _ := ulid.New(ms, entropy)

	ret := strings.ToLower(rand.String())
	ret = strings.Replace(ret, "-", "", -1)

	return ret
}
