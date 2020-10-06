package controller

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
)

func getFromCache(cacheURL, key string) (io.ReadCloser, bool) {
	//resp, err := http.Get(cachServiceURL + "/?key=" + key)
	resp, err := http.Get(cacheURL + key)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, false
	}
	return resp.Body, true
}

func saveToCache(cacheURL, key string, duration int64, data []byte) {
	//req, _ := http.NewRequest(http.MethodPost, cachServiceURL+"/?key="+key,
	req, err := http.NewRequest(http.MethodPost, cacheURL+key,
		bytes.NewBuffer(data))
	if err != nil {
		log.Println("failed to save to cache:", err)
		return
	}
	req.Header.Add("cache-control", "maxage="+strconv.FormatInt(duration, 10))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error while htpp post:", err)
	}
	log.Printf("Successfuly POST on %s with response %+v", cacheURL+key, resp)
}

/*func invalidateCacheEntry(cacheURL, key string) {
	http.Get(cachServiceURL + "/invalidate" + key)
}*/
