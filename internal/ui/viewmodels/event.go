package viewmodels

type Event struct {
	ID        int64  `json:"id"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
}

