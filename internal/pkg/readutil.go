package internal

import (
	mk "StakeBackendGoTest/internal/entity/mock"

	"encoding/json"
	"io/ioutil"
)

func ReadMockPositionData() mk.Positions {
	file, _ := ioutil.ReadFile("internal/mockdata/mockpositions.json")
	data := mk.Positions{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func ReadMockPriceData() mk.Prices {
	file, _ := ioutil.ReadFile("internal/mockdata/mockprices.json")
	data := mk.Prices{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}
