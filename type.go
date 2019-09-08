package gorest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

// RestRequest struct
type RestRequest struct {
	Body            []byte
	Method          string
	Headers         http.Header
	CookieJar       *cookiejar.Jar
	URL             url.URL
	FollowRedirects bool
	Delay           time.Duration
}

// NewRestRequest factory
func NewRestRequest() RestRequest {
	restRequest := RestRequest{}
	restRequest.Headers = http.Header{}
	restRequest.FollowRedirects = true
	restRequest.CookieJar, _ = cookiejar.New(nil)

	return restRequest
}

func (r RestRequest) String() string {
	s := fmt.Sprintf("Request URL: %s\n", r.URL.String())
	s += fmt.Sprintf("Request Method: %s\n", r.Method)
	s += fmt.Sprintf("Request Body: %s\n", r.Body)

	for k, v := range r.Headers {
		s += fmt.Sprintln("Request Header Key: ", k, "Value: ", v)
	}

	s += fmt.Sprintf("Request Cookies: %v\n", r.CookieJar)

	return s
}

// RestResponse struct
type RestResponse struct {
	Body       []byte
	Headers    http.Header
	CookieJar  *cookiejar.Jar
	StatusCode int
}

func (r RestResponse) String() string {
	s := fmt.Sprintf("Response status code: %d\n", r.StatusCode)
	s += fmt.Sprintf("Response Body: %s\n", string(r.Body))

	for k, v := range r.Headers {
		s += fmt.Sprintln("Response Header Key: ", k, "Value: ", v)
	}

	s += fmt.Sprintf("Response Cookies: %v\n", r.CookieJar)

	return s
}

// A Generator creates new tests from responses from existing tests
type Generator func(restTestResponse RestTest) (newTests []RestTest)

// RestTest struct
type RestTest struct {
	RestRequest        RestRequest
	RestResponse       RestResponse
	RestTestResult     RestTestResult
	Generator          Generator
	Description        string
	Values             map[string]string
	ExpectedStatusCode int
	BodyExpectations   []*regexp.Regexp
	BodyRefusals       []*regexp.Regexp
}

// RestTestResult struct
type RestTestResult struct {
	Description      string    `json:"description"`
	URLUnescaped     string    `json:"urlUnescaped"`
	RequestTimeStart time.Time `json:"requestTimeStart"`
	RequestTimeEnd   time.Time `json:"requestTimeEnd"`
	RequestDuration  float64   `json:"requestTimeDuration"`
	StatusCode       int       `json:"statusCode"`
	Errors           []error   `json:"errors"`
}

func (r RestTestResult) String() string {
	json, _ := json.Marshal(r)
	return fmt.Sprintln(string(json))
}

// ResultTally is a summary recalculated after each request
type ResultTally struct {
	ShortestRequestDuration float64 `json:"shortestRequestDuration"`
	LongestRequestDuration  float64 `json:"longestRequestDuration"`
	TotalElapsedDuration    float64 `json:"totalElapsedDuration"`
	TotalRequests           int     `json:"totalRequests"`
	TotalErrors             int     `json:"totalErrors"`
}

func (r ResultTally) String() string {
	json, _ := json.Marshal(r)
	return fmt.Sprintln(string(json))
}

//ResultTallys Is all the result tallys for each type of test
type ResultTallys map[string]*ResultTally

//Add Add
func (rs ResultTallys) Add(restTest RestTest) {
	if _, ok := rs[restTest.Description]; !ok {
		var result ResultTally
		result.ShortestRequestDuration = restTest.RestTestResult.RequestDuration
		result.LongestRequestDuration = restTest.RestTestResult.RequestDuration
		result.TotalElapsedDuration = restTest.RestTestResult.RequestDuration
		result.TotalRequests = 1
		result.TotalErrors = len(restTest.RestTestResult.Errors)

		rs[restTest.Description] = &result
	} else {

		if restTest.RestTestResult.RequestDuration < rs[restTest.Description].ShortestRequestDuration {
			rs[restTest.Description].ShortestRequestDuration = restTest.RestTestResult.RequestDuration
		}
		if restTest.RestTestResult.RequestDuration > rs[restTest.Description].LongestRequestDuration {
			rs[restTest.Description].LongestRequestDuration = restTest.RestTestResult.RequestDuration
		}
		rs[restTest.Description].TotalElapsedDuration += restTest.RestTestResult.RequestDuration
		rs[restTest.Description].TotalRequests++
		rs[restTest.Description].TotalErrors += len(restTest.RestTestResult.Errors)
	}
}

func (rs ResultTallys) String() string {
	var totalResultTally ResultTally
	var s string

	firstResult := true
	for key, result := range rs {
		if firstResult {
			totalResultTally.ShortestRequestDuration = result.ShortestRequestDuration
			totalResultTally.LongestRequestDuration = result.LongestRequestDuration
			firstResult = false
		} else {
			if result.ShortestRequestDuration < totalResultTally.ShortestRequestDuration {
				totalResultTally.ShortestRequestDuration = result.ShortestRequestDuration
			}
			if result.LongestRequestDuration > totalResultTally.LongestRequestDuration {
				totalResultTally.LongestRequestDuration = result.LongestRequestDuration
			}
		}

		totalResultTally.TotalElapsedDuration += result.TotalElapsedDuration
		totalResultTally.TotalErrors += result.TotalErrors
		totalResultTally.TotalRequests += result.TotalRequests

		s += fmt.Sprintf("%s: %s", key, result)
	}

	s += fmt.Sprintln("TOTAL RESULT TALLY:", totalResultTally)

	return s
}
