package api

import (
	"github.com/gin-gonic/gin"
	"github.com/godprogrammer3/lmwn-dev-test/controller"
	"github.com/godprogrammer3/lmwn-dev-test/dto"
)

type CovidAPI struct {
	covidController controller.CovidController
}

func NewCovidAPI(covidController controller.CovidController) *CovidAPI {
	return &CovidAPI{
		covidController: covidController,
	}
}

// GetCovidSummary godoc
// @Summary Get covid summary
// @Description Get covid summary include case per province and case per group age
// @Tags covid
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.CovidSummary
// @Failure 500 {object} dto.Response
// @Router /covid/summary [get]
func (api *CovidAPI) GetCovidSummary(ctx *gin.Context) {
	res, err := api.covidController.GetSummary()
	if err == nil {
		ctx.JSON(200, res)
	} else {
		ctx.JSON(500, &dto.Response{
			Message: err.Error(),
		})
	}

}
