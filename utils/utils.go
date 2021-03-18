package utils

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func SaveImage(resp *http.Response, imagePath string) error {
	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func ParseSkuID(skuIDs string) map[string]string {
	skuIDs = strings.TrimSpace(skuIDs)
	skuIDList := strings.Split(skuIDs, ",")
	ret := make(map[string]string)
	for _, v := range skuIDList {
		if strings.Contains(v, ":") {
			skuIDsSplit := strings.Split(v, ":")
			ret[skuIDsSplit[0]] = skuIDsSplit[1]
			continue
		}
		ret[v] = "1"
	}
	return ret
}

func ParseAreaID(area string) string {
	area = strings.TrimSpace(area)
	areas := strings.FieldsFunc(area, func(r rune) bool {
		return r == '-' || r == '_'
	})
	lenA := len(areas)
	if lenA < 4 {
		for i := 0; i < 4-lenA; i++ {
			areas = append(areas, "0")
		}
	}
	return strings.Join(areas, "_")
}
