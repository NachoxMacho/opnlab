package opnsense

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type OPNSenseResponse interface {
	DHCPLease | Interface
}

func getOPNSenseData[ResponseType OPNSenseResponse](path string) ([]ResponseType, error) {

	log.Println("https://" + opnsenseURL + "/api" + path)
	request, err := http.NewRequest("GET", "https://"+opnsenseURL+"/api"+path, nil)
	if err != nil {
		return nil, err
	}

	defaultTransport := http.DefaultTransport.(*http.Transport)

	// Create new Transport that ignores self-signed SSL
	customTransport := &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: customTransport}
	request.SetBasicAuth(opnsenseKey, opnsenseSecret)
	request.Header.Set("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	type OPNSenseData struct {
		Rows []ResponseType `json:"rows"`
	}

	res := OPNSenseData{}
	json.Unmarshal(bodyBytes, &res)
	return res.Rows, nil
}
