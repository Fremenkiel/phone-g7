package g7

// CommandType is the command name ("Cmd") of a message. A constant without a remark is a request
// from the host to the device.
type CommandType string

// General commands.
const (
	// CmdHeartbeat is a request from the host or from the device.
	CmdHeartbeat         CommandType = "Heartbeat"
	CmdLogin             CommandType = "Login"
	CmdSystemSetHostInfo CommandType = "System.SetHostInfo"
	CmdGetDeviceInfo     CommandType = "GetDeviceInfo"
	CmdRegisterEvent     CommandType = "RegisterEvent"
	CmdUnregisterEvent   CommandType = "UnregisterEvent"
)

// Identification.
const (
	// CmdAppCardIdent is an event.
	CmdAppCardIdent CommandType = "App.CardIdent"
)

// Access.
const (
	CmdAppStartUnlockProcess CommandType = "App.StartUnlockProcess"
	// CmdAppFeedback is an event: the door feedback triggered.
	CmdAppFeedback CommandType = "App.Feedback"
	// CmdAppUnlockProcessResult is an event.
	CmdAppUnlockProcessResult CommandType = "App.UnlockProcessResult"
	CmdAppStartDenyProcess    CommandType = "App.StartDenyProcess"
)

// Locker check-in.
const (
	CmdAppStartLockerCheckInProcess CommandType = "App.StartLockerCheckInProcess"
	// CmdAppLockerCheckInProcessResult is an event.
	CmdAppLockerCheckInProcessResult CommandType = "App.LockerCheckInProcessResult"
)

// Fingerprint.
const (
	CmdAppStartFingerVerifyProcess CommandType = "App.StartFingerVerifyProcess"
	// CmdAppFingerVerifyProcessResult is an event.
	CmdAppFingerVerifyProcessResult      CommandType = "App.FingerVerifyProcessResult"
	CmdAppStartFingerEnrollmentProcess   CommandType = "App.StartFingerEnrollmentProcess"
	CmdAppStartDisableFingerCheckProcess CommandType = "App.StartDisableFingerCheckProcess"
	// CmdAppEnrollmentProcessResult is an event.
	CmdAppEnrollmentProcessResult CommandType = "App.EnrollmentProcessResult"
)

// PIN and number input.
const (
	CmdAppStartEnterNumberProcess CommandType = "App.StartEnterNumberProcess"
	// CmdAppEnterNumberProcessResult is an event.
	CmdAppEnterNumberProcessResult CommandType = "App.EnterNumberProcessResult"
	CmdAppStartPinVerifyProcess    CommandType = "App.StartPinVerifyProcess"
	// CmdAppPinProcessResult is an event.
	CmdAppPinProcessResult CommandType = "App.PinProcessResult"
)

// Confirm request.
const (
	CmdAppStartConfirmProcess CommandType = "App.StartConfirmProcess"
	// CmdAppConfirmProcessResult is an event.
	CmdAppConfirmProcessResult CommandType = "App.ConfirmProcessResult"
)

// Membership selection request.
const (
	CmdAppStartSelectionProcess CommandType = "App.StartSelectionProcess"
	// CmdAppSelectionProcessResult is an event.
	CmdAppSelectionProcessResult CommandType = "App.SelectionProcessResult"
)

// QR code and background image.
const (
	CmdAppSetBackgroundImage CommandType = "App.SetBackgroundImage"
	CmdAppSetQRCode          CommandType = "App.SetQRCode"
	CmdAppSetOfflineQRCode   CommandType = "App.SetOfflineQRCode"
)

// App state.
const (
	CmdAppGetAppState CommandType = "App.GetAppState"
	CmdAppSetAppState CommandType = "App.SetAppState"
	// CmdAppAppStateChanged is an event.
	CmdAppAppStateChanged CommandType = "App.AppStateChanged"
)

// Camera.
const (
	CmdAppGetCameraImage CommandType = "App.GetCameraImage"
)

// Authorization list.
const (
	CmdAppAuthorizationsAdd    CommandType = "App.Authorizations.Add"
	CmdAppAuthorizationsUpdate CommandType = "App.Authorizations.Update"
	CmdAppAuthorizationsDelete CommandType = "App.Authorizations.Delete"
	CmdAppAuthorizationsGet    CommandType = "App.Authorizations.Get"
)

// Bookings of the Access, Time, Enrollment and Return Apps.
const (
	CmdAppBookingsGet    CommandType = "App.Bookings.Get"
	CmdAppBookingsDelete CommandType = "App.Bookings.Delete"
	// CmdAppNewBookingAvailable is an event.
	CmdAppNewBookingAvailable CommandType = "App.NewBookingAvailable"
)

// Optocoupler state.
const (
	CmdAppGetOptoStates CommandType = "App.GetOptoStates"
	// CmdAppOptoStateChanged is an event.
	CmdAppOptoStateChanged CommandType = "App.OptoStateChanged"
)

// G7 Info App.
const (
	CmdAppStartInfoProcess CommandType = "App.StartInfoProcess"
)

// G7 Time App and G7 Countdown App.
const (
	CmdAppStartIdentProcess          CommandType = "App.StartIdentProcess"
	CmdAppStartDurationSelectProcess CommandType = "App.StartDurationSelectProcess"
	// CmdAppDurationSelected is an event.
	CmdAppDurationSelected  CommandType = "App.DurationSelected"
	CmdAppStartTimerProcess CommandType = "App.StartTimerProcess"
	// CmdAppTimerStarted is an event.
	CmdAppTimerStarted CommandType = "App.TimerStarted"
	// CmdAppTimerElapsed is an event.
	CmdAppTimerElapsed CommandType = "App.TimerElapsed"
)

// G7 Return App.
const (
	CmdAppStartReturnProcess CommandType = "App.StartReturnProcess"
	// CmdAppReturnProcessResult is an event.
	CmdAppReturnProcessResult CommandType = "App.ReturnProcessResult"
	CmdAppGetLevel            CommandType = "App.GetLevel"
	// CmdAppLevelChanged is an event.
	CmdAppLevelChanged CommandType = "App.LevelChanged"
)

// G7 Writer App.
const (
	CmdAppStartSuccessProcess CommandType = "App.StartSuccessProcess"
)

// G7 Vending App (cashless device).
const (
	CmdAppCDStartSession CommandType = "App.CD.StartSession"
	// CmdAppCDEndSession is an event.
	CmdAppCDEndSession CommandType = "App.CD.EndSession"
	// CmdAppCDVendRequest is an event.
	CmdAppCDVendRequest CommandType = "App.CD.VendRequest"
	// CmdAppCDVendResponse is a request from the host. It answers the event CmdAppCDVendRequest.
	CmdAppCDVendResponse CommandType = "App.CD.VendResponse"
	// CmdAppCDRevalueRequest is an event.
	CmdAppCDRevalueRequest CommandType = "App.CD.RevalueRequest"
	// CmdAppCDRevalueResult is a request from the host. It answers the event
	// CmdAppCDRevalueRequest.
	CmdAppCDRevalueResult CommandType = "App.CD.RevalueResult"
)

// G7 Time and Attendance App.
const (
	// CmdAppPersonIdent is an event.
	CmdAppPersonIdent CommandType = "App.PersonIdent"
	// CmdAppBalance is an event.
	CmdAppBalance                 CommandType = "App.Balance"
	CmdAppStartGrantProcess       CommandType = "App.StartGrantProcess"
	CmdAppStartBalancesProcess    CommandType = "App.StartBalancesProcess"
	CmdAppGetBookings             CommandType = "App.GetBookings"
	CmdAppDeleteBooking           CommandType = "App.DeleteBooking"
	CmdAppSavePersons             CommandType = "App.SavePersons"
	CmdAppSaveDataCarriers        CommandType = "App.SaveDataCarriers"
	CmdAppSaveFingerprints        CommandType = "App.SaveFingerprints"
	CmdAppGetPersons              CommandType = "App.GetPersons"
	CmdAppGetDataCarriers         CommandType = "App.GetDataCarriers"
	CmdAppGetFingerprints         CommandType = "App.GetFingerprints"
	CmdAppDeletePerson            CommandType = "App.DeletePerson"
	CmdAppDeleteDataCarrier       CommandType = "App.DeleteDataCarrier"
	CmdAppDeleteFingerprint       CommandType = "App.DeleteFingerprint"
	CmdAppStartPersonalizeProcess CommandType = "App.StartPersonalizeProcess"
)
