# go-gfx [![Go-Linux](https://github.com/taylorza/go-gfx/actions/workflows/linux-build.yml/badge.svg?branch=master)](https://github.com/taylorza/go-gfx/actions/workflows/linux-build.yml)[![Go Reference](https://pkg.go.dev/badge/github.com/taylorza/go-gfx.svg)](https://pkg.go.dev/github.com/taylorza/go-gfx)[![Go-Windows](https://github.com/taylorza/go-gfx/actions/workflows/windows-build.yml/badge.svg)](https://github.com/taylorza/go-gfx/actions/workflows/windows-build.yml)

Package **gfx** provides a simple native Go 2d graphics library for Windows and Linux.

## Installation

Use the 'go' command:

    $ go get github.com/taylorza/go-gfx
    
## Examples

```go
  package main
  
  import(
    "github.com/taylorza/go-gfx/pkg/gfx"
  )
  
  // Define a type that represents your application. This type must satisfy the Application interface
  // and provide the Load, Update and Unload methods. This struct can maintain any additional state you 
  // require for your graphics application
  type myapp struct {
  }
  
  // Load called once when the application is initialized
  func (app *myapp) Load() {}
  
  // Unload called when the application ends
  func (app *myapp) Unload() {}
  
  // Update called for every frame and provides the time in seconds since the last frame update
  func (app *myapp) Update(delta float64) {
    // Update and draw your frame
    gfx.Clear(gfx.Cyan)
    gfx.DrawCircle(gfx.Width() / 2, gfx.Height() / 2, gfx.Height() / 3, gfx.Red)
  }

  func main() {
    if gfx.Init("GFX Example", 0, 0, 320, 240, 2, 2) {
      gfx.Rin(&myapp{})
    }
  }
```
