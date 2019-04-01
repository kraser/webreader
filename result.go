// result.go
package gocurl

import (
	"io"
	"io/ioutil"
	"net/http"

	errs "github.com/kraser/errorshandler"
)

type RequestResult struct {
	response *http.Response
	Text     string
	stream   io.Reader
}

func (res *RequestResult) ProcessResult(response *http.Response) RequestError {
	res.response = response
	defer response.Body.Close()
	res.stream = response.Body
	body, err := ioutil.ReadAll(response.Body)
	errs.ErrorHandle(err)
	res.Text = string(body)
	reqErrror := RequestError{}
	return reqErrror
}

func (res *RequestResult) Reset() {
	res.Text = ""
	res.stream = nil
}
