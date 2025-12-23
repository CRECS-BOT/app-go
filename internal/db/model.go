package db

import "github.com/kamva/mgm/v3"

type User struct {
	mgm.DefaultModel `bson:",inline"`

	TelegramID int64  `bson:"telegram_id"`
	Username   string `bson:"username"`
	FirstName  string `bson:"first_name"`

	// FSM state: current flow step ("" if none)
	State string `bson:"state"`

	// Example domain data
	IsAdmin bool `bson:"is_admin"`
}
