package vec

import "math"

func Clamp(value, min, max float64) float64 { return math.Max(math.Min(value, max), min) }
func Sign(value float64) float64 {
	if value < 0 {
		return -1
	}
	if value > 0 {
		return 1
	}
	return 0
}

func Step(edge, x float64) float64 {
	if x > edge {
		return 1
	}
	return 0
}

func Reflect(rd, n Vec3) Vec3 { return rd.Sub(n.Mul(Vec3FromScalar(2 * n.Dot(rd)))) }
func RotateX(v Vec3, angle float64) Vec3 {
	r := v
	r.Z = v.Z*math.Cos(angle) - v.Y*math.Sin(angle)
	r.Y = v.Z*math.Sin(angle) + v.Y*math.Cos(angle)
	return r
}
func RotateY(v Vec3, angle float64) Vec3 {
	r := v
	r.X = v.X*math.Cos(angle) - v.Z*math.Sin(angle)
	r.Z = v.X*math.Sin(angle) + v.Z*math.Cos(angle)
	return r
}
func RotateZ(v Vec3, angle float64) Vec3 {
	r := v
	r.X = v.X*math.Cos(angle) - v.Y*math.Sin(angle)
	r.Y = v.X*math.Sin(angle) + v.Y*math.Cos(angle)
	return r
}

func Sphere(ro, rd Vec3, r float64) Vec2 {
	b := ro.Dot(rd)
	c := ro.Dot(ro) - r*r
	h := b*b - c
	if h < 0 {
		return Vec2FromScalar(-1)
	}
	h = math.Sqrt(h)
	return Vec2FromXY(-h-b, h-b)
}
func Box(ro, rd, boxSize Vec3) (Vec2, Vec3) {
	m := Vec3FromScalar(1).Div(rd)
	n := m.Mul(ro)
	k := m.Abs().Mul(boxSize)
	t1 := n.Neg().Sub(k)
	t2 := n.Neg().Add(k)
	tN := math.Max(math.Max(t1.X, t1.Y), t1.Z)
	tF := math.Min(math.Min(t2.X, t2.Y), t2.Z)
	if tN > tF || tF < 0 {
		return Vec2FromScalar(-1), Vec3FromScalar(0)
	}
	yzx := Vec3FromXYZ(t1.Y, t1.Z, t1.X)
	zxy := Vec3FromXYZ(t1.Z, t1.X, t1.Y)
	outNormal := rd.Sign().Neg().Mul(yzx.Step(t1)).Mul(zxy.Step(t1))
	return Vec2FromXY(tN, tF), outNormal
}

func Plane(ro, rd, p Vec3, w float64) float64 {
	return -(ro.Dot(p) + w) / rd.Dot(p)
}
