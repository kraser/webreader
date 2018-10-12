package webreader

import (
	errs "errorshandler"
	"io"
	"io/ioutil"

	//"log"
	"errors"
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

func processResponse(response *http.Response) string {
	logger.Debug("RESPONSE_STATUS:", response.StatusCode)
	logger.Debug("RESPONSE:", response.Status)
	logger.Debug("RESPONSE_HEADERS")
	for name, value := range response.Header {
		logger.Debug(name, value)
	}
	cookieHandler.SaveCookies(response)

	body, err := ioutil.ReadAll(response.Body)
	errs.ErrorHandle(err)

	//result.Stream = response.Body
	//result.Text = string(body)

	return string(body)
}

func PrepareRequestParameters() (*http.Request, error) {
	myReq, err := http.NewRequest(currentOptions.Method, currentUrl, nil)
	errs.ErrorHandle(err)
	logger.Debug("REQUEST_HEADERS")
	myReq.Header.Add("Host", myReq.Host)
	logger.Debug("Host", myReq.Host)
	myReq.Header.Add("User-Agent", currentOptions.UserAgent)
	logger.Debug("User-Agent", currentOptions.UserAgent)
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

func DoRequest(url string, options *RequestOptions) (string, error) {
	result.Reset()

	currentUrl = url
	processOptions()
	req, err := PrepareRequestParameters()
	errs.ErrorHandle(err)
	client := &http.Client{}
	toDoReq := true
	var html string
	var respErr error = nil
	var trials int8
	trials = 0
	for toDoReq {
		logger.Debug("TRY: ", trials)
		resp, err := client.Do(req)
		errs.ErrorHandle(err)
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			respErr = errors.New(resp.Status)
			logger.Error(respErr)
			trials++
			toDoReq = trials <= options.Trials
		} else {
			toDoReq = false
			html = processResponse(resp)
		}
	}

	return html, respErr
}
