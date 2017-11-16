package gorest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

// DoAndVerify DoAndVerify
func DoAndVerify(restTest RestTest) RestTest {
	restTest, err := Do(restTest)

	request := restTest.RestRequest
	response := restTest.RestResponse

	if err != nil {
		fmt.Println("FATAL: There was an error excuting the rest request: ", request, "With Error: ", err.Error())
		os.Exit(1)
	}

	if response.StatusCode != restTest.ExpectedStatusCode {
		fmt.Println("** WARNING: with ", request.URL.RequestURI(), " Expecting Status Code: ", restTest.ExpectedStatusCode, ", Recieved: ", response.StatusCode, " **")
	}

	fmt.Println("LOG: with ", request.URL.RequestURI(), " Elasped Time ", restTest.ElapsedTime)

	return restTest
}

// Do execute the HTTP request
func Do(restTest RestTest) (RestTest, error) {
	restRequest := restTest.RestRequest

	var client *http.Client

	if restRequest.FollowRedirects {
		client = &http.Client{Jar: restRequest.Cookies}
	} else {
		client = &http.Client{
			Jar: restRequest.Cookies,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	httpRequest, err := http.NewRequest(restRequest.Method, restRequest.URL.String(), bytes.NewReader(restRequest.Body))

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
	restTest.ElapsedTime = time.Since(start).Seconds()

	contents, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return restTest, err
	}

	restResponse.Cookies, _ = cookiejar.New(nil)
	restResponse.Cookies.SetCookies(&restRequest.URL, httpResponse.Cookies())

	restResponse.Body = contents
	restResponse.Headers = httpResponse.Header
	restResponse.StatusCode = httpResponse.StatusCode

	restTest.RestResponse = restResponse

	return restTest, nil
}
