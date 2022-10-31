package base

import (
	"fmt"
	"log"
	"crypto/rand"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"encoding/base64"
)

const (
	availabilityTopic            = "mode/availability"
	powerCommandTopic            = "power/set"
	opModeCommandTopic           = "mode/set"
	opModeStateTopic             = "mode/state"
	actionTopic                  = "action"
	currentTemperatureStateTopic = "current_temperature/state"
	temperatureCommandTopic      = "temperature/set"
	temperatureStateTopic        = "temperature/state"
	fanModeCommandTopic          = "fan_mode/set"
	fanModeStateTopic            = "fan_mode/state"
	presetModeCommandTopic       = "preset_mode/set"
	presetModeStateTopic         = "preset_mode/state"
	swingModeCommandTopic        = "swing_mode/set"
	swingModeStateTopic          = "swing_mode/state"
	purifyCommandTopic           = "purify/set"
	purifyStateTopic             = "purify/state"
	automaticCleanCommandTopic   = "automatic_clean/set"
	automaticCleanStateTopic     = "automatic_clean/state"
)

type MQTT struct {
	clientId string
	prefix   string
	haPrefix string

	client      mqtt.Client
	controllers map[string]Controller
	prefixes    map[string]string
}

type MQTTNotifier struct {
	mqtt   *MQTT
	prefix string
}

func (m *MQTTNotifier) UpdateAction(action string) {
	m.mqtt.updateAction(m.prefix, action)
}
func (m *MQTTNotifier) UpdateOpMode(opMode string) {
	m.mqtt.updateOpMode(m.prefix, opMode)
}
func (m *MQTTNotifier) UpdateFanMode(fanMode string) {
	m.mqtt.updateFanMode(m.prefix, fanMode)
}
func (m *MQTTNotifier) UpdatePresetMode(presetMode string) {
	m.mqtt.updatePresetMode(m.prefix, presetMode)
}
func (m *MQTTNotifier) UpdateSwingMode(swingMode string) {
	m.mqtt.updateSwingMode(m.prefix, swingMode)
}
func (m *MQTTNotifier) UpdatePurify(status string) {
	m.mqtt.updatePurify(m.prefix, status)
}
func (m *MQTTNotifier) UpdateAutomaticClean(status string) {
	m.mqtt.updateAutomaticClean(m.prefix, status)
}
func (m *MQTTNotifier) UpdateTemperature(temperature string) {
	m.mqtt.updateTemperature(m.prefix, temperature)
}
func (m *MQTTNotifier) UpdateCurrentTemperature(temperature string) {
	m.mqtt.updateCurrentTemperature(m.prefix, temperature)
}
func (m *MQTTNotifier) UpdateAttribute(topic string, attribute string) {
	m.mqtt.updateAttribute(m.prefix, topic, attribute)
}
func (m *MQTTNotifier) UpdateHomeAssistantConfig(topic string, config string) {
	m.mqtt.updateHomeAssistantConfig(m.mqtt.haPrefix, topic, config)
}

func NewMQTT(broker string, username string, password string, clientId string) *MQTT {
	log.Printf("Connecting to MQTT broker %s for %s", broker, clientId)
	m := &MQTT{
		clientId:    clientId,
		controllers: make(map[string]Controller),
		prefixes:    make(map[string]string),
	}

	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	if username != "" {
		options.SetUsername(username)
	}
	if password != "" {
		options.SetPassword(password)
	}
	random_id := make([]byte, 8)
	log.Printf("Reading random")
	_, err := rand.Read(random_id)
	if err != nil {
		log.Fatal("Cannot get random for client ID")
	}
	clientId = fmt.Sprintf("%s_%s", clientId, base64.StdEncoding.EncodeToString(random_id))
	log.Printf("MQTT Client ID %s", clientId)

	options.SetClientID(clientId)
	options.SetOnConnectHandler(func(client mqtt.Client) {
		log.Printf("Connection established to %s:%s", clientId, broker)
		m.subscribeTopics()
	})
	options.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection lost to %s:%s %s", clientId, broker, err)
	})
	options.SetAutoReconnect(true)

	m.client = mqtt.NewClient(options)
	return m
}

func (m *MQTT) RegisterController(id string, prefix string, haPrefix string, controller Controller) StateNotifier {
	m.haPrefix = haPrefix
	m.controllers[id] = controller
	m.prefixes[id] = prefix
	return &MQTTNotifier{
		mqtt:   m,
		prefix: prefix,
	}
}

func (m *MQTT) Connect() {
	token := m.client.Connect()
	if token.Wait() && token.Error() == nil {
		log.Println("MQTT Connection succeeded:", m.client.IsConnectionOpen())
	} else {
		log.Println("MQTT Connection failed:", token.Error())
	}
}

func (m *MQTT) subscribeTopics() {
	for controllerId := range m.controllers {
		prefix := m.prefixes[controllerId]
		key := fmt.Sprintf("%s", controllerId)
		log.Printf("subscribing to prefix %s for %s", prefix, controllerId)
		tokens := []mqtt.Token{
			m.client.Subscribe(prefix+"/"+powerCommandTopic, 0,
				func(client mqtt.Client, message mqtt.Message) {
					log.Printf("Received %s:%s:%s", key, message.Topic(), string(message.Payload()))
					m.controllers[key].SetPowerMode(string(message.Payload()))
				}),
			m.client.Subscribe(prefix+"/"+opModeCommandTopic, 0,
				func(client mqtt.Client, message mqtt.Message) {
					log.Println("Received %s:%s:%s", key, message.Topic(), string(message.Payload()))
					m.controllers[key].SetOpMode(string(message.Payload()))
				}),
			m.client.Subscribe(prefix+"/"+fanModeCommandTopic, 0,
				func(client mqtt.Client, message mqtt.Message) {
					log.Println("Received %s:%s:%s", key, message.Topic(), string(message.Payload()))
					m.controllers[key].SetFanMode(string(message.Payload()))
				}),
			m.client.Subscribe(prefix+"/"+presetModeCommandTopic, 0,
				func(client mqtt.Client, message mqtt.Message) {
					log.Println("Received %s:%s:%s", key, message.Topic(), string(message.Payload()))
					m.controllers[key].SetPresetMode(string(message.Payload()))
				}),
			m.client.Subscribe(prefix+"/"+swingModeCommandTopic, 0,
				func(client mqtt.Client, message mqtt.Message) {
					log.Println("Received %s:%s:%s", key, message.Topic(), string(message.Payload()))
					m.controllers[key].SetSwingMode(string(message.Payload()))
				}),
			m.client.Subscribe(prefix+"/"+purifyCommandTopic, 0,
				func(client mqtt.Client, message mqtt.Message) {
					log.Println("Received %s:%s:%s", key, message.Topic(), string(message.Payload()))
					m.controllers[key].SetPurify(string(message.Payload()))
				}),
			m.client.Subscribe(prefix+"/"+automaticCleanCommandTopic, 0,
				func(client mqtt.Client, message mqtt.Message) {
					log.Println("Received %s:%s:%s", key, message.Topic(), string(message.Payload()))
					m.controllers[key].SetAutomaticClean(string(message.Payload()))
				}),
			m.client.Subscribe(prefix+"/"+temperatureCommandTopic, 0,
				func(client mqtt.Client, message mqtt.Message) {
					log.Println("Received %s:%s:%s", key, message.Topic(), string(message.Payload()))
					m.controllers[key].SetTemperature(string(message.Payload()))
				}),
			m.client.Subscribe(m.haPrefix+"/status", 0,
				func(client mqtt.Client, message mqtt.Message) {
					log.Println("Received %s:%s:%s", key, message.Topic(), string(message.Payload()))
					if string(message.Payload()) == "online" {
						m.controllers[key].SetHomeAssistantConfig()
					}
				}),
			// TODO(gsasha): subscribe to more commands.
		}
		for _, token := range tokens {
			if token.Wait() && token.Error() != nil {
				log.Printf("Error subscribing to topics %s: %s", controllerId, token.Error())
				return
			}
		}
		log.Printf("Subscribed to topics for %s", controllerId)
	}
}

func (m *MQTT) updateAction(prefix string, action string) {
	m.publish(prefix, actionTopic, action)
}
func (m *MQTT) updateOpMode(prefix string, opMode string) {
	m.publish(prefix, opModeStateTopic, opMode)
}
func (m *MQTT) updateFanMode(prefix string, fanMode string) {
	m.publish(prefix, fanModeStateTopic, fanMode)
}
func (m *MQTT) updatePresetMode(prefix string, presetMode string) {
	m.publish(prefix, presetModeStateTopic, presetMode)
}
func (m *MQTT) updateSwingMode(prefix string, swingMode string) {
	m.publish(prefix, swingModeStateTopic, swingMode)
}
func (m *MQTT) updatePurify(prefix string, status string) {
	m.publish(prefix, purifyStateTopic, status)
}
func (m *MQTT) updateAutomaticClean(prefix string, status string) {
	m.publish(prefix, automaticCleanStateTopic, status)
}
func (m *MQTT) updateTemperature(prefix string, temperature string) {
	m.publish(prefix, temperatureStateTopic, temperature)
}
func (m *MQTT) updateCurrentTemperature(prefix string, temperature string) {
	m.publish(prefix, currentTemperatureStateTopic, temperature)
}
func (m *MQTT) updateAttribute(prefix string, topic string, attribute string) {
	m.publish(prefix, topic, attribute)
}
func (m *MQTT) updateHomeAssistantConfig(prefix string, topic string, config string) {
	m.publish(prefix, topic, config)
}
func (m *MQTT) publish(prefix string, topic string, message string) {
	log.Println("mqtt publishing", prefix+"/"+topic, message)
	m.client.Publish(prefix+"/"+topic, 0, false, message)
}
