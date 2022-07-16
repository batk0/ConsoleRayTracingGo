package main

import (
	"ConsolRayTracingGo/vec"
	"fmt"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	width := 120.0
	height := 30.0
	if ws, err := unix.IoctlGetWinsize(syscall.Stdout, unix.TIOCGWINSZ); err == nil {
		// get terminal size if possible
		width = float64(ws.Col)
		height = float64(ws.Row)
	}
	aspect := float64(width) / float64(height)
	pixelAspect := 11.0 / 24
	gradient := " .:!/r(l1Z4H9W8$@"
	gradientSize := len(gradient) - 1

	screen := make([]byte, int(width*height+height))
	ts := time.Now().UnixMilli()
	tsOld := ts
	for t := 0; t < 10000; t++ {
		// Main loop
		light := vec.Vec3FromXYZ(-0.5, 0.5, -1.0).Norm()
		spherePos := vec.Vec3FromXYZ(0, 3, 0)
		for i := 0; i < int(width); i++ {
			for j := 0; j < int(height); j++ {
				uv := vec.Vec2FromXY(float64(i), float64(j)).Div(vec.Vec2FromXY(width, height)).Mul(vec.Vec2FromScalar(2)).Sub(vec.Vec2FromScalar(1))
				uv.X *= aspect * pixelAspect
				ro := vec.Vec3FromXYZ(-6, 0, 0)
				rd := vec.Vec3FromVec2(2, uv).Norm()
				ro = vec.RotateY(ro, 0.25)
				rd = vec.RotateY(rd, 0.25)
				ro = vec.RotateZ(ro, float64(t)*0.01)
				rd = vec.RotateZ(rd, float64(t)*0.01)
				diff := 1.0
				for k := 0; k < 5; k++ {
					minIt := 99999.0
					intersection := vec.Sphere(ro.Sub(spherePos), rd, 1)
					n := vec.Vec3FromScalar(0)
					albedo := 1.0
					if intersection.X > 0 {
						itPoint := ro.Sub(spherePos).Add(rd.Mul(vec.Vec3FromScalar(intersection.X)))
						minIt = intersection.X
						n = itPoint.Norm()
					}
					// boxN := vec.Vec3FromScalar(0)
					intersection, boxN := vec.Box(ro, rd, vec.Vec3FromScalar(1))
					if intersection.X > 0 && intersection.X < minIt {
						minIt = intersection.X
						n = boxN
					}
					intersection = vec.Vec2FromScalar(vec.Plane(ro, rd, vec.Vec3FromXYZ(0, 0, -1), 1))
					if intersection.X > 0 && intersection.X < minIt {
						minIt = intersection.X
						n = vec.Vec3FromXYZ(0, 0, -1)
						albedo = 0.5
					}
					if minIt < 99999 {
						diff *= (n.Dot(light)*0.5 + 0.5) * albedo
						ro = ro.Add(rd.Mul(vec.Vec3FromScalar(minIt - 0.01)))
						rd = vec.Reflect(rd, n)
					} else {
						break
					}
				}
				color := int(diff * 20)
				color = int(vec.Clamp(float64(color), 0, float64(gradientSize)))
				pixel := gradient[color]
				screen[i+j*int(width)] = pixel
				screen[(j+1)*int(width)] = '\n'
			}
		}
		// screen[width * height +height- 1] = '\0';
		// fmt.Printf("\e[2J\e[0;0H");
		fmt.Printf("%s", string(screen))
		tsOld = ts
		ts = time.Now().UnixMilli()
		fmt.Printf("FPS: %v", 1000.0/(ts-tsOld))
	}

}
