package g7

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// The card types ("CardType") of the events that the emulator makes.
const (
	CardTypeMifareClassic1K = "MIFARE_CLASSIC_1K"
	CardTypeBarcode         = "BARCODE"
	CardTypePIN             = "PIN"
)

const maxIdentLength = 128

var (
	ErrIdentEmpty   = errors.New("the id is empty")
	ErrIdentTooLong = fmt.Errorf("the id is longer than %d characters", maxIdentLength)
	ErrPINNotDigits = errors.New("a PIN has digits only")
)

func cleanIdent(id string) (string, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return "", ErrIdentEmpty
	}
	if len(id) > maxIdentLength {
		return "", ErrIdentTooLong
	}
	return id, nil
}

// NewCardReadEvent makes the event App.CardIdent of a data carrier at the internal reader. The
// id is the UID of the data carrier.
func NewCardReadEvent(id string) (Command, error) {
	id, err := cleanIdent(id)
	if err != nil {
		return Command{}, err
	}

	return NewCommand(CmdAppCardIdent, MessageTypeEvent, CardIdentEventBody{
		UID:        new(id),
		RfStandard: new("ISO14443A"),
		CardType:   new(CardTypeMifareClassic1K),
		CardTypeID: new(1),
		CardName:   new("Mifare Classic Datacarrier"),
		Segments:   []Segment{},
		Language:   new("en"),
		Device:     new(DeviceInternal),
	})
}

// NewBarcodeEvent makes the event App.CardIdent of a barcode at the first barcode reader. The id
// is the raw barcode data and the barcode ID. If the id is a decimal number, the UID is that
// number as a hex string, as a device converts it.
func NewBarcodeEvent(id string) (Command, error) {
	id, err := cleanIdent(id)
	if err != nil {
		return Command{}, err
	}

	body := CardIdentEventBody{
		RfStandard: new("UNKNOWN"),
		CardType:   new(CardTypeBarcode),
		CardTypeID: new(2),
		CardName:   new("Barcode"),
		Segments: []Segment{{
			SegmentType: new("BARCODE_DATA"),
			Data:        new(id),
			BarcodeID:   new(id),
		}},
		Language: new("en"),
		Device:   new(DeviceBarcode(1)),
	}
	if number, err := strconv.ParseUint(id, 10, 64); err == nil {
		body.UID = new(fmt.Sprintf("%08X", number))
	}

	return NewCommand(CmdAppCardIdent, MessageTypeEvent, body)
}

// NewPinEvent makes the event App.CardIdent of a PIN that a person entered at the device. The id
// is the PIN.
func NewPinEvent(id string) (Command, error) {
	id, err := cleanIdent(id)
	if err != nil {
		return Command{}, err
	}
	for _, r := range id {
		if r < '0' || r > '9' {
			return Command{}, ErrPINNotDigits
		}
	}

	return NewCommand(CmdAppCardIdent, MessageTypeEvent, CardIdentEventBody{
		UID:        new("0"),
		RfStandard: new("UNKNOWN"),
		CardType:   new(CardTypePIN),
		Segments: []Segment{{
			SegmentType: new("PIN_DATA"),
			Data:        new(id),
			PINType:     new("IDENTIFICATION"),
		}},
		Language: new("en"),
		Device:   new(DeviceInternal),
	})
}
