package domain

type EventType string

func (e EventType) String() string {
	return string(e)
}
