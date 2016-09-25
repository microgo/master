package helper

import (
	"master/model/postgres"
	"master/resource"
)

type HelperInterface interface {
	GetUserByID(int) (*postgres.User, error)
}

type Helper struct {
	*resource.Resource
}

func NewResourceHelper(r *resource.Resource) HelperInterface {
	return &Helper{Resource: r}
}
