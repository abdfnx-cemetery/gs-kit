package data

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	typeDuration      = reflect.TypeOf(time.Duration(5))
	typeTime          = reflect.TypeOf(time.Time{})
	typeStringDecoder = reflect.TypeOf((*StringDecoder)(nil)).Elem()
)

type StringDecoder interface {
	DecodeString(value string) error
}
