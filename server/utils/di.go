package utils

import (
	"fmt"
)

var ErrKeyAlreadyExists = fmt.Errorf("map: key already exists")

type CreateFunc func(args ...interface{}) any

type di_config struct {
	is_only     bool
	instance    interface{}
	create_func CreateFunc
}

type DI interface {
	Set(key string, create_func CreateFunc) (e error)
	SetShare(key string, create_func CreateFunc) (e error)
	MustGet(key string, args ...interface{}) any
}

//-----------------------------------------------

type di_store map[string]di_config

func (s di_store) Set(key string, create_func CreateFunc) (e error) {
	_, ok := s[key]
	if ok {
		return ErrKeyAlreadyExists
	}
	s[key] = di_config{
		is_only:     false,
		create_func: create_func,
	}
	return
}

func (s di_store) SetShare(key string, create_func CreateFunc) (e error) {
	_, ok := s[key]
	if ok {
		return ErrKeyAlreadyExists
	}
	s[key] = di_config{
		is_only:     true,
		create_func: create_func,
	}
	return
}

func (s di_store) MustGet(key string, args ...interface{}) any {
	config, ok := s[key]
	if !ok {
		panic(fmt.Errorf("service(%s) unset", key))
	}

	if config.is_only {
		if config.instance == nil {
			config.instance = config.create_func(args...)
		}
		return config.instance
	}

	return config.create_func(args...)
}

//-----------------------------------------------

var di = &di_store{}

func GetDI() DI {
	return di
}
