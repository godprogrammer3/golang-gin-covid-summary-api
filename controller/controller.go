package controller

import (
	"github.com/godprogrammer3/lmwn-dev-test/entity"
	"github.com/godprogrammer3/lmwn-dev-test/service"
)

type CovidController interface {
	GetSummary() (entity.CovidSummary, error)
}

type controller struct {
	service service.CovidService
}

func New(service service.CovidService) CovidController {
	return &controller{
		service: service,
	}
}

func (controller *controller) GetSummary() (entity.CovidSummary, error) {
	return controller.service.GetSummary()
}
