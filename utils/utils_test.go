package utils_test

import (
	"JD_Purchase/utils"
	"log"
	"testing"
)

func TestParseAreaID(t *testing.T) {
	ret := utils.ParseAreaID("28_2487_")
	log.Println(ret)
}

func TestParseSkuID(t *testing.T) {
	ret := utils.ParseSkuID("123456:2,123789")
	log.Println(ret)
}
