package command

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/thatpix3l/persephone/pkg/zeropad"
	"golang.org/x/exp/constraints"
)

type Response struct {
	Shutter           bool // `commandID:"1"`
	Sleep             bool // `commandID:"5"`
	SetDateTime       bool // `commandID:"13"`
	SetLocalDateTime  bool // `commandID:"15"`
	SetLivestreamMode bool // `commandID:"21"`
	WifiAP            bool // `commandID:"23"`
	HiLightMoment     bool // `commandID:"24"`
	LoadPresetGroup   bool // `commandID:"62"`
	LoadPreset        bool // `commandID:"64"`
	Analytics         bool // `commandID:"80"`
}

func intToBool[Integer constraints.Integer](i Integer) (bool, error) {

	if i == 0 {
		return false, nil

	} else if i == 1 {
		return true, nil

	}

	return false, errors.New("integer is not 0 or 1, cannot convert to bool")

}

func (r *Response) Unmarshal(data []byte) error {

	if len(data) < 3 {
		return fmt.Errorf("data length is less than minimum of 3: %v", data)
	}

	suggestedPacketLength := data[0]
	if len(data)-1 != int(suggestedPacketLength) {
		return fmt.Errorf("packet length (index 0) does not match suggested packet length: %v", data)
	}

	id := data[1]
	valueBuf := data[2:]
	valueZeroPadded := zeropad.BigEndian64(valueBuf)
	valueInt := binary.BigEndian.Uint64(valueZeroPadded)

	var funcError error = nil

	switch id {

	case 0x01:
		r.Shutter, funcError = intToBool(valueInt)

	case 0x05:
		r.Sleep, funcError = intToBool(valueInt)

	case 0x0d:
		r.SetDateTime, funcError = intToBool(valueInt)

	case 0x0f:
		r.SetLocalDateTime, funcError = intToBool(valueInt)

	case 0x15:
		r.SetLivestreamMode, funcError = intToBool(valueInt)

	case 0x17:
		r.WifiAP, funcError = intToBool(valueInt)

	case 0x18:
		r.HiLightMoment, funcError = intToBool(valueInt)

	case 0x3e:
		r.LoadPresetGroup, funcError = intToBool(valueInt)

	case 0x40:
		r.LoadPreset, funcError = intToBool(valueInt)

	case 0x50:
		r.Analytics, funcError = intToBool(valueInt)

	default:
		funcError = fmt.Errorf("status id does not exist: %v (%x)", valueInt, valueInt)

	}

	return funcError

}
