// cookies
package gocurl

import (
	errs "errorshandler"
	"logger"
	"net/http"
	"os"
)

type Cookies struct {
	Cookies            []*http.Cookie
	cookiesFileName    string
	cookiesFileHandler *os.File
}

func (handler *Cookies) SaveCookies(cookies []*http.Cookie) {
	handler.Cookies = cookies
	logger.Debug("COOKIEFILE", handler.cookiesFileName)
	if len(handler.cookiesFileName) != 0 {
		cookiesFileHandler, err := os.OpenFile(handler.cookiesFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		errs.ErrorHandle(err)
		defer cookiesFileHandler.Close()
		cookiesFileHandler.Truncate(0)
		logger.Debug("RESPONSE_COOKIES")
		for _, value := range handler.Cookies {
			logger.Debug(value.String())
			cookiesFileHandler.WriteString(value.String() + "\n")
		}
	}
}

func (handler *Cookies) SetCookieFileName(fileName string) {
	handler.cookiesFileName = fileName
}
