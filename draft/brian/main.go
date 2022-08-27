package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	{
		var status int32
		gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
		if status == gl.FALSE {
			var logLen int32
			gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLen)

			log := strings.Repeat("\x00", int(logLen+1))
			gl.GetProgramInfoLog(shaderProgram, logLen, nil, gl.Str(log))

			return fmt.Errorf("failed to link: %v", log)
		}
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	points := []float32{
		// positions // colors
		+.5, +.5, 0, 0.0, 1.0, 0.0,
		+.5, -.5, 0, 0.0, 1.0, 1.0,
		-.5, -.5, 0, 0.0, 0.0, 1.0,
		-.5, +.5, 0, 1.0, 0.0, 0.0,
	}
	const vertexSize = 6
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

	// vertexColorLocation := gl.GetUniformLocation(shaderProgram, gl.Str("ourColor"+"\x00"))
	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shaderProgram)
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

func compileShader(source string, shaderType uint32) (shader uint32, err error) {
	shader = gl.CreateShader(shaderType)

	cstrs, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, cstrs, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(shader, logLen, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return
}

var vertexShaderSource = `
	#version 330 core

	layout (location = 0) in vec3 aPos;
	layout (location = 1) in vec3 aColor;
	out vec3 ourColor;

	void main()
	{
		gl_Position = vec4(aPos, 1.0);
		ourColor = aColor;
	}
` + "\x00"

var fragmentShaderSource = `
	#version 330

	in vec3 ourColor;
	out vec4 FragColor;
	//uniform vec4 ourColor;

	void main()
	{
		FragColor = vec4(ourColor, 1.0f);
	}
` + "\x00"

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
