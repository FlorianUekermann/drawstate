package drawstate

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/all-core/gl"
)

func compileProgram(vertexShader, fragmentShader string) uint32 {
	var vs = makeShader(vertexShader, gl.VERTEX_SHADER)
	var fs = makeShader(fragmentShader, gl.FRAGMENT_SHADER)
	p := gl.CreateProgram()
	gl.AttachShader(p, vs)
	gl.AttachShader(p, fs)
	infoLogBuffer := make([]uint8, 10000)
	var infoLogLength int32
	gl.GetProgramInfoLog(p, int32(len(infoLogBuffer)), &infoLogLength, &infoLogBuffer[0])
	if infoLogLength != 0 {
		log.Fatalln("Traditional pipeline:\n" + string(infoLogBuffer[0:infoLogLength]))
	}
	gl.LinkProgram(p)
	gl.DetachShader(p, fs)
	gl.DeleteShader(fs)
	gl.DetachShader(p, vs)
	gl.DeleteShader(vs)
	glCheckError()
	return p
}

func makeShader(shaderSource string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	var cShaderSource, cFreeShaderSource = gl.Strs(shaderSource)
	gl.ShaderSource(shader, 1, cShaderSource, nil)
	cFreeShaderSource()
	gl.CompileShader(shader)
	infoLogBuffer := make([]uint8, 10000)
	var infoLogLength int32
	gl.GetShaderInfoLog(shader, int32(len(infoLogBuffer)), &infoLogLength, &infoLogBuffer[0])
	if infoLogLength != 0 {
		var typeString = "unknown shader"
		switch shaderType {
		case gl.VERTEX_SHADER:
			typeString = "vertex shader"
		case gl.FRAGMENT_SHADER:
			typeString = "fragment shader"
		}
		panic(fmt.Sprintln(typeString + ":\n" + string(infoLogBuffer[0:infoLogLength])))
	}
	return shader
}

func glCheckError() {
	if e := gl.GetError(); e != 0 {
		str := fmt.Sprintf("0x%04X\n", e)
		switch e {
		case 0x0502:
			str = "Invalid Operation"
		case 0x0501:
			str = "Invalid Value"
		case 0x0500:
			str = "Invalid Enum"
		}
		panic("GlError: " + str)
	}
}
