package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"go-hue-light/config"
)

type LightState struct {
	Brightness float64
	ColorX     float64
	ColorY     float64
}

type Lights map[string]string

func (l *Lights) String() string {
	var str string

	for name, rid := range *l {
		str += fmt.Sprintf("Name: %s, ResourceId: %s\n", name, rid)
	}

	return str
}

func CreateNewClientRequestBody() []byte {
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

func CreateLightStateRequestBody(state LightState) []byte {
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
		On On `json:"on"`
        Dimming Dimming `json:"dimming"`
        Color Color `json:"color"`
	}

	opt := Option{}

    if state.Brightness > 0 {
        opt.On.On = true
    }

    if state.Brightness == 0 {
        opt.On.On = false
    }

    opt.Dimming.Brightness = state.Brightness
    opt.Color.Xy.X = state.ColorX
    opt.Color.Xy.Y = state.ColorY

	requestBody, err := json.Marshal(opt)
	if err != nil {
		log.Println(err)
	}

    return requestBody
}

func SetLightState(c config.BridgeConfig, rId string, state LightState) {
	req, err := http.NewRequest("PUT",
		fmt.Sprintf("https://%s/clip/v2/resource/light/%s", c.Ip, rId),
		bytes.NewBuffer(CreateLightStateRequestBody(state)))
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

func GetDevices(c config.BridgeConfig) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/clip/v2/resource/device", c.Ip), nil)
	req.Header.Add("hue-application-key", c.ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	resp.Body.Close()

	return respBody
}

func SendNewClientRequest(ip string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Post(fmt.Sprintf("https://%s/%s", ip, "api"), "", bytes.NewBuffer(CreateNewClientRequestBody()))
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

func RegisterApiKey(ip string) (string, error) {
	resp, err := SendNewClientRequest(ip)
	if err != nil {
		return "", err
	}

	if resp != nil && IsLinkButtonResponse(resp) {
		log.Println("Please press the link button on the HUE Bridge, then press any button...")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		resp, err := SendNewClientRequest(ip)
		if err != nil {
			return "", err
		}

		return ParseNewUserResult(resp), nil
	}

	return "", errors.New("Error registering new api key")
}

func ParseNewUserResult(resp []byte) string {
	type SuccessData struct {
		Username  string `json:"username"`
		Clientkey string `json:"clientkey"`
	}

	type Success struct {
		Success SuccessData
	}

	var success []Success
	err := json.Unmarshal(resp, &success)
	if err != nil {
		log.Println("error parsing new user request", err)
		return ""
	}

	return success[0].Success.Username
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
		log.Fatal("error parsing error response")
		return false
	}

	if error[0].Error.Type == 101 || error[0].Error.Description == "link button not pressed" {
		return true
	}

	return false
}

func HandleEmptyConfig(c *config.BridgeConfig) {
	fmt.Print("Please input HUE Bridge IP: ")

	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
    ip := net.ParseIP(stdin.Text())
    if ip == nil {
        log.Fatal("IP out of range")
    }

    c.Ip = ip.To4().String()

	apiKey, err := RegisterApiKey(c.Ip)
	if err != nil {
		log.Fatal(err)
	}

	c.ApiKey = apiKey
}

func GetLights(config config.BridgeConfig) Lights {
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

	resp := GetDevices(config)

	err := json.Unmarshal(resp, &res)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	lights := Lights{}

	for _, light := range res.Data {
		for _, service := range light.Services {
			if service.Rtype == "light" {
				lights[light.Metadata.Name] = service.Rid
			}
		}
	}

	return lights
}

func HandleCommand(config config.BridgeConfig, lightName string, lightState LightState) {
	lights := GetLights(config)

	if len(lightName) == 0 {
		fmt.Println("Registered Lights:")
		fmt.Println(lights.String())
		return
	}

	rid, ok := lights[lightName]
	if !ok {
		log.Fatalf("Light with name \"%s\" is not registered", lightName)
	}

	SetLightState(config, rid, lightState)
}

func main() {
	config := config.BridgeConfig{}
	config.Load()

	if len(config.Ip) == 0 || len(config.ApiKey) == 0 {
		HandleEmptyConfig(&config)

		config.Save()
	}

	brightness := flag.Float64("br", 100.0, "Controls the brightness of the given light. [0 - 100]")
	colorX := flag.Float64("colorx", 0.0, "Controls the X Coordinate in the color diagram. [0.0 - 1.0]")
	colorY := flag.Float64("colory", 0.0, "Controls the Y Coordinate in the color diagram. [0.0 - 1.0]")
	lightName := flag.String("light", "", "Name of the light to control")

	flag.Parse()

	HandleCommand(config, *lightName, LightState{*brightness, *colorX, *colorY})
}
