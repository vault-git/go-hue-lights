package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
    deviceResourceUrl = "https://%s/clip/v2/resource/device"
    lightResourceUrl = "https://%s/clip/v2/resource/light/%s"
)

func getHueApiResource(c BridgeConfig, url string) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	req.Header.Add("hue-application-key", c.ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}

    defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	return respBody
}

func GetDeviceResource(c BridgeConfig) []byte {
    return getHueApiResource(c, fmt.Sprintf(deviceResourceUrl, c.Ip))
}

func GetLightResource(c BridgeConfig, rId string) []byte {
    return getHueApiResource(c, fmt.Sprintf(lightResourceUrl, c.Ip, rId))
}

func PostNewClient(ip string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Post(fmt.Sprintf("https://%s/api", ip), "", bytes.NewBuffer(CreateNewClientRequest()))
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, err
	}

	return respBody, nil
}

func PutLightResource(c BridgeConfig, rId string, lp LightProps) {
    lightProps := bytes.NewBuffer(CreateLightPropertiesRequest(lp))
	req, err := http.NewRequest("PUT", fmt.Sprintf(lightResourceUrl, c.Ip, rId), lightProps)
    if err != nil {
        log.Println(err)
        return
    }

	req.Header.Add("hue-application-key", c.ApiKey)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	_, err = client.Do(req)
	if err != nil {
		log.Println(err)
	}
}
