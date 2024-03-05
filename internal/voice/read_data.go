package voice

import (
	"bufio"
	"errors"
	"log"
	"os"
	assert "skillbox_final/internal/assertions"
	"strconv"
	"strings"
)

const CSVRowLength = 8

func isValidVoice(data VoiceData) bool {
	if assert.Alpha2Map[data.Country] == "" {
		return false
	} else if !assert.CheckValueInArray(data.Provider, assert.VoiceProviders[:]) {
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

func getVoiceData(params []string) (VoiceData, error) {
	ConnectionStability, err := strconv.ParseFloat(params[4], 32)
	if err != nil {
		return VoiceData{}, errors.New("bad data format")
	}
	TTFB, err := strconv.Atoi(params[5])
	if err != nil {
		return VoiceData{}, errors.New("bad data format")
	}
	VoicePurity, err := strconv.Atoi(params[6])
	if err != nil {
		return VoiceData{}, errors.New("bad data format")
	}
	MedianOfCallsTime, err := strconv.Atoi(params[7])
	if err != nil {
		return VoiceData{}, errors.New("bad data format")
	}
	return VoiceData{
		Country:             params[0],
		Bandwidth:           params[1],
		ResponseTime:        params[2],
		Provider:            params[3],
		ConnectionStability: float32(ConnectionStability),
		TTFB:                TTFB,
		VoicePurity:         VoicePurity,
		MedianOfCallsTime:   MedianOfCallsTime,
	}, nil

}

func GetVoiceDataSlice(path string) ([]VoiceData, error) {
	rows, err := readFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result []VoiceData
	for _, row := range rows {
		params := strings.Split(row, ";")
		if len(params) != CSVRowLength {
			continue
		}
		voice, err := getVoiceData(params)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		if !isValidVoice(voice) {
			continue
		}
		result = append(result, voice)
	}
	return result, nil
}
