package repository

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/godprogrammer3/lmwn-dev-test/entity"
)

type CovidRepository interface {
	GetCovidSummary() (entity.CovidSummary, error)
	getCovidStat() (entity.CovidPublicData, error)
}

type HttpClienter interface {
	Do(req *http.Request) (*http.Response, error)
}

type covidRepository struct {
	publicCovidStatURL string
	client             HttpClienter
}

func NewCovidRepository(client HttpClienter, publicCovidStatURL string) CovidRepository {
	return &covidRepository{
		publicCovidStatURL: publicCovidStatURL,
		client:             client,
	}
}

func (covidRepository *covidRepository) GetCovidSummary() (entity.CovidSummary, error) {
	covidPublicData, err := covidRepository.getCovidStat()
	if err != nil {
		return entity.CovidSummary{}, err
	}
	covidSummary := entity.CovidSummary{Province: make(map[string]int)}
	for _, data := range covidPublicData.Data {

		if data.ProvinceEn == "" {
			data.ProvinceEn = "N/A"
		}
		covidSummary.Province[data.ProvinceEn]++
		if data.Age == nil {
			covidSummary.AgeGroup.AgeNA++
		} else if *data.Age <= 30 {
			covidSummary.AgeGroup.Age0to30++
		} else if *data.Age <= 60 {
			covidSummary.AgeGroup.Age31to60++
		} else {
			covidSummary.AgeGroup.Age61plus++
		}
	}
	return covidSummary, nil
}

func (covidRepository *covidRepository) getCovidStat() (entity.CovidPublicData, error) {
	request, err := http.NewRequest(http.MethodGet, covidRepository.publicCovidStatURL, nil)
	if err != nil {
		print("Can not make request to public api: ")
		println(err.Error())
		return entity.CovidPublicData{}, errors.New("can not make request to public api")
	}
	response, err := covidRepository.client.Do(request)
	if err != nil {
		print("Can not get covid stats from public api: ")
		println(err.Error())
		return entity.CovidPublicData{}, errors.New("can not get covid stats from public api")
	}

	var covidPublicData entity.CovidPublicData

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&covidPublicData)
	if err != nil {
		print("Can not parse covid stats from public api: ")
		println(err.Error())
		return entity.CovidPublicData{}, errors.New("can not parse covid stats from public api")
	}

	return covidPublicData, nil
}
