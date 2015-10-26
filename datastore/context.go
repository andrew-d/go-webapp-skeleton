package datastore

import (
	"golang.org/x/net/context"
)

type private struct{}

var contextKey private

func NewContext(parent context.Context, ds Datastore) context.Context {
	return context.WithValue(parent, contextKey, ds)
}

func FromContext(c context.Context) Datastore {
	return c.Value(contextKey).(Datastore)
}
