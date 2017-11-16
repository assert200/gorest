package gorest

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// RestRequest struct
type RestRequest struct {
	Body            []byte
	Method          string
	Headers         http.Header
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

	return s
}

// RestTest struct
type RestTest struct {
	RestRequest        RestRequest
	RestResponse       RestResponse
	Description        string
	Values             map[string]string
	ElapsedTime        float64
	ExpectedStatusCode int
	//ExpectedBodyBlackList      []regexp.Regexp
	//ExpectedBodyWhiteList      []regexp.Regexp
}

// Result Result
type Result struct {
	TotalElaspedTime float64
	TotalRequests    float64
}

func (r Result) String() string {
	s := fmt.Sprintf("Total Elapsed Time: %f\n", r.TotalElaspedTime)
	s += fmt.Sprintf("Total Requests: %f\n", r.TotalRequests)
	s += fmt.Sprintf("Avg Request Time: %f\n", r.TotalElaspedTime/r.TotalRequests)

	return s
}

//Results Results
type Results map[string]*Result

//Add Add
func (rs Results) Add(restTest RestTest) {
	if _, ok := rs[restTest.Description]; !ok {
		var result Result
		result.TotalElaspedTime = restTest.ElapsedTime
		result.TotalRequests = 1

		rs[restTest.Description] = &result
	} else {
		rs[restTest.Description].TotalElaspedTime += restTest.ElapsedTime
		rs[restTest.Description].TotalRequests++
	}
}

func (rs Results) String() string {
	var s string
	for key, tally := range rs {
		s += fmt.Sprintln(key, " ", tally)
	}

	return s
}
