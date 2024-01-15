package main

import (
	"testing"
)

func TestParseDeviceResource(t *testing.T) {
	testBuf := []byte(`{"errors":[],"data":[{"id":"33f98776-9d82-49cd-8a2e-181a150b8415","product_data":{"model_id":"BSB002","manufacturer_name":"Signify Netherlands B.V.","product_name":"Hue Bridge","product_archetype":"bridge_v2","certified":true,"software_version":"1.62.1962097030"},"metadata":{"name":"Hue Bridge","archetype":"bridge_v2"},"identify":{},"services":[{"rid":"a845c168-8043-451a-b644-6b7930c3aa2c","rtype":"bridge"},{"rid":"3422e9cd-2c2b-48c2-94a7-9bf4a9afc8b8","rtype":"zigbee_connectivity"},{"rid":"6a415ae5-96fc-43a0-bf30-5bd136c014e7","rtype":"entertainment"},{"rid":"e2fdb013-a3bd-4c20-8e62-e7ae8deed27d","rtype":"zigbee_device_discovery"}],"type":"device"},{"id":"49e3d193-c4c6-4525-8bc9-456511567bf5","id_v1":"/lights/3","product_data":{"model_id":"LCA004","manufacturer_name":"Signify Netherlands B.V.","product_name":"Hue color lamp","product_archetype":"sultan_bulb","certified":true,"software_version":"1.104.2","hardware_platform_type":"100b-114"},"metadata":{"name":"hue_color_1","archetype":"table_shade"},"identify":{},"services":[{"rid":"b0d8ca80-f421-4976-ba8c-ba2026ed4224","rtype":"zigbee_connectivity"},{"rid":"8535734d-7136-45e6-a551-70b5f2782fa8","rtype":"light"},{"rid":"9bca5549-4d9f-49d9-9f33-5761b616306a","rtype":"entertainment"},{"rid":"8ed00696-944e-46a1-b30d-5fe6285468e3","rtype":"taurus_7455"},{"rid":"4415a2f4-893e-4fb0-82f7-d8b98b36f138","rtype":"device_software_update"}],"type":"device"},{"id":"6542ce87-379d-4c45-8548-2992c6251533","id_v1":"/sensors/15","product_data":{"model_id":"ROM001","manufacturer_name":"Signify Netherlands B.V.","product_name":"Hue Smart button","product_archetype":"unknown_archetype","certified":true,"software_version":"2.47.8","hardware_platform_type":"100b-116"},"metadata":{"name":"Hue Smart button 1","archetype":"unknown_archetype"},"services":[{"rid":"63346ff3-83d7-4686-abf8-fcf3d6248679","rtype":"zigbee_connectivity"},{"rid":"dc825a6e-1331-46e6-9e55-6f9f47ad413f","rtype":"button"},{"rid":"91c50913-8eae-4848-b585-ac666464731b","rtype":"device_power"},{"rid":"9b359495-e2a1-4296-b213-9075ee9c8a5b","rtype":"device_software_update"}],"type":"device"},{"id":"9195f3ad-48e7-44e5-8894-4c8855eae5c2","id_v1":"/lights/1","product_data":{"model_id":"LWA017","manufacturer_name":"Signify Netherlands B.V.","product_name":"Hue white lamp","product_archetype":"sultan_bulb","certified":true,"software_version":"1.104.2","hardware_platform_type":"100b-114"},"metadata":{"name":"white_lamp_1","archetype":"single_spot"},"identify":{},"services":[{"rid":"7db433d3-1da0-41a7-8d4f-2ac2bb8c55aa","rtype":"zigbee_connectivity"},{"rid":"c873dbb6-aae7-44b2-b0b9-e1ef3992cec7","rtype":"light"},{"rid":"4649a11a-05ba-45bb-92a3-465f6e7f326e","rtype":"taurus_7455"},{"rid":"b6988c09-ce55-47c1-8971-e4bd0359375b","rtype":"device_software_update"}],"type":"device"},{"id":"b20d9992-18cb-407b-a4f9-23afd3296fe9","id_v1":"/lights/2","product_data":{"model_id":"LWA017","manufacturer_name":"Signify Netherlands B.V.","product_name":"Hue white lamp","product_archetype":"sultan_bulb","certified":true,"software_version":"1.104.2","hardware_platform_type":"100b-114"},"metadata":{"name":"white_lamp_2","archetype":"single_spot"},"identify":{},"services":[{"rid":"0114d615-5605-4e72-8fa7-69404e79ebc8","rtype":"zigbee_connectivity"},{"rid":"2ea13e3b-0401-4396-bcb7-d58d2c542388","rtype":"light"},{"rid":"228e9079-f7de-4bb1-8e43-8c9312c32cc9","rtype":"taurus_7455"},{"rid":"07f2d8a0-f842-43ab-a2ea-850c11c7bea2","rtype":"device_software_update"}],"type":"device"}]}`)

	lights := ParseDeviceResource(testBuf)

	if len(lights) != 3 {
		t.Errorf("ParseDeviceResource(%q) = %v", testBuf, len(lights))
	}

	for _, lightName := range []string{"white_lamp_1", "white_lamp_2", "hue_color_1"} {
		rid, ok := lights[lightName]
		if !ok {
			t.Errorf("ParseDeviceResource(%q) lamp does not exist", testBuf)
		}

		if rid != "8535734d-7136-45e6-a551-70b5f2782fa8" &&
           rid != "c873dbb6-aae7-44b2-b0b9-e1ef3992cec7" &&
           rid != "2ea13e3b-0401-4396-bcb7-d58d2c542388" {
			t.Errorf("ParseDeviceResource(%q) wrong resource id for light", testBuf)
		}
	}
}

func TestParseLightResource(t *testing.T) {
	var tests = []struct {
		input []byte
		want  LightProps
	}{
		// color lamp
		{[]byte(`{"errors":[],"data":[{"id":"8535734d-7136-45e6-a551-70b5f2782fa8","id_v1":"/lights/3","owner":{"rid":"49e3d193-c4c6-4525-8bc9-456511567bf5","rtype":"device"},"metadata":{"name":"hue_color_1","archetype":"table_shade","function":"mixed"},"product_data":{"function":"mixed"},"identify":{},"on":{"on":true},"dimming":{"brightness":50.2,"min_dim_level":0.20000000298023225},"dimming_delta":{},"color_temperature":{"mirek":497,"mirek_valid":true,"mirek_schema":{"mirek_minimum":153,"mirek_maximum":500}},"color_temperature_delta":{},"color":{"xy":{"x":0.5249,"y":0.4136},"gamut":{"red":{"x":0.6915,"y":0.3083},"green":{"x":0.17,"y":0.7},"blue":{"x":0.1532,"y":0.0475}},"gamut_type":"C"},"dynamics":{"status":"none","status_values":["none","dynamic_palette"],"speed":0.0,"speed_valid":false},"alert":{"action_values":["breathe"]},"signaling":{"signal_values":["no_signal","on_off","on_off_color","alternating"]},"mode":"normal","effects":{"status_values":["no_effect","candle","fire","prism"],"status":"no_effect","effect_values":["no_effect","candle","fire","prism"]},"powerup":{"preset":"safety","configured":true,"on":{"mode":"on","on":{"on":true}},"dimming":{"mode":"dimming","dimming":{"brightness":100.0}},"color":{"mode":"color_temperature","color_temperature":{"mirek":366}}},"type":"light"}]}`),
			LightProps{true, 50.2, 0.5249, 0.4136}},
		// white lamp
		{[]byte(`{"errors":[],"data":[{"id":"c873dbb6-aae7-44b2-b0b9-e1ef3992cec7","id_v1":"/lights/1","owner":{"rid":"9195f3ad-48e7-44e5-8894-4c8855eae5c2","rtype":"device"},"metadata":{"name":"white_lamp_1","archetype":"single_spot","fixed_mired":366,"function":"functional"},"product_data":{"function":"functional"},"identify":{},"on":{"on":false},"dimming":{"brightness":100.0,"min_dim_level":5.0},"dimming_delta":{},"dynamics":{"status":"none","status_values":["none"],"speed":0.0,"speed_valid":false},"alert":{"action_values":["breathe"]},"signaling":{"signal_values":["no_signal","on_off"]},"mode":"normal","effects":{"status_values":["no_effect","candle"],"status":"no_effect","effect_values":["no_effect","candle"]},"powerup":{"preset":"safety","configured":true,"on":{"mode":"on","on":{"on":true}},"dimming":{"mode":"dimming","dimming":{"brightness":100.0}}},"type":"light"}]}`),
			LightProps{false, 100.0, 0.0, 0.0}},
	}

	for _, test := range tests {
		if got := ParseLightResource(test.input); got != test.want {
			t.Errorf("ParseLightResource(%q) == %v", test.input, test.want)
		}
	}
}

func TestIsLinkButtonResponse(t *testing.T) {
	var tests = []struct {
		input []byte
		want  bool
	}{
        {[]byte(`[{"error":{"type": 101, "address": "", "description": "link button not pressed"}}]`), true},
        {[]byte(`[{"error":{"type": 100, "address": "", "description": "other error"}}]`), false},
        {[]byte(`{"invalid":"stuff"}`), false},
	}

	for _, test := range tests {
		if got := IsLinkButtonResponse(test.input); got != test.want {
			t.Errorf("IsLinkButtonResponse(%q) == %v", test.input, test.want)
		}
	}
}

func TestParseNewUserResponse(t *testing.T) {
	var tests = []struct {
		input []byte
		want  string
	}{
        {[]byte(`[{"success": {"username": "new_user_name","clientkey": "new_client_key"}}]`), "new_user_name"},
        {[]byte(`[{"other":"json"}]`), ""},
	}

	for _, test := range tests {
		if got, _ := ParseNewUserResult(test.input); got != test.want {
			t.Errorf("ParseNewUserResult(%q) == %v", test.input, test.want)
		}
	}
}
