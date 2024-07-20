package models

import "gorm.io/gorm"

type FriendRequest struct {
	gorm.Model
	FromUserName string   `json:"from_user_name"`
	ToUserName   string   `json:"to_user_name"`
	Status     string // "pending", "accepted", "rejected"
}

type FriendshipRequestDao struct {
	gorm.Model
	FromUserName string   `json:"from_user_name"`
	ToUserName   string   `json:"to_user_name"`
	Status     string `json:"status"` // e.g., "pending", "accepted", "rejected"
	RequestedBy uint `gorm:"not null"` // User who sent the request
	RequestedTo uint `gorm:"not null"` // User who received the request
}

