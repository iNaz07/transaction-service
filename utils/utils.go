package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// TODO: generate uniq number
func GenerateNumber() string {
	control := strconv.Itoa(10 + rand.Intn(89)) //контрольный номер
	randNum := RandomNumberGenerator()
	nuniq := strconv.Itoa(randNum)        // уникальный номер счета
	return "KZ" + control + "777" + nuniq //777 - код банка
}

func RandomNumberGenerator() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Int() / 1000000
}
