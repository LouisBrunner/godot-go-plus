package wrapper

import (
	"github.com/godot-go/godot-go/pkg/builtin"
	"github.com/godot-go/godot-go/pkg/core"
	"github.com/godot-go/godot-go/pkg/ffi"
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
		for _, method := range me.info.methods {
			if method.isVirtual {
				core.ClassDBBindMethodVirtual(t, method.goName, method.gdName, method.gdArgs, nil)
			} else {
				core.ClassDBBindMethod(t, method.goName, method.gdName, method.gdArgs, nil)
			}
		}

		for _, prop := range me.info.properties {
			core.ClassDBBindMethod(t, prop.goGetter, prop.gdGetter, nil, nil)
			core.ClassDBBindMethod(t, prop.goSetter, prop.gdSetter, prop.gdSetterArgs, nil)
			core.ClassDBAddProperty(t, prop.gdTyp, prop.gdName, prop.goSetter, prop.goGetter)
		}

		for _, signal := range me.info.signals {
			core.ClassDBAddSignal(t, signal.gdName, signal.gdArgs...)
		}
	})
}

func (me *classRegister) Terminate() {
}
