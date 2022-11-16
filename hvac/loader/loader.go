package loader

import (
	"fmt"
	yaml "github.com/goccy/go-yaml"
	"github.com/gsasha/hvac_ip_mqtt_bridge/hvac/base"
	"github.com/gsasha/hvac_ip_mqtt_bridge/hvac/models"
	"io/ioutil"
	"log"
)

type Config struct {
	MQTT    *MQTTConfig    `yaml:"mqtt"`
	Devices []DeviceConfig `yaml:"devices"`
}

type MQTTConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type DeviceConfig struct {
	Name       string `yaml:"name"`
	Model      string `yaml:"model"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	MQTTPrefix string `yaml:"mqtt_prefix"`
	HAPrefix   string `yaml:"ha_prefix"`
	DUID       string `yaml:"duid"`
	AuthToken  string `yaml:"auth_token"`
}

type Device struct {
	mqtt       *base.MQTT
	controller base.Controller
}

func NewDevice(mqtt *base.MQTT, deviceConfig DeviceConfig) (*Device, error) {
	controller, err := models.NewController(
		deviceConfig.Model,
		deviceConfig.Name,
		deviceConfig.Host,
		deviceConfig.Port,
		deviceConfig.DUID,
		deviceConfig.AuthToken,
		deviceConfig.MQTTPrefix,
		(deviceConfig.HAPrefix != ""),
		)

	log.Printf("Registering controller %s %s", deviceConfig.Name, deviceConfig.MQTTPrefix)
	notifier := mqtt.RegisterController(deviceConfig.Name, deviceConfig.MQTTPrefix, deviceConfig.HAPrefix, controller)
	controller.SetStateNotifier(notifier)

	if err != nil {
		return nil, err
	}
	return &Device{
		mqtt:       mqtt,
		controller: controller,
	}, nil
}

func (device *Device) Run() {
	device.controller.Connect()
}

func Load(configFile string) ([]*Device, error) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, err
	}
	if config.MQTT == nil {
		return nil, fmt.Errorf("mqtt missing in configuration")
	}
	protocol := config.MQTT.Protocol
	if protocol == "" {
		protocol = "tcp"
	}
	host := config.MQTT.Host
	if host == "" {
		return nil, fmt.Errorf("MQTT host not given")
	}
	port := config.MQTT.Port
	if port == "" {
		port = "1883"
	}
	username := config.MQTT.Username
	password := config.MQTT.Password
	mqttBroker := fmt.Sprintf("%s://%s:%s", protocol, config.MQTT.Host, port)
	mqtt := base.NewMQTT(mqttBroker, username, password, "hvac_ip_mqtt_bridge")

	var devices []*Device
	for _, deviceConfig := range config.Devices {
		device, err := NewDevice(mqtt, deviceConfig)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	mqtt.Connect()
	return devices, nil
}
