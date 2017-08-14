package main

import (
	"testing"
	"time"

	
)

type mockOutputter struct {
	called bool
}

func (mo *mockOutputter) output(_ time.Time, _ []logEntry) {
	mo.called = true

}

func TestAddBucket(t *testing.T) {
	mo := new(mockOutputter)
	bck := newBucket(mo)
	entries := []string{
		`2016-05-07T09:07:00.002979+00:00 heroku[router]: at=info method=HEAD path="/api" host="debora.com" request_id=0e787a0d-9cba-451c-9a44-b4009009dd60 fwd="10.29.10.29" dyno=web.4 connect=1ms service=282ms status=200 bytes=22`,
		`2016-05-07T09:07:00.001490+00:00 heroku[router]: at=info method=GET path="/blog" host="brs.org" request_id=fc693802-8851-484e-aab4-1d013714b68b fwd="10.29.10.29" dyno=web.3 connect=2ms service=994ms status=200 bytes=552`,
		`2016-05-07T21:06:59.000000+00:00 heroku[router]: host="feds.net" status=200 fwd="10.29.10.29" dyno=web.3 at=info connect=0ms bytes=726 request_id=21c36d4b-ba45-4093-96be-8aa0e76526a1 service=342ms method=GET path="/api`,
	}
	for i, ent := range entries {
		err := bck.add(ent)
		if err != nil {
			t.Fatal("unexpected error")
		}
		if i < 2 {
			if mo.called {
				t.Fatal("output shouldn't have happened")
			}
			
			if len(bck.entries) != (i + 1) {
				t.Fatal("entry count wrong")
			}
		} else {
			if !mo.called {
				t.Fatal("output should have happened")
			}
			if len(bck.entries) != 3 {
				t.Fatal("entry should reset to 3")
			}

		}
	}
}
