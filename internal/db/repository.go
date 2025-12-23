package db

import (
	"errors"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUserByTelegramID(tgID int64) (*User, error) {
	u := &User{}
	coll := mgm.Coll(u)

	err := coll.First(bson.M{"telegram_id": tgID}, u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func CreateUser(tgID int64, username, firstName string) (*User, error) {
	u := &User{
		TelegramID: tgID,
		Username:   username,
		FirstName:  firstName,
		State:      "",
	}
	if err := mgm.Coll(u).Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func UpsertUserBasic(tgID int64, username, firstName string) (*User, error) {
	u, err := FindUserByTelegramID(tgID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return CreateUser(tgID, username, firstName)
	}
	changed := false
	if u.Username != username {
		u.Username = username
		changed = true
	}
	if u.FirstName != firstName {
		u.FirstName = firstName
		changed = true
	}
	if changed {
		if err := mgm.Coll(u).Update(u); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func SetUserState(tgID int64, state string) error {
	u, err := FindUserByTelegramID(tgID)
	if err != nil {
		return err
	}
	if u == nil {
		// in un bot pro di solito crei user al primo contatto.
		u, err = CreateUser(tgID, "", "")
		if err != nil {
			return err
		}
	}
	u.State = state
	return mgm.Coll(u).Update(u)
}
