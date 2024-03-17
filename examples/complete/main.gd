extends Node2D

var time: float = 0

@export var factor: float = 100

@onready var my: MyNode2D = $MyNode2D

func _ready() -> void:
  my.monitor($ColorRect)
  
func _process(delta: float) -> void:
  time += delta
  my.move(Vector2(sin(time) * factor, cos(time) * factor) + Vector2(500, 300))
