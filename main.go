package main

import (
	"ConsolRayTracingGo/vec"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var width uint16 = 120
	var height uint16 = 30
	if ws, err := unix.IoctlGetWinsize(syscall.Stdout, unix.TIOCGWINSZ); err == nil {
		// get terminal size if possible
		width = ws.Col
		height = ws.Row
	}
	aspect := float64(width) / float64(height)
	pixelAspect := 11.0 / 24
	gradient := " .:!/r(l1Z4H9W8$@"
	gradientSize := len(gradient) - 1

	screen := newScreen(height, width)
	ts_start := float64(time.Now().UnixMilli())
	ts := ts_start
	tsOld := ts
	hideCursor()
	defer showCursor()
	go func() {
		<-sigs
		showCursor()
		fmt.Println("Exiting")
		os.Exit(0)
	}()
	for {
		// Main loop
		t := float64(time.Now().UnixMilli()) - ts_start
		light := vec.Vec3FromXYZ(-0.5, 0.5, -1.0).Norm()
		spheres := []vec.Vec3{
			vec.Vec3FromXYZ(0, 3, 0),
			vec.Vec3FromXYZ(3, 0, 0),
			vec.Vec3FromXYZ(0, -3, 0),
			vec.Vec3FromXYZ(-3, 0, 0),
		}
		for j := 0; j < int(height); j++ {
			for i := 0; i < int(width); i++ {
				uv := vec.Vec2FromXY(float64(i), float64(j)).Div(vec.Vec2FromXY(float64(width), float64(height))).Mul(vec.Vec2FromScalar(2)).Sub(vec.Vec2FromScalar(1))
				uv.X *= aspect * pixelAspect
				ro := vec.Vec3FromXYZ(-6, 0, 0)
				rd := vec.Vec3FromVec2(2, uv).Norm()
				ro = vec.RotateY(ro, 0.25)
				rd = vec.RotateY(rd, 0.25)
				ro = vec.RotateZ(ro, float64(t)*0.001)
				rd = vec.RotateZ(rd, float64(t)*0.001)
				diff := 1.0
				for k := 0; k < 5; k++ {
					minIt := 99999.0
					n := vec.Vec3FromScalar(0)
					albedo := 1.0
					for _, spherePos := range spheres {
						intersectSphere(ro, spherePos, rd, &minIt, &n)
					}
					intersectCube(ro, rd, &minIt, &n)
					intersectPlane(ro, rd, &minIt, &n, &albedo)
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
				screen[j].line[i] = pixel
			}
			movePrint(0, j, string(screen[j].line))
		}
		tsOld = ts
		ts = float64(time.Now().UnixMilli())
		movePrint(0, int(height)-1, fmt.Sprintf("FPS: %v ", 1000.0/(ts-tsOld)))
	}

}

func intersectPlane(ro vec.Vec3, rd vec.Vec3, minIt *float64, n *vec.Vec3, albedo *float64) {
	intersection := vec.Vec2FromScalar(vec.Plane(ro, rd, vec.Vec3FromXYZ(0, 0, -1), 1))
	if intersection.X > 0 && intersection.X < *minIt {
		*minIt = intersection.X
		*n = vec.Vec3FromXYZ(0, 0, -1)
		*albedo = 0.5
	}
}

func intersectCube(ro vec.Vec3, rd vec.Vec3, minIt *float64, n *vec.Vec3) {
	intersection, boxN := vec.Box(ro, rd, vec.Vec3FromScalar(1))
	if intersection.X > 0 && intersection.X < *minIt {
		*minIt = intersection.X
		*n = boxN
	}
}

func intersectSphere(ro vec.Vec3, spherePos vec.Vec3, rd vec.Vec3, minIt *float64, n *vec.Vec3) {
	intersection := vec.Sphere(ro.Sub(spherePos), rd, 1)
	if intersection.X > 0 && intersection.X < *minIt {
		itPoint := ro.Sub(spherePos).Add(rd.Mul(vec.Vec3FromScalar(intersection.X)))
		*minIt = intersection.X
		*n = itPoint.Norm()
	}
}

type row struct {
	line []byte
	n    uint16
}

func newRow(width uint16, n uint16) row {
	r := row{make([]byte, width), n}
	return r
}

func newScreen(height, width uint16) []row {
	screen := make([]row, height)
	for i := uint16(0); i < height; i++ {
		screen[i] = newRow(width, i)
	}
	return screen
}

func movePrint(x, y int, str string) {
	fmt.Printf("\033[%d;%dH%s", y, x, str)
}

func hideCursor() {
	fmt.Printf("\033[?25l")
}

func showCursor() {
	fmt.Printf("\033[?25h")
}
