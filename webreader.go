package webreader

import (
	"io"
	"io/ioutil"

	//"log"
	"logger"
	"net/http"
)

type Dyad struct {
	KeyName string
	Value   string
}

type RequestResult struct {
	Text   string
	Stream io.Reader
}

func (res *RequestResult) Reset() {
	res.Text = ""
	res.Stream = nil
}

var cookieHandler = new(Cookies)
var currentUrl string
var client http.Client
var myReq http.Request
var result = new(RequestResult)

func errorHandle(e error) {
	if e != nil {
		panic(e)
	}
}

func processResponse(response *http.Response) string {
	logger.Debug("RESPONSE_HEADERS")
	for name, value := range response.Header {
		logger.Debug(name, value)
	}

	body, err := ioutil.ReadAll(response.Body)
	errorHandle(err)

	//result.Stream = response.Body
	//result.Text = string(body)

	return string(body)
}

func PrepareRequestParameters() (*http.Request, error) {
	myReq, err := http.NewRequest(currentOptions.Method, currentUrl, nil)
	errorHandle(err)
	logger.Debug("REQUEST_HEADERS")
	myReq.Header.Add("Host", myReq.Host)
	for _, value := range currentOptions.HttpHeaders {
		logger.Debug(value.KeyName, value.Value)
		myReq.Header.Add(value.KeyName, value.Value)
	}
	if len(cookieHandler.Cookies) == 0 {
		logger.Debug("NO_REQUEST_COOKIES")
	} else {
		logger.Debug("REQUEST_COOKIES")
		for _, cookie := range cookieHandler.Cookies {
			logger.Debug(cookie.String())
		}
	}
	return myReq, err
}

func DoRequest(url string, options *RequestOptions) string {
	result.Reset()

	currentUrl = url
	processOptions()
	req, err := PrepareRequestParameters()
	client := &http.Client{}

	resp, err := client.Do(req)
	errorHandle(err)
	defer resp.Body.Close()

	cookieHandler.SaveCookies(resp)

	return processResponse(resp)
}
