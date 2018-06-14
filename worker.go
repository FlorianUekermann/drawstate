package drawstate

import (
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func worker() {
	runtime.LockOSThread()

	var config = <-configChannel
	close(configChannel)

	// open a window
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	var window, err = glfw.CreateWindow(config.windowHeight, config.windowWidth, "window title", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(0)
	go func() {
		for {
			drawChannelLock.Lock()
			var shouldClose = window.ShouldClose() || drawChannelClosed
			glfw.PollEvents()
			drawChannelLock.Unlock()
			if shouldClose {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		Close()
	}()

	// initialize opengl program and texture
	gl.Init()
	var openGLProgram = compileProgram(vertexShader, fragmentShader)
	var openGLVertexArray uint32
	gl.GenVertexArrays(1, &openGLVertexArray)
	var openGLTexture uint32
	gl.GenTextures(1, &openGLTexture)
	gl.BindTexture(gl.TEXTURE_2D, openGLTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.R32UI, config.stateWidth, config.stateHeight, 0, gl.RED_INTEGER, gl.UNSIGNED_INT, nil)
	glCheckError()

	for state := range drawChannel {
		if len(state) != int(config.stateWidth)*int(config.stateHeight) {
			panic("state has wrong size")
		}

		// draw frame
		gl.UseProgram(openGLProgram)
		gl.BindVertexArray(openGLVertexArray)
		gl.BindTexture(gl.TEXTURE_2D, openGLTexture)
		gl.BindImageTexture(0, openGLTexture, 0, false, 0, gl.READ_ONLY, gl.R32UI)
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, config.stateWidth, config.stateHeight, gl.RED_INTEGER, gl.UNSIGNED_INT, unsafe.Pointer(&state[0]))
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		glCheckError()
		window.SwapBuffers()
	}
}
