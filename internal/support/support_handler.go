package support

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func GetSupport(path string) ([]SupportData, error) {
	resp, err := http.Get(path)
	if err != nil {
		log.Println(err)
		return []SupportData{}, err
	}
	if resp.StatusCode != 200 {
		log.Println("Status code is not 200")
		var list []SupportData
		return list, errors.New("status code is not 200")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	return getSupportStruct(body)

}

func getSupportStruct(body []byte) ([]SupportData, error) {
	var list []SupportData
	err := json.Unmarshal(body, &list)
	if err != nil {
		log.Println(err)
	}
	return list, err
}
