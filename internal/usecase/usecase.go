package usecase

import "context"

type Input interface{}

type Output interface{}

type Usecase[I Input, O Output] interface {
	Execute(context.Context, I) (O, error)
}
