package command

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

// Semantic versioning
type semVer struct {
	Major int
	Minor int
	Patch int
}

// Return a string in the form of vX.Y.Z, where X, Y and Z corresponds to the major, minor and patch version, respectively
func (s *semVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
}

type hardware struct {
	ModelNumber     string
	ModelName       string
	Board           string
	FirmwareVersion string
	SerialNumber    string
	SSID            string
	SSIDMacAddress  string
}

type response struct {
	Shutter           bool      // `commandID:"1"`
	Sleep             bool      // `commandID:"5"`
	SetDateTime       bool      // `commandID:"13"`
	DateTime          time.Time //
	SetLocalDateTime  bool      // `commandID:"15"`
	LocalDateTime     time.Time //
	SetLivestreamMode bool      // `commandID:"21"`
	WifiAP            bool      // `commandID:"23"`
	HiLightMoment     bool      // `commandID:"24"`
	Hardware          hardware  // `commandID:"60"`
	LoadPresetGroup   bool      // `commandID:"62"`
	LoadPreset        bool      // `commandID:"64"`
	Analytics         bool      // `commandID:"80"`
	OpenGoProVersion  semVer    // `commandID:"81"`
}

func NewResponse() response {
	return response{}
}

// Return true or false if "i" is 1 or 0, error if any other number
func intToBool[Integer constraints.Integer](i Integer) (bool, error) {

	if i == 0 {
		return false, nil

	} else if i == 1 {
		return true, nil

	}

	return false, errors.New("integer is not 0 or 1, cannot convert to bool")

}

// Return the hexadecimal string representation of a byte slice, separated by colons
func bytesToHexString(buf []byte) string {
	hexString := ""
	for _, b := range buf {
		hexString += fmt.Sprintf("%x:", b)
	}
	hexString = strings.TrimSuffix(hexString, ":")
	return hexString
}

// Unmarshal a byte slice into *response.
//
// Errors if not minimum length, length does not match expected length, or has unknown status ID.
func (r *response) Unmarshal(data []byte) error {

	if len(data) < 3 {
		return fmt.Errorf("data length is less than minimum of 3: %v", data)
	}

	suggestedPacketLength := data[0]
	if len(data)-1 != int(suggestedPacketLength) {
		return fmt.Errorf("packet length (index 0) does not match suggested packet length: %v", data)
	}

	id := data[1]
	successCode := data[2]
	valueBuf := data[3:]

	shiftValueBuf := func() {
		valueBuf = valueBuf[valueBuf[0]+1:]
	}

	var funcError error = nil

	switch id {

	case 0x01:
		r.Shutter, funcError = intToBool(successCode)

	case 0x05:
		r.Sleep, funcError = intToBool(successCode)

	case 0x0d:
		r.SetDateTime, funcError = intToBool(successCode)

	case 0x0e:
		r.DateTime = time.Date(
			int(binary.BigEndian.Uint16(valueBuf[1:2])), // Two bytes for year
			time.Month(valueBuf[3]),                     // Byte for current month
			int(valueBuf[4]),                            // day
			int(valueBuf[5]),                            // hour
			int(valueBuf[6]),                            // minute
			int(valueBuf[7]),                            // second
			0,                                           // "DateTime" doesn't provide nanoseconds, so ignore
			time.Local,                                  // "DateTime"
		)

	case 0x0f:
		r.SetLocalDateTime, funcError = intToBool(successCode)

	// case 0x10:
	// 	r.LocalDateTime = time.Date(
	// 		int(binary.BigEndian.Uint16(valueBuf[1:2])), // Two bytes for year
	// 		time.Month(valueBuf[3]),                     // Byte for current month
	// 		int(valueBuf[4]),                            // day
	// 		int(valueBuf[5]),                            // hour
	// 		int(valueBuf[6]),                            // minute
	// 		int(valueBuf[7]),                            // second
	// 		0,                                           // "DateTime" doesn't provide nanoseconds, so ignore
	// 		time.Local,                                  // "DateTime" doesn't provide a timezone, so add local
	// 	)

	case 0x15:
		r.SetLivestreamMode, funcError = intToBool(successCode)

	case 0x17:
		r.WifiAP, funcError = intToBool(successCode)

	case 0x18:
		r.HiLightMoment, funcError = intToBool(successCode)

	case 0x3c:
		r.Hardware.ModelNumber = fmt.Sprintf("%x:%x:%x:%x", valueBuf[1], valueBuf[2], valueBuf[3], valueBuf[4])
		shiftValueBuf()
		r.Hardware.ModelName = string(valueBuf[1 : valueBuf[0]+1])
		shiftValueBuf()
		r.Hardware.Board = string(valueBuf[1 : valueBuf[0]+1])
		shiftValueBuf()
		r.Hardware.FirmwareVersion = string(valueBuf[1 : valueBuf[0]+1])
		shiftValueBuf()
		r.Hardware.SerialNumber = string(valueBuf[1 : valueBuf[0]+1])
		shiftValueBuf()
		r.Hardware.SSID = string(valueBuf[1 : valueBuf[0]+1])
		shiftValueBuf()
		r.Hardware.SSIDMacAddress = bytesToHexString(valueBuf[1 : valueBuf[0]+1])
		shiftValueBuf()

	case 0x3e:
		r.LoadPresetGroup, funcError = intToBool(successCode)

	case 0x40:
		r.LoadPreset, funcError = intToBool(successCode)

	case 0x50:
		r.Analytics, funcError = intToBool(successCode)

	case 0x51:
		r.OpenGoProVersion.Major = int(valueBuf[1])
		r.OpenGoProVersion.Minor = int(valueBuf[3])

	default:
		funcError = fmt.Errorf("status id does not exist: %v (%x)", id, id)

	}

	return funcError

}
