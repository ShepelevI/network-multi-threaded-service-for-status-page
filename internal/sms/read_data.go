package sms

import (
	"bufio"
	"log"
	"os"
	assert "skillbox_final/internal/assertions"
	"strings"
)

const CSVRowLength = 4

func isValidSMS(data SMSData) bool {

	if assert.Alpha2Map[data.Country] == "" {
		return false
	} else if !assert.CheckValueInArray(data.Provider, assert.Providers[:]) {
		return false
	}
	return true
}

func readFile(path string) ([]string, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("File does not exist")
			return nil, err
		}
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}(file)

	var rows []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		rows = append(rows, sc.Text())
	}
	return rows, nil

}

func GetSMSDataSlice(path string) ([]SMSData, error) {
	rows, err := readFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result []SMSData
	for _, row := range rows {
		params := strings.Split(row, ";")
		if len(params) != CSVRowLength {
			continue
		}
		sms := SMSData{
			Country:      params[0],
			Bandwidth:    params[1],
			ResponseTime: params[2],
			Provider:     params[3],
		}
		if !isValidSMS(sms) {
			continue
		}
		result = append(result, sms)
	}
	return result, nil
}
