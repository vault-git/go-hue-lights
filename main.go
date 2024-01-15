package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func RegisterApiKey(ip string) (string, error) {
	resp, err := PostNewClient(ip)
	if err != nil {
		return "", err
	}

	if resp != nil && IsLinkButtonResponse(resp) {
		log.Println("Please press the link button on the HUE Bridge, then press any button...")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		resp, err := PostNewClient(ip)
		if err != nil {
			return "", err
		}

		return ParseNewUserResult(resp)
	}

	return "", errors.New("Error registering new api key")
}

func HandleEmptyConfig(c *BridgeConfig) {
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

func SetNewLightProps(currentProps LightProps, newProps *LightProps) {
	if newProps.Brightness == -1.0 {
		newProps.Brightness = currentProps.Brightness
	}

	if newProps.Brightness >= 5.0 {
		newProps.On = true
	}

	if newProps.ColorX == -1.0 {
		newProps.ColorX = currentProps.ColorX
	}

	if newProps.ColorY == -1.0 {
		newProps.ColorY = currentProps.ColorY
	}
}

func PrintAllLights(config BridgeConfig) {
	lights := ParseDeviceResource(GetDeviceResource(config))

	for name, rid := range lights {
		light := ParseLightResource(GetLightResource(config, rid))

		fmt.Printf("%s: %s", name, light.String())
	}
}

func ControlLights(config BridgeConfig, lightName string, newLightProps LightProps) {
	lights := ParseDeviceResource(GetDeviceResource(config))

	rid, ok := lights[lightName]
	if !ok {
		log.Fatalf("Light with name \"%s\" is not registered", lightName)
	}

	currentLightProps := ParseLightResource(GetLightResource(config, rid))

	SetNewLightProps(currentLightProps, &newLightProps)

	PutLightResource(config, rid, newLightProps)
}

func main() {
	config := BridgeConfig{}
	config.Load()

	if len(config.Ip) == 0 || len(config.ApiKey) == 0 {
		HandleEmptyConfig(&config)

		config.Save()
	}

	listLights := flag.Bool("list", false, "Lists all registered lights.")
	lightName := flag.String("light", "", "Name of the light to control.")
	brightness := flag.Float64("br", -1.0, "Controls the brightness of the given light. [0 - 100]")
	colorX := flag.Float64("colorx", -1.0, "Controls the X Coordinate in the color diagram. [0.0 - 1.0]")
	colorY := flag.Float64("colory", -1.0, "Controls the Y Coordinate in the color diagram. [0.0 - 1.0]")

	flag.Parse()

	if *listLights {
        PrintAllLights(config)
        return
	}

    if len(*lightName) == 0 {
        flag.Usage()
        return
    }

	ControlLights(config, *lightName, LightProps{false, *brightness, *colorX, *colorY})
}
