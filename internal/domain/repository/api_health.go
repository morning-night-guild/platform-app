package repository

import "context"

type APIHealth interface {
	Check(context.Context) error
}
