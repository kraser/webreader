// options
package webreader

import (
	"log"
	//"os"
)

type RequestOptions struct {
	Url         string
	Method      string
	PostFields  map[string]Dyad
	CookieFile  string
	HttpHeaders map[string]Dyad
	Test        string
}

var currentOptions = new(RequestOptions)

func (options *RequestOptions) AddPostField(fieldName string, fieldValue string) {
	field := Dyad{KeyName: fieldName, Value: fieldValue}
	options.PostFields[fieldName] = field
}

func (options *RequestOptions) AddHeader(headerName string, headerValue string) {
	header := Dyad{KeyName: headerName, Value: headerValue}
	options.HttpHeaders[headerName] = header
}

func GetOptions() *RequestOptions {
	currentOptions.Test = "THE TEST"
	currentOptions.PostFields = make(map[string]Dyad)
	currentOptions.HttpHeaders = make(map[string]Dyad)
	return currentOptions
}

func processOptions() {
	if len(currentOptions.Method) == 0 {
		currentOptions.Method = "GET"
	}

	if len(currentOptions.CookieFile) != 0 {
		log.Println("COOKIEFILE: ", currentOptions.CookieFile)
		cookieHandler.SetCookieFileName(currentOptions.CookieFile)
	}
}
