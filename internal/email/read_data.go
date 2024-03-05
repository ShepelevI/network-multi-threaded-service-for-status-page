package email

import (
	"bufio"
	"errors"
	"log"
	"os"
	assert "skillbox_final/internal/assertions"
	"strconv"
	"strings"
)

const CSVRowLength = 3

func isValidEmail(data EmailData) bool {
	if assert.Alpha2Map[data.Country] == "" {
		return false
	} else if !assert.CheckValueInArray(data.Provider, assert.EmailProviders[:]) {
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

func getEmailData(params []string) (EmailData, error) {
	DeliveryTime, err := strconv.Atoi(params[2])
	if err != nil {
		return EmailData{}, errors.New("bad data format")
	}
	return EmailData{
		Country:      params[0],
		Provider:     params[1],
		DeliveryTime: DeliveryTime,
	}, nil

}

func GetEmailDataSlice(path string) ([]EmailData, error) {
	rows, err := readFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result []EmailData
	for _, row := range rows {
		params := strings.Split(row, ";")
		if len(params) != CSVRowLength {
			continue
		}
		email, err := getEmailData(params)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		if !isValidEmail(email) {
			continue
		}
		result = append(result, email)
	}
	return result, nil
}
