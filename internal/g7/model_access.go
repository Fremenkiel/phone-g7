package g7

// StartUnlockProcessRequestBody is the data of App.StartUnlockProcess. Each field replaces a
// configured value.
type StartUnlockProcessRequestBody struct {
	// Device is a reader name, see DeviceInternal. The host sets the reader of the App.CardIdent
	// event, or the wrong door can open.
	Device *string `json:"Device,omitzero"`

	UnlockingTimeMs *int `json:"UnlockingTime_ms,omitzero"`

	// DisplayText can be HTML. The wildcards @Count and @MaxCount are available.
	DisplayText *string `json:"DisplayText,omitzero"`

	// MaxCount is the number of persons of a group entrance. The default is 1.
	MaxCount *int `json:"MaxCount,omitzero"`

	// PersonTimeoutMs is the time that one person has to enter. Group entrance only.
	PersonTimeoutMs *int `json:"PersonTimeout_ms,omitzero"`

	// IconName is none, check, user-check, users, wheelchair, male, female, arrow-up, arrow-left,
	// arrow-right, level-up, level-down, sign-in, sign-out, thumbs-up, smile, portal-enter or
	// portal-exit.
	IconName *string `json:"IconName,omitzero"`

	// IconURL is a URL or a data URL of a custom icon or image.
	IconURL *string `json:"IconUrl,omitzero"`

	// IconSize is 1x to 10x.
	IconSize *string `json:"IconSize,omitzero"`

	// LedBarColor has the form #RRGGBB. The default is #00FF00.
	LedBarColor *string `json:"LedBarColor,omitzero"`

	Relays []Relay `json:"Relays,omitzero"`

	// DoLockerCheckIn: if the locker check-in fails, the device denies the access.
	DoLockerCheckIn *bool `json:"DoLockerCheckIn,omitzero"`

	// DoFingerprintVerify: if the fingerprint check fails, the device denies the access.
	DoFingerprintVerify *bool `json:"DoFingerprintVerify,omitzero"`

	// Challenge is the one-time challenge of a QR code template.
	Challenge *string `json:"Challenge,omitzero"`

	// AppState is the App state after the process: IDLE, UNLOCKED or LOCKED.
	AppState *string `json:"AppState,omitzero"`
}

type Relay struct {
	// Device is INTERNAL, EXPANDER1 to EXPANDER8 or READER1 to READER8.
	Device string `json:"Device"`

	// ID is 1 to 8, as the device has.
	ID int `json:"Id"`

	// OnTimeMs is the pulse time. 0 turns off a configured relay.
	OnTimeMs int `json:"OnTime_ms"`
}

// FeedbackEventBody is the data of App.Feedback: the door is unlocked and the door feedback
// triggered.
type FeedbackEventBody struct {
	// Count is the number of accesses.
	Count *int `json:"Count,omitzero"`

	Device *string `json:"Device,omitzero"`
}

// UnlockProcessResultEventBody is the data of App.UnlockProcessResult: the unlock process is
// complete and the App goes to idle.
type UnlockProcessResultEventBody struct {
	AccessGranted bool `json:"AccessGranted"`

	// Feedback true: the door feedback triggered.
	Feedback bool `json:"Feedback"`

	// Count is the number of accesses.
	Count *int `json:"Count,omitzero"`

	// Error is CHECK_IN_FAILED or FINGERPRINT_VERIFY_FAILED, if the access was not granted.
	Error *string `json:"Error,omitzero"`

	Device     *string `json:"Device,omitzero"`
	UID        *string `json:"UID,omitzero"`
	PIN        *string `json:"PIN,omitzero"`
	Barcode    *string `json:"Barcode,omitzero"`
	CardTypeID *int    `json:"CardTypeId,omitzero"`
}

// StartDenyProcessRequestBody is the data of App.StartDenyProcess: show the deny screen, then go
// to idle.
type StartDenyProcessRequestBody struct {
	Device *string `json:"Device,omitzero"`

	// DisplayText can be HTML.
	DisplayText *string `json:"DisplayText,omitzero"`

	TimeoutMs *int `json:"Timeout_ms,omitzero"`

	// IconName is none, times, exclamation-triangle, frown, user-slash or users-slash.
	IconName *string `json:"IconName,omitzero"`

	// LedBarColor has the form #RRGGBB. The default is #FF0000.
	LedBarColor *string `json:"LedBarColor,omitzero"`
}

// StartLockerCheckInProcessRequestBody is the data of App.StartLockerCheckInProcess.
type StartLockerCheckInProcessRequestBody struct {
	Device *string `json:"Device,omitzero"`

	// Segments replace the configured check-in, or are segments of the host to write.
	Segments []Segment `json:"Segments,omitzero"`
}

type LockerCheckInProcessResultEventBody struct {
	Result bool `json:"Result"`

	// SuccessCount is the number of locker segments that the device wrote.
	SuccessCount *int `json:"SuccessCount,omitzero"`
}

// StartFingerVerifyProcessRequestBody is the data of App.StartFingerVerifyProcess.
type StartFingerVerifyProcessRequestBody struct {
	Device *string `json:"Device,omitzero"`

	// Templates are the templates to compare with. By default the device reads them from the
	// data carrier.
	Templates []string `json:"Templates,omitzero"`

	DisplayText *string `json:"DisplayText,omitzero"`
}

// MatchResultEventBody is the data of App.FingerVerifyProcessResult and App.PinProcessResult.
type MatchResultEventBody struct {
	Match  bool    `json:"Match"`
	Device *string `json:"Device,omitzero"`
}

// StartEnterNumberProcessRequestBody is the data of App.StartEnterNumberProcess.
type StartEnterNumberProcessRequestBody struct {
	Device      *string `json:"Device,omitzero"`
	DisplayText *string `json:"DisplayText,omitzero"`

	// ShowNumber false: the device shows * for each digit. The default is true.
	ShowNumber *bool `json:"ShowNumber,omitzero"`

	MinLength *int `json:"MinLength,omitzero"`
	MaxLength *int `json:"MaxLength,omitzero"`

	// TimeoutMs has the default 30000. Each key press starts the time again.
	TimeoutMs *int `json:"Timeout_ms,omitzero"`
}

type EnterNumberProcessResultEventBody struct {
	IsSuccess bool    `json:"IsSuccess"`
	Number    *string `json:"Number,omitzero"`
	Device    *string `json:"Device,omitzero"`
}

// StartPinVerifyProcessRequestBody is the data of App.StartPinVerifyProcess.
type StartPinVerifyProcessRequestBody struct {
	Device *string `json:"Device,omitzero"`

	// PIN is the PIN that the person must enter.
	PIN string `json:"PIN"`

	DisplayText *string `json:"DisplayText,omitzero"`

	// TimeoutMs has the default 30000. Each key press starts the time again.
	TimeoutMs *int `json:"Timeout_ms,omitzero"`
}

// StartConfirmProcessRequestBody is the data of App.StartConfirmProcess.
type StartConfirmProcessRequestBody struct {
	Device    *string `json:"Device,omitzero"`
	TimeoutMs *int    `json:"Timeout_ms,omitzero"`

	// DisplayText can be HTML. Time App: the wildcards @Price and @Duration are available.
	DisplayText *string `json:"DisplayText,omitzero"`

	// ConfirmByButton has the default true.
	ConfirmByButton *bool `json:"ConfirmByButton,omitzero"`

	// ConfirmByCard has the default false.
	ConfirmByCard *bool `json:"ConfirmByCard,omitzero"`
}

type ConfirmProcessResultEventBody struct {
	// Result true: confirmed. False: abort or timeout.
	Result bool `json:"Result"`

	Device *string `json:"Device,omitzero"`
}

// StartSelectionProcessRequestBody is the data of App.StartSelectionProcess: show a table of
// buttons and texts.
type StartSelectionProcessRequestBody struct {
	TimeoutMs *int            `json:"Timeout_ms,omitzero"`
	Items     []SelectionItem `json:"Items"`
}

type SelectionItem struct {
	// ID is the ID that App.SelectionProcessResult returns for this button.
	ID *string `json:"Id,omitzero"`

	// Row is the row number, from 1. One row can hold more than one item.
	Row int `json:"Row"`

	// Type is BUTTON (default) or TEXT.
	Type *string `json:"Type,omitzero"`

	HTML    string `json:"Html"`
	Colspan *int   `json:"Colspan,omitzero"`
	Rowspan *int   `json:"Rowspan,omitzero"`
}

type SelectionProcessResultEventBody struct {
	// Result true: a button was pressed. False: timeout.
	Result bool `json:"Result"`

	ID *string `json:"Id,omitzero"`
}

// SetBackgroundImageRequestBody is the data of App.SetBackgroundImage. Without a URL, the
// configured image applies again.
type SetBackgroundImageRequestBody struct {
	// Screen is default, idle, deny, wait, locked, unlock, info, time-select, timer-pre,
	// timer-active, timer-post, timer-pause, enrollment-select, enrollment or return.
	Screen string `json:"Screen"`

	// URL is a URL or a data URL. Wildcards: @Theme and @Random.
	URL *string `json:"Url,omitzero"`

	// RefreshIntervalMs is the interval in which the device replaces @Random. More than 500 is
	// recommended. The JSON name has the spelling of GANTNER.
	RefreshIntervalMs *int `json:"RefreshIntervall_ms,omitzero"`

	// Color is the design color.
	Color *string `json:"Color,omitzero"`
}

// SetQRCodeRequestBody is the data of App.SetQRCode: the QR code of the idle screen while the
// host is connected.
type SetQRCodeRequestBody struct {
	Data string `json:"Data"`
}

// SetOfflineQRCodeRequestBody is the data of App.SetOfflineQRCode: the QR code of the emergency
// mode. The device stores it.
type SetOfflineQRCodeRequestBody struct {
	Data string `json:"Data"`

	// OfflinePINs are the PIN codes that are valid in emergency mode.
	OfflinePINs []string `json:"OfflinePINs,omitzero"`
}

// GetAppStateRequestBody is the data of App.GetAppState.
type GetAppStateRequestBody struct {
	Device *string `json:"Device,omitzero"`
}

// AppStateBody is the response of App.GetAppState and the data of the event App.AppStateChanged.
type AppStateBody struct {
	// AppState is IDLE, BUSY, UNLOCKED (permanent) or LOCKED (permanent).
	AppState *string `json:"AppState,omitzero"`

	Device *string `json:"Device,omitzero"`
}

// SetAppStateRequestBody is the data of App.SetAppState. It cancels all active processes.
type SetAppStateRequestBody struct {
	// Device is a reader name, see DeviceInternal. The default is ALL.
	Device *string `json:"Device,omitzero"`

	// AppState is IDLE, UNLOCKED (permanent) or LOCKED (permanent).
	AppState string `json:"AppState"`

	// TimeoutMs is the time after which the device goes back to IDLE. The default is no timeout.
	TimeoutMs *int `json:"Timeout_ms,omitzero"`

	DisplayText *string `json:"DisplayText,omitzero"`

	// IconName is an icon name of App.StartUnlockProcess or of App.StartDenyProcess.
	IconName *string `json:"IconName,omitzero"`
}

// GetCameraImageRequestBody is the data of App.GetCameraImage.
type GetCameraImageRequestBody struct {
	Mirror *bool `json:"Mirror,omitzero"`
}

type GetCameraImageResponseBody struct {
	// Image is the JPEG image, Base64 encoded.
	Image *string `json:"Image,omitzero"`

	// Error is CAMERA_NOT_AVAILABLE, CAMERA_DISABLED, CAMERA_ACCESS_LIMITED, IMAGE_READ_ERROR or
	// UNKNOWN_ERROR.
	Error *string `json:"Error,omitzero"`
}

// GetOptoStatesRequestBody is the data of App.GetOptoStates.
type GetOptoStatesRequestBody struct {
	// EnableStateChangedEvent true: the device sends App.OptoStateChanged events. The default is
	// false.
	EnableStateChangedEvent *bool `json:"EnableStateChangedEvent,omitzero"`
}

type GetOptoStatesResponseBody struct {
	Optos []OptoState `json:"Optos,omitzero"`
}

// OptoState is one optocoupler state. It is also the data of the event App.OptoStateChanged.
type OptoState struct {
	// Device is INTERNAL, EXPANDER1 to EXPANDER8 or READER1 to READER8.
	Device *string `json:"Device,omitzero"`

	ID    *int `json:"Id,omitzero"`
	State bool `json:"State"`
}
