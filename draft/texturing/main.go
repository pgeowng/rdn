package main

import (
	"fmt"
	"image"
	"runtime"

	"neilpa.me/go-stbi"
	_ "neilpa.me/go-stbi/jpeg"
	_ "neilpa.me/go-stbi/png"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/pgeowng/rende/draft/texturing/glm"
	"github.com/pgeowng/rende/draft/texturing/shader"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() (err error) {
	err = glfw.Init()
	if err != nil {
		return
	}

	// glfw.WindowHint(glfw.VersionMajor, 3)
	// glfw.WindowHint(glfw.VersionMinor, 3)
	// glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		return
	}

	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(framebufferSizeCallback)

	if err = gl.Init(); err != nil {
		return
	}

	sh := shader.New("./vertex.glsl", "./fragment.glsl")
	prog, err := sh.Compile()
	if err != nil {
		return
	}

	points := []float32{
		// positions // colors     // texture coords
		+.5, +.5, 0, 0.0, 1.0, 0.0, 1.0, 0.0,
		+.5, -.5, 0, 0.0, 1.0, 1.0, 1.0, 1.0,
		-.5, -.5, 0, 0.0, 0.0, 1.0, 0.0, 1.0,
		-.5, +.5, 0, 1.0, 0.0, 0.0, 0.0, 0.0,
	}
	const vertexSize = 8
	const f32size = 4
	var vbo, vao uint32
	{
		gl.GenVertexArrays(1, &vao)
		gl.GenBuffers(1, &vbo)
		gl.BindVertexArray(vao)

		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, f32size*len(points), gl.Ptr(points), gl.STATIC_DRAW)

		// position attribute
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, f32size*vertexSize, nil)
		gl.EnableVertexAttribArray(0)

		// color attribute
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, f32size*vertexSize, gl.PtrOffset(f32size*3))
		gl.EnableVertexAttribArray(1)

		// tex attribute
		gl.VertexAttribPointer(2, 2, gl.FLOAT, false, f32size*vertexSize, gl.PtrOffset(f32size*6))
		gl.EnableVertexAttribArray(2)

		defer gl.DeleteVertexArrays(1, &vao)
		defer gl.DeleteBuffers(1, &vbo)
	}

	indices := []uint32{
		0, 1, 3,
		1, 2, 3,
	}

	var ebo uint32
	{
		gl.GenBuffers(1, &ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
		defer gl.DeleteBuffers(1, &ebo)
	}

	var texture1 uint32
	{
		gl.GenTextures(1, &texture1)
		gl.BindTexture(gl.TEXTURE_2D, texture1)

		// repeat texture1
		gl.TexParameteri(texture1, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(texture1, gl.TEXTURE_WRAP_T, gl.REPEAT)

		// texture1 interpolation
		gl.TexParameteri(texture1, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(texture1, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		var filePath string = "./tex.png"

		var rgba *image.RGBA
		rgba, err = stbi.Load(filePath)
		if err != nil {
			return
		}
		// if rgba.Stride != rgba.Rect.Size().X*4 {
		// 	return fmt.Errorf("image: unsupported stride: %d", rgba.Stride)
		// }

		fmt.Println(rgba.Rect.Size())

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
		gl.GenerateMipmap(gl.TEXTURE_2D)
		defer gl.DeleteTextures(1, &texture1)
	}

	var texture2 uint32
	{
		gl.GenTextures(1, &texture2)
		gl.BindTexture(gl.TEXTURE_2D, texture2)

		// repeat texture2
		gl.TexParameteri(texture2, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(texture2, gl.TEXTURE_WRAP_T, gl.REPEAT)

		// texture2 interpolation
		gl.TexParameteri(texture2, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(texture2, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		var filePath string = "./lumi.jpg"

		var rgba *image.RGBA
		rgba, err = stbi.Load(filePath)
		if err != nil {
			return
		}
		// if rgba.Stride != rgba.Rect.Size().X*4 {
		// 	return fmt.Errorf("image: unsupported stride: %d", rgba.Stride)
		// }

		fmt.Println(rgba.Rect.Size())

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
		gl.GenerateMipmap(gl.TEXTURE_2D)
		defer gl.DeleteTextures(1, &texture2)
	}

	gl.UseProgram(prog)

	gl.Uniform1i(gl.GetUniformLocation(prog, gl.Str("texture1"+"\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(prog, gl.Str("texture2"+"\x00")), 1)

	// vertexColorLocation := gl.GetUniformLocation(shaderProgram, gl.Str("ourColor"+"\x00"))
	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)

		translation := glm.Identity().Translate(glm.Vec3{.5, -.5, 0})
		rotation := glm.RotationZ(float32(glfw.GetTime()))
		transform := translation.Times(rotation)

		gl.UseProgram(prog)
		gl.UniformMatrix4fv(gl.GetUniformLocation(prog, gl.Str("transform"+"\x00")), 1, false, &transform[0])
		gl.BindVertexArray(vao)

		// timeValue := glfw.GetTime()
		// greenValue := float32(math.Sin(timeValue)/2.0 + 0.5)
		// gl.Uniform4f(vertexColorLocation, 0, greenValue, 0, 1)

		// gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	return
}

func framebufferSizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func processInput(w *glfw.Window) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
}

func makeVao(points []float32) (vbo uint32, vao uint32) {
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	return
}
