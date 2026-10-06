package model

import (
	"time"

	"github.com/google/uuid"
)

// EntityType defines the allowed types for an entity.
type EntityType string

const (
	TypeVehicle   EntityType = "vehicle"
	TypeIoTDevice EntityType = "iot_device"
	TypeFacility  EntityType = "facility"
	TypeOther     EntityType = "other"
)

// EntityStatus defines the allowed statuses for an entity.
type EntityStatus string

const (
	StatusActive      EntityStatus = "active"
	StatusInactive    EntityStatus = "inactive"
	StatusMaintenance EntityStatus = "maintenance"
	StatusOffline     EntityStatus = "offline"
)

// Entity represents a geographical entity stored in the system.
type Entity struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Type        EntityType   `json:"type"`
	Status      EntityStatus `json:"status"`
	Description string       `json:"description"`
	Latitude    float64      `json:"latitude"`
	Longitude   float64      `json:"longitude"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
