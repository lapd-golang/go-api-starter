package models

import (
	"go-admin-starter/database"
	"time"
)

type Base struct {
	ID        int        `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

var db = database.GetDB()
