package entry

import "github.com/LouisBrunner/godot-go-plus/pkg/wrapper"

var globalRegisteredClasses []wrapper.ClassConstructor

func Register(constructors ...wrapper.ClassConstructor) {
	globalRegisteredClasses = append(globalRegisteredClasses, constructors...)
}
