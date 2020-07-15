package models

import (
	"context"
	"github.com/satori/go.uuid"
)

const tasksCollection = "tasks"

type TaskDatabase interface {
	FindOne(context.Context, interface{}) (*Task, error)
	Create(context.Context, *Task) error
	Delete(context.Context, uuid.UUID) error
	Update(context.Context, interface{}, interface{}) (error)
}

type taskDatabase struct {
	db DatabaseHelper
}

func NewTaskDatabase(db DatabaseHelper) TaskDatabase {
	return &taskDatabase{
		db: db,
	}
}

func (u *taskDatabase) FindOne(ctx context.Context, filter interface{}) (*Task, error) {
	task := &Task{}
	err := u.db.Collection(tasksCollection).FindOne(ctx, filter).Decode(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (u *taskDatabase) Create(ctx context.Context, usr *Task) error {
	_, err := u.db.Collection(tasksCollection).InsertOne(ctx, usr)
	if err != nil {
		return err
	}
	return nil
}

func (u *taskDatabase) Delete(ctx context.Context, uid uuid.UUID) error {
	// In this case it is possible to use bson.M{"username":username} but I tend
	// to avoid another dependency in this layer and for demonstration purposes
	// used omitempty in the model
	task := &Task{
		UUID: uid,
	}
	_, err := u.db.Collection(tasksCollection).DeleteOne(ctx, task)
	return err
}

func (u *taskDatabase) Update(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := u.db.Collection(tasksCollection).UpdateOne(ctx, filter, update)
	return err
}
