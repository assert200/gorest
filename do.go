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

// DoAndVerify DoAndVerify
func DoAndVerify(restTest RestTest) RestTest {
	var verifyErrors []error

	restTest, err := Do(restTest)

	request := restTest.RestRequest
	response := restTest.RestResponse

	if err != nil {
		verifyErrors = append(verifyErrors, err)
	} else {
		if response.StatusCode != restTest.ExpectedStatusCode {
			errorMsg := fmt.Sprintf("Expecting Status Code: %d Recieved: %d", restTest.ExpectedStatusCode, response.StatusCode)
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

	fmt.Printf("LOG: %s %s Elasped Time: %f Errors: %v \n", restTest.Description, request.URL.RequestURI(), restTest.ElapsedTime, verifyErrors)

	restTest.Errors = verifyErrors
	return restTest
}

// Do execute the HTTP request
func Do(restTest RestTest) (RestTest, error) {
	restRequest := restTest.RestRequest

	var client *http.Client

	if restRequest.FollowRedirects {
		client = &http.Client{ /*Jar: restRequest.Cookies*/ }
	} else {
		client = &http.Client{
			/* Jar: restRequest.Cookies, */
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

	if err != nil {
		return restTest, err
	}

	defer httpResponse.Body.Close()

	var restResponse RestResponse

	contents, err := ioutil.ReadAll(httpResponse.Body)
	restTest.ElapsedTime = time.Since(start).Seconds()
	if err != nil {
		return restTest, err
	}

	//restResponse.Cookies, _ = cookiejar.New(nil)
	//restResponse.Cookies.SetCookies(&restRequest.URL, httpResponse.Cookies())

	restResponse.Body = contents
	restResponse.Headers = httpResponse.Header
	restResponse.StatusCode = httpResponse.StatusCode

	restTest.RestResponse = restResponse

	return restTest, nil
}
