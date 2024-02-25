package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ProxmoxResponses interface {
	[]APINode | []APIVM | []APIStorage | APIVMConfig
}

func getProxmoxData[ResponseType ProxmoxResponses](path string) (ResponseType, error) {

	var resp ResponseType

	log.Println("https://" + proxmoxURL + "/api2/json" + path)
	request, err := http.NewRequest("GET", "https://"+proxmoxURL+"/api2/json"+path, nil)
	if err != nil {
		return resp, err
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
	request.Header.Set("Authorization", "PVEAPIToken="+proxmoxToken)
	request.Header.Set("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return resp, err
	}

	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return resp, err
	}

	type proxmoxData struct {
		Data ResponseType
	}

	res := proxmoxData{}
	json.Unmarshal(bodyBytes, &res)
	return res.Data, nil
}
