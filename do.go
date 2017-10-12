package gorest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

// Do execute the HTTP request
func Do(session *Session, request Request) (Response, error) {
	return do(session, request, true)
}

// DoWithoutFollowRedirects execute the HTTP request without following redirects
func DoWithoutFollowRedirects(session *Session, request Request) (Response, error) {
	return do(session, request, false)
}

func do(session *Session, request Request, followRedirects bool) (Response, error) {
	if session.Cookie == nil {
		session.Cookie, _ = cookiejar.New(nil)
	}

	var client *http.Client

	if followRedirects {
		client = &http.Client{Jar: session.Cookie}
	} else {
		client = &http.Client{
			Jar: session.Cookie,
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

	resp, err := client.Do(req)

	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	session.Cookie.SetCookies(&request.URL, resp.Cookies())

	var response Response
	response.Body = contents
	response.Header = resp.Header
	response.StatusCode = resp.StatusCode

	return response, nil
}
