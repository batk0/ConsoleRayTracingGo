package vec

import "math"

type Vec2 struct {
	X, Y float64
}

func Vec2FromScalar(value float64) Vec2 { return Vec2{value, value} }
func Vec2FromXY(x, y float64) Vec2      { return Vec2{x, y} }

func (v1 Vec2) Add(v2 Vec2) Vec2 { return Vec2{v1.X + v2.X, v1.Y + v2.Y} }
func (v1 Vec2) Sub(v2 Vec2) Vec2 { return Vec2{v1.X - v2.X, v1.Y - v2.Y} }
func (v1 Vec2) Mul(v2 Vec2) Vec2 { return Vec2{v1.X * v2.X, v1.Y * v2.Y} }
func (v1 Vec2) Div(v2 Vec2) Vec2 { return Vec2{v1.X / v2.X, v1.Y / v2.Y} }

func (v Vec2) Length() float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y) }
