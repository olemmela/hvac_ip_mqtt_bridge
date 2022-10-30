package samsung

import (
	"strings"
)

type translationEntry struct {
	mqtt string
	ac   string
}

func toAc(value string, table []translationEntry) string {
	for _, e := range table {
		if strings.ToLower(value) == strings.ToLower(e.mqtt) {
			return e.ac
		}
	}
	return value
}

func fromAc(value string, table []translationEntry) string {
	for _, e := range table {
		if strings.ToLower(value) == strings.ToLower(e.ac) {
			return e.mqtt
		}
	}
	return strings.ToLower(value)
}

var booleanTable = []translationEntry{
	{"on", "On"},
	{"off", "Off"},
}

func booleanToAC(mode string) string   { return toAc(mode, booleanTable) }
func booleanFromAC(mode string) string { return fromAc(mode, booleanTable) }

var opModeTable = []translationEntry{
	{"cool", "Cool"},
	{"heat", "Heat"},
	{"dry", "Dry"},
	{"auto", "Auto"},
	{"fan_only", "Wind"},
	{"off", "Off"},
}

func OpModeToAC(mode string) string   { return toAc(mode, opModeTable) }
func OpModeFromAC(mode string) string { return fromAc(mode, opModeTable) }

var fanModeTable = []translationEntry{
	{"auto", "Auto"},
	{"low", "Low"},
	{"medium", "Mid"},
	{"high", "High"},
	{"turbo", "Turbo"},
}

func FanModeToAC(mode string) string   { return toAc(mode, fanModeTable) }
func FanModeFromAC(mode string) string { return fromAc(mode, fanModeTable) }

var presetModeTable = []translationEntry{
	{"none", "Off"},
	{"eco", "Quiet"},
	{"sleep", "Sleep"},
	{"smart", "Smart"},
	{"comfort", "SoftCool"},
	{"boost", "TurboMode"},
}

func PresetModeToAC(mode string) string   { return toAc(mode, presetModeTable) }
func PresetModeFromAC(mode string) string { return fromAc(mode, presetModeTable) }

var swingModeTable = []translationEntry{
	{"horizontal", "SwingLR"},
	{"vertical", "SwingUD"},
	{"off", "Fixed"},
	{"both", "Rotation"},
}

func SwingModeToAC(mode string) string   { return toAc(mode, swingModeTable) }
func SwingModeFromAC(mode string) string { return fromAc(mode, swingModeTable) }
