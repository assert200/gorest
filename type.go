package gorest

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

// RestRequest struct
type RestRequest struct {
	Body    []byte
	Method  string
	Headers http.Header
	Cookies         *cookiejar.Jar 
	URL             url.URL
	FollowRedirects bool
}

// NewRestRequest factory
func NewRestRequest() RestRequest {
	restRequest := RestRequest{}
	restRequest.Headers = http.Header{}
	restRequest.FollowRedirects = true
	restRequest.Cookies, _ = cookiejar.New(nil)

	return restRequest
}

func (r RestRequest) String() string {
	s := fmt.Sprintf("Request URL: %s\n", r.URL.String())
	s += fmt.Sprintf("Request Method: %s\n", r.Method)
	s += fmt.Sprintf("Response Body: %s\n", r.Body)

	for k, v := range r.Headers {
		s += fmt.Sprintln("Request Header Key: ", k, "Value: ", v)
	}
	
	s += fmt.Sprintf("Response Cookies: %v\n", r.Cookies)

	return s
}

// RestResponse struct
type RestResponse struct {
	Body       []byte
	Headers    http.Header
	Cookies    *cookiejar.Jar
	StatusCode int
}

func (r RestResponse) String() string {
	s := fmt.Sprintf("Response status code: %d\n", r.StatusCode)
	s += fmt.Sprintf("Response Body: %s\n", string(r.Body))

	for k, v := range r.Headers {
		s += fmt.Sprintln("Response Header Key: ", k, "Value: ", v)
	}
	
	s += fmt.Sprintf("Response Cookies: %v\n", r.Cookies)

	return s
}

// A Generator creates new tests from responses from existing tests
type Generator func(restTestResponse RestTest) (newTests []RestTest)

// RestTest struct
type RestTest struct {
	RestRequest        RestRequest
	RestResponse       RestResponse
	Generator          Generator
	Description        string
	Values             map[string]string
	RequestTime        float64
	ExpectedStatusCode int
	BodyExpectations   []*regexp.Regexp
	BodyRefusals       []*regexp.Regexp
	Errors             []error
}

// Result Result
type Result struct {
	ShortestRequestTime float64
	LongestRequestTime  float64
	TotalElapsedTime    float64
	TotalRequests       int
	TotalErrors         int
}

func (r Result) String() string {
	return fmt.Sprintf("AvgReq, %f, Variance, %f, ShortestReq, %f, LongestReq, %f, TotalElapsed, %f, TotalReqs, %d, TotalErrors, %d", r.TotalElapsedTime/float64(r.TotalRequests), (r.LongestRequestTime - r.ShortestRequestTime), r.ShortestRequestTime, r.LongestRequestTime, r.TotalElapsedTime, r.TotalRequests, r.TotalErrors)
}

//Results Results
type Results map[string]*Result

//Add Add
func (rs Results) Add(restTest RestTest) {
	if _, ok := rs[restTest.Description]; !ok {
		var result Result
		result.ShortestRequestTime = restTest.RequestTime
		result.LongestRequestTime = restTest.RequestTime
		result.TotalElapsedTime = restTest.RequestTime
		result.TotalRequests = 1
		result.TotalErrors = len(restTest.Errors)

		rs[restTest.Description] = &result
	} else {

		if restTest.RequestTime < rs[restTest.Description].ShortestRequestTime {
			rs[restTest.Description].ShortestRequestTime = restTest.RequestTime
		}
		if restTest.RequestTime > rs[restTest.Description].LongestRequestTime {
			rs[restTest.Description].LongestRequestTime = restTest.RequestTime
		}
		rs[restTest.Description].TotalElapsedTime += restTest.RequestTime
		rs[restTest.Description].TotalRequests++
		rs[restTest.Description].TotalErrors += len(restTest.Errors)
	}
}

func (rs Results) String() string {
	var totalResult Result
	var s string

	firstResult := true
	for key, result := range rs {
		if firstResult {
			totalResult.ShortestRequestTime = result.ShortestRequestTime
			totalResult.LongestRequestTime = result.LongestRequestTime
			firstResult = false
		} else {
			if result.ShortestRequestTime < totalResult.ShortestRequestTime {
				totalResult.ShortestRequestTime = result.ShortestRequestTime
			}
			if result.LongestRequestTime > totalResult.LongestRequestTime {
				totalResult.LongestRequestTime = result.LongestRequestTime
			}
		}

		totalResult.TotalElapsedTime += result.TotalElapsedTime
		totalResult.TotalErrors += result.TotalErrors
		totalResult.TotalRequests += result.TotalRequests

		s += fmt.Sprintln(key, ",", result)
	}

	s += fmt.Sprintln("TOTAL RESULT,", totalResult)

	return s
}
