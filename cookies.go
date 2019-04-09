// cookies
package gocurl

import (
	errs "errorshandler"
	"io/ioutil"
	"logger"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
		//handler.cookiesFileHandler = cookiesFileHandler
		cookiesFileHandler.Truncate(0)
		logger.Debug("RESPONSE_COOKIES")
		for _, value := range handler.Cookies {

			logger.Debug(value.String())
			cookiesFileHandler.WriteString(handler.CookieToString(value))
		}
	}
}

func (handler *Cookies) SetCookieFileName(fileName string) {
	handler.cookiesFileName = fileName
}

func (handler *Cookies) CookieToString(cookie *http.Cookie) string {
	var b strings.Builder
	b.WriteString(cookie.Name)
	b.WriteString("\t")
	b.WriteString(cookie.Value)
	b.WriteString("\t")
	b.WriteString(cookie.Path)
	b.WriteString("\t")
	b.WriteString(cookie.Domain)
	b.WriteString("\t")
	b.WriteString(cookie.Expires.String())
	b.WriteString("\t")
	b.WriteString(cookie.RawExpires)
	b.WriteString("\t")
	b.WriteString(strconv.Itoa(cookie.MaxAge))
	b.WriteString("\t")
	b.WriteString(strconv.FormatBool(cookie.Secure))
	b.WriteString("\t")
	b.WriteString(strconv.FormatBool(cookie.HttpOnly))
	b.WriteString("\t")
	b.WriteString(strconv.Itoa(int(cookie.SameSite)))
	b.WriteString("\t")
	b.WriteString(cookie.Raw)
	b.WriteString("\n")

	return b.String()
}

func (handler *Cookies) ReadCookies() {
	if _, err := os.Stat(handler.cookiesFileName); os.IsNotExist(err) {
		return
	}
	data, err := ioutil.ReadFile(handler.cookiesFileName)
	errs.ErrorHandle(err)
	cookieStrings := strings.Split(string(data), "\n")
	handler.Cookies = make([]*http.Cookie, len(cookieStrings)-1)
	for i, cookieString := range cookieStrings {
		if len(cookieString) == 0 {
			continue
		}
		parts := strings.Split(cookieString, "\t")
		cookie := &http.Cookie{}
		cookie.Name = parts[0]
		cookie.Value = parts[1]
		cookie.Path = parts[2]
		cookie.Domain = parts[3]
		cookie.Expires, _ = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", parts[4])
		cookie.RawExpires = parts[5]
		cookie.MaxAge, _ = strconv.Atoi(parts[6])
		cookie.Secure, _ = strconv.ParseBool(parts[7])
		cookie.HttpOnly, _ = strconv.ParseBool(parts[8])
		sameSite, _ := strconv.Atoi(parts[9])
		cookie.SameSite = http.SameSite(sameSite)
		cookie.Raw = parts[10]
		handler.Cookies[i] = cookie
	}
}

func (handler *Cookies) ActualCookiesRaw() string {
	cookies := make([]string, 0 /*len(handler.Cookies)*/)
	for _, cookie := range handler.Cookies {
		var raw strings.Builder
		raw.WriteString(cookie.Name)
		raw.WriteString("=")
		raw.WriteString(cookie.Value)
		cookies = append(cookies, raw.String())
	}
	return strings.Join(cookies, ";")
}
