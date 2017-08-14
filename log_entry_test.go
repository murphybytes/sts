package main

import (
	"testing"
)

func TestLogEntry(t *testing.T) {
	line := `2016-05-07T09:07:00.001490+00:00 heroku[router]: at=info method=GET path="/blog" host="brs.org" request_id=fc693802-8851-484e-aab4-1d013714b68b fwd="10.29.10.29" dyno=web.3 connect=2ms service=994ms status=200 bytes=552`
	le, err := newLogEntry(line)
	if err != nil {
		t.Fatal("error on valid log line")
	}
	if le == nil {
		t.Fatal("expected logEntry not nil")
	}
	if le.timeStamp.Hour() != 9 {
		t.Fatal("timestamp parse failed")
	}
	if le.serviceTime != 994 {
		t.Fatal("service parse failed")
	}
	if le.host != "brs.org" {
		t.Fatal("host parse failed")
	}
}
