// options
package gocurl

import (
	log "logger"
	"math/rand"

	//"net/http"
	"strings"
	"time"
)

type RequestOptions struct {
	Url            string            //URL запроса
	method         string            //Метод запроса
	PostFields     map[string]string //Параметры POST-запроса
	CookieFile     string            //Имя файла с Cookies
	HttpHeaders    map[string]string //Заголовки запроса
	UserAgent      string            //User-Agent
	Trials         int               //Кол-во попыток
	Interval       float64           //Интервал между попытками
	timeout        time.Duration     //Время ожидания ответа на запрос
	FollowLocation bool              //Флаг управления редиректом на Location
	//Preprocess  func(req *http.Request)
}

var methods = [7]string{
	"GET",     //Запрашивает данные из указанного ресурса
	"POST",    // Отправляет данные для обработки указанному ресурсу
	"HEAD",    // То же, что и GET, но возвращает только заголовки HTTP, а не тело документа
	"PUT",     // Загружает представление указанного URI
	"DELETE",  // Удаляет указанный ресурс
	"OPTIONS", // Возвращает методы HTTP, поддерживаемые сервером
	"CONNECT", // Преобразует соединение запроса в прозрачный туннель TCP/IP
}

/**
 * Adds new field to POST-parameters
 */
func (options *RequestOptions) AddPostField(fieldName string, fieldValue string) {
	options.PostFields[fieldName] = fieldValue
}

/**
 * Adds new header to request
 */
func (options *RequestOptions) AddHeader(headerName string, headerValue string) {
	log.Debug(headerName+":", headerValue)
	options.HttpHeaders[headerName] = headerValue
}

/**
 * Adds headers to request
 */
func (options *RequestOptions) AddHeaders(headers map[string]string) {
	log.Debug("HEADERS")
	for name, value := range headers {
		options.AddHeader(name, value)
	}
}

/**
 * Set random User-Agent header to request
 */
func (options *RequestOptions) SetRandUserAgent() {
	options.UserAgent = useragents[rand.Intn(len(useragents))]
	log.Info("UA:", options.UserAgent)
}

/**
 * Set request method
 */
func (options *RequestOptions) SetMethod(method string) {
	method = strings.ToUpper(method)
	if !validMethod(method) {
		panic("Method " + method + " don't exists")
	}
	options.method = method
}

/**
 * Return request method
 */
func (options *RequestOptions) GetMethod() string {
	if len(options.method) == 0 {
		options.method = "GET"
	}
	return options.method
}

/**
 * Устанавливает время ожидания ответа на запрос
 */
func (options *RequestOptions) SetTimeout(timeout string) {
	value, _ := time.ParseDuration(timeout)
	options.timeout = value
}

/**
 * Возвращает время ожидания ответа на запрос
 */
func (options *RequestOptions) GetTimeout() time.Duration {
	return options.timeout
}

func validMethod(method string) bool {
	valid := false
	method = strings.ToUpper(method)
	for _, v := range methods {
		if v == method {
			valid = true
			break
		}
	}
	return valid
}

/* Methods *RequestOptions end */

/**
 * Возвращает инициализированные по умолчанию опции запроса
 */
func GetOptions() *RequestOptions {
	rand.Seed(time.Now().UnixNano())
	options := new(RequestOptions)
	options.Interval = 0
	options.PostFields = make(map[string]string)
	options.HttpHeaders = make(map[string]string)
	options.SetRandUserAgent()
	options.SetTimeout("60s")
	options.FollowLocation = false
	return options
}
