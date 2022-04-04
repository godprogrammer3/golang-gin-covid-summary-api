package entity

type CovidSummary struct {
	Province map[string]int
	AgeGroup Group
}

type Group struct {
	Age0to30  int `json:"0-30"`
	Age31to60 int `json:"31-60"`
	Age61plus int `json:"61+"`
	AgeNA     int `json:"N/A"`
}

type CovidPublicData struct {
	Data []Data `json:"Data"`
}

type Data struct {
	Age        *int   `json:"Age"`
	ProvinceEn string `json:"ProvinceEn"`
}
