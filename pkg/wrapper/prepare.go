package wrapper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/godot-go/godot-go/pkg/core"
	"github.com/godot-go/godot-go/pkg/ffi"
	"github.com/iancoleman/strcase"
)

type classProperty struct {
	gdName       string
	goGetter     string
	gdGetter     string
	goSetter     string
	gdSetter     string
	gdTyp        ffi.GDExtensionVariantType
	gdSetterArgs []string
}

type classMethod struct {
	gdName    string
	goName    string
	gdArgs    []string
	isVirtual bool
}

type classSignal struct {
	gdName string
	gdArgs []core.SignalParam
}

type classInfo struct {
	name       string
	instance   Class
	properties []classProperty
	methods    []classMethod
	signals    []classSignal
}

func prepareClass(ctr ClassConstructor) (*classInfo, error) {
	instance := ctr()
	// val := reflect.ValueOf(instance)
	typ := reflect.TypeOf(instance)
	if typ.Kind() != reflect.Pointer && typ.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a pointer to struct but got %s (%T)", typ.Kind(), instance)
	}

	name := typ.Elem().Name()
	properties := make([]classProperty, 0, typ.Elem().NumField())
	methods := make([]classMethod, 0, typ.NumMethod())

	methodsFromParent := map[string]struct{}{}
	for _, field := range reflect.VisibleFields(typ.Elem()) {
		ftyp := field.Type
		if field.Anonymous {
			for m := 0; m < ftyp.NumMethod(); m++ {
				method := ftyp.Method(m)
				if method.IsExported() {
					methodsFromParent[method.Name] = struct{}{}
				}
			}
			ftyp = reflect.PointerTo(ftyp)
			for m := 0; m < ftyp.NumMethod(); m++ {
				method := ftyp.Method(m)
				if method.IsExported() {
					methodsFromParent[method.Name] = struct{}{}
				}
			}
		}

		if !field.IsExported() || field.Anonymous {
			continue
		}

		// TODO: allow override through tags
		// fieldName := strcase.ToSnake(field.Name)
		// getterName := fmt.Sprintf("_get_property_%s", fieldName)
		// setterName := fmt.Sprintf("_set_property_%s", fieldName)

		// // prop, err := infoFromType(fieldName, val.Elem().FieldByName(field.Name))
		// // if err != nil {
		// // 	return nil, fmt.Errorf("class %q: %w", name, err)
		// // }
		// // propCpy := gdc.CNewPropertyInfo()
		// // *propCpy = *prop
		// methods = append(methods, classMethod{
		// 	name: getterName,
		// 	// method: gdc.ClassMethodInfo{
		// 	// 	Name:            getter.AsPtr(),
		// 	// 	MethodUserdata:  *(*unsafe.Pointer)(unsafe.Pointer(utils.ToPointer[int](len(methods)))),
		// 	// 	PtrcallFunc:     gdc.Callbacks.GetClassMethodInfoPtrcallFuncCallback(),
		// 	// 	CallFunc:        gdc.Callbacks.GetClassMethodInfoCallFuncCallback(),
		// 	// 	HasReturnValue:  gdc.Bool(1),
		// 	// 	ReturnValueInfo: propCpy,
		// 	// 	MethodFlags:     uint(gdapi.MethodFlagsDefault | gdapi.MethodFlagConst),
		// 	// 	// FIXME: metadata missing
		// 	// },
		// 	fn: func(me interface{}) interface{} {
		// 		return reflect.ValueOf(me).Elem().FieldByName(field.Name).Interface()
		// 	},
		// })
		// methods = append(methods, classMethod{
		// 	name: setterName,
		// 	// method: gdc.ClassMethodInfo{
		// 	// 	Name:           setter.AsPtr(),
		// 	// 	MethodUserdata: *(*unsafe.Pointer)(unsafe.Pointer(utils.ToPointer[int](len(methods)))),
		// 	// 	PtrcallFunc:    gdc.Callbacks.GetClassMethodInfoPtrcallFuncCallback(),
		// 	// 	CallFunc:       gdc.Callbacks.GetClassMethodInfoCallFuncCallback(),
		// 	// 	MethodFlags:    uint(gdapi.MethodFlagsDefault),
		// 	// 	ArgumentCount:  1,
		// 	// 	ArgumentsInfo:  propCpy,
		// 	// 	// FIXME: metadata missing
		// 	// },
		// 	fn: func(me interface{}, value interface{}) {
		// 		reflect.ValueOf(me).Elem().FieldByName(field.Name).Set(reflect.ValueOf(value))
		// 	},
		// })
		// properties = append(properties, classProperty{
		// 	name: field.Name,
		// 	// property: *prop,
		// 	// getter:   getter,
		// 	// setter:   setter,
		// })
	}

	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		if !method.IsExported() {
			continue
		}
		if _, found := methodsFromParent[method.Name]; found {
			continue
		}
		methodName := strcase.ToSnake(method.Name)
		isVirtual := false
		if strings.HasPrefix(methodName, "x_") {
			methodName = strings.TrimPrefix(methodName, "x")
			isVirtual = true
		}
		args := make([]string, 0, method.Type.NumIn()-1)
		for i := 1; i < method.Type.NumIn(); i += 1 {
			args = append(args, method.Type.In(i).Name())
		}
		methods = append(methods, classMethod{
			goName:    method.Name,
			gdName:    methodName,
			isVirtual: isVirtual,
			gdArgs:    args,
		})
	}

	return &classInfo{
		name:       name,
		instance:   instance,
		properties: properties, // TODO: properties
		methods:    methods,
		signals:    []classSignal{}, // TODO: signals
	}, nil
}
