package gocurl

import (
	"logger"

	"net/http"
	"net/url"
	"sync"
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
	client.httpClient = &http.Client{Jar: NewCurlJar()}

	return client
}

/* CurlClient methods */
func (c *CurlClient) DoRequest(url string) string {
	c.result.Reset()
	c.url = url
	req, err := c.PrepareRequestParameters()
	errs.ErrorHandle(err)

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
		resp, err := c.httpClient.Do(req)
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
	//myReq.Header.Set("Cookie", c.Cookies.ActualRaws())
	logger.Debug("COOKIES:", c.Cookies.ActualCookiesRaw())
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

type CurlJar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewCurlJar() *CurlJar {
	curlJar := new(CurlJar)
	curlJar.cookies = make(map[string][]*http.Cookie)
	return curlJar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *CurlJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *CurlJar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}
