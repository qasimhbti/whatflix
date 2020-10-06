package loghelper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/whatflix/entity"
)

var tr = http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}
var client = &http.Client{Transport: &tr}

func WriteEntry(logServiceURL string, entry *entity.LogEntry) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	_ = enc.Encode(entry)
	req, _ := http.NewRequest(http.MethodPost, logServiceURL, &buf)
	_, _ = client.Do(req)
}
