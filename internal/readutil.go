package internal

import (
	"encoding/json"
	"io/ioutil"
)

func ReadMockPositionData() Positions {
	file, _ := ioutil.ReadFile("internal/mockdata/mockpositions.json")
	data := Positions{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func ReadMockPriceData() Prices {
	file, _ := ioutil.ReadFile("internal/mockdata/mockprices.json")
	data := Prices{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}
