package service

import (
	"log"
	"math/rand"
	"time"
	config "turbo-converter/config"

	hashids "github.com/speps/go-hashids"
)

var hashGen *hashids.HashID

func InitHashGen() {
	hd := hashids.NewData()
	hd.Salt = config.GlobalConfig().Salt
	hg, err := hashids.NewWithData(hd)

	if err != nil {
		log.Fatalf("Cannot create hashids [%s]", err.Error())
		panic("Cannot create hashids!")
	}

	rand.Seed(time.Now().UTC().UnixNano())
	hashGen = hg
}

func GetHash() string {

	year := time.Now().Year() % 100
	yearDay := time.Now().YearDay()
	nano := time.Now().Nanosecond() % 10000
	randomNum := rand.Intn(99)

	resultHash, _ := hashGen.Encode([]int{year, yearDay, nano, randomNum})

	return resultHash
}
