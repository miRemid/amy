package event

// CQRequest is the CQRequest event struct
type CQRequest struct {
	*CQEvent

	Type string
}
