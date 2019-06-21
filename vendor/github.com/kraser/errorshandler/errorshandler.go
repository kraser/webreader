// errorshandler project errorshandler.go
package errorshandler

//log "logger"

func ErrorHandle(e error) {
	//log.Debug(e)
	if e != nil {
		panic(e)
	}
}
