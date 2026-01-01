package model

import "github.com/google/uuid"

type Activity struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Summary     string    `json:"summary"`
	Location    string    `json:"location"` // should be lat long object
	Detail      string    `json:"detail"`   // should be activity detail object
	ItineraryId uuid.UUID `json:"itineraryId"`
}
