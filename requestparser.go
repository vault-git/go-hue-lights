package main

import (
	"encoding/json"
	"fmt"
)

type LightProps struct {
	On         bool
	Brightness float64
	ColorX     float64
	ColorY     float64
}

func (lp *LightProps) String() string {
	state := func() string {
		if lp.On {
			return "On"
		}

		return "Off"
	}()

	return fmt.Sprintf("State: %s, Brightness: %0.1f, ColorCoords: %0.4f, %0.4f\n", state, lp.Brightness, lp.ColorX, lp.ColorY)
}

type LightDevices map[string]string

func ParseLightResource(buf []byte) LightProps {
	type On struct {
		On bool
	}

	type Dimming struct {
		Brightness float64
	}

	type Xy struct {
		X float64
		Y float64
	}

	type Color struct {
		Xy Xy
	}

	type Data struct {
		Dimming Dimming
		On      On
		Color   Color
	}

	type Result struct {
		Errors []string
		Data   []Data
	}

	var res = Result{}

	err := json.Unmarshal(buf, &res)
	if err != nil {
		fmt.Println(err)
		return LightProps{}
	}

	lightData := res.Data[0]

	return LightProps{lightData.On.On, lightData.Dimming.Brightness, lightData.Color.Xy.X, lightData.Color.Xy.Y}
}

func ParseDeviceResource(buf []byte) LightDevices {
	type Service struct {
		Rid   string
		Rtype string
	}

	type Data struct {
		Services []Service
		Metadata struct {
			Name string
		}
	}

	type Result struct {
		Errors []string
		Data   []Data
	}

	var res = Result{}

	err := json.Unmarshal(buf, &res)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	lights := LightDevices{}

	for _, light := range res.Data {
		for _, service := range light.Services {
			if service.Rtype == "light" {
				lights[light.Metadata.Name] = service.Rid
			}
		}
	}

	return lights
}

func IsLinkButtonResponse(resp []byte) bool {
	type ErrorData struct {
		Type        int
		Address     string
		Description string
	}

	type Error struct {
		Error ErrorData
	}

	var error []Error
	err := json.Unmarshal(resp, &error)
	if err != nil {
		return false
	}

	if error[0].Error.Type == 101 || error[0].Error.Description == "link button not pressed" {
		return true
	}

	return false
}

func ParseNewUserResult(resp []byte) (string, error) {
	type SuccessData struct {
		Username  string
		Clientkey string
	}

	type Success struct {
		Success SuccessData
	}

	var success []Success
	err := json.Unmarshal(resp, &success)
	if err != nil {
        return "", err
	}

	return success[0].Success.Username, nil
}

