package g7

import "time"

// AuthorizationRecord is one record of the authorization list of the device. Exactly one of UID,
// PIN and Barcode has a value.
type AuthorizationRecord struct {
	// ID is a record identifier of the host. Advanced Access App 1.7.0.
	ID *string `json:"Id,omitzero"`

	// UID is a hex string, most significant byte first, with an even number of characters.
	UID *string `json:"UID,omitzero"`

	PIN *string `json:"PIN,omitzero"`

	// Barcode is the raw barcode data. GANTNER recommends the UID where possible.
	Barcode *string `json:"Barcode,omitzero"`

	CardTypeID *int      `json:"CardTypeId,omitzero"`
	ValidFrom  *DateTime `json:"ValidFrom,omitzero"`
	ValidUntil *DateTime `json:"ValidUntil,omitzero"`
}

// AuthorizationsRequestBody is the data of App.Authorizations.Add and App.Authorizations.Update.
// The host sends about 1000 records at most in one command. Add does not look for duplicates.
type AuthorizationsRequestBody struct {
	Authorizations []AuthorizationRecord `json:"Authorizations"`
}

// AuthorizationsDeleteRequestBody is the data of App.Authorizations.Delete.
type AuthorizationsDeleteRequestBody struct {
	// ClearAll true: delete all records.
	ClearAll *bool `json:"ClearAll,omitzero"`

	Authorizations []AuthorizationRecord `json:"Authorizations,omitzero"`
}

// AuthorizationsGetRequestBody is the data of App.Authorizations.Get. The host reads in pages, or
// finds records by one identifier.
type AuthorizationsGetRequestBody struct {
	Skip    *int    `json:"Skip,omitzero"`
	Count   *int    `json:"Count,omitzero"`
	ID      *string `json:"Id,omitzero"`
	UID     *string `json:"UID,omitzero"`
	PIN     *string `json:"PIN,omitzero"`
	Barcode *string `json:"Barcode,omitzero"`
}

type AuthorizationsGetResponseBody struct {
	Authorizations []AuthorizationRecord `json:"Authorizations,omitzero"`
}

// BookingsGetRequestBody is the data of App.Bookings.Get: the bookings that the device stored in
// emergency mode.
type BookingsGetRequestBody struct {
	// SinceID: the device returns the bookings after this ID.
	SinceID *int64 `json:"SinceId,omitzero"`

	PageSize *int `json:"PageSize,omitzero"`
}

type BookingsGetResponseBody struct {
	Bookings []AccessBooking `json:"Bookings,omitzero"`

	// LastID is the ID of the last booking in this response.
	LastID *int64 `json:"LastId,omitzero"`

	MoreBookingsAvailable bool `json:"MoreBookingsAvailable"`
}

type AccessBooking struct {
	ID           int64              `json:"Id"`
	UTCTimestamp *time.Time         `json:"UTCTimestamp,omitzero"`
	BookingData  *AccessBookingData `json:"BookingData,omitzero"`
}

type AccessBookingData struct {
	// Booking is AccessWithFeedback, AccessNoFeedback, TimerUsed or NotAuthorized.
	Booking *string `json:"Booking,omitzero"`

	Data *AccessBookingDetails `json:"Data,omitzero"`
}

type AccessBookingDetails struct {
	UID        *string `json:"UID,omitzero"`
	PIN        *string `json:"PIN,omitzero"`
	Barcode    *string `json:"Barcode,omitzero"`
	CardTypeID *int    `json:"CardTypeId,omitzero"`

	// Duration is the selected duration in seconds. TimerUsed only.
	Duration *int `json:"Duration,omitzero"`

	Device *string `json:"Device,omitzero"`
}
