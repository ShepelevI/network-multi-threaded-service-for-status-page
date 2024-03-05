package billing

import (
	"errors"
	"io"
	"log"
	"os"
)

const ByteMaskLength = 6

func readFile(path string) ([]byte, error) {
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
	data := make([]byte, 6)
	for {
		_, err := file.Read(data)
		if err == io.EOF {
			break
		}
	}
	return data, nil
}

func GetBillingData(path string) (BillingData, error) {
	buff, err := readFile(path)
	if err != nil {
		log.Println(err)
		return BillingData{}, err
	}
	if len(buff) != ByteMaskLength {
		log.Printf("Bad byte mask")
		return BillingData{}, errors.New("bad byte mask")
	}
	return BillingData{
		CreateCustomer: buff[5] == '1',
		Purchase:       buff[4] == '1',
		Payout:         buff[3] == '1',
		Recurring:      buff[2] == '1',
		FraudControl:   buff[1] == '1',
		CheckoutPage:   buff[0] == '1',
	}, nil
}
