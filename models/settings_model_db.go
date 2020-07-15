package models

import (
	"context"
	"github.com/satori/go.uuid"
)

const settingsCollection = "settings"

type SettingDatabase interface {
	FindOne(context.Context, interface{}) (*Setting, error)
	Create(context.Context, *Setting) error
	Delete(context.Context, uuid.UUID) error
	Update(context.Context, interface{}, interface{}) (error)
}

type settingDatabase struct {
	db DatabaseHelper
}

func NewSettingDatabase(db DatabaseHelper) SettingDatabase {
	return &settingDatabase{
		db: db,
	}
}

func (u *settingDatabase) FindOne(ctx context.Context, filter interface{}) (*Setting, error) {
	setting := &Setting{}
	err := u.db.Collection(settingsCollection).FindOne(ctx, filter).Decode(setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func (u *settingDatabase) Create(ctx context.Context, usr *Setting) error {
	_, err := u.db.Collection(settingsCollection).InsertOne(ctx, usr)
	if err != nil {
		return err
	}
	return nil
}

func (u *settingDatabase) Delete(ctx context.Context, uid uuid.UUID) error {
	// In this case it is possible to use bson.M{"username":username} but I tend
	// to avoid another dependency in this layer and for demonstration purposes
	// used omitempty in the model
	setting := &Setting{
		UUID: uid,
	}
	_, err := u.db.Collection(settingsCollection).DeleteOne(ctx, setting)
	return err
}

func (u *settingDatabase) Update(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := u.db.Collection(settingsCollection).UpdateOne(ctx, filter, update)
	return err
}
