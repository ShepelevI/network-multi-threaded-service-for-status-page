package json

import (
	"path/filepath"
	assert "skillbox_final/internal/assertions"
	"skillbox_final/internal/billing"
	"skillbox_final/internal/email"
	incidentData "skillbox_final/internal/incidents"
	"skillbox_final/internal/mms"
	"skillbox_final/internal/sms"
	"skillbox_final/internal/support"
	"skillbox_final/internal/voice"
	"time"
)

const speedSupport = 18
const dataPath = "simulator"

var BufferedDataT ResultSetT

type ResultSetT struct {
	SMS       [][]sms.SMSData                `json:"sms"`
	MMS       [][]mms.MMSData                `json:"mms"`
	VoiceCall []voice.VoiceData              `json:"voice_call"`
	Email     map[string][][]email.EmailData `json:"email"`
	Billing   billing.BillingData            `json:"billing"`
	Support   []int                          `json:"support"`
	Incidents []incidentData.IncidentData    `json:"incident"`
}

type ResultT struct {
	Status bool       `json:"status"`
	Data   ResultSetT `json:"data"`
	Error  string     `json:"error"`
}
type ResultTErr struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

const gate = "http://127.0.0.1:8383"

type messageData interface {
	SetCountry(new string)
	GetCountry() string
	GetProvider() string
}

func sortMessageByDelivery(list []email.EmailData) []email.EmailData {
	length := len(list)
	for i := 0; i < (length - 1); i++ {
		for j := 0; j < ((length - 1) - i); j++ {
			if list[j].DeliveryTime > list[j+1].DeliveryTime {
				list[j], list[j+1] = list[j+1], list[j]
			}
		}
	}
	return list
}

func sortMessageByProvider(list []messageData) []messageData {
	length := len(list)
	for i := 0; i < (length - 1); i++ {
		for j := 0; j < ((length - 1) - i); j++ {
			if list[j].GetProvider() > list[j+1].GetProvider() {
				list[j], list[j+1] = list[j+1], list[j]
			}
		}
	}
	return list
}

func sortMessageByCountry(list []messageData) []messageData {
	length := len(list)
	for i := 0; i < (length - 1); i++ {
		for j := length - 1; j > i; j-- {
			if list[j].GetCountry() < list[j-1].GetCountry() {
				list[j], list[j-1] = list[j-1], list[j]
			}
		}
	}
	return list
}

func convertCodeToNameSMS(list []sms.SMSData) []sms.SMSData {
	result := make([]sms.SMSData, 0, len(list))
	for _, data := range list {
		data.Country = assert.Alpha2Map[data.Country]
		result = append(result, data)
	}
	return result
}

func convertCodeToNameMMS(list []mms.MMSData) []mms.MMSData {
	result := make([]mms.MMSData, 0, len(list))
	for _, data := range list {
		data.Country = assert.Alpha2Map[data.Country]
		result = append(result, data)
	}
	return result
}

func getDataSMS(list []sms.SMSData) [][]sms.SMSData {
	resultSMS := convertCodeToNameSMS(list)
	generalMessages := make([]messageData, 0, len(resultSMS))
	for i := range resultSMS {
		generalMessages = append(generalMessages, &resultSMS[i])
	}
	sortedSMSByProvider := sortMessageByProvider(generalMessages)
	sliceSMS1 := make([]sms.SMSData, len(sortedSMSByProvider))
	for i, message := range sortedSMSByProvider {
		smsMessage, ok := message.(*sms.SMSData)
		if !ok {
		}
		sliceSMS1[i] = *smsMessage
	}
	sliceSMS2 := make([]sms.SMSData, len(sortedSMSByProvider))
	sortedSMSByCountry := sortMessageByCountry(generalMessages)
	for i, message := range sortedSMSByCountry {
		smsMessage, ok := message.(*sms.SMSData)
		if !ok {
		}
		sliceSMS2[i] = *smsMessage
	}
	return [][]sms.SMSData{
		sliceSMS1,
		sliceSMS2,
	}

}

func getDataMMS(list []mms.MMSData, ch chan [][]mms.MMSData) {
	resultMMS := convertCodeToNameMMS(list)
	generalMessages := make([]messageData, 0, len(resultMMS))
	for i := range resultMMS {
		generalMessages = append(generalMessages, &resultMMS[i])
	}
	sortedSMSByProvider := sortMessageByProvider(generalMessages)
	sliceMMS1 := make([]mms.MMSData, len(sortedSMSByProvider))
	for i, message := range sortedSMSByProvider {
		smsMessage, ok := message.(*mms.MMSData)
		if !ok {
		}
		sliceMMS1[i] = *smsMessage
	}
	sliceMMS2 := make([]mms.MMSData, len(sortedSMSByProvider))
	sortedSMSByCountry := sortMessageByCountry(generalMessages)
	for i, message := range sortedSMSByCountry {
		smsMessage, ok := message.(*mms.MMSData)
		if !ok {
		}
		sliceMMS2[i] = *smsMessage
	}

	ch <- [][]mms.MMSData{
		sliceMMS1,
		sliceMMS2,
	}

}

func getMinMaxEmail(list []email.EmailData) (result [][]email.EmailData) {
	resultSort := sortMessageByDelivery(list)
	min3 := make([]email.EmailData, 3)
	max3 := make([]email.EmailData, 3)
	min3 = resultSort[:3]
	max3 = resultSort[len(resultSort)-3:]
	result = [][]email.EmailData{
		min3,
		max3,
	}
	return result
}

func transformCountryNameForEmailData(m map[string][][]email.EmailData) map[string][][]email.EmailData {
	result := make(map[string][][]email.EmailData, len(m))
	for k, v := range m {
		result[assert.Alpha2Map[k]] = v
	}
	return result
}

func getEmailData(list []email.EmailData, ch chan map[string][][]email.EmailData) {
	result := make(map[string][][]email.EmailData)
	CountryProvider := make(map[string][]email.EmailData)
	for _, data := range list {
		CountryProvider[data.Country] = append(CountryProvider[data.Country], data)
	}
	for country, provider := range CountryProvider {
		sorted := getMinMaxEmail(provider)
		result[country] = sorted
	}
	ch <- transformCountryNameForEmailData(result)
}

func getSupportData(list []support.SupportData, ch chan []int) {
	result := make([]int, 2)
	hour := time.Now().Hour()
	if hour < 9 {
		result[0] = 1
	} else if hour > 16 {
		result[0] = 3
	} else {
		result[0] = 2
	}
	waitTime := func() int {
		result := 0
		for _, data := range list {
			result += data.ActiveTickets
		}
		return result * speedSupport
	}()
	result[1] = waitTime
	ch <- result

}

func getIncidentData(list []incidentData.IncidentData, ch chan []incidentData.IncidentData) {
	length := len(list)
	for i := 0; i < (length - 1); i++ {
		for j := 0; j < ((length - 1) - i); j++ {
			if list[j+1].Status == "active" && list[j].Status != "active" {
				list[j], list[j+1] = list[j+1], list[j]
			}
		}
	}
	ch <- list
}

func sortVoice(list []voice.VoiceData, ch chan []voice.VoiceData) {
	result := make([]voice.VoiceData, len(list))
	for i, data := range list {
		data.Country = assert.Alpha2Map[data.Country]
		result[i] = data
	}
	ch <- result
}

func GetResultData() (ResultSetT, error) {
	if BufferedDataT.Support != nil {
		return BufferedDataT, nil // Если в буфере есть значение, которое не успело очиститься клинером, то сразу отдаем его
	}
	resultSMS, errSMS := sms.GetSMSDataSlice(filepath.Join(dataPath, "sms.data"))
	if errSMS != nil {
		return ResultSetT{}, errSMS
	}
	resultMMS, errMMS := mms.GetMMS(gate + "/mms")
	if errMMS != nil {
		return ResultSetT{}, errSMS
	}
	resultVoice, errVoice := voice.GetVoiceDataSlice(filepath.Join(dataPath, "voice.data"))
	if errVoice != nil {
		return ResultSetT{}, errSMS
	}
	resultEmail, errEmail := email.GetEmailDataSlice(filepath.Join(dataPath, "email.data"))
	if errEmail != nil {
		return ResultSetT{}, errSMS
	}
	resultBilling, errBilling := billing.GetBillingData(filepath.Join(dataPath, "billing.data"))
	if errBilling != nil {
		return ResultSetT{}, errSMS
	}
	resultSupport, errSupport := support.GetSupport(gate + "/support")
	if errSupport != nil {
		return ResultSetT{}, errSMS
	}
	resultIncident, errIncident := incidentData.GetIncident(gate + "/accendent")
	if errIncident != nil {
		return ResultSetT{}, errSMS
	}
	dataSMS := getDataSMS(resultSMS)
	go func() {
		getDataSMS(resultSMS)
	}()
	chanDataVoice := make(chan []voice.VoiceData, 1)
	sortVoice(resultVoice, chanDataVoice)
	dataVoice := <-chanDataVoice

	chanDataMMS := make(chan [][]mms.MMSData, 1)
	getDataMMS(resultMMS, chanDataMMS)
	dataMMS := <-chanDataMMS

	chanDataEmail := make(chan map[string][][]email.EmailData, 1)
	getEmailData(resultEmail, chanDataEmail)
	dataEmail := <-chanDataEmail

	chanDataSupport := make(chan []int, 1)
	getSupportData(resultSupport, chanDataSupport)
	dataSupport := <-chanDataSupport

	chanDataIncident := make(chan []incidentData.IncidentData, 1)
	getIncidentData(resultIncident, chanDataIncident)
	dataIncidents := <-chanDataIncident
	data := ResultSetT{
		SMS:       dataSMS,
		MMS:       dataMMS,
		VoiceCall: dataVoice,
		Email:     dataEmail,
		Billing:   resultBilling,
		Support:   dataSupport,
		Incidents: dataIncidents,
	}
	BufferedDataT = data
	return data, nil

}

func StartBufferCleaner(second int) {
	go func() {
		for {
			select {
			case <-time.Tick(time.Second * time.Duration(second)):
				BufferedDataT = ResultSetT{}
			}
		}
	}()
}
