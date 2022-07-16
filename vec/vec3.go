package vec

import "math"

type Vec3 struct {
	X, Y, Z float64
}

func Vec3FromXYZ(x, y, z float64) Vec3    { return Vec3{x, y, z} }
func Vec3FromScalar(value float64) Vec3   { return Vec3{value, value, value} }
func Vec3FromVec2(x float64, v Vec2) Vec3 { return Vec3{x, v.X, v.Y} }

func (v1 Vec3) Add(v2 Vec3) Vec3 { return Vec3{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z} }
func (v1 Vec3) Sub(v2 Vec3) Vec3 { return Vec3{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z} }
func (v1 Vec3) Mul(v2 Vec3) Vec3 { return Vec3{v1.X * v2.X, v1.Y * v2.Y, v1.Z * v2.Z} }
func (v1 Vec3) Div(v2 Vec3) Vec3 { return Vec3{v1.X / v2.X, v1.Y / v2.Y, v1.Z / v2.Z} }
func (v Vec3) Neg() Vec3         { return Vec3{-v.X, -v.Y, -v.Z} }

func (v Vec3) Length() float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z) }

func (v Vec3) Norm() Vec3           { return v.Div(Vec3FromScalar(v.Length())) }
func (v1 Vec3) Dot(v2 Vec3) float64 { return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z }
func (v Vec3) Abs() Vec3            { return Vec3{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)} }
func (v Vec3) Sign() Vec3           { return Vec3{Sign(v.X), Sign(v.Y), Sign(v.Z)} }
func (edge Vec3) Step(v Vec3) Vec3 {
	return Vec3{Step(edge.X, v.X), Step(edge.Y, v.Y), Step(edge.Z, v.Z)}
}
