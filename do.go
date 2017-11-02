package gorest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Do execute the HTTP request
func Do(session *Session, request Request) (Response, error) {
	return do(session, request, true)
}

// DoWithoutFollowRedirects execute the HTTP request without following redirects
func DoWithoutFollowRedirects(session *Session, request Request) (Response, error) {
	return do(session, request, false)
}

// NewSession generates blank session
func NewSession() *Session {
	session := &Session{}
	session.Cookies, _ = cookiejar.New(nil)
	return session
}

func do(session *Session, request Request, followRedirects bool) (Response, error) {
	if session.Cookies == nil {
		session.Cookies, _ = cookiejar.New(nil)
	}

	var client *http.Client

	if followRedirects {
		client = &http.Client{Jar: session.Cookies}
	} else {
		client = &http.Client{
			Jar: session.Cookies,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	req, err := http.NewRequest(request.Method, request.URL.String(), bytes.NewReader(request.Body))

	if err != nil {
		return Response{}, err
	}

	req.Close = true
	req.Header = request.Header

	start := time.Now()
	resp, err := client.Do(req)

	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()

	var response Response
	response.ElapsedTime = time.Since(start).Seconds()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	session.Cookies.SetCookies(&request.URL, resp.Cookies())

	response.Body = contents
	response.Header = resp.Header
	response.StatusCode = resp.StatusCode

	return response, nil
}
