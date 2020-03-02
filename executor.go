package gorest

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// ExecuteAndVerify ExecuteAndVerify
func ExecuteAndVerify(restTest RestTest) RestTest {
	var verifyErrors []error

	restTest, err := Execute(restTest)

	response := restTest.RestResponse

	if err != nil {
		verifyErrors = append(verifyErrors, err)
	} else {
		if response.StatusCode != restTest.ExpectedStatusCode {
			errorMsg := fmt.Sprintf("Expecting Status Code: %d Received: %d", restTest.ExpectedStatusCode, response.StatusCode)
			verifyErrors = append(verifyErrors, errors.New(errorMsg))
		}

		for _, bodyExpectation := range restTest.BodyExpectations {
			if !bodyExpectation.MatchString(string(restTest.RestResponse.Body)) {
				errorMsg := fmt.Sprintf("Body expectation %v was not met: %s", bodyExpectation, string(restTest.RestResponse.Body))
				verifyErrors = append(verifyErrors, errors.New(errorMsg))
			}
		}

		for _, bodyRefusal := range restTest.BodyRefusals {
			if bodyRefusal.MatchString(string(restTest.RestResponse.Body)) {
				errorMsg := fmt.Sprintf("Body refusal %v was detected", bodyRefusal)
				verifyErrors = append(verifyErrors, errors.New(errorMsg))
			}
		}
	}

	restTest.RestTestResult.Errors = verifyErrors

	//log.Println(restTest.RestTestResult)
	return restTest
}

// Execute execute the HTTP request
func Execute(restTest RestTest) (RestTest, error) {
	restRequest := restTest.RestRequest

	time.Sleep(restRequest.Delay)

	var client *http.Client

	if restRequest.FollowRedirects {
		client = &http.Client{Jar: restRequest.CookieJar}
	} else {
		client = &http.Client{
			Jar: restRequest.CookieJar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	URLUnescaped, _ := url.PathUnescape(restRequest.URL.String())
	httpRequest, err := http.NewRequest(restRequest.Method, URLUnescaped, bytes.NewReader(restRequest.Body))

	if err != nil {
		return restTest, err
	}

	httpRequest.Close = true
	httpRequest.Header = restRequest.Headers

	start := time.Now()
	httpResponse, err := client.Do(httpRequest)
	end := time.Now()
	if err != nil {
		return restTest, err
	}

	defer httpResponse.Body.Close()

	var restResponse RestResponse

	contents, err := ioutil.ReadAll(httpResponse.Body)

	if err != nil {
		return restTest, err
	}

	restResponse.CookieJar = restRequest.CookieJar

	restResponse.Body = contents
	restResponse.Headers = httpResponse.Header
	restResponse.StatusCode = httpResponse.StatusCode
	restTest.RestResponse = restResponse

	var restTestResult RestTestResult
	restTestResult.Description = restTest.Description
	restTestResult.URLUnescaped, err = url.PathUnescape(restRequest.URL.String())
	if err != nil {
		panic(err)
	}
	restTestResult.RequestTimeStart = start
	restTestResult.RequestTimeEnd = end
	restTestResult.RequestDuration = end.Sub(start).Seconds()
	restTestResult.StatusCode = httpResponse.StatusCode
	restTest.RestTestResult = restTestResult

	return restTest, nil
}
