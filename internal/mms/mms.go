package mms

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

func (s *MMSData) SetCountry(new string) {
	s.Country = new
}
func (s *MMSData) GetCountry() string {
	return s.Country
}
func (s *MMSData) GetProvider() string {
	return s.Provider
}
