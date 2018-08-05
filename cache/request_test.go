package cache

import (
	"testing"
)

func Test_Request_Log(t *testing.T) {
	t.Parallel()

	fakeReq := []byte("ok this is\nthe\nrequest\tlogged")
	if err := LogWebRequest("testing", fakeReq); err != nil {
		t.Fatal(err)
	}

	// we want the last one inserted so we can test
	reqID, b, err := GetWebRequest(false)
	if err != nil {
		t.Error(err)
	} else if reqID != "testing" {
		t.Errorf("reqID was %s and we were expecting testing", reqID)
	} else if string(b) != string(fakeReq) {
		t.Errorf("the returned request was %s and we were looking for %s", string(b), string(fakeReq))
	}
}
