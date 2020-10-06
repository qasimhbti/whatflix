package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/whatflix/entity"
	"github.com/whatflix/logservice/loghelper"
)

var (
	appservers   = []string{}
	currentIndex = 0
	client       = http.Client{Transport: &transport}
)

func processRequests(logServiceURL string) {
	for {
		select {
		case request := <-requestCh:
			log.Println("request")
			if len(appservers) == 0 {
				request.w.WriteHeader(http.StatusInternalServerError)
				_, _ = request.w.Write([]byte("No app servers found"))
				request.doneCh <- struct{}{}
				continue
			}
			currentIndex++
			if currentIndex == len(appservers) {
				currentIndex = 0
			}
			host := appservers[currentIndex]
			go processRequest(host, request)
		case host := <-registerCh:
			log.Println("register " + host)
			go loghelper.WriteEntry(logServiceURL, &entity.LogEntry{
				Level:     entity.LogLevelInfo,
				Timestamp: time.Now(),
				Source:    "load balancer",
				Message:   "Registering application server with address: " + host,
			})
			isFound := false
			for _, h := range appservers {
				if host == h {
					isFound = true
					break
				}
			}

			if !isFound {
				appservers = append(appservers, host)
			}
		case host := <-unregisterCh:
			log.Println("unregister " + host)
			go loghelper.WriteEntry(logServiceURL, &entity.LogEntry{
				Level:     entity.LogLevelInfo,
				Timestamp: time.Now(),
				Source:    "load balancer",
				Message:   "Unregistering application server with address: " + host,
			})
			for i := len(appservers) - 1; i >= 0; i-- {
				if appservers[i] == host {
					appservers = append(appservers[:i], appservers[i+1:]...)
				}
			}
		case <-heartbeatCh:
			log.Println("heartbeat...")
			servers := appservers[:]
			go func(servers []string) {
				for _, server := range servers {
					//resp, err := http.Get("https://" + server + "/ping")
					resp, err := http.Get("http://" + server + "/ping")
					if err != nil || resp.StatusCode != 200 {
						unregisterCh <- server
					}
					log.Println("success")
				}
			}(servers)
		}
	}
}

func processRequest(host string, request *webRequest) {
	hostURL, _ := url.Parse(request.r.URL.String())
	//hostURL.Scheme = "https"
	hostURL.Scheme = "http"
	hostURL.Host = host
	log.Println(hostURL.String())
	req, _ := http.NewRequest(request.r.Method, hostURL.String(), request.r.Body)
	for k, v := range request.r.Header {
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}
		req.Header.Add(k, values)
	}

	resp, err := client.Do(req)
	if err != nil {
		request.w.WriteHeader(http.StatusInternalServerError)
		request.doneCh <- struct{}{}
		return
	}

	for k, v := range resp.Header {
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}
		request.w.Header().Add(k, values)
	}
	_, err = io.Copy(request.w, resp.Body)
	if err != nil {
		log.Printf("error while copy: %v\n", err)
	}
	request.doneCh <- struct{}{}
}
