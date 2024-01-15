package main

import (
	"encoding/json"
	"log"
)

func CreateLightPropertiesRequest(lp LightProps) []byte {
	type On struct {
		On bool `json:"on"`
	}

	type Dimming struct {
		Brightness float64 `json:"brightness"`
	}

	type Xy struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	type Color struct {
		Xy Xy `json:"xy"`
	}

	type Option struct {
		On      On      `json:"on"`
		Dimming Dimming `json:"dimming"`
		Color   Color   `json:"color"`
	}

	opt := Option{On{lp.On}, Dimming{lp.Brightness}, Color{Xy{lp.ColorX, lp.ColorY}}}

	requestBody, err := json.Marshal(opt)
	if err != nil {
		log.Println(err)
	}

	return requestBody
}

func CreateNewClientRequest() []byte {
	type Client struct {
		DeviceType        string `json:"devicetype"`
		GenerateClientKey bool   `json:"generateclientkey"`
	}

	requestBody, err := json.Marshal(Client{"go-hue-lights", true})
	if err != nil {
		return nil
	}

	return requestBody
}
