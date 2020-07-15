package models

import (
	"context"
	"log"
)

const usersCollection = "users"

type UserDatabase interface {
	Find(context.Context, interface{}) (*Users, error)
	FindOne(context.Context, interface{}) (*User, error)
	Create(context.Context, *User) error
	Delete(context.Context, string) error
	Update(context.Context, interface{}, interface{}) (error)
}

type userDatabase struct {
	db DatabaseHelper
}

type Users []User

func NewUserDatabase(db DatabaseHelper) UserDatabase {
	return &userDatabase{
		db: db,
	}
}

func (u *userDatabase) Find(ctx context.Context, filter interface{}) (*Users, error) {
	users := Users{}
	user := User{}

	cur, err := u.db.Collection(usersCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	for cur.Next(context.Background()) {
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return &users, nil
}

func (u *userDatabase) FindOne(ctx context.Context, filter interface{}) (*User, error) {
	user := &User{}
	err := u.db.Collection(usersCollection).FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userDatabase) Create(ctx context.Context, usr *User) error {
	_, err := u.db.Collection(usersCollection).InsertOne(ctx, usr)
	if err != nil {
		return err
	}
	return nil
}

func (u *userDatabase) Delete(ctx context.Context, username string) error {
	// In this case it is possible to use bson.M{"username":username} but I tend
	// to avoid another dependency in this layer and for demonstration purposes
	// used omitempty in the model
	user := &User{
		UserName: username,
	}
	_, err := u.db.Collection(usersCollection).DeleteOne(ctx, user)
	return err
}

func (u *userDatabase) Update(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := u.db.Collection(usersCollection).UpdateOne(ctx, filter, update)
	return err
}
