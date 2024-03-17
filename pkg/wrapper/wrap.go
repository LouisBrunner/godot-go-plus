package wrapper

import (
	"fmt"

	"github.com/godot-go/godot-go/pkg/builtin"
	"github.com/godot-go/godot-go/pkg/core"
	"github.com/godot-go/godot-go/pkg/ffi"
	"github.com/godot-go/godot-go/pkg/log"
)

type ClassWrapper interface {
	Initialize()
	Terminate()
}

func Wrap(ctr ClassConstructor) (ClassWrapper, error) {
	info, err := prepareClass(ctr)
	if err != nil {
		return nil, err
	}

	return &classRegister{
		ctr:  ctr,
		info: info,
	}, nil
}

type classRegister struct {
	ctr  ClassConstructor
	info *classInfo
}

func (me *classRegister) Initialize() {
	core.ClassDBRegisterClass(me.info.instance, []ffi.GDExtensionPropertyInfo{}, nil, func(t builtin.GDClass) {
		log.Debug(fmt.Sprintf("Registering class %s", me.info.name))

		for _, method := range me.info.methods {
			log.Debug(fmt.Sprintf("Registering method %s::%s", me.info.name, method.goName))
			if method.isVirtual {
				core.ClassDBBindMethodVirtual(t, method.goName, method.gdName, method.gdArgs, nil)
			} else {
				core.ClassDBBindMethod(t, method.goName, method.gdName, method.gdArgs, nil)
			}
		}

		for _, prop := range me.info.properties {
			log.Debug(fmt.Sprintf("Registering property %s::%s", me.info.name, prop.gdName))
			core.ClassDBBindMethod(t, prop.goGetter, prop.gdGetter, nil, nil)
			if prop.goSetter != "" {
				core.ClassDBBindMethod(t, prop.goSetter, prop.gdSetter, prop.gdSetterArgs, nil)
			}
			core.ClassDBAddProperty(t, prop.gdTyp, prop.gdName, prop.gdSetter, prop.gdGetter)
		}

		for _, signal := range me.info.signals {
			log.Debug(fmt.Sprintf("Registering signal %s::%s", me.info.name, signal.gdName))
			core.ClassDBAddSignal(t, signal.gdName, signal.gdArgs...)
		}
	})
}

func (me *classRegister) Terminate() {
}
