package gocurl

import (
	"logger"

	"net/http"
	//"net/url"
	//"sync"
	"time"

	errs "errorshandler"
)

var client = new(CurlClient)

type CurlClient struct {
	url     string
	Options *RequestOptions
	request http.Request
	result  *RequestResult
	*Cookies
	requestTime time.Time
	httpClient  *http.Client
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
		client.Cookies.ReadCookies()
	}
	client.httpClient = &http.Client{
		Jar: NewCurlJar(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return client
}

/* CurlClient methods */
func (c *CurlClient) DoRequest(url string) string {
	c.result.Reset()
	c.url = url
	req, err := c.PrepareRequestParameters()
	errs.ErrorHandle(err)

	toContinue := true
	var requestErr RequestError
	var result string
	var trials int
	trials = 0
	for toContinue {
		current := time.Now()
		elapsed := current.Sub(c.requestTime).Seconds()
		if elapsed < c.Options.Interval {
			time.Sleep(time.Duration(c.Options.Interval-elapsed) * time.Second)
		}
		c.requestTime = time.Now()
		logger.Debug("TRY: ", trials)
		resp, err := c.httpClient.Do(req)
		errs.ErrorHandle(err)
		defer resp.Body.Close()
		result, requestErr = c.processResponse(resp)
		if requestErr.IsNull() {
			trials++
			toContinue = trials <= c.Options.Trials

			//вынести в ParserHandleError
			//c.Options.SetRandUserAgent()
			//req.Header.Set("User-Agent", c.Options.UserAgent)
			//вынести в ParserHandleError

		} else {
			toContinue = false
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
	logger.Debug("COOKIES:", c.Cookies.ActualCookiesRaw())
	myReq.Header.Set("Cookie", c.Cookies.ActualCookiesRaw())

	return myReq, err
}

func (c *CurlClient) processResponse(response *http.Response) (string, RequestError) {
	logger.Debug("RESPONSE_STATUS:", response.StatusCode)
	logger.Debug("RESPONSE:", response.Status)
	logger.Debug(response.Request.URL)
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
