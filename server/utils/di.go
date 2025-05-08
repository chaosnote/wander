package utils

import (
	"fmt"
)

const (
	SERVICE_LOGGER  = "logger"
	SERVICE_MARIADB = "mariadb"
	SERVICE_NATS    = "nats"
	SERVICE_REDIS   = "redis"
)

var ErrKeyAlreadyExists = fmt.Errorf("map: key already exists")

type CreateFunc func() any

type DI interface {
	Set(key string, create CreateFunc) (e error)
	Get(key string) (service any, ok bool)
	MustGet(key string) (service any)
}

//-----------------------------------------------

type di_store map[string]any

func (s di_store) Set(key string, create CreateFunc) (e error) {
	_, ok := s[key]
	if ok {
		return ErrKeyAlreadyExists
	}
	s[key] = create()
	return
}

func (s di_store) Get(key string) (service any, ok bool) {
	service, ok = s[key]
	return
}

func (s di_store) MustGet(key string) (service any) {
	var ok bool
	service, ok = s[key]
	if !ok {
		panic(fmt.Errorf("service(%s) unset", key))
	}
	return
}

//-----------------------------------------------

var di = di_store{}

func GetDI() DI {
	return di
}
