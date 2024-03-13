package wrapper

import (
	"reflect"

	"github.com/godot-go/godot-go/pkg/builtin"
)

type Class = builtin.Object

type ClassConstructor func() Class

func (me *ClassConstructor) Name() string {
	return reflect.TypeOf(me).Elem().Name()
}
