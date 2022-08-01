package main

import (
	"ConsolRayTracingGo/vec"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const gradient = " .:!/r(l1Z4H9W8$@"
const gradientSize = len(gradient) - 1

func main() {
	// go routines setup
	var workers uint16 = 16
	done := make(chan struct{})
	defer close(done)
	fOut := make([]<-chan row, workers)

	// Signals handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	hideCursor()
	defer showCursor()
	go func() {
		<-sigs
		close(done)
		showCursor()
		fmt.Println("Exiting")
		os.Exit(0)
	}()

	// Screen setup
	var width uint16 = 120
	var height uint16 = 30
	if ws, err := unix.IoctlGetWinsize(syscall.Stdout, unix.TIOCGWINSZ); err == nil {
		// get terminal size if possible
		width = ws.Col
		height = ws.Row
	}
	aspect := float64(width) / float64(height)
	pixelAspect := 11.0 / 24

	// Setup common params for rows
	light := vec.Vec3FromXYZ(-0.5, 0.5, -1.0).Norm()
	objects := []vec.Object{
		vec.Box{Size: vec.Vec3FromScalar(1), Position: vec.Vec3FromXYZ(0, 3, 0)},
		// vec.Sphere{Radius: 1, Position: vec.Vec3FromXYZ(0, 3, 0)},
		vec.Sphere{Radius: 1, Position: vec.Vec3FromXYZ(3, 0, 0)},
		vec.Sphere{Radius: 1, Position: vec.Vec3FromXYZ(0, -3, 0)},
		vec.Sphere{Radius: 1, Position: vec.Vec3FromXYZ(-3, 0, 0)},
		vec.Box{Size: vec.Vec3FromScalar(1), Position: vec.Vec3FromXYZ(0, 0, -1)},
		vec.Plane{Normal: vec.Vec3FromXYZ(0, 0, 1), Position: vec.Vec3FromXYZ(0, 0, 2)},
	}
	params := rowParams{width, height, 0, aspect, pixelAspect, 0, objects, light}

	// Start timer
	ts_start := float64(time.Now().UnixMilli())
	ts := ts_start
	tsOld := ts

	for {
		// Main loop
		t := float64(time.Now().UnixMilli()) - ts_start
		params.t = t

		// Start workers
		for w := uint16(0); w < workers; w++ {
			// fmt.Println(w)
			fOut[w] = fanOut(done, genRowParams(done, params, height, workers, w))
		}
		// Print results
		for r := range fanIn(done, fOut...) {
			movePrint(0, int(r.n), string(r.line))
		}

		// Reset timer and print FPS
		tsOld = ts
		ts = float64(time.Now().UnixMilli())
		movePrint(0, int(height)-1, fmt.Sprintf("FPS: %v ", 1000.0/(ts-tsOld)))
	}

}

type rowParams struct {
	width, height, j       uint16
	aspect, pixelAspect, t float64
	objects                []vec.Object
	light                  vec.Vec3
}

type row struct {
	line []byte
	n    uint16
}

func newRow(width uint16, n uint16) row {
	r := row{make([]byte, width), n}
	return r
}

// Generates rowParams and stream them to the chan
func genRowParams(done chan struct{}, params rowParams, height, workers, w uint16) <-chan rowParams {
	r := make(chan rowParams)
	go func() {
		defer close(r)
		for j := w; j < height; j += workers {
			params.j = j
			select {
			case <-done:
				return
			case r <- params:
			}
		}

	}()
	return r
}

func fanOut(done chan struct{}, in <-chan rowParams) <-chan row {
	resultRow := make(chan row)
	go func() {
		defer close(resultRow)
		for rp := range in {
			select {
			case <-done:
				return
			case resultRow <- computeRow(rp):
			}
		}
	}()
	return resultRow
}

func fanIn(done chan struct{}, cs ...<-chan row) <-chan row {
	var wg sync.WaitGroup
	resultRow := make(chan row)

	mx := func(c <-chan row) {
		defer wg.Done()
		for r := range c {
			select {
			case <-done:
				return
			case resultRow <- r:
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go mx(c)
	}

	go func() {
		wg.Wait()
		close(resultRow)
	}()
	return resultRow
}

func computeRow(params rowParams) row {
	screenRow := newRow(params.width, params.j)
	for i := 0; i < int(params.width); i++ {
		uv := vec.Vec2FromXY(float64(i), float64(params.j)).Div(vec.Vec2FromXY(float64(params.width), float64(params.height))).Mul(vec.Vec2FromScalar(2)).Sub(vec.Vec2FromScalar(1))
		uv.X *= params.aspect * params.pixelAspect
		ro := vec.Vec3FromXYZ(-10, 0, 0)
		rd := vec.Vec3FromVec2(2, uv).Norm()
		ro = vec.RotateY(ro, 0.25)
		rd = vec.RotateY(rd, 0.25)
		ro = vec.RotateZ(ro, float64(params.t)*0.001)
		rd = vec.RotateZ(rd, float64(params.t)*0.001)
		diff := 1.0
		for k := 0; k < 5; k++ {
			minIt := 99999.0
			n := vec.Vec3FromScalar(0)
			albedo := 1.0
			for _, o := range params.objects {
				o.GetReflection(ro, rd, &minIt, &n, &albedo)
			}
			if minIt < 99999 {
				diff *= (n.Dot(params.light)*0.5 + 0.5) * albedo
				ro = ro.Add(rd.Mul(vec.Vec3FromScalar(minIt - 0.01)))
				rd = vec.Reflect(rd, n)
			} else {
				break
			}
		}
		color := int(diff * 20)
		color = int(vec.Clamp(float64(color), 0, float64(gradientSize)))
		pixel := gradient[color]
		screenRow.line[i] = pixel
	}
	return screenRow
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
