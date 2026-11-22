package g7

import (
	"encoding/json"
	"fmt"
)

// CardIdentEventBody is the data of App.CardIdent: a data carrier, a barcode or a PIN was
// presented. The device waits for a command, for example App.StartUnlockProcess or
// App.StartDenyProcess.
type CardIdentEventBody struct {
	// UID is the UID of the data carrier as a hex string, most significant byte first.
	UID *string `json:"UID,omitzero"`

	// RfStandard is LEGIC, ISO15693, ISO14443A, ISO14443B, INSIDE_SECURE, SONY_FELICA, LF_FSK,
	// LF_ASK, LF_PSK or UNKNOWN.
	RfStandard *string `json:"RfStandard,omitzero"`

	// CardType is, for example, MIFARE_DESFIRE, LEGIC_ADVANT, APPLE_PASS, WIEGAND, BARCODE, PIN
	// or UNKNOWN.
	CardType *string `json:"CardType,omitzero"`

	// CardTypeID is the data carrier type as configured, 0 to 15.
	CardTypeID *int `json:"CardTypeId,omitzero"`

	CardName *string   `json:"CardName,omitzero"`
	Segments []Segment `json:"Segments,omitzero"`

	// Language is ISO 639-1, for example "en".
	Language *string `json:"Language,omitzero"`

	// Device is a reader name, see DeviceInternal. Advanced Access App only.
	Device *string `json:"Device,omitzero"`

	// Group is the door or group number, if configured.
	Group *int `json:"Group,omitzero"`

	// Direction is Undefined, Entry or Exit.
	Direction *string `json:"Direction,omitzero"`

	// EnrollmentState is ENROLLED, DISBALED (spelling of GANTNER), NOT_ENROLLED or NOT_CODED.
	// Enrollment App only.
	EnrollmentState *string `json:"EnrollmentState,omitzero"`
}

// The reader names of the Advanced Access App ("Device"). The default is INTERNAL.
const (
	DeviceInternal = "INTERNAL"
	DeviceAll      = "ALL"
)

// DeviceReader returns READER1 to READER8.
func DeviceReader(number int) string { return fmt.Sprintf("READER%d", number) }

// DeviceBarcode returns BARCODE1 or BARCODE2.
func DeviceBarcode(number int) string { return fmt.Sprintf("BARCODE%d", number) }

// DeviceExpander returns EXPANDER1 to EXPANDER8. For relays and optocouplers only.
func DeviceExpander(number int) string { return fmt.Sprintf("EXPANDER%d", number) }

// Segment is one segment of a data carrier. SegmentType tells which fields have a value: LOCKER,
// FIU, CUSTOM, WIEGAND_DATA, BARCODE_DATA, PIN_DATA, GANTNER_TOKEN, TIME_ACCESS,
// TIME_ACCESS_CARDNET_G1, TIME_ACCESS_CARDNET_G2, REPORT or EU_COVID_CERT_DATA. Times are strings
// as the card holds them, in the form 2007-01-01T00:00:00.
type Segment struct {
	SegmentType   *string `json:"SegmentType,omitzero"`
	ChecksumError *bool   `json:"ChecksumError,omitzero"`

	// FID is the company ID.
	FID    *int64 `json:"FID,omitzero"`
	SubFID *int64 `json:"SubFID,omitzero"`

	// CardType is a number in a LOCKER segment (0: user card) and a string in a TIME_ACCESS
	// segment ("USER").
	CardType json.RawMessage `json:"CardType,omitzero"`

	ValidFrom  *string `json:"ValidFrom,omitzero"`
	ValidUntil *string `json:"ValidUntil,omitzero"`

	// LOCKER

	LockerSegmentNumber *int    `json:"LockerSegmentNumber,omitzero"`
	TimeStampFirstUse   *string `json:"TimeStampFirstUse,omitzero"`
	TimeStampLastUse    *string `json:"TimeStampLastUse,omitzero"`
	UseCounter          *int    `json:"UseCounter,omitzero"`

	// BatteryInfo is the number of a lock with an empty battery in a LOCKER segment, and the
	// number of actions left in a REPORT segment.
	BatteryInfo *int `json:"BatteryInfo,omitzero"`

	TimeStampBattery *string `json:"TimeStampBattery,omitzero"`

	// BatteryStatus 0: no information. 1: battery changed. 2: battery nearly empty. 3: battery
	// empty.
	BatteryStatus *int `json:"BatteryStatus,omitzero"`

	Locked        *bool `json:"Locked,omitzero"`
	LockerNumber  *int  `json:"LockerNumber,omitzero"`
	Locked2       *bool `json:"Locked2,omitzero"`
	LockerNumber2 *int  `json:"LockerNumber2,omitzero"`
	Index         *int  `json:"Index,omitzero"`
	WasLocked     *bool `json:"WasLocked,omitzero"`
	WasLocked2    *bool `json:"WasLocked2,omitzero"`
	Disabled      *bool `json:"Disabled,omitzero"`
	CheckUseTime  *bool `json:"CheckUseTime,omitzero"`
	CertificateOK *bool `json:"CertificateOK,omitzero"`

	// FIU (fingerprint)

	Templates []string `json:"Templates,omitzero"`

	// VerificationLevel 0: no fingerprint check. 1 to 10: the level. 255: the default level.
	VerificationLevel *int `json:"VerificationLevel,omitzero"`

	// CUSTOM, BARCODE_DATA and PIN_DATA

	CustomSegmentID *int `json:"CustomSegmentId,omitzero"`

	// Data is the segment data as a hex string (CUSTOM), the raw barcode (BARCODE_DATA) or the
	// digits (PIN_DATA).
	Data *string `json:"Data,omitzero"`

	// DataOffset is the byte offset at which the device writes Data. CUSTOM, write only.
	DataOffset *int `json:"DataOffset,omitzero"`

	// BarcodeID is the extracted barcode ID. BARCODE_DATA only. The JSON name is "Id".
	BarcodeID *string `json:"Id,omitzero"`

	// PINType is IDENTIFICATION or VERIFICATION. PIN_DATA only.
	PINType *string `json:"PINType,omitzero"`

	// WIEGAND_DATA

	// SiteCode is a string by the table of GANTNER and a number in its example.
	SiteCode json.RawMessage `json:"SiteCode,omitzero"`

	// GANTNER_TOKEN (Apple Pass)

	// TokenID is the token ID as a hex string, most significant byte first. The JSON name is
	// "ID".
	TokenID *string `json:"ID,omitzero"`

	MessageTimeStamp *string `json:"MessageTimeStamp,omitzero"`

	// TIME_ACCESS, TIME_ACCESS_CARDNET_G1 and TIME_ACCESS_CARDNET_G2

	// CardNumber is the card ID as a hex string, most significant byte first.
	CardNumber *string `json:"CardNumber,omitzero"`

	CompareKey *int `json:"CompareKey,omitzero"`
	DataLength *int `json:"DataLength,omitzero"`

	// REPORT

	LockNumber           *int    `json:"LockNumber,omitzero"`
	CardNETChecksumError *bool   `json:"CardNETChecksumError,omitzero"`
	Date                 *string `json:"Date,omitzero"`

	// EU_COVID_CERT_DATA

	FirstName     *string `json:"FirstName,omitzero"`
	LastName      *string `json:"LastName,omitzero"`
	DayOfBirth    *string `json:"DayOfBirth,omitzero"`
	CovidStatusOk *bool   `json:"CovidStatusOk,omitzero"`
	SignatureOk   *bool   `json:"SignatureOk,omitzero"`
	PassValid     *bool   `json:"PassValid,omitzero"`

	// CovidStatusAddInfo is INVALID_CERT, RECOVERED_OK, TEST_NEGATIVE, TEST_POSITIVE, VACCINATED
	// or UNKNOWN_VACCINE.
	CovidStatusAddInfo *string `json:"CovidStatusAddInfo,omitzero"`

	CovidStatusTimestamp *string `json:"CovidStatusTimestamp,omitzero"`
}
