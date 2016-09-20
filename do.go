package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

func Do(session *Session, request Request) (Response, error) {
	if session.Cookie == nil {
		session.Cookie, _ = cookiejar.New(nil)
	}

	client := &http.Client{Jar: session.Cookie}

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
