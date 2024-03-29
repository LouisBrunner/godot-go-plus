package main

import (
	"fmt"

	ggp "github.com/LouisBrunner/godot-go-plus"
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

type MyNode2D struct {
	Node2DImpl

	Speed  int
	secret string
}

func (n *MyNode2D) SetSpeed(speed int) {
	n.Speed = speed
}

func (n *MyNode2D) GetSpeed() int {
	return n.Speed
}

func (n *MyNode2D) Move(vec Vector2) {
	n.Node2DImpl.SetPosition(vec.Multiply_int(int64(n.Speed)))
}

func (n *MyNode2D) X_Ready() {
	n.printSecret()
	n.Speed *= 10
}

func (n *MyNode2D) printSecret() {
	fmt.Println(n.secret)
}

func newMyNode2D() ggp.Class {
	return &MyNode2D{
		secret: "123",
		Speed:  10,
	}
}

func init() {
	ggp.Register(newMyNode2D)
}

func main() {}
