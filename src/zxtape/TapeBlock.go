package zxtape

import (
	"bytes"
	"io"
	"wav"
	"zx"
)

type TapeBlock struct {
	Data     *[]byte
	Checksum byte
}

func writeDataByte(data byte, hi byte, lo byte, writer *bytes.Buffer, freq int) error {
	const (
		PULSELEN_ZERO = 855
		PULSELEN_ONE  = 1710
	)

	var mask byte = 0x80
	for mask != 0 {
		var len int
		if (data & mask) == 0 {
			len = PULSELEN_ZERO
		} else {
			len = PULSELEN_ONE
		}

		if err := wav.DoSignal(writer, hi, len, freq); err != nil {
			return err
		}
		if err := wav.DoSignal(writer, lo, len, freq); err != nil {
			return err
		}
		mask >>= 1
	}
	return nil
}

func (t *TapeBlock) SaveSoundData(amplify bool, soundBuffer *bytes.Buffer, freq int) error {
	const (
		PULSELEN_PILOT            = 2168
		PULSELEN_SYNC1            = 667
		PULSELEN_SYNC2            = 735
		PULSELEN_SYNC3            = 954
		IMPULSNUMBER_PILOT_HEADER = 8063
		IMPULSNUMBER_PILOT_DATA   = 3223
	)

	var err error

	var pilotImpulses int
	if (*t.Data)[0] < 128 {
		pilotImpulses = IMPULSNUMBER_PILOT_HEADER
	} else {
		pilotImpulses = IMPULSNUMBER_PILOT_DATA
	}

	var HI, LO byte
	if amplify {
		HI = 0xFF
		LO = 0x00
	} else {
		HI = 0xC0
		LO = 0x40
	}

	var signalState = HI

	for i := 0; i < pilotImpulses; i++ {
		if err = wav.DoSignal(soundBuffer, signalState, PULSELEN_PILOT, freq); err != nil {
			return err
		}

		if signalState == HI {
			signalState = LO
		} else {
			signalState = HI
		}
	}

	if signalState == LO {
		if err = wav.DoSignal(soundBuffer, LO, PULSELEN_PILOT, freq); err != nil {
			return err
		}
	}

	if err = wav.DoSignal(soundBuffer, HI, PULSELEN_SYNC1, freq); err != nil {
		return err
	}

	if err = wav.DoSignal(soundBuffer, LO, PULSELEN_SYNC2, freq); err != nil {
		return err
	}

	for _, d := range *t.Data {
		if err = writeDataByte(d, HI, LO, soundBuffer, freq); err != nil {
			return err
		}
	}

	if err = writeDataByte(t.Checksum, HI, LO, soundBuffer, freq); err != nil {
		return err
	}

	if err = wav.DoSignal(soundBuffer, HI, PULSELEN_SYNC3, freq); err != nil {
		return err
	}

	return nil
}

func ReadTapeBlock(reader io.Reader) (*TapeBlock, error) {
	var length int
	var err error
	var checksum byte

	length, err = zx.ReadZxShort(reader)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length-1)

	_, err = io.ReadAtLeast(reader, data, len(data))
	if err != nil {
		return nil, err
	}

	checksum, err = zx.ReadZxByte(reader)
	if err != nil {
		return nil, err
	}

	return &TapeBlock{&data, checksum}, nil
}
