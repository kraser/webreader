package webreader

import (
	"errors"
	errs "errorshandler"
	"io"
	"io/ioutil"
	"logger"
	"net/http"
	"strings"
	"time"
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

	//вынести в ParserHandleError
	if len(req.Cookies()) == 0 {
		logger.Debug("NO_REQUEST_COOKIES")
	} else {
		logger.Debug("REQUEST_COOKIES")
		myReq.Header.Del("Cookie")
	}
	//вынести в ParserHandleError

	toDoReq := true
	var html string
	var respErr error = nil
	var trials int
	trials = 0
	for toDoReq {
		logger.Debug("TRY: ", trials)
		resp, err := client.Do(req)
		errs.ErrorHandle(err)
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			ResponseHeaders(resp)
			respErr = errors.New(strings.Join([]string{resp.Status, url}, " AT "))
			logger.Error(respErr)
			trials++
			toDoReq = trials <= options.Trials

			//вынести в ParserHandleError
			currentOptions.SetRandUserAgent()
			req.Header.Set("User-Agent", currentOptions.UserAgent)
			//вынести в ParserHandleError

			time.Sleep(options.Interval * time.Second)
		} else {
			toDoReq = false
			html = processResponse(resp)
		}
	}

	return html, respErr
}

func ResponseHeaders(response *http.Response) {
	for key, headers := range response.Header {
		logger.Debug(key)
		logger.Debug(headers)
	}
}
