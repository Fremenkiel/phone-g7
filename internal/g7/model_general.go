package g7

import "time"

// LoginRequestBody is the data of Login (server mode only). The host gives the user and the
// password, or the one time token, or the access token.
type LoginRequestBody struct {
	// User is system, admin, user or app.
	User *string `json:"User,omitzero"`

	// Password is Base64 encoded.
	Password *string `json:"Password,omitzero"`

	OneTimeToken *string `json:"OneTimeToken,omitzero"`

	// AccessToken is the JSON web token of an earlier login. It is valid for 12 hours.
	AccessToken *string `json:"AccessToken,omitzero"`

	// OnlineDecision true: the host answers the requests of the App. Only one host can decide.
	OnlineDecision *bool `json:"OnlineDecision,omitzero"`

	// Info is the name and the version of the host software.
	Info *string `json:"Info,omitzero"`
}

type LoginResponseBody struct {
	AccessToken *string `json:"AccessToken,omitzero"`
	User        *string `json:"User,omitzero"`
}

// SetHostInfoRequestBody is the data of System.SetHostInfo (client mode, firmware 3.2.0). It
// replaces the login. The host sends it before the registration to the App.* events.
type SetHostInfoRequestBody struct {
	// OnlineDecision true: the host answers the requests of the App. Only one host can decide.
	OnlineDecision *bool `json:"OnlineDecision,omitzero"`

	// Info is the name and the version of the host software.
	Info *string `json:"Info,omitzero"`

	// FrameSize is the supported frame size in bytes. 32768 for Amazon API Gateway.
	FrameSize *int `json:"FrameSize,omitzero"`

	// MessageSize is the maximum message size in bytes. 131072 for Amazon API Gateway.
	MessageSize *int `json:"MessageSize,omitzero"`

	// ReconnectTimeout is the seconds until the device reconnects silently. 6600 for Amazon API
	// Gateway.
	ReconnectTimeout *int `json:"ReconnectTimeout,omitzero"`
}

// RegisterEventRequestBody is the data of RegisterEvent and UnregisterEvent.
type RegisterEventRequestBody struct {
	// Event is the event name. "App.*" selects all events of the active App.
	Event string `json:"Event"`
}

type GetDeviceInfoResponseBody struct {
	DeviceModel *string `json:"DeviceModel,omitzero"`

	// DeviceTypeID 113: GT7 terminal. 114: GC7 controller. 115: GRT7. 116: G7 Base. 117: G7 Vending.
	DeviceTypeID *int `json:"DeviceTypeId,omitzero"`

	DeviceName           *string          `json:"DeviceName,omitzero"`
	SerialNumber         *string          `json:"SerialNumber,omitzero"`
	ArticleNumber        *string          `json:"ArticleNumber,omitzero"`
	FirmwareVersion      *string          `json:"FirmwareVersion,omitzero"`
	HardwareVersion      *string          `json:"HardwareVersion,omitzero"`
	CoproFirmwareVersion *string          `json:"CoproFirmwareVersion,omitzero"`
	DetectedRfidReader   *string          `json:"DetectedRfidReader,omitzero"`
	ActiveApp            *DeviceActiveApp `json:"ActiveApp,omitzero"`
	BLE                  *DeviceBLE       `json:"BLE,omitzero"`
	IO                   *DeviceIO        `json:"IO,omitzero"`
	License              *DeviceLicense   `json:"License,omitzero"`
	Network              *DeviceNetwork   `json:"Network,omitzero"`
	Wlan                 *DeviceNetwork   `json:"Wlan,omitzero"`
	Routing              []DeviceRoute    `json:"Routing,omitzero"`
}

type DeviceActiveApp struct {
	// Name is the technical name, for example "access".
	Name *string `json:"Name,omitzero"`

	DisplayName   *string       `json:"DisplayName,omitzero"`
	Version       *string       `json:"Version,omitzero"`
	Author        *string       `json:"Author,omitzero"`
	Date          *string       `json:"Date,omitzero"`
	LastChangedAt *time.Time    `json:"LastChangedAt,omitzero"`
	Level         *int          `json:"Level,omitzero"`
	Plugin        *DevicePlugin `json:"Plugin,omitzero"`
}

type DevicePlugin struct {
	DisplayName *string `json:"DisplayName,omitzero"`
	Version     *string `json:"Version,omitzero"`
	Author      *string `json:"Author,omitzero"`
	Date        *string `json:"Date,omitzero"`
}

type DeviceBLE struct {
	// DeviceTypeID 0: not specified. 1: QS9322. 3: BLEM_SL22.
	DeviceTypeID *int `json:"DeviceTypeId,omitzero"`

	FirmwareVersion      *string `json:"FirmwareVersion,omitzero"`
	BootloaderVersion    *string `json:"BootloaderVersion,omitzero"`
	SecureElementVersion *string `json:"SEVersion,omitzero"`
	MAC                  *string `json:"MAC,omitzero"`
}

type DeviceIO struct {
	NumberRelays *int `json:"NumberRelays,omitzero"`
	NumberOptos  *int `json:"NumberOptos,omitzero"`
}

type DeviceLicense struct {
	// Level is the available license points.
	Level *int `json:"Level,omitzero"`

	OpenCard      *bool   `json:"OpenCard,omitzero"`
	EmbeddedEloxx *bool   `json:"EmbeddedEloxx,omitzero"`
	FIDList       []int64 `json:"FIDList,omitzero"`
	SubDevices    *int    `json:"SubDevices,omitzero"`
}

type DeviceNetwork struct {
	IP     *string  `json:"Ip,omitzero"`
	IPList []string `json:"IpList,omitzero"`
	MAC    *string  `json:"MAC,omitzero"`
}

type DeviceRoute struct {
	Dest *string `json:"Dest,omitzero"`
	Dev  *string `json:"Dev,omitzero"`
	Gw   *string `json:"Gw,omitzero"`
}
