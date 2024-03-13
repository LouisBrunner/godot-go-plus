package entry

import "C"

import (
	"fmt"
	"unsafe"

	"github.com/LouisBrunner/godot-go-plus/pkg/wrapper"
	"github.com/godot-go/godot-go/pkg/core"
	"github.com/godot-go/godot-go/pkg/ffi"
	"github.com/godot-go/godot-go/pkg/log"
)

//export godot_go_plus_entry
func godot_go_plus_entry(p_get_proc_address unsafe.Pointer, p_library unsafe.Pointer, r_initialization unsafe.Pointer) bool {
	initObj := core.NewInitObject(
		(ffi.GDExtensionInterfaceGetProcAddress)(p_get_proc_address),
		(ffi.GDExtensionClassLibraryPtr)(p_library),
		(*ffi.GDExtensionInitialization)(unsafe.Pointer(r_initialization)),
	)

	registeredClasses := make([]wrapper.ClassWrapper, 0, len(globalRegisteredClasses))
	for _, constructor := range globalRegisteredClasses {
		wrapper, err := wrapper.Wrap(constructor)
		if err != nil {
			log.Error(fmt.Sprintf("Error wrapping class %q: %s", constructor.Name(), err))
			continue
		}
		registeredClasses = append(registeredClasses, wrapper)
	}

	initObj.RegisterSceneInitializer(func() {
		for _, class := range registeredClasses {
			class.Initialize()
		}
	})

	initObj.RegisterSceneTerminator(func() {
		for _, class := range registeredClasses {
			class.Terminate()
		}
	})

	return initObj.Init()
}
