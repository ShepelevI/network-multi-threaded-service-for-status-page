package mms

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	assert "skillbox_final/internal/assertions"
)

func isValidMMS(data MMSData) bool {

	if assert.Alpha2Map[data.Country] == "" {
		return false
	} else if !assert.CheckValueInArray(data.Provider, assert.Providers[:]) {
		return false
	}
	return true
}

func GetMMS(path string) ([]MMSData, error) {
	resp, err := http.Get(path)
	if err != nil {
		log.Println(err)
		log.Println("http-get response send error", path)
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Println("status code is not 200")
		return nil, errors.New("status code is not 200")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	return getStructMMS(body), nil

}

func getStructMMS(body []byte) []MMSData {
	var list []MMSData
	err := json.Unmarshal(body, &list)
	if err != nil {
		log.Println(err)
	}
	for i, v := range list {
		if !isValidMMS(v) {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}
