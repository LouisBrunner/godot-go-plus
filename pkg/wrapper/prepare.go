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

const tagName = "godot"

type tagData struct {
	name   string
	getter string
	setter string
}

func parseTag(tag string) *tagData {
	if tag == "" {
		return nil
	}
	fields := strings.Split(tag, ",")
	data := &tagData{}
	for _, field := range fields {
		parts := strings.Split(field, "=")
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "name":
			data.name = parts[1]
		case "get":
			data.getter = parts[1]
		case "set":
			data.setter = parts[1]
		}
	}
	return data
}

var reservedMethods = []string{
	"Destroy",
}

func prepareClass(ctr ClassConstructor) (*classInfo, error) {
	instance := ctr()
	typ := reflect.TypeOf(instance)
	name := typ.Elem().Name()

	if typ.Kind() != reflect.Pointer && typ.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("class %q: expected a pointer to struct but got %s (%T)", name, typ.Kind(), instance)
	}

	properties := make([]classProperty, 0, typ.Elem().NumField())
	methods := make([]classMethod, 0, typ.NumMethod())
	signals := make([]classSignal, 0)

	ignoredMethods := map[string]struct{}{}
	for _, method := range reservedMethods {
		ignoredMethods[method] = struct{}{}
	}

	for _, field := range reflect.VisibleFields(typ.Elem()) {
		ftyp := field.Type
		if field.Anonymous {
			for m := 0; m < ftyp.NumMethod(); m++ {
				method := ftyp.Method(m)
				if method.IsExported() {
					ignoredMethods[method.Name] = struct{}{}
				}
			}
			ftyp = reflect.PointerTo(ftyp)
			for m := 0; m < ftyp.NumMethod(); m++ {
				method := ftyp.Method(m)
				if method.IsExported() {
					ignoredMethods[method.Name] = struct{}{}
				}
			}
		}
	}

	for f := 0; f < typ.Elem().NumField(); f++ {
		field := typ.Elem().Field(f)

		if !field.IsExported() || field.Anonymous {
			continue
		}

		fieldName := strcase.ToSnake(field.Name)
		info := parseTag(field.Tag.Get(tagName))
		if info != nil && info.name != "" {
			fieldName = info.name
		}

		fmt.Printf("field: %s\n", fieldName)
		fmt.Printf("field type: %s\n", field.Type)
		fmt.Printf("field type: %T\n", reflect.New(field.Type).Elem().Interface())
		isSignal := field.Type.Implements(reflect.TypeFor[Signal]())
		fmt.Printf("isSignal: %v\n", isSignal)
		if isSignal {
			signals = append(signals, classSignal{
				gdName: fieldName,
				// TODO: no idea how to get those
				gdArgs: []core.SignalParam{},
			})
			continue
		}

		getterName := fmt.Sprintf("_get_property_%s", fieldName)
		goGetter := fmt.Sprintf("Get%s", strcase.ToCamel(field.Name))
		if info != nil && info.getter != "" {
			goGetter = info.getter
		}
		method, found := typ.MethodByName(goGetter)
		if !found {
			return nil, fmt.Errorf("class %q: getter %q not found", name, goGetter)
		}
		if method.Type.NumIn() != 1 || method.Type.NumOut() != 1 {
			return nil, fmt.Errorf("class %q: getter %q has wrong signature", name, goGetter)
		}
		ignoredMethods[goGetter] = struct{}{}

		setterName := fmt.Sprintf("_set_property_%s", fieldName)
		goSetter := fmt.Sprintf("Set%s", strcase.ToCamel(field.Name))
		if info != nil && info.setter != "" {
			goSetter = info.setter
			if goSetter == "nil" {
				goSetter = ""
				setterName = ""
			}
		}
		if goSetter != "" {
			method, found = typ.MethodByName(goSetter)
			if !found {
				return nil, fmt.Errorf("class %q: setter %q not found", name, goSetter)
			}
			if method.Type.NumIn() != 2 || method.Type.NumOut() != 0 {
				return nil, fmt.Errorf("class %q: setter %q has wrong signature", name, goSetter)
			}
			ignoredMethods[goSetter] = struct{}{}
		}

		properties = append(properties, classProperty{
			gdName:       fieldName,
			gdGetter:     getterName,
			gdSetter:     setterName,
			goGetter:     goGetter,
			goSetter:     goSetter,
			gdSetterArgs: []string{"value"},
			gdTyp:        core.ReflectTypeToGDExtensionVariantType(field.Type),
		})
	}

	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		if !method.IsExported() {
			continue
		}
		if _, found := ignoredMethods[method.Name]; found {
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
			args = append(args, fmt.Sprintf("arg%d", i-1))
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
		properties: properties,
		methods:    methods,
		signals:    signals,
	}, nil
}
