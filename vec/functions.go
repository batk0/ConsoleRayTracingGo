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

type Object interface {
	// intersect(ro, rd Vec3) float64
	GetReflection(ro, rd Vec3, minIt *float64, normal *Vec3, albedo *float64)
}
type Sphere struct {
	Radius   float64
	Position Vec3
}
type Box struct {
	Size     Vec3
	Position Vec3
}
type Plane struct {
	Normal   Vec3
	Position Vec3
}

func (s Sphere) intersect(ro, rd Vec3) float64 {
	// Solving quadratic equation where a=1
	b := ro.Dot(rd)
	c := ro.Dot(ro) - s.Radius*s.Radius
	h := b*b - c
	if h < 0 {
		return -1
	}
	h = math.Sqrt(h)
	return -h - b // We don't use second root h-b
}

func (b Box) intersect(ro, rd Vec3) (float64, Vec3) {
	m := Vec3FromScalar(1).Div(rd)
	n := m.Mul(ro)
	k := m.Abs().Mul(b.Size)
	t1 := n.Neg().Sub(k)
	t2 := n.Neg().Add(k)
	tN := math.Max(math.Max(t1.X, t1.Y), t1.Z) // near point of intersection
	tF := math.Min(math.Min(t2.X, t2.Y), t2.Z) // far point of intersection
	if tN > tF || tF < 0 {
		return -1, Vec3FromScalar(0)
	}
	yzx := Vec3FromXYZ(t1.Y, t1.Z, t1.X)
	zxy := Vec3FromXYZ(t1.Z, t1.X, t1.Y)
	outNormal := rd.Sign().Neg().Mul(yzx.Step(t1)).Mul(zxy.Step(t1))
	return tN, outNormal // We don't use tF
}

func (p Plane) intersect(ro, rd Vec3) float64 {
	return ro.Neg().Dot(p.Normal) / rd.Dot(p.Normal)
}

func (p Plane) GetReflection(ro Vec3, rd Vec3, minIt *float64, normal *Vec3, albedo *float64) {
	intersection := p.intersect(ro.Sub(p.Position), rd)
	if intersection > 0 && intersection < *minIt {
		*minIt = intersection
		*normal = Vec3FromXYZ(0, 0, -1)
		*albedo = 0.5
	}
}

func (b Box) GetReflection(ro Vec3, rd Vec3, minIt *float64, normal *Vec3, albedo *float64) {
	intersection, boxN := b.intersect(ro.Sub(b.Position), rd)
	if intersection > 0 && intersection < *minIt {
		*minIt = intersection
		*normal = boxN
	}
}

func (s Sphere) GetReflection(ro Vec3, rd Vec3, minIt *float64, normal *Vec3, albedo *float64) {
	intersection := s.intersect(ro.Sub(s.Position), rd)
	if intersection > 0 && intersection < *minIt {
		itPoint := ro.Sub(s.Position).Add(rd.Mul(Vec3FromScalar(intersection)))
		*minIt = intersection
		*normal = itPoint.Norm()
	}
}
