package shader

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Shader struct {
	vertexPath   string
	fragmentPath string
	program      uint32
}

func New(vertexPath, fragmentPath string) *Shader {
	return &Shader{
		vertexPath:   vertexPath,
		fragmentPath: fragmentPath,
	}
}

func (s *Shader) Compile() (prog uint32, err error) {

	vertexBody, err := os.ReadFile(s.vertexPath)
	if err != nil {
		return
	}

	fragmentBody, err := os.ReadFile(s.fragmentPath)
	if err != nil {
		return
	}

	vertexShader, err := compileShader(string(vertexBody)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		return
	}

	fragmentShader, err := compileShader(string(fragmentBody)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		return
	}

	prog = gl.CreateProgram()
	s.program = prog
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	{
		var status int32
		gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
		if status == gl.FALSE {
			var logLen int32
			gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLen)

			log := strings.Repeat("\x00", int(logLen+1))
			gl.GetProgramInfoLog(prog, logLen, nil, gl.Str(log))

			err = fmt.Errorf("failed to link: %v", log)
			return
		}
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return
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

func (s *Shader) Uniform(label string) int32 {
	return gl.GetUniformLocation(s.program, gl.Str(label+"\x00"))
}

func (s *Shader) UseProgram() {
	gl.UseProgram(s.program)
}
