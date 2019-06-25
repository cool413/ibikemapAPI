package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ibikemap")
}

func TaichungStation(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]string)
	body["action"] = "ub_site_by_sno_class"
	body["datas[lang]"] = "tw"
	body["datas[loc]"] = "taichung"

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Referer"] = "https://i.youbike.com.tw/station/map"

	resp, err := Post("https://apis.youbike.com.tw/useAPI", nil, body, headers)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	w.Write(sitemap)
}

//Post http post method
func Post(url string, body map[string]string, params map[string]string, headers map[string]string) (*http.Response, error) {
	//add post body
	var bodyJson []byte
	var req *http.Request
	if body != nil {
		var err error
		bodyJson, err = json.Marshal(body)
		if err != nil {
			log.Println(err)
			return nil, errors.New("http post body to json failed")
		}
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail: %v \n")
	}
	req.Header.Set("Content-type", "application/json")
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{}
	log.Printf("Go POST URL : %s \n", req.URL.String())
	return client.Do(req)
}
