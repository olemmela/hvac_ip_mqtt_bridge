package samsung

import (
	"bytes"
	"encoding/xml"
	"encoding/json"
	"fmt"
	"github.com/gsasha/hvac_ip_mqtt_bridge/hvac/base"
	"log"
	"strings"
	"strconv"
	"text/template"
	"time"
)

type homeAssistantDevice struct {
	Identifiers string `json:"identifiers"`
	Name string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Model string `json:"model"`
}

type homeAssistantSwitchConfig struct {
	Name string `json:"name"`
	Id string `json:"unique_id"`
	PayloadOn string `json:"payload_on"`
	PayloadOff string `json:"payload_off"`
	StateTopic string `json:"state_topic"`
	CommandTopic string `json:"command_topic"`
	Device homeAssistantDevice `json:"device"`
}

type homeAssistantSensorConfig struct {
	Name string `json:"name"`
	Id string `json:"unique_id"`
	DeviceClass string `json:"device_class"`
	StateTopic string `json:"state_topic"`
	StateClass string `json:"state_class,omitempty"`
	ValueTemplate string `json:"value_template,omitempty"`
	Unit string `json:"unit_of_measurement"`
	Icon string `json:"icon"`
	Device homeAssistantDevice `json:"device"`
}

type homeAssistantClimateConfig struct {
	Name string `json:"name"`
	Id string `json:"unique_id"`
	PowerCommandTopic string `json:"power_command_topic"`
	PayloadOn string `json:"payload_on"`
	PayloadOff string `json:"payload_off"`
	ModeStateTopic string `json:"mode_state_topic"`
	ModeCommandTopic string `json:"mode_command_topic"`
	ActionTopic string `json:"action_topic"`
	FanModeStateTopic string `json:"fan_mode_state_topic"`
	FanModeCommandTopic string `json:"fan_mode_command_topic"`
	FanModes []string `json:"fan_modes"`
	TemperatureStateTopic	string `json:"temperature_state_topic"`
	TemperatureCommandTopic string `json:"temperature_command_topic"`
	CurrentTemperatureTopic string `json:"current_temperature_topic"`
	PresetModeStateTopic string `json:"preset_mode_state_topic"`
	PresetModeCommandTopic string `json:"preset_mode_command_topic"`
	PresetModes []string `json:"preset_modes"`
	SwingModeStateTopic string `json:"swing_mode_state_topic"`
	SwingModeCommandTopic string `json:"swing_mode_command_topic"`
	SwingModes []string `json:"swing_modes"`
	Precision int `json:"precision"`
	MinTemp int `json:"min_temp"`
	MaxTemp int `json:"max_temp"`
	Device homeAssistantDevice `json:"device"`
}

type SamsungAC2878 struct {
	name      string
	host      string
	port      string
	authToken string
	duid      string

	connection    base.Connection
	stateNotifier base.StateNotifier

	online             bool
	homeassistant      bool
	prefix             string
	err                string
	powerMode          string
	opMode             string
	attrs              map[string]string
}

func NewSamsungAC2878(name string, host, port, duid, authToken string, prefix string, homeassistant bool) (*SamsungAC2878, error) {
	if port == "" {
		port = "2878"
	}
	conn, err := base.NewTLSSocketConnection()
	return &SamsungAC2878{
		name:       name,
		host:       host,
		port:       port,
		authToken:  authToken,
		duid:       duid,
		connection: conn,
		prefix:     prefix,
		homeassistant: homeassistant,
		attrs:      make(map[string]string),
	}, err
}

func (c *SamsungAC2878) SetStateNotifier(stateNotifier base.StateNotifier) {
	c.stateNotifier = stateNotifier
}

func (c *SamsungAC2878) SetHomeAssistantConfig() {
	device := homeAssistantDevice{
		Identifiers: c.duid,
		Name: c.name,
		Manufacturer: "Samsung",
		Model: "Smart A/C",
	}

	climate := homeAssistantClimateConfig{
		Name: c.name,
		Id: c.duid,
		PowerCommandTopic: c.prefix+"/power/set",
		PayloadOn: "on",
		PayloadOff: "off",
		ModeStateTopic: c.prefix+"/mode/state",
		ModeCommandTopic: c.prefix+"/mode/set",
		ActionTopic: c.prefix+"/action",
		FanModeStateTopic: c.prefix+"/fan_mode/state",
		FanModeCommandTopic: c.prefix+"/fan_mode/set",
		FanModes: []string{"auto", "low", "medium", "high", "turbo"},
		TemperatureStateTopic: c.prefix+"/temperature/state",
		TemperatureCommandTopic: c.prefix+"/temperature/set",
		CurrentTemperatureTopic: c.prefix+"/current_temperature/state",
		PresetModeStateTopic: c.prefix+"/preset_mode/state",
		PresetModeCommandTopic: c.prefix+"/preset_mode/set",
		PresetModes: []string{"eco", "sleep", "smart", "comfort", "boost"},
		SwingModeStateTopic: c.prefix+"/swing_mode/state",
		SwingModeCommandTopic: c.prefix+"/swing_mode/set",
		SwingModes: []string{"horizontal", "vertical", "both", "off"},
		Precision: 1,
		MinTemp: 8,
		MaxTemp: 30,
		Device: device,
	}

	config, err := json.Marshal(climate)
	if err == nil {
		c.stateNotifier.UpdateHomeAssistantConfig("climate/"+c.duid+"/config", string(config))
	}

	Sensor := homeAssistantSensorConfig{
		Name: c.name+" Outdoor Temperature",
		Id: c.duid+".outdoor_temp",
		StateTopic: c.prefix+"/outdoor_temperature",
		DeviceClass: "temperature",
		Icon: "mdi:temperature-celsius",
		Unit: "°C",
		Device: device,
	}

	config, err = json.Marshal(Sensor)
	if err == nil {
		c.stateNotifier.UpdateHomeAssistantConfig("sensor/"+c.duid+"/outdoor_temp/config", string(config))
	}

	Sensor = homeAssistantSensorConfig{
		Name: c.name+" Energy",
		Id: c.duid+".user_power",
		StateTopic: c.prefix+"/used_power",
		StateClass: "total_increasing",
		DeviceClass: "energy",
		ValueTemplate: "{{ value | float / 10 }}",
		Icon: "mdi:lightning-bolt",
		Unit: "kWh",
		Device: device,
	}

	config, err = json.Marshal(Sensor)
	if err == nil {
		c.stateNotifier.UpdateHomeAssistantConfig("sensor/"+c.duid+"/used_power/config", string(config))
	}

	Switch := homeAssistantSwitchConfig{
		Name: c.name+" Automatic Clean",
		Id: c.duid+".automatic_clean",
		StateTopic: c.prefix+"/automatic_clean/state",
		CommandTopic: c.prefix+"/automatic_clean/set",
		PayloadOn: "on",
		PayloadOff: "off",
		Device: device,
	}

	config, err = json.Marshal(Switch)
	if err == nil {
		c.stateNotifier.UpdateHomeAssistantConfig("switch/"+c.duid+"/automatic_clean/config", string(config))
	}

	Switch = homeAssistantSwitchConfig{
		Name: c.name+" Purify",
		Id: c.duid+".purify",
		StateTopic: c.prefix+"/purify/state",
		CommandTopic: c.prefix+"/purify/set",
		PayloadOn: "on",
		PayloadOff: "off",
		Device: device,
	}

	config, err = json.Marshal(Switch)
	if err == nil {
		c.stateNotifier.UpdateHomeAssistantConfig("switch/"+c.duid+"/purify/config", string(config))
	}
}

func (c *SamsungAC2878) Connect() {
	c.connection.Connect(c.host, c.port, c)
	if c.homeassistant {
		c.SetHomeAssistantConfig()
	}
	go func() {
		for range time.Tick(time.Second * 60) {
			c.sendDeviceStateRequest()
		}
	}()
}

var (
	authenticateTemplate = template.Must(template.New("authenticate").Parse(
		`<Request Type="AuthToken"><User Token="{{.token}}" /></Request>
`))
	deviceStateTemplate = template.Must(template.New("deviceState").Parse(
		`<Request Type="DeviceState" DUID="{{.duid}}"></Request>
`))
	setPowerModeTemplate = template.Must(template.New("setPowerMode").Parse(
		`<Request Type="DeviceControl"><Control CommandID="AC_FUN_POWER" DUID="{{.duid}}"><Attr ID="AC_FUN_POWER" Value="{{.value}}" /></Control></Request>
`))
	setModeTemplate = template.Must(template.New("setMode").Parse(
		`<Request Type="DeviceControl"><Control CommandID="AC_FUN_OPMODE" DUID="{{.duid}}"><Attr ID="AC_FUN_OPMODE" Value="{{.value}}" /></Control></Request>
`))
	setFanModeTemplate = template.Must(template.New("setFanMode").Parse(
		`<Request Type="DeviceControl"><Control CommandID="AC_FUN_WINDLEVEL" DUID="{{.duid}}"><Attr ID="AC_FUN_WINDLEVEL" Value="{{.value}}" /></Control></Request>
`))
	setPresetModeTemplate = template.Must(template.New("setPresetMode").Parse(
		`<Request Type="DeviceControl"><Control CommandID="AC_FUN_COMODE" DUID="{{.duid}}"><Attr ID="AC_FUN_COMODE" Value="{{.value}}" /></Control></Request>
`))
	setSwingModeTemplate = template.Must(template.New("setSwingMode").Parse(
		`<Request Type="DeviceControl"><Control CommandID="AC_FUN_DIRECTION" DUID="{{.duid}}"><Attr ID="AC_FUN_DIRECTION" Value="{{.value}}" /></Control></Request>
`))
	setTemperatureTemplate = template.Must(template.New("setTemperature").Parse(
		`<Request Type="DeviceControl"><Control CommandID="AC_FUN_TEMPSET" DUID="{{.duid}}"><Attr ID="AC_FUN_TEMPSET" Value="{{.value}}" /></Control></Request>
`))
	setPurifyTemplate = template.Must(template.New("setPurify").Parse(
		`<Request Type="DeviceControl"><Control CommandID="AC_ADD_SPI" DUID="{{.duid}}"><Attr ID="AC_ADD_SPI" Value="{{.value}}" /></Control></Request>
`))
	setAutomaticCleanTemplate = template.Must(template.New("setAutomaticClean").Parse(
		`<Request Type="DeviceControl"><Control CommandID="AC_ADD_AUTOCLEAN" DUID="{{.duid}}"><Attr ID="AC_ADD_AUTOCLEAN" Value="{{.value}}" /></Control></Request>
`))
)

func (c *SamsungAC2878) SetPowerMode(powerMode string) {
	c.sendMessage(setPowerModeTemplate, map[string]string{
		"value": booleanToAC(powerMode),
		"duid":  c.duid,
	})
}

func (c *SamsungAC2878) SetOpMode(mode string) {
	if mode == "off" {
		c.sendMessage(setPowerModeTemplate, map[string]string{
			"value": "Off",
			"duid":  c.duid,
		})
	} else {
		c.sendMessage(setModeTemplate, map[string]string{
			"value": OpModeToAC(mode),
			"duid":  c.duid,
		})
	}
}

func (c *SamsungAC2878) SetFanMode(fanMode string) {
	c.sendMessage(setFanModeTemplate, map[string]string{
		"value": FanModeToAC(fanMode),
		"duid":  c.duid,
	})
}

func (c *SamsungAC2878) SetPresetMode(presetMode string) {
	c.sendMessage(setPresetModeTemplate, map[string]string{
		"value": PresetModeToAC(presetMode),
		"duid":  c.duid,
	})
}

func (c *SamsungAC2878) SetSwingMode(swingMode string) {
	c.sendMessage(setSwingModeTemplate, map[string]string{
		"value": SwingModeToAC(swingMode),
		"duid":  c.duid,
	})
}

func (c *SamsungAC2878) SetPurify(status string) {
	c.sendMessage(setPurifyTemplate, map[string]string{
		"value": booleanToAC(status),
		"duid":  c.duid,
	})
}

func (c *SamsungAC2878) SetAutomaticClean(status string) {
	c.sendMessage(setAutomaticCleanTemplate, map[string]string{
		"value": booleanToAC(status),
		"duid":  c.duid,
	})
}

func (c *SamsungAC2878) SetTemperature(temperature string) {
	c.sendMessage(setTemperatureTemplate, map[string]string{
		"value": temperature,
		"duid":  c.duid,
	})
}

type Response struct {
	XMLName     xml.Name `xml:"Response"`
	Type        string   `xml:"Type,attr"`
	Status      string   `xml:"Status,attr"`
	DeviceState DeviceState
	Inner       []byte `xml:",innerxml"`
}

type Update struct {
	XMLName xml.Name `xml:"Update"`
	Type    string   `xml:"Type,attr"`
	Status  Status
}
type Attr struct {
	XMLName xml.Name `xml:"Attr"`
	ID      string   `xml:"ID,attr"`
	Type    string   `xml:"Type,attr"`
	Value   string   `xml:"Value,attr"`
}
type Status struct {
	XMLName xml.Name `xml:"Status"`
	DUID    string   `xml:"DUID"`
	GroupID string   `xml:GroupID,attr"`
	ModelID string   `xml:ModelID,attr"`
	Attr    []Attr
}
type Device struct {
	XMLName xml.Name `xml:"Device"`
	DUID    string   `xml:"DUID,attr"`
	GroupID string   `xml:"GroupID,attr"`
	ModelID string   `xml:"ModelID,attr"`
	Attr    []Attr
}
type DeviceState struct {
	XMLName xml.Name `xml:"DeviceState"`
	Device  Device
}

func (c *SamsungAC2878) OnConnectionEstablished() {
	log.Printf("Established connection to %s", c.name)
	c.connection.ExpectRead()
}

func (c *SamsungAC2878) HandleMessage(message []byte) {
	log.Printf("Received message from %s: %s", c.name, string(message))

	if string(message) == "DPLUG-1.6\n" {
		log.Printf("Connection hello received from %s", c.name)
		c.connection.ExpectRead()
	}
	var update Update
	if err := xml.Unmarshal(message, &update); err == nil {
		c.handleUpdate(&update)
		return
	}
	var response Response
	if err := xml.Unmarshal(message, &response); err == nil {
		c.handleResponse(&response)
		return
	}
}

func (c *SamsungAC2878) handleUpdate(update *Update) error {
	switch update.Type {
	case "InvalidateAccount":
		c.handleInvalidateAccount()
	case "Status":
		c.handleUpdateStatus(&update.Status)
	default:
		log.Println("Error: %s unknown update type", c.name, update.Type)
		return nil
	}
	return nil
}

func (c *SamsungAC2878) handleResponse(response *Response) error {
	switch response.Type {
	case "AuthToken":
		c.handleAuthToken(response.Status)
	case "DeviceState":
		c.handleDeviceState(&response.DeviceState)
	case "DeviceControl":
		c.handleDeviceControl(response.Status)
	default:
		fmt.Println("Error: %s got unknown response", c.name, response.Type)
	}
	return nil
}

func (c *SamsungAC2878) handleInvalidateAccount() {
	c.sendMessage(authenticateTemplate, map[string]string{
		"duid":  c.duid,
		"token": c.authToken,
	})
}

func (c *SamsungAC2878) sendDeviceStateRequest() {
	c.sendMessage(deviceStateTemplate, map[string]string{
		"duid": c.duid,
	})
}

func (c *SamsungAC2878) handleAuthToken(status string) {
	if status == "Okay" {
		c.online = true
	} else {
		c.online = false
	}
	c.sendDeviceStateRequest()
}

func (c *SamsungAC2878) handleDeviceControl(status string) {
	if status == "Okay" {
		c.err = ""
	} else {
		c.err = status
	}
}

func (c *SamsungAC2878) handleUpdateStatus(status *Status) {
	if status == nil {
		fmt.Println("Error: No status")
		return
	}
	c.handleAttributes(status.Attr)
}

func (c *SamsungAC2878) handleDeviceState(deviceState *DeviceState) {
	c.handleAttributes(deviceState.Device.Attr)
}

func (c *SamsungAC2878) handleOpModeUpdate() {
	if strings.ToLower(c.powerMode) == "off" {
		c.stateNotifier.UpdateOpMode(OpModeFromAC("Off"))
	} else {
		c.stateNotifier.UpdateOpMode(OpModeFromAC(c.opMode))
	}
}

func (c *SamsungAC2878) handleAttributes(attrs []Attr) {
	if c.stateNotifier == nil {
		fmt.Println("Error: want to notify state, but no notifer defined")
		return
	}
	for _, attr := range attrs {
		c.attrs[attr.Type] = attr.Value
		switch attr.ID {
		case "AC_FUN_POWER":
			c.powerMode = attr.Value
			c.handleOpModeUpdate()
		case "AC_FUN_OPMODE":
			c.opMode = attr.Value
			c.handleOpModeUpdate()
		case "AC_FUN_COMODE":
			c.stateNotifier.UpdatePresetMode(PresetModeFromAC(attr.Value))
		case "AC_FUN_DIRECTION":
			c.stateNotifier.UpdateSwingMode(SwingModeFromAC(attr.Value))
		case "AC_FUN_TEMPSET":
			c.stateNotifier.UpdateTemperature(attr.Value)
		case "AC_FUN_TEMPNOW":
			c.stateNotifier.UpdateCurrentTemperature(attr.Value)
		case "AC_ADD_SPI":
			c.stateNotifier.UpdatePurify(booleanFromAC(attr.Value))
		case "AC_ADD_AUTOCLEAN":
			c.stateNotifier.UpdateAutomaticClean(booleanFromAC(attr.Value))
		case "AC_OUTDOOR_TEMP":
			value, err := strconv.Atoi(attr.Value)
			if err == nil {
				c.stateNotifier.UpdateAttribute("outdoor_temperature", strconv.Itoa(value - 55))
			}
		case "AC_FUN_WINDLEVEL":
			c.stateNotifier.UpdateFanMode(FanModeFromAC(attr.Value))
		case "AC_ADD2_USEDPOWER":
			c.stateNotifier.UpdateAttribute("used_power", attr.Value)
		}
	}
}

func (c *SamsungAC2878) sendMessage(messageTemplate *template.Template, data map[string]string) {
	var buf bytes.Buffer
	messageTemplate.Execute(&buf, data)
	log.Printf("sending request to %s [%s]\n", c.name, buf.String())
	c.connection.SendMessage(buf.Bytes())
}
