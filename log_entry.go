package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type logEntry struct {
	timeStamp   time.Time
	host        string
	serviceTime int
}

type invalidLogFormat struct {
	logLine string
}

func (ilf *invalidLogFormat) Error() string {
	return fmt.Sprintf("log format invalid: %q", ilf.logLine)
}

const (
	hostKey    = "host"
	serviceKey = "service"
)

func newLogEntry(logLine string) (*logEntry, error) {
	var le logEntry
	parts := strings.Split(logLine, " ")
	if len(parts) < 3 {
		return nil, &invalidLogFormat{logLine}
	}
	ts, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, &invalidLogFormat{logLine}
	}
	le.timeStamp = ts

	for _, part := range parts {
		pair := strings.Split(part, "=")
		if len(pair) == 2 {
			switch pair[0] {
			case hostKey:
				if len(pair) < 2 {
					return nil, &invalidLogFormat{logLine}
				}
				// trim quotes
				le.host = pair[1][1 : len(pair[1])-1]
			case serviceKey:
				// get rid of the 'ms' suffix and convert to int
				strval := pair[1]
				strval = strval[:len(strval)-2]
				val, err := strconv.Atoi(strval)
				if err != nil {
					return nil, &invalidLogFormat{logLine}
				}
				le.serviceTime = val
			}
		}
	}
	return &le, nil

}
