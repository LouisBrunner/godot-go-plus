# godot-go-plus

**WARNING: this project is still in early development and not currently functional.**

This project allows you to easily wrap Go code into Godot classes and use them in your Godot project.

The API is heavily inspired by https://github.com/ShadowApex/godot-go and uses https://github.com/godot-go/godot-go under the hood.

## Usage

In order to use Go in your Godot project, you need to compile the Go code into a shared library and load it in Godot.
This is done by using the Go `-buildmode=c-shared` flag on a `main` package with a `main` function.

### Boilerplate code

Create a go file containing the following:

```go
package main

func init() {
}

func main() {
}
```

We will be using the `init` function to register our custom Godot classes, implemented in Go.

## Create a custom Godot class

A basic Godot class is defined as a struct with an embbeded `gdapi` struct. For example, to inherit from `Node2D`:

```go
import (
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

type MyNode2D struct {
  Node2DImpl
}
```

Any public method receiver defined on `MyNode2D` will be available as a method on the Godot class and any public field will be available as a property.

```go
import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

type MyNode2D struct {
  Node2DImpl

  Speed int
}

func (n *MyNode2D) Move(vec builtin.Vector2) {
  n.Node2DImpl.SetPosition(vec.Multiply_int(int64(n.Speed)))
}
```

## Register the custom Godot class

In the `init` function, register the custom Godot class:

```go
package main

import (
  ggp "github.com/LouisBrunner/godot-go-plus"
)

func NewMyNode2D() ggp.Class {
  return &MyNode2D{}
}

func init() {
  ggp.Register(NewMyNode2D)
}

func main() {
}
```

### Compiling the shared library

Your package is now ready to be built as a shared library.

```bash
# Linux
go build -buildmode=c-shared -o libmyextension.so pkg_folder
# macOS
go build -buildmode=c-shared -o libmyextension.dylib pkg_folder
# Windows
go build -buildmode=c-shared -o libmyextension.dll pkg_folder
```

You can name the shared library whatever you want, the example above uses `libmyextension`.

### Including the shared library in Godot

You will need to create a file with the extension `.gdextension` and place it in your project for Godot to be able to load the shared library.

```gdscript
[configuration]
entry_symbol = "godot_go_plus_entry"
compatibility_minimum = "4.2"

[libraries]
macos.debug.arm64 = "res://your_folder/libmyextension-darwin-arm64.dylib"
macos.release.arm64 = "res://your_folder/libmyextension-darwin-arm64.dylib"
macos.debug.amd64 = "res://your_folder/libmyextension-darwin-amd64.dylib"
macos.release.amd64 = "res://your_folder/libmyextension-darwin-amd64.dylib"
windows.debug.amd64 = "res://your_folder/libmyextension-windows-amd64.dll"
windows.release.amd64 = "res://your_folder/libmyextension-windows-amd64.dll"
linux.debug.amd64 = "res://your_folder/libmyextension-linux-amd64.so"
linux.release.amd64 = "res://your_folder/libmyextension-linux-amd64.so"
```

You will need to replace `your_folder` with the path to the folder containing the shared library and `libmyextension` with the name of the shared library.

### Caveats

- Due to the way `godot-go` works, you will need to set the environment variable `GODEBUG=cgocheck=0` anytime you use the library as Go will panic otherwise.

### Examples

You can find examples in the `examples` folder.

- `simple`: a simple example of a custom Godot class, which basically reproduces the above usage
- `complete`: a more complete example showcasing custom getter/setter names, signals, GDScript usage, etc

## Acknowledgements

- https://github.com/ShadowApex/godot-go: for their great work and straight-forward API which inspired this project
- https://github.com/godot-go/godot-go: for creating Go bindings for Godot without which this project would not be possible
