package gocurl

import (
	"logger"

	"net/http"
	"time"

	errs "github.com/kraser/errorshandler"
)

var client = new(CurlClient)

type CurlClient struct {
	url     string
	Options *RequestOptions
	request http.Request
	result  *RequestResult
	*Cookies
	requestTime time.Time
}

func GetCurl() *CurlClient {
	return client
}

func InitCurl(options *RequestOptions) *CurlClient {
	client.Options = options
	client.requestTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	client.result = new(RequestResult)
	client.Cookies = &Cookies{}
	if len(options.CookieFile) != 0 {
		client.Cookies.SetCookieFileName(options.CookieFile)
	}

	return client
}

/* CurlClient methods */
func (c *CurlClient) DoRequest(url string) string {
	c.result.Reset()
	c.url = url
	req, err := c.PrepareRequestParameters()
	errs.ErrorHandle(err)
	httpClient := &http.Client{}

	toDoRequest := true
	var requestErr RequestError
	var result string
	var trials int
	trials = 0
	for toDoRequest {
		current := time.Now()
		elapsed := current.Sub(c.requestTime).Seconds()
		if elapsed < c.Options.Interval {
			time.Sleep(time.Duration(c.Options.Interval-elapsed) * time.Second)
		}
		c.requestTime = time.Now()
		logger.Debug("TRY: ", trials)
		resp, err := httpClient.Do(req)
		errs.ErrorHandle(err)
		defer resp.Body.Close()
		result, requestErr = c.processResponse(resp)
		if requestErr.IsNull() {
			trials++
			toDoRequest = trials <= c.Options.Trials

			//вынести в ParserHandleError
			c.Options.SetRandUserAgent()
			req.Header.Set("User-Agent", c.Options.UserAgent)
			//вынести в ParserHandleError

		} else {
			toDoRequest = false
		}
	}

	return result
}

func (c *CurlClient) PrepareRequestParameters() (*http.Request, error) {
	myReq, err := http.NewRequest(c.Options.GetMethod(), c.url, nil)
	errs.ErrorHandle(err)
	logger.Debug("REQUEST_HEADERS")
	myReq.Header.Add("Host", myReq.Host)
	logger.Debug("Host", myReq.Host)
	myReq.Header.Add("User-Agent", c.Options.UserAgent)
	logger.Debug("User-Agent", c.Options.UserAgent)
	for name, value := range c.Options.HttpHeaders {
		logger.Debug(name, value)
		myReq.Header.Add(name, value)
	}

	return myReq, err
}

func (c *CurlClient) processResponse(response *http.Response) (string, RequestError) {
	logger.Debug("RESPONSE_STATUS:", response.StatusCode)
	logger.Debug("RESPONSE:", response.Status)
	logger.Debug("RESPONSE_HEADERS")
	for name, value := range response.Header {
		logger.Debug(name, value)
	}

	c.Cookies.SaveCookies(response.Cookies())
	responseError := c.result.ProcessResult(response)
	logger.Info(response.Cookies())
	//result.Stream = response.Body
	//result.Text = string(body)
	return c.result.Text, responseError
}

/* CurlClient methods */
/*
type Dyad struct {
	KeyName string
	Value   string
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

	currentOptions.Preprocess(myReq)

	return myReq, err
}

func DoRequest(url string, options *RequestOptions) (string, error) {
	result.Reset()

	currentUrl = url
	options.Process()
	req, err := PrepareRequestParameters()
	errs.ErrorHandle(err)
	client := &http.Client{}

	toDoReq := true
	var html string
	var reqErr error = nil
	var trials int
	trials = 0
	for toDoReq {
		logger.Debug("TRY: ", trials)
		resp, err := client.Do(req)
		errs.ErrorHandle(err)
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			ResponseHeaders(resp)
			reqErr = NewRequestError(resp, url)
			logger.Error(reqErr)
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

	return html, reqErr
}

func ResponseHeaders(response *http.Response) {
	for key, headers := range response.Header {
		logger.Debug(key)
		logger.Debug(headers)
	}
}
*/
