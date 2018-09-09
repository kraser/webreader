package webreader

import (
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	//"github.com/opesun/goquery"
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

	body, err := ioutil.ReadAll(response.Body)
	errorHandle(err)

	return string(body)
}

func PrepareRequestParameters() (*http.Request, error) {
	myReq, err := http.NewRequest(currentOptions.Method, currentUrl, nil)
	errorHandle(err)
	log.Println("HEADER_ACCEPT", currentOptions.HttpHeaders["Accept"].Value)
	myReq.Header.Add("Accept", currentOptions.HttpHeaders["Accept"].Value)
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
