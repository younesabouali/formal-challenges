// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventNames string

const (
	EventNamesSighted     EventNames = "sighted"
	EventNamesAreaChecked EventNames = "area_checked"
	EventNamesFound       EventNames = "found"
	EventNamesRewarded    EventNames = "rewarded"
)

func (e *EventNames) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EventNames(s)
	case string:
		*e = EventNames(s)
	default:
		return fmt.Errorf("unsupported scan type for EventNames: %T", src)
	}
	return nil
}

type NullEventNames struct {
	EventNames EventNames
	Valid      bool // Valid is true if EventNames is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEventNames) Scan(value interface{}) error {
	if value == nil {
		ns.EventNames, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EventNames.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEventNames) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EventNames), nil
}

type MissingPetsStatus string

const (
	MissingPetsStatusMissing MissingPetsStatus = "missing"
	MissingPetsStatusFound   MissingPetsStatus = "found"
)

func (e *MissingPetsStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MissingPetsStatus(s)
	case string:
		*e = MissingPetsStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for MissingPetsStatus: %T", src)
	}
	return nil
}

type NullMissingPetsStatus struct {
	MissingPetsStatus MissingPetsStatus
	Valid             bool // Valid is true if MissingPetsStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMissingPetsStatus) Scan(value interface{}) error {
	if value == nil {
		ns.MissingPetsStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MissingPetsStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMissingPetsStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MissingPetsStatus), nil
}

type Roles string

const (
	RolesAdmin Roles = "admin"
	RolesUser  Roles = "user"
)

func (e *Roles) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Roles(s)
	case string:
		*e = Roles(s)
	default:
		return fmt.Errorf("unsupported scan type for Roles: %T", src)
	}
	return nil
}

type NullRoles struct {
	Roles Roles
	Valid bool // Valid is true if Roles is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRoles) Scan(value interface{}) error {
	if value == nil {
		ns.Roles, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Roles.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRoles) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Roles), nil
}

type MissingPet struct {
	ID          uuid.UUID
	Createdat   time.Time
	Updatedat   time.Time
	PetName     string
	Description sql.NullString
	ImageUrl    sql.NullString
	Status      MissingPetsStatus
	LostIn      interface{}
	LostAt      time.Time
	UserID      uuid.UUID
}

type PetsRecord struct {
	ID          uuid.UUID
	Createdat   time.Time
	Updatedat   time.Time
	PetID       uuid.UUID
	UserID      uuid.UUID
	EventName   EventNames
	ImageUrl    sql.NullString
	Description sql.NullString
}

type User struct {
	ID        uuid.UUID
	Createdat time.Time
	Updatedat time.Time
	Name      string
	Email     string
	Password  string
	Role      Roles
	ApiKey    string
}
