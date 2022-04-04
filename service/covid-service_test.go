package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/godprogrammer3/lmwn-dev-test/entity"
	"github.com/godprogrammer3/lmwn-dev-test/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ClientMockResponse struct {
	response string
}

func (c *ClientMockResponse) Do(req *http.Request) (*http.Response, error) {
	mockResponse := c.response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockResponse)
	}))
	defer ts.Close()
	return http.Get(ts.URL)
}

type ClientMockStatus struct {
	status int
	err    error
}

func (c *ClientMockStatus) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: c.status}, c.err
}

const (
	WITHOUT_NULL_RESPONSE_FROM_API          = `{"Data":[{"ProvinceEn":"Phrae", "Age":0},{"ProvinceEn":"Suphan Buri", "Age":31}, {"ProvinceEn":"Roi Et", "Age":61} ]}`
	NULL_PROVINCE_REPONSE_FROM_API          = `{"Data":[{"ProvinceEn":"Phrae", "Age":0},{"ProvinceEn":"Suphan Buri", "Age":31}, {"ProvinceEn":"Roi Et", "Age":61}, {"ProvinceEn": null, "Age":61}, {"ProvinceEn": null, "Age":61} ]}`
	NULL_AGE_RESPONSE_FROM_API              = `{"Data":[{"ProvinceEn":"Phrae", "Age":0},{"ProvinceEn":"Suphan Buri", "Age":31}, {"ProvinceEn":"Roi Et", "Age":61}, {"ProvinceEn": "Roi Et", "Age": null}, {"ProvinceEn": "Roi Et", "Age":null} ]}`
	NULL_PROVINCE_AND_AGE_RESPONSE_FROM_API = `{"Data":[{"ProvinceEn":"Phrae", "Age":0},{"ProvinceEn":"Suphan Buri", "Age":31}, {"ProvinceEn":"Roi Et", "Age":61}, {"ProvinceEn": null, "Age":61}, {"ProvinceEn": null, "Age":61}, {"ProvinceEn": "Roi Et", "Age": null}, {"ProvinceEn": "Roi Et", "Age":null} ]}`
	OVER_61_AGE_RESPONSE_FROM_API           = `{"Data":[{"ProvinceEn":"Phrae", "Age":0},{"ProvinceEn":"Suphan Buri", "Age":31}, {"ProvinceEn":"Roi Et", "Age":61},  {"ProvinceEn":"Roi Et", "Age":62}, {"ProvinceEn":"Roi Et", "Age":62}]}`
)

var (
	EXPECTED_WITHOUT_NULL_RESPONSE          = entity.CovidSummary{Province: map[string]int{"Phrae": 1, "Suphan Buri": 1, "Roi Et": 1}, AgeGroup: entity.Group{Age0to30: 1, Age31to60: 1, Age61plus: 1, AgeNA: 0}}
	EXPECTED_NULL_PROVINCE_REPONSE          = entity.CovidSummary{Province: map[string]int{"Phrae": 1, "Suphan Buri": 1, "Roi Et": 1, "N/A": 2}, AgeGroup: entity.Group{Age0to30: 1, Age31to60: 1, Age61plus: 3, AgeNA: 0}}
	EXPECTED_NULL_AGE_RESPONSE              = entity.CovidSummary{Province: map[string]int{"Phrae": 1, "Suphan Buri": 1, "Roi Et": 3}, AgeGroup: entity.Group{Age0to30: 1, Age31to60: 1, Age61plus: 1, AgeNA: 2}}
	EXPEXTED_NULL_PROVINCE_AND_AGE_RESPONSE = entity.CovidSummary{Province: map[string]int{"Phrae": 1, "Suphan Buri": 1, "Roi Et": 3, "N/A": 2}, AgeGroup: entity.Group{Age0to30: 1, Age31to60: 1, Age61plus: 3, AgeNA: 2}}
	EXPEXTED_OVER_61_AGE__RESPONSE          = entity.CovidSummary{Province: map[string]int{"Phrae": 1, "Suphan Buri": 1, "Roi Et": 3}, AgeGroup: entity.Group{Age0to30: 1, Age31to60: 1, Age61plus: 3, AgeNA: 0}}
)

var _ = Describe("Covid Service", func() {
	var (
		covidRepository repository.CovidRepository
		covidService    CovidService
	)

	BeforeSuite(func() {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	})
	Describe("Get Covid Summary", func() {
		Context("If public url is incorrect", func() {
			BeforeEach(func() {
				covidStatsURL := "**incorect url**\n"
				client := &http.Client{}
				covidRepository = repository.NewCovidRepository(client, covidStatsURL)
				covidService = New(covidRepository)
			})

			It("Should error", func() {
				_, err := covidService.GetSummary()
				Ω(err).Should(HaveOccurred())
			})

			It("Shold have error message = \"can not make request to public api\"", func() {
				_, err := covidService.GetSummary()
				Ω(err.Error()).Should(Equal("can not make request to public api"))
			})
		})

		Context("If cant connect to public api", func() {
			BeforeEach(func() {
				covidStatsURL := "https://static.wongnai.com/devinterview/covid-cases.json"
				client := &ClientMockStatus{status: http.StatusInternalServerError, err: errors.New(http.StatusText(http.StatusInternalServerError))}
				covidRepository = repository.NewCovidRepository(client, covidStatsURL)
				covidService = New(covidRepository)
			})

			It("Should error", func() {
				_, err := covidService.GetSummary()
				Ω(err).Should(HaveOccurred())
			})

			It("Shold have error message = \"can not get covid stats from public api\"", func() {
				_, err := covidService.GetSummary()
				Ω(err.Error()).Should(Equal("can not get covid stats from public api"))
			})
		})

		Context("If response from public api is invalid format", func() {
			BeforeEach(func() {
				covidStatsURL := "https://static.wongnai.com/devinterview/covid-cases.json"
				client := &ClientMockResponse{response: `***invalid response***\n`}
				covidRepository = repository.NewCovidRepository(client, covidStatsURL)
				covidService = New(covidRepository)
			})

			It("Should error", func() {
				res, err := covidService.GetSummary()
				fmt.Printf("\nres:%v\n", res)
				Ω(err).Should(HaveOccurred())
			})

			It("Shold have error message = \"can not parse covid stats from public api\"", func() {
				res, err := covidService.GetSummary()
				if err != nil {
					fmt.Println(res)
				}
				Ω(err.Error()).Should(Equal("can not parse covid stats from public api"))
			})
		})

		Context("If public url correct and parse data correctly", func() {
			BeforeEach(func() {
				covidStatsURL := "https://static.wongnai.com/devinterview/covid-cases.json"
				client := &http.Client{}
				covidRepository = repository.NewCovidRepository(client, covidStatsURL)
				covidService = New(covidRepository)
			})

			It("Should not error", func() {
				_, err := covidService.GetSummary()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Should return CovidSummary entity", func() {
				res, _ := covidService.GetSummary()
				Ω(fmt.Sprintf("%T", res)).Should(Equal("entity.CovidSummary"))
			})

			It("Should not return empty province", func() {
				res, _ := covidService.GetSummary()
				Ω(len(res.Province)).ShouldNot(Equal(0))
			})

			It("Should return AgeGroup property", func() {
				res, _ := covidService.GetSummary()
				Ω(fmt.Sprintf("%T", res.AgeGroup)).Should(Equal("entity.Group"))
			})
		})

		Context("If response from api is WITHOUT_NULL_RESPONSE_FROM_API", func() {
			BeforeEach(func() {
				covidStatsURL := "https://static.wongnai.com/devinterview/covid-cases.json"
				client := &ClientMockResponse{response: WITHOUT_NULL_RESPONSE_FROM_API}
				covidRepository = repository.NewCovidRepository(client, covidStatsURL)
				covidService = New(covidRepository)
			})

			It("Should not error", func() {
				_, err := covidService.GetSummary()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Should  match EXPECTED_WITHOUT_NULL_RESPONSE", func() {
				res, _ := covidService.GetSummary()
				fmt.Println(res)
				Ω(res).Should(Equal(EXPECTED_WITHOUT_NULL_RESPONSE))

			})
		})

		Context("If response from api is NULL_PROVINCE_REPONSE_FROM_API", func() {
			BeforeEach(func() {
				covidStatsURL := "https://static.wongnai.com/devinterview/covid-cases.json"
				client := &ClientMockResponse{response: NULL_PROVINCE_REPONSE_FROM_API}
				covidRepository = repository.NewCovidRepository(client, covidStatsURL)
				covidService = New(covidRepository)
			})

			It("Should not error", func() {
				_, err := covidService.GetSummary()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Should  match EXPECTED_NULL_PROVINCE_REPONSE", func() {
				res, _ := covidService.GetSummary()
				fmt.Println(res)
				Ω(res).Should(Equal(EXPECTED_NULL_PROVINCE_REPONSE))

			})
		})

		Context("If response from api is NULL_AGE_RESPONSE_FROM_API", func() {
			BeforeEach(func() {
				covidStatsURL := "https://static.wongnai.com/devinterview/covid-cases.json"
				client := &ClientMockResponse{response: NULL_AGE_RESPONSE_FROM_API}
				covidRepository = repository.NewCovidRepository(client, covidStatsURL)
				covidService = New(covidRepository)
			})

			It("Should not error", func() {
				_, err := covidService.GetSummary()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Should  match EXPECTED_NULL_AGE_RESPONSE", func() {
				res, _ := covidService.GetSummary()
				fmt.Println(res)
				Ω(res).Should(Equal(EXPECTED_NULL_AGE_RESPONSE))

			})
		})

		Context("If response from api is NULL_PROVINCE_AND_AGE_RESPONSE_FROM_API", func() {
			BeforeEach(func() {
				covidStatsURL := "https://static.wongnai.com/devinterview/covid-cases.json"
				client := &ClientMockResponse{response: NULL_PROVINCE_AND_AGE_RESPONSE_FROM_API}
				covidRepository = repository.NewCovidRepository(client, covidStatsURL)
				covidService = New(covidRepository)
			})

			It("Should not error", func() {
				_, err := covidService.GetSummary()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Should  match EXPEXTED_NULL_PROVINCE_AND_AGE_RESPONSE", func() {
				res, _ := covidService.GetSummary()
				fmt.Println(res)
				Ω(res).Should(Equal(EXPEXTED_NULL_PROVINCE_AND_AGE_RESPONSE))

			})
		})

		Context("If response from api is OVER_61_AGE_RESPONSE_FROM_API", func() {
			BeforeEach(func() {
				covidStatsURL := "https://static.wongnai.com/devinterview/covid-cases.json"
				client := &ClientMockResponse{response: OVER_61_AGE_RESPONSE_FROM_API}
				covidRepository = repository.NewCovidRepository(client, covidStatsURL)
				covidService = New(covidRepository)
			})

			It("Should not error", func() {
				_, err := covidService.GetSummary()
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Should  match EXPEXTED_OVER_61_AGE__RESPONSE", func() {
				res, _ := covidService.GetSummary()
				fmt.Println(res)
				Ω(res).Should(Equal(EXPEXTED_OVER_61_AGE__RESPONSE))

			})
		})
	})
})
