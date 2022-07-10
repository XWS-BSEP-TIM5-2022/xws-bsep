package domain

type EventStore interface {
	GetAllEvents() ([]*Event, error)
	NewEvent(event *Event) (*Event, error)
}
