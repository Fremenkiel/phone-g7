package g7

import "time"

// The data objects of the Time and Attendance App.

// BalanceEventBody is the data of App.Balance: a person selected the balance function and was
// identified.
type BalanceEventBody struct {
	// Card is the data carrier. It is nil if the person used a PIN or a fingerprint.
	Card *Card `json:"Card,omitzero"`

	// IDType 0: data carrier. 1: PIN pad. 2: fingerprint.
	IDType *int `json:"IdType,omitzero"`

	// IDValue is the UID for ID type 0, and the personal number for the ID types 1 and 2.
	IDValue *string `json:"IdValue,omitzero"`

	IDValueType *int       `json:"IdValueType,omitzero"`
	Language    *string    `json:"Language,omitzero"`
	Timestamp   *time.Time `json:"Timestamp,omitzero"`
}

// PersonIdentEventBody is the data of App.PersonIdent: a person was identified. The codes are as
// configured in the App. The JSON holds the fields of BalanceEventBody at the same level.
type PersonIdentEventBody struct {
	BalanceEventBody

	// BookingCode is, for example, coming, leaving or pause.
	BookingCode *string `json:"BookingCode,omitzero"`

	ReasonCode  *string `json:"ReasonCode,omitzero"`
	AccountCode *string `json:"AccountCode,omitzero"`
}

type Card struct {
	CardType   *string   `json:"CardType,omitzero"`
	RfStandard *string   `json:"RfStandard,omitzero"`
	Segments   []Segment `json:"Segments,omitzero"`
	UID        *string   `json:"UID,omitzero"`
}

// StartGrantProcessRequestBody is the data of App.StartGrantProcess: tell the person that the
// request was handled.
type StartGrantProcessRequestBody struct {
	TimeoutMs   *int    `json:"Timeout_ms,omitzero"`
	DisplayText *string `json:"DisplayText,omitzero"`
	IconName    *string `json:"IconName,omitzero"`

	// IconClass is a CSS class for the icon.
	IconClass *string `json:"IconClass,omitzero"`

	// Autonomous true: the App handles the booking as an offline booking.
	Autonomous *bool `json:"Autonomous,omitzero"`

	// SoundFile is a sound file of the theme.
	SoundFile *string `json:"SoundFile,omitzero"`
}

// StartBalancesProcessRequestBody is the data of App.StartBalancesProcess.
type StartBalancesProcessRequestBody struct {
	TimeoutMs   *int      `json:"Timeout_ms,omitzero"`
	DisplayText *string   `json:"DisplayText,omitzero"`
	Balances    []Balance `json:"Balances"`
}

type Balance struct {
	Name  string  `json:"Name"`
	Value *string `json:"Value,omitzero"`
}

// GetBookingsRequestBody is the data of App.GetBookings: the bookings that the App made offline.
type GetBookingsRequestBody struct {
	Skip  *int `json:"Skip,omitzero"`
	Count *int `json:"Count,omitzero"`
}

type GetBookingsResponseBody struct {
	TotalCount *int                `json:"TotalCount,omitzero"`
	Bookings   []AttendanceBooking `json:"Bookings,omitzero"`
}

type AttendanceBooking struct {
	ID         int64   `json:"Id"`
	TerminalID *string `json:"TerminalId,omitzero"`
	Card       *Card   `json:"Card,omitzero"`

	// IDType 0: data carrier. 1: PIN pad. 2: fingerprint.
	IDType *int `json:"IdType,omitzero"`

	IDValue     *string    `json:"IdValue,omitzero"`
	IDValueType *int       `json:"IdValueType,omitzero"`
	Timestamp   *time.Time `json:"Timestamp,omitzero"`
	BookingCode *string    `json:"BookingCode,omitzero"`
	ReasonCode  *string    `json:"ReasonCode,omitzero"`
	AccountCode *string    `json:"AccountCode,omitzero"`

	// FiuInfo 0: no verification. 1: hit. 2: no hit. 3: timeout. 4: no template. 5: error.
	FiuInfo *int `json:"FiuInfo,omitzero"`

	// FiuScore 0: best. 1: good. 2: OK. 3: poor. 4: very poor.
	FiuScore *int `json:"FiuScore,omitzero"`

	// PresenceState 0: not used. 1: present. 2: absent. 3: paused.
	PresenceState *int `json:"PresenceState,omitzero"`

	// Exported true: the booking is in the bookings.exp FTP file.
	Exported *bool `json:"Exported,omitzero"`
}

// DeleteBookingRequestBody is the data of App.DeleteBooking.
type DeleteBookingRequestBody struct {
	// ID is the booking to delete. Not used if ClearAll is true or Count is more than 0.
	ID int64 `json:"Id"`

	ClearAll bool `json:"ClearAll"`

	// Count is the number of bookings to delete, from the first.
	Count int `json:"Count"`
}

// SavePersonsRequestBody is the data of App.SavePersons: add or update persons for offline
// operation.
type SavePersonsRequestBody struct {
	// Clear true: delete all persons first.
	Clear bool `json:"Clear"`

	Persons []Person `json:"Persons"`
}

type Person struct {
	ID string `json:"Id"`

	// PersonalNumber is the number with which the person can identify at the PIN pad.
	PersonalNumber *string `json:"PersonalNumber,omitzero"`

	FirstName *string `json:"FirstName,omitzero"`
	LastName  *string `json:"LastName,omitzero"`
	Language  *string `json:"Language,omitzero"`

	// FiuVerification 0: not used. 1: from the database. 2: from the card.
	FiuVerification *int `json:"FiuVerification,omitzero"`

	// PresenceState 0: not used. 1: present. 2: absent. 3: paused.
	PresenceState *int `json:"PresenceState,omitzero"`

	Balances []Balance `json:"Balances,omitzero"`

	// DataCarriers is for App.GetPersons only.
	DataCarriers []DataCarrier `json:"DataCarriers,omitzero"`

	// FiuTemplates is for App.GetPersons only.
	FiuTemplates []Fingerprint `json:"FiuTemplates,omitzero"`
}

// SaveRecordsResponseBody is the response of App.SavePersons, App.SaveDataCarriers and
// App.SaveFingerprints.
type SaveRecordsResponseBody struct {
	AffectedRecords *int    `json:"AffectedRecords,omitzero"`
	Error           *string `json:"Error,omitzero"`

	// FailedPersons are the person IDs of the records that the device did not save.
	FailedPersons []string `json:"FailedPersons,omitzero"`

	// FiuRemainingPersons is the number of persons that the fingerprint reader can still store.
	// App.SaveFingerprints only.
	FiuRemainingPersons *int `json:"FiuRemainingPersons,omitzero"`
}

// SaveDataCarriersRequestBody is the data of App.SaveDataCarriers.
type SaveDataCarriersRequestBody struct {
	// Clear true: delete all data carriers first.
	Clear bool `json:"Clear"`

	DataCarriers []DataCarrier `json:"DataCarriers"`
}

type DataCarrier struct {
	// PersonID is not present in the data carriers of a Person.
	PersonID *string `json:"PersonId,omitzero"`

	// UID is a hex string, most significant byte first. The JSON name is "Uid" in this App.
	UID string `json:"Uid"`

	CardType *int `json:"CardType,omitzero"`
}

// SaveFingerprintsRequestBody is the data of App.SaveFingerprints.
type SaveFingerprintsRequestBody struct {
	// Clear true: delete all fingerprints first.
	Clear bool `json:"Clear"`

	// AutoAdjustFingerTemplates true: the device corrects template headers that it does not
	// support. Usually false.
	AutoAdjustFingerTemplates bool `json:"AutoAdjustFingerTemplates"`

	Fingerprints []Fingerprint `json:"Fingerprints"`
}

type Fingerprint struct {
	// PersonID is not present in the templates of a Person.
	PersonID *string `json:"PersonId,omitzero"`

	// No is the finger number, 1 to 10.
	No int `json:"No"`

	// Data is the template data.
	Data string `json:"Data"`
}

// GetPersonsRequestBody is the data of App.GetPersons. With an ID, the device ignores Skip and
// Count.
type GetPersonsRequestBody struct {
	ID    *string `json:"Id,omitzero"`
	Skip  *int    `json:"Skip,omitzero"`
	Count *int    `json:"Count,omitzero"`
}

type GetPersonsResponseBody struct {
	TotalCount *int     `json:"TotalCount,omitzero"`
	Persons    []Person `json:"Persons,omitzero"`
}

// GetPersonRecordsRequestBody is the data of App.GetDataCarriers and App.GetFingerprints. With a
// person ID, the device ignores Skip and Count.
type GetPersonRecordsRequestBody struct {
	PersonID *string `json:"PersonId,omitzero"`
	Skip     *int    `json:"Skip,omitzero"`
	Count    *int    `json:"Count,omitzero"`
}

type GetDataCarriersResponseBody struct {
	TotalCount   *int          `json:"TotalCount,omitzero"`
	DataCarriers []DataCarrier `json:"DataCarriers,omitzero"`
}

type GetFingerprintsResponseBody struct {
	TotalCount   *int          `json:"TotalCount,omitzero"`
	Fingerprints []Fingerprint `json:"Fingerprints,omitzero"`
}

// DeletePersonRequestBody is the data of App.DeletePerson.
type DeletePersonRequestBody struct {
	ID *string `json:"Id,omitzero"`

	// ClearAll true: delete all persons.
	ClearAll bool `json:"ClearAll"`
}

// DeleteDataCarrierRequestBody is the data of App.DeleteDataCarrier. Without a person ID, the UID
// selects the data carriers of all persons.
type DeleteDataCarrierRequestBody struct {
	PersonID *string `json:"PersonId,omitzero"`

	// UID has the JSON name "Uid" in this App.
	UID *string `json:"Uid,omitzero"`

	// ClearAll true: delete all data carriers.
	ClearAll bool `json:"ClearAll"`
}

// DeleteFingerprintRequestBody is the data of App.DeleteFingerprint. Without a person ID, the
// number selects the fingerprints of all persons.
type DeleteFingerprintRequestBody struct {
	PersonID *string `json:"PersonId,omitzero"`
	No       *int    `json:"No,omitzero"`

	// ClearAll true: delete all fingerprints.
	ClearAll bool `json:"ClearAll"`
}

// StartPersonalizeProcessRequestBody is the data of App.StartPersonalizeProcess: change the
// language until the welcome screen shows again.
type StartPersonalizeProcessRequestBody struct {
	// Language is ISO 639-1, for example "en".
	Language string `json:"Language"`
}
