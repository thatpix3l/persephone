// Utilities for building byte sequence commands, as well as (un)marshaling responses
package command

import (
	"encoding/binary"
	"time"
)

const (
	Action actionT = iota // Root of all functions for generating command byte sequences
)

type actionT int

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	} else {
		return 0
	}
}

func buildAction(id byte, parameters ...byte) []byte {

	finalPacket := []byte{1, id}
	if len(parameters) == 0 {
		return finalPacket
	}

	finalPacket = append(finalPacket, byte(len(parameters)))
	finalPacket = append(finalPacket, parameters...)
	finalPacket[0] = finalPacket[2] + 2

	return finalPacket

}

func (a actionT) TurnShutterOn() []byte {
	return buildAction(0x01, 0x01)
}

func (a actionT) TurnShutterOff() []byte {
	return buildAction(0x01, 0x00)
}

func (a actionT) Sleep() []byte {
	return buildAction(0x05)
}

func (a actionT) SetDateTime(t time.Time) []byte {

	dateBuf := []byte{}

	year := make([]byte, 2)
	binary.BigEndian.PutUint16(year, uint16(t.Year()))

	dateBuf = append(dateBuf, year...)
	dateBuf = append(dateBuf, byte(t.Month()))
	dateBuf = append(dateBuf, byte(t.Day()))
	dateBuf = append(dateBuf, byte(t.Hour()))
	dateBuf = append(dateBuf, byte(t.Minute()))
	dateBuf = append(dateBuf, byte(t.Second()))

	return buildAction(0x0d, dateBuf...)

}

func (a actionT) GetDateTime() []byte {
	return buildAction(0x0e)
}

func (a actionT) SetLocalDateTime(t time.Time) []byte {

	dateBuf := []byte{}

	// Convert the year into a 16-bit byte array
	year := make([]byte, 2)
	binary.BigEndian.PutUint16(year, uint16(t.Year()))

	// Append to the date buffer stuff related to the date
	dateBuf = append(dateBuf, year...)
	dateBuf = append(dateBuf, byte(t.Month()))
	dateBuf = append(dateBuf, byte(t.Day()))
	dateBuf = append(dateBuf, byte(t.Hour()))
	dateBuf = append(dateBuf, byte(t.Minute()))
	dateBuf = append(dateBuf, byte(t.Second()))

	// Extract the zone offset duration
	_, offset := t.Zone()
	offsetDuration := time.Second * time.Duration(offset)

	// Append to the date buffer the offset in hours and minutes
	dateBuf = append(dateBuf, byte(offsetDuration.Hours()))
	dateBuf = append(dateBuf, byte(offsetDuration.Minutes()))

	// Check if the given time was configured for DST
	var isDSTInt byte = 0
	if t.IsDST() {
		isDSTInt = 1
	}

	// Append DST switch byte
	dateBuf = append(dateBuf, isDSTInt)

	return buildAction(0x0f, dateBuf...)

}

func (a actionT) GetLocalDateTime() []byte {
	return buildAction(0x10)
}

func (a actionT) TurnAccessPointOff() []byte {
	return buildAction(0x17, 0x00)
}

func (a actionT) TurnAccessPointOn() []byte {
	return buildAction(0x17, 0x01)
}

func (a actionT) HilightMoment() []byte {
	return buildAction(0x18)
}

func (a actionT) GetHardwareInfo() []byte {
	return buildAction(0x3c)
}

func (a actionT) LoadPresetGroupVideo() []byte {
	return buildAction(0x3e, 0x03, 0xe8)
}

func (a actionT) LoadPresetGroupPhoto() []byte {
	return buildAction(0x3e, 0x03, 0xe9)
}

func (a actionT) LoadPresetGroupTimelapse() []byte {
	return buildAction(0x3e, 0x03, 0xea)
}

func (a actionT) Analytics() []byte {
	return buildAction(0x50)
}

func (a actionT) GetVersion() []byte {
	return buildAction(0x51)
}
