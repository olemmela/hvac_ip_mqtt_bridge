package models

import (
	"fmt"
	"github.com/gsasha/hvac_ip_mqtt_bridge/hvac/base"
	"github.com/gsasha/hvac_ip_mqtt_bridge/hvac/models/samsung"
)

func NewController(
	model string, name string,
	host, port, duid, authToken string,
	prefix string, homeassistant bool) (base.Controller, error) {
	switch model {
	case "samsungac2878":
		return samsung.NewSamsungAC2878(name, host, port, duid, authToken, prefix, homeassistant)
	}
	return nil, fmt.Errorf("Model not supported: %s", model)
}
