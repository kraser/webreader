package webreader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//"github.com/opesun/goquery"
)

type RequestOptions struct {
	Method     string
	PosrFields map[string]string
}

var client http.Client
var myReq http.Request

//var request http.Request

func errorHandle(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckRequestUrl(response *http.Response) {

	body, err := ioutil.ReadAll(response.Body)
	errorHandle(err)
	defer response.Body.Close()
	fmt.Println(string(body))
}

func PrepareRequestParameters(url string, options RequestOptions) (*http.Request, error) {
	if len(options.Method) == 0 {
		options.Method = "get"
	}
	myReq, err := http.NewRequest(options.Method, url, nil)
	fmt.Println(myReq.Method)
	fmt.Println(myReq.URL)
	errorHandle(err)
	myReq.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	return myReq, err
}

func DoRequest(url string, options RequestOptions) string {
	fmt.Println(url)
	client := &http.Client{}
	req, err := PrepareRequestParameters(url, options)
	fmt.Println(req.Method)
	fmt.Println("URL", req.URL)
	resp, err := client.Do(req)
	errorHandle(err)
	CheckRequestUrl(resp)
	return "test"
}
