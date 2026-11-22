package g7

// The data objects of the Info, Time, Countdown, Enrollment, Return and Writer Apps.

// StartInfoProcessRequestBody is the data of App.StartInfoProcess: show up to 4 information
// screens.
type StartInfoProcessRequestBody struct {
	// DisplayText1 can be HTML. An empty string hides a configured screen.
	DisplayText1 *string `json:"DisplayText1,omitzero"`
	DisplayText2 *string `json:"DisplayText2,omitzero"`
	DisplayText3 *string `json:"DisplayText3,omitzero"`
	DisplayText4 *string `json:"DisplayText4,omitzero"`

	// Accepted: the device beeps and signals accepted (true) or denied (false).
	Accepted *bool `json:"Accepted,omitzero"`

	// TimeoutMs is the time for each screen.
	TimeoutMs *int `json:"Timeout_ms,omitzero"`

	// IconName is a Font Awesome 5 icon name. The default is "info". An empty string shows no
	// icon.
	IconName *string `json:"IconName,omitzero"`

	// OnlineData are the rows of the @Table wildcard. The keys are wildcard names, for example
	// OnlineGroup, OnlineLabel, OnlineValidFrom and OnlineValidTo, or custom names.
	OnlineData []map[string]string `json:"OnlineData,omitzero"`

	Challenge *string `json:"Challenge,omitzero"`
}

// StartIdentProcessRequestBody is the data of App.StartIdentProcess.
type StartIdentProcessRequestBody struct {
	TimeoutMs   *int    `json:"Timeout_ms,omitzero"`
	DisplayText *string `json:"DisplayText,omitzero"`
}

// StartDurationSelectProcessRequestBody is the data of App.StartDurationSelectProcess. Times are
// in seconds.
type StartDurationSelectProcessRequestBody struct {
	DisplayText *string  `json:"DisplayText,omitzero"`
	BaseTime    *int     `json:"BaseTime,omitzero"`
	BasePrice   *float64 `json:"BasePrice,omitzero"`
	StepTime    *int     `json:"StepTime,omitzero"`
	MaxSteps    *int     `json:"MaxSteps,omitzero"`
	StepPrice   *float64 `json:"StepPrice,omitzero"`

	// Currency is the currency symbol.
	Currency *string `json:"Currency,omitzero"`
}

type DurationSelectedEventBody struct {
	// Result true: a duration was selected. False: timeout or abort.
	Result bool `json:"Result"`

	// Step is the selected step. The base step is 0.
	Step *int `json:"Step,omitzero"`

	Price *float64 `json:"Price,omitzero"`

	// Time is the selected duration in seconds.
	Time *int `json:"Time,omitzero"`
}

// StartTimerProcessRequestBody is the data of App.StartTimerProcess. Times are in seconds.
type StartTimerProcessRequestBody struct {
	PreTime        *int    `json:"PreTime,omitzero"`
	PreDisplayText *string `json:"PreDisplayText,omitzero"`

	// Time is the active time. The default is the selected duration.
	Time *int `json:"Time,omitzero"`

	DisplayText     *string `json:"DisplayText,omitzero"`
	PostTime        *int    `json:"PostTime,omitzero"`
	PostDisplayText *string `json:"PostDisplayText,omitzero"`
	Challenge       *string `json:"Challenge,omitzero"`
}

// StartFingerEnrollmentProcessRequestBody is the data of App.StartFingerEnrollmentProcess.
type StartFingerEnrollmentProcessRequestBody struct {
	// NumberOfFingers is 1 or 2.
	NumberOfFingers *int `json:"NumberOfFingers,omitzero"`

	TimeoutMs *int `json:"Timeout_ms,omitzero"`
}

// StartDisableFingerCheckProcessRequestBody is the data of App.StartDisableFingerCheckProcess.
type StartDisableFingerCheckProcessRequestBody struct {
	TimeoutMs *int `json:"Timeout_ms,omitzero"`
}

type EnrollmentProcessResultEventBody struct {
	Result bool `json:"Result"`
}

// StartReturnProcessRequestBody is the data of App.StartReturnProcess: take back a data carrier.
type StartReturnProcessRequestBody struct {
	TimeoutMs *int `json:"Timeout_ms,omitzero"`

	// DoLockerCheckIn has the default false.
	DoLockerCheckIn *bool `json:"DoLockerCheckIn,omitzero"`

	DisplayText *string `json:"DisplayText,omitzero"`
}

type ReturnProcessResultEventBody struct {
	// Result true: the device has the data carrier.
	Result bool `json:"Result"`

	// Error is CHECK_IN_FAILED, NO_TAG_IN_READER, WRONG_TAG_IN_READER, LEVEL_FULL, TIMEOUT or
	// MANIPULATION.
	Error *string `json:"Error,omitzero"`
}

// LevelBody is the response of App.GetLevel and the data of the event App.LevelChanged.
type LevelBody struct {
	// Level is the fill level in percent.
	Level *int `json:"Level,omitzero"`

	// State is LEVEL_OK, LEVEL_WARNING or LEVEL_FULL.
	State *string `json:"State,omitzero"`
}

// StartSuccessProcessRequestBody is the data of App.StartSuccessProcess: signal success, then go
// to idle.
type StartSuccessProcessRequestBody struct {
	TimeoutMs *int `json:"Timeout_ms,omitzero"`

	// LedBarColor has the form #RRGGBB. The default is #00FF00.
	LedBarColor *string `json:"LedBarColor,omitzero"`
}
