package webreader

import (
	"io/ioutil"
	//"log"
	"logger"
	"net/http"
)

type Dyad struct {
	KeyName string
	Value   string
}

var cookieHandler = new(Cookies)
var currentUrl string
var client http.Client
var myReq http.Request

func errorHandle(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckRequestUrl(response *http.Response) string {
	logger.Debug("RESPONSE_HEADERS")
	for name, value := range response.Header {
		logger.Debug(name, value)
	}

	body, err := ioutil.ReadAll(response.Body)
	errorHandle(err)

	return string(body)
}

func PrepareRequestParameters() (*http.Request, error) {
	myReq, err := http.NewRequest(currentOptions.Method, currentUrl, nil)
	errorHandle(err)
	logger.Debug("REQUEST_HEADERS")
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
	currentUrl = url
	processOptions()
	req, err := PrepareRequestParameters()
	client := &http.Client{}

	resp, err := client.Do(req)
	errorHandle(err)
	defer resp.Body.Close()

	cookieHandler.SaveCookies(resp)

	return CheckRequestUrl(resp)
}

/*
Cache-Control
no-store, no-cache, must-reval…te, post-check=0, pre-check=0
Connection
keep-alive
Content-Encoding
gzip
Content-Type
text/html; UTF-8;charset=UTF-8
Date
Sun, 09 Sep 2018 08:51:55 GMT
Expires
Thu, 19 Nov 1981 08:52:00 GMT
Keep-Alive
timeout=30
Pragma
no-cache
Server
nginx-reuseport/1.13.4
Set-Cookie
CMS_SESSION=8b0f9f070aa7ecc72d…55 GMT; Max-Age=36000; path=/
Set-Cookie
authorized=989ff9c4541fd56fe2d…:55 GMT; Max-Age=3600; path=/
Set-Cookie
user=deleted; expires=Thu, 01-…:00:01 GMT; Max-Age=0; path=/
Transfer-Encoding
chunked
Vary
Accept-Encoding
X-Powered-By
PHP/5.6.30
*/
