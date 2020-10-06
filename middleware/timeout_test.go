package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTimeoutMiddleware(t *testing.T) {
	server := httptest.NewServer(&TimeoutMiddleware{})
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Get / err = %s; want nil", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil ReadAll / err = %s; want nil", err)
	}
	//t.Log("res :", string(res))
}
