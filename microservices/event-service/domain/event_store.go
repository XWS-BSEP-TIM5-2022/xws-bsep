package domain

type EventStore interface {
	GetAllEvents() ([]*Event, error)
}
