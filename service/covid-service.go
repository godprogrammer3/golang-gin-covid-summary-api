package service

import (
	"github.com/godprogrammer3/lmwn-dev-test/entity"
	"github.com/godprogrammer3/lmwn-dev-test/repository"
)

type CovidService interface {
	GetSummary() (entity.CovidSummary, error)
}

type covidService struct {
	repository repository.CovidRepository
}

func New(covidRepository repository.CovidRepository) CovidService {
	return &covidService{
		repository: covidRepository,
	}
}

func (covidService *covidService) GetSummary() (entity.CovidSummary, error) {
	return covidService.repository.GetCovidSummary()
}
