package app

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

type App struct {
	*app.Application
	currentDemo IGame
	dirData     string
	scene       *core.Node
	demoScene   *core.Node
	ambLight    *light.Ambient

	camera *camera.Camera       // Camera
	orbit  *camera.OrbitControl // Orbit control
}

type IGame interface {
	Start(*App)                 // Called once at the start of the demo
	Update(*App, time.Duration) // Called every frame
	Cleanup(*App)               // Called once at the end of the demo
}

var GameMap = map[string]IGame{}

const (
	progName = "DotGo Alpha"
	execName = "dotgo"
	vmajor   = 0
	vminor   = 1
)

func Create() *App {

	a := new(App)
	a.Application = app.App()

	// Create scenes
	a.demoScene = core.NewNode() // demoScene will be cleared before a new demo is started
	a.scene = core.NewNode()
	a.scene.Add(a.demoScene)

	width, height := a.GetSize()
	aspect := float32(width) / float32(height)
	a.camera = camera.New(aspect)
	a.scene.Add(a.camera) // Add camera to scene (important for audio demos)
	a.orbit = camera.NewOrbitControl(a.camera)

	// Create and add ambient light to scene
	a.ambLight = light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.5)
	a.scene.Add(a.ambLight)

	a.Subscribe(window.OnWindowSize, func(evname string, ev interface{}) { a.OnWindowResize() })
	a.OnWindowResize()

	a.Subscribe(window.OnKeyDown, func(evname string, ev interface{}) {
		kev := ev.(*window.KeyEvent)
		if kev.Key == window.KeyEscape {
			a.Exit()
		}
	})

	a.setupScene()

	return a
}

func (a *App) setupScene() {

	if a.currentDemo != nil {
		a.currentDemo.Cleanup(a)
	}

	a.UnsubscribeAllID(a)

	a.DisposeAllCustomCursors()
	a.SetCursor(window.ArrowCursor)

	a.Gls().ClearColor(0, 0, 0, 1.0)

	a.Renderer().SetObjectSorting(true)

	a.ambLight.SetColor(&math32.Color{1.0, 1.0, 1.0})
	a.ambLight.SetIntensity(0.5)

	a.camera.SetPosition(0, 0, 5)
	a.camera.UpdateSize(5)
	a.camera.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 1, 0})
	a.camera.SetProjection(camera.Perspective)

	a.orbit.Reset()

}

func (a *App) AmbLight() *light.Ambient {

	return a.ambLight
}

func (a *App) Scene() *core.Node {

	return a.demoScene
}

func (a *App) Camera() *camera.Camera {

	return a.camera
}

func (a *App) Orbit() *camera.OrbitControl {

	return a.orbit
}

func (a *App) OnWindowResize() {

	width, height := a.GetFramebufferSize()
	a.Gls().Viewport(0, 0, int32(width), int32(height))

	a.camera.SetAspect(float32(width) / float32(height))
}

func (a *App) Run() {

	a.Application.Run(a.Update)
}

func (a *App) Update(rend *renderer.Renderer, deltaTime time.Duration) {

	a.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT)

	if a.currentDemo != nil {
		a.currentDemo.Update(a, deltaTime)
	}

	err := rend.Render(a.scene, a.camera)
	if err != nil {
		panic(err)
	}

	gui.Manager().TimerManager.ProcessTimers()

}
