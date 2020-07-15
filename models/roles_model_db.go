package models

import (
	"context"
	"github.com/satori/go.uuid"
)

const rolesCollection = "roles"

type RoleDatabase interface {
	FindOne(context.Context, interface{}) (*Role, error)
	Create(context.Context, *Role) error
	Delete(context.Context, uuid.UUID) error
	Update(context.Context, interface{}, interface{}) (error)
}

type roleDatabase struct {
	db DatabaseHelper
}

func NewRoleDatabase(db DatabaseHelper) RoleDatabase {
	return &roleDatabase{
		db: db,
	}
}

func (u *roleDatabase) FindOne(ctx context.Context, filter interface{}) (*Role, error) {
	role := &Role{}
	err := u.db.Collection(rolesCollection).FindOne(ctx, filter).Decode(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (u *roleDatabase) Create(ctx context.Context, usr *Role) error {
	_, err := u.db.Collection(rolesCollection).InsertOne(ctx, usr)
	if err != nil {
		return err
	}
	return nil
}

func (u *roleDatabase) Delete(ctx context.Context, uid uuid.UUID) error {
	// In this case it is possible to use bson.M{"username":username} but I tend
	// to avoid another dependency in this layer and for demonstration purposes
	// used omitempty in the model
	role := &Role{
		//UUID: uid,
	}
	_, err := u.db.Collection(rolesCollection).DeleteOne(ctx, role)
	return err
}

func (u *roleDatabase) Update(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := u.db.Collection(rolesCollection).UpdateOne(ctx, filter, update)
	return err
}
