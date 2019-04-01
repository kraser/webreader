// options
package gocurl

import (
	log "logger"
	"math/rand"
	//"net/http"
	//"time"
)

type RequestOptions struct {
	Url         string
	method      string
	PostFields  map[string]string
	CookieFile  string
	HttpHeaders map[string]string
	UserAgent   string
	Trials      int
	Interval    float64
	//Preprocess  func(req *http.Request)
}

/* Methods *RequestOptions start */
//Adds new field to POST-parameters
func (options *RequestOptions) AddPostField(fieldName string, fieldValue string) {
	options.PostFields[fieldName] = fieldValue
}

func (options *RequestOptions) AddHeader(headerName string, headerValue string) {
	log.Debug(headerName+":", headerValue)
	options.HttpHeaders[headerName] = headerValue
}

func (options *RequestOptions) AddHeaders(headers map[string]string) {
	log.Debug("HEADERS")
	for name, value := range headers {
		options.AddHeader(name, value)
	}
}

func (options *RequestOptions) SetRandUserAgent() {
	options.UserAgent = useragents[rand.Intn(len(useragents))]
	log.Info("UA:", options.UserAgent)
}

func (options *RequestOptions) SetMethod(method string) {
	options.method = method
}

func (options *RequestOptions) GetMethod() string {
	if len(options.method) == 0 {
		options.method = "GET"
	}
	return options.method
}

/* Methods *RequestOptions end */

func GetOptions() *RequestOptions {
	options := new(RequestOptions)
	options.Interval = 0
	options.PostFields = make(map[string]string)
	options.HttpHeaders = make(map[string]string)
	options.SetRandUserAgent()
	return options
}
