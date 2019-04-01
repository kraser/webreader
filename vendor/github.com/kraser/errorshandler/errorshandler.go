// errorshandler project errorshandler.go
package errorshandler
import (
	log "logger"
)
func ErrorHandle(e error) {
	log.Debug(e)
	if e != nil {
		panic(e)
	}
}
