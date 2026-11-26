package action

// IdentRequest is the body of the card read, barcode and PIN actions.
type IdentRequest struct {
	ID string `json:"id"`
}

type IdentResponse struct {
	// ID is the ID of the event in the database.
	ID int64 `json:"id"`

	// Subscribed false: no host is registered to the event. The event stays in the queue.
	Subscribed bool `json:"subscribed"`
}
