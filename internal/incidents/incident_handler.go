package incidents

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func GetIncident(path string) ([]IncidentData, error) {
	resp, err := http.Get(path)
	if err != nil {
		log.Println(err)
		log.Println("Http-get response send error", path)
		return []IncidentData{}, errors.New("http-get response send error")
	}
	if resp.StatusCode != 200 {
		log.Println("Status code is not 200")
		var list []IncidentData
		return list, errors.New("status code is not 200")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	return getIncidentStruct(body), nil

}

func getIncidentStruct(body []byte) []IncidentData {
	var list []IncidentData
	err := json.Unmarshal(body, &list)
	if err != nil {
		log.Println(err)
	}
	return list
}
