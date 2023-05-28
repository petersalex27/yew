package error

import (
	"fmt"
	"os"
	"time"
	"yew/info"
)

// Log logs messages to arbitrary locations
type Log struct {
	message primaryMessageStruct
	includeTime bool
	file *os.File
}

func (log Log) GetLocation() info.Location {
	return log.message.GetLocation()
}

func (log Log) shouldAbort() bool {
	return log.message.shouldAbort()
}

func (log Log) ToString() (msg string) {
	msg = log.message.ToString()
	if log.includeTime {
		msg = "(" + time.Now().UTC().String() + ") " + msg
	}
	return msg
}

func (log Log) Print() int {
	res, err := fmt.Fprintf(log.file, "%s\n", log.ToString())
	if err != nil {
		SystemError(err.Error()).Print()
		return 0
	}
	return res
}