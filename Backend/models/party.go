package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type Party struct {
	gorm.Model
	Name         string `json:"name"`
	PartyID      uint   `json:"party_id"`
	HostUsername string `json:"host_username"`
	Host         User   `gorm:"foreignKey:HostUsername;references:Username"`
	Users        Array  `gorm:"type:jsonb"`
	Invites      Array  `gorm:"type:jsonb"`
}

type PartyInvitation struct {
	gorm.Model
	PartyID      uint   `json:"party_id"`
	Status       string `json:"status"` // "pending", "accepted", "rejected"
	FromUserName string `json:"from_user_name"`
	ToUserName   string `json:"to_user_name"`
}

type Array []interface{}

// Value Marshal
func (a Array) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *Array) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (a Array) ToStringSlice() ([]string, error) {
	strSlice := make([]string, len(a))
	for i, v := range a {
		str, ok := v.(string)
		if !ok {
			return nil, errors.New("element in Array is not a string")
		}
		strSlice[i] = str
	}
	return strSlice, nil
}
