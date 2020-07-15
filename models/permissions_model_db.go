package models

import (
	"context"
	"github.com/satori/go.uuid"
)

const permissionsCollection = "permissions"

type PermissionDatabase interface {
	FindOne(context.Context, interface{}) (*Permission, error)
	Create(context.Context, *Permission) error
	Delete(context.Context, uuid.UUID) error
	Update(context.Context, interface{}, interface{}) (error)
}

type permissionDatabase struct {
	db DatabaseHelper
}

func NewPermissionDatabase(db DatabaseHelper) PermissionDatabase {
	return &permissionDatabase{
		db: db,
	}
}

func (u *permissionDatabase) FindOne(ctx context.Context, filter interface{}) (*Permission, error) {
	permission := &Permission{}
	err := u.db.Collection(permissionsCollection).FindOne(ctx, filter).Decode(permission)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (u *permissionDatabase) Create(ctx context.Context, usr *Permission) error {
	_, err := u.db.Collection(permissionsCollection).InsertOne(ctx, usr)
	if err != nil {
		return err
	}
	return nil
}

func (u *permissionDatabase) Delete(ctx context.Context, uid uuid.UUID) error {
	// In this case it is possible to use bson.M{"username":username} but I tend
	// to avoid another dependency in this layer and for demonstration purposes
	// used omitempty in the model
	permission := &Permission{
		UUID: uid,
	}
	_, err := u.db.Collection(permissionsCollection).DeleteOne(ctx, permission)
	return err
}

func (u *permissionDatabase) Update(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := u.db.Collection(permissionsCollection).UpdateOne(ctx, filter, update)
	return err
}
