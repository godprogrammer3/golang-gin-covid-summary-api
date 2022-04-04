package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godprogrammer3/lmwn-dev-test/api"
	"github.com/godprogrammer3/lmwn-dev-test/controller"
	"github.com/godprogrammer3/lmwn-dev-test/docs"
	"github.com/godprogrammer3/lmwn-dev-test/repository"
	"github.com/godprogrammer3/lmwn-dev-test/service"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	intialConfig()
	r := setupRouter()
	port, ok := viper.Get("APP_PORT").(string)
	if !ok {
		port = "3000"
	}
	r.Run(fmt.Sprintf(":" + port))
}

func intialConfig() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
}
func setupRouter() *gin.Engine {
	covidStatsURL, ok := viper.Get("COVID_STATS_URL").(string)
	if !ok {
		panic("can not get COVID_STATS_URL env")
	}
	client := &http.Client{}
	var (
		covidRepository repository.CovidRepository = repository.NewCovidRepository(client, covidStatsURL)
		covidService    service.CovidService       = service.New(covidRepository)
		covidController controller.CovidController = controller.New(covidService)
	)

	docs.SwaggerInfo.Title = "Covid Summary API"
	docs.SwaggerInfo.Description = "This API use to get summary of covid case"
	docs.SwaggerInfo.Version = "1.0"
	port, ok := viper.Get("APP_PORT").(string)
	if !ok {
		port = "3000"
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	server := gin.Default()

	covidApi := api.NewCovidAPI(covidController)

	apiRoutes := server.Group(docs.SwaggerInfo.BasePath)
	{
		covid := apiRoutes.Group("/covid")
		{
			covid.GET("/summary", covidApi.GetCovidSummary)
		}
	}
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return server
}
