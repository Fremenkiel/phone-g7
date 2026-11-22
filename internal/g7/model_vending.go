package g7

// The data objects of the Vending App. Money values are in cents.

// CDStartSessionRequestBody is the data of App.CD.StartSession.
type CDStartSessionRequestBody struct {
	SessionID string `json:"SessionId"`

	// CashValue is the available cash. The host omits it for an unknown or hidden value.
	CashValue *int `json:"CashValue,omitzero"`

	PaymentInfos *PaymentInfos `json:"PaymentInfos,omitzero"`

	// RemoteVend: the terminal starts a sale. The vending machine must permit it.
	RemoteVend *RemoteVend `json:"RemoteVend,omitzero"`

	Challenge *string `json:"Challenge,omitzero"`

	// DisplayText can use the wildcard @CashValue.
	DisplayText *string `json:"DisplayText,omitzero"`

	TimeoutMs *int `json:"Timeout_ms,omitzero"`
}

// CDSessionBody is the response of App.CD.StartSession, the data of the event App.CD.EndSession
// and the response of App.CD.RevalueResult.
type CDSessionBody struct {
	SessionID *string `json:"SessionId,omitzero"`
}

type PaymentInfos struct {
	// UID is the MDB payment media ID as a hex string. The device uses 4 bytes.
	UID string `json:"UID"`

	IsTestMedia    *bool `json:"IsTestMedia,omitzero"`
	IsFreeVendCard *bool `json:"IsFreeVendCard,omitzero"`
	UserGroup      *int  `json:"UserGroup,omitzero"`
	DiscountGroup  *int  `json:"DiscountGroup,omitzero"`
	PriceList      *int  `json:"PriceList,omitzero"`

	// DiscountPercentage -100 to 0: a discount in percent. 0 to 100: a surcharge in percent.
	DiscountPercentage *int `json:"DiscountPercentage,omitzero"`
}

type RemoteVend struct {
	ItemNumber int   `json:"ItemNumber"`
	IsLastVend *bool `json:"IsLastVend,omitzero"`
}

// CDVendRequestEventBody is the data of App.CD.VendRequest. The host answers it with the request
// App.CD.VendResponse.
type CDVendRequestEventBody struct {
	SessionID *string `json:"SessionId,omitzero"`

	// Count is the number of items. More than 1 only with the basket option.
	Count *int `json:"Count,omitzero"`

	// Price is in cents.
	Price *int `json:"Price,omitzero"`

	ItemNumber *int `json:"ItemNumber,omitzero"`
}

// CDVendResponseRequestBody is the data of the request App.CD.VendResponse.
type CDVendResponseRequestBody struct {
	SessionID string `json:"SessionId"`
	Approved  bool   `json:"Approved"`

	// OnlinePrice is a price in cents that replaces the price of the machine.
	OnlinePrice *int `json:"OnlinePrice,omitzero"`

	DisplayText *string `json:"DisplayText,omitzero"`
}

// CDVendResponseResponseBody is the response to the request App.CD.VendResponse.
type CDVendResponseResponseBody struct {
	SessionID *string    `json:"SessionId,omitzero"`
	Items     []VendItem `json:"Items,omitzero"`
}

type VendItem struct {
	// Result is VEND_SUCCESS or VEND_FAILED.
	Result *string `json:"Result,omitzero"`

	// Reason is UNKNOWN_OR_OUT_OF_ORDER, SELECTION_NOT_EXIST, SELECTION_EMPTY,
	// SELECTION_DEFECTIVE, MISSING_INGRIEDIENT (spelling of GANTNER), SELECTION_BLOCKED,
	// SELECTION_INHIBITED, PRODUCT_EXPIRED, TEMPERATURE_OUT_OF_RANGE, MACHINE_BUSY_OR_WRONG_MODE,
	// INSUFFICIENT_CREDIT or TEMPORARY_BUSY.
	Reason *string `json:"Reason,omitzero"`

	// Price is in cents.
	Price *int `json:"Price,omitzero"`

	ItemNumber *int `json:"ItemNumber,omitzero"`
}

// CDRevalueRequestEventBody is the data of App.CD.RevalueRequest. The host answers it with the
// request App.CD.RevalueResult.
type CDRevalueRequestEventBody struct {
	SessionID *string `json:"SessionId,omitzero"`

	// ValueToAdd is the cash value to add, in cents.
	ValueToAdd *int `json:"ValueToAdd,omitzero"`
}

// CDRevalueResultRequestBody is the data of the request App.CD.RevalueResult.
type CDRevalueResultRequestBody struct {
	SessionID string `json:"SessionId"`
	Success   bool   `json:"Success"`

	// CashValue is the new cash value in cents.
	CashValue *int `json:"CashValue,omitzero"`
}
