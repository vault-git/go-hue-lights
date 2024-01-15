package main

import (
	"testing"
)

func TestCreateLightPropertiesRequest(t *testing.T) {
    var tests = []struct{
        input LightProps
        want string
    }{
        {LightProps{true, 100.0, 0.7512, 0.7541}, `{"on":{"on":true},"dimming":{"brightness":100},"color":{"xy":{"x":0.7512,"y":0.7541}}}`},
        {LightProps{false, 0.0, 0.0, 0.0}, `{"on":{"on":false},"dimming":{"brightness":0},"color":{"xy":{"x":0,"y":0}}}`},
    }

    for _, test := range tests {
        if got := CreateLightPropertiesRequest(test.input); string(got) != test.want {
            t.Errorf("CreateLightPropertiesRequest(%v): got %v wanted %v", test.input, string(got), test.want)
        }
    }
}

func TestCreateNewClientRequest(t *testing.T) {
    correctOutput := string([]byte(`{"devicetype":"go-hue-lights","generateclientkey":true}`))
    if string(CreateNewClientRequest()) != correctOutput {
        t.Errorf("CreateNewClientRequest(): wrong output")
    }
}
