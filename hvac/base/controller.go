package base

type StateNotifier interface {
	UpdateAction(action string)
	UpdateOpMode(mode string)
	UpdateFanMode(fanMode string)
	UpdatePresetMode(presetMode string)
	UpdateSwingMode(swingMode string)
	UpdatePurify(status string)
	UpdateAutomaticClean(status string)
	UpdateTemperature(temperature string)
	UpdateCurrentTemperature(temperature string)
	UpdateAttribute(topic string, attribute string)
}

type Controller interface {
	SetStateNotifier(stateNotifier StateNotifier)
	Connect()
	SetPowerMode(powerMode string)
	SetOpMode(mode string)
	SetFanMode(fanMode string)
	SetPresetMode(presetMode string)
	SetSwingMode(SwingMode string)
	SetPurify(status string)
	SetAutomaticClean(status string)
	SetTemperature(temperature string)
}
