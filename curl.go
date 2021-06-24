package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	URL      string
	URI      string
	Method   string
	Headers  map[string]string
	Params   url.Values
	Json     interface{}
	IsJson   bool
	Byte     []byte `json:"-"`
	ByteStr  string
	IsByte   bool
	IsNotify bool
	Client   http.Client `json:"-"`
	Xml      string
	IsXml    bool
}

var client http.Client

func init() {
	client = http.Client{
		Timeout: time.Second * 5,
	}
}

func NewRequest() *Request {
	return &Request{
		Headers: map[string]string{},
		Client:  client,
	}
}

func (r *Request) validateURL() error {
	if r.URL == "" || len(r.URL) == 0 {
		return errors.New("URL is required")
	}
	return nil
}

func (r *Request) validateMethod() error {
	switch r.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return nil
	}

	return errors.New("Unsupported method " + r.Method)
}

func (r *Request) doRequest() (*http.Response, []byte, error) {
	var request *http.Request
	var response *http.Response

	if err := r.validateURL(); err != nil {
		return response, nil, err
	}

	r.Method = strings.ToUpper(r.Method)
	if err := r.validateMethod(); err != nil {
		return response, nil, err
	}

	u, err := url.Parse(r.URL)
	if err != nil {
		return response, nil, err
	}

	if r.Headers == nil {
		r.Headers = map[string]string{}
	}

	if r.Method == "GET" {
		request, err = r.Get(u)
	} else if r.Method == "POST" && r.IsJson {
		request, err = r.PostJSON(u)
	} else if r.Method == "POST" && r.IsByte {
		request, err = r.PostByte(u)
	} else if r.Method == "POST" && r.IsXml {
		request, err = r.PostXML(u)
	} else if r.Method == "POST" {
		request, err = r.Post(u)
	} else if r.Method == "DELETE" && r.IsJson {
		request, err = r.PostJSON(u)
	} else if r.Method == "DELETE" {
		request, err = r.Post(u)
	} else if r.Method == "PATCH" && r.IsJson {
		request, err = r.PostJSON(u)
	} else if r.Method == "PATCH" {
		request, err = r.Post(u)
	} else if r.Method == "PUT" {
		request, err = r.Post(u)
	}

	if request == nil {
		return response, nil, errors.New("Failed create new request")
	}

	if err != nil {
		return response, nil, err
	}

	for key, value := range r.Headers {
		request.Header.Set(key, value)
	}
	// request.Close = true

	response, err = r.Client.Do(request)
	if err != nil {
		return response, nil, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response, nil, err
	}

	return response, contents, nil
}

func (r *Request) Get(u *url.URL) (*http.Request, error) {
	u.RawQuery = r.Params.Encode()
	req, err := http.NewRequest(r.Method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *Request) Post(u *url.URL) (*http.Request, error) {
	form := strings.NewReader(r.Params.Encode())
	req, err := http.NewRequest(r.Method, u.String(), form)
	if err != nil {
		return nil, err
	}

	r.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	return req, nil
}

func (r *Request) PostJSON(u *url.URL) (*http.Request, error) {
	body, err := json.Marshal(r.Json)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	r.Headers["Content-Type"] = "application/json"

	return req, nil
}

func (r *Request) PostXML(u *url.URL) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, u.String(), bytes.NewReader([]byte(r.Xml)))
	if err != nil {
		return nil, err
	}

	r.Headers["Content-Type"] = "application/xml"

	return req, nil
}

func (r *Request) PostByte(u *url.URL) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, u.String(), bytes.NewBuffer(r.Byte))
	if err != nil {
		return nil, err
	}

	return req, nil
}
