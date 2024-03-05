package sms

type SMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

func (s *SMSData) SetCountry(new string) {
	s.Country = new
}
func (s *SMSData) GetCountry() string {
	return s.Country
}
func (s *SMSData) GetProvider() string {
	return s.Provider
}
