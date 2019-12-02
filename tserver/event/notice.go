package event

// CQNotice is the CQNotice event struct
type CQNotice struct {
	*CQEvent

	Type string
}
