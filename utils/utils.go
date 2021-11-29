package utils

import (
	"math/rand"
	"strconv"
)

func GenerateNumber() string {
	control := strconv.Itoa(rand.Intn(99))          //контрольный номер
	nuniq := strconv.Itoa(rand.Intn(9999999999999)) // уникальный номер счета
	return "KZ" + control + "777" + nuniq           //777 - код банка
}
