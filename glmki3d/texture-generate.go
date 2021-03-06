package glmki3d

import (
	// "fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	// "github.com/mki1967/go-mki3d/glmki3d"
	"github.com/mki1967/go-mki3d/mki3d"
	"math"
	"strconv"
	"strings"
)

const texSize = 256

func MakeGeneratorVertexShader(def mki3d.TexturionDefType) string {
	test := strings.Join([]string{def.R, def.G, def.B, def.A}, "")
	if strings.ContainsAny(test, ";}") { // security check failed -> replace with black.
		def.R = "0.0"
		def.G = "0.0"
		def.B = "0.0"
		def.A = "1.0"
	}
	return "" +
		"#version 330\n" +
		"const float PI = " + strconv.FormatFloat(math.Pi, 'f', -1, 64) + ";\n" +
		"const int texSize= " + strconv.Itoa(texSize) + ";\n" +
		"float G(float x,float y);\n" +
		"float B(float x,float y);\n" +
		"float A(float x,float y);\n" +
		"float R(float x,float y){ return  " + def.R + "; }\n" +
		"float G(float x,float y){ return  " + def.G + "; }\n" +
		"float B(float x,float y){ return  " + def.B + "; }\n" +
		"float A(float x,float y){ return  " + def.A + "; }\n" +
		"layout (location = 0) in float h;\n" +
		"uniform float v;\n" +
		"out vec4 color;\n" +
		"void main()\n" +
		"{\n" +
		"  float  args[6];\n" +
		"  float h=h-float(texSize)/2.0;\n" +
		"  float v=v-float(texSize)/2.0;\n" +
		"  float x= 2.0*h/float(texSize); \n" +
		"  float y= 2.0*v/float(texSize); \n" +
		"  color= vec4( R(x,y), G(x,y), B(x,y), A(x,y) );\n" +
		"  gl_Position = vec4( x, y, 0.0, 1.0 );\n" + /// w=0.5 for perspective division
		"  gl_PointSize=1.0;\n" + /// test it
		"}\n" +
		"\x00"
}

var GeneratorFragmentShader = `
#version 330
in vec4 color;
out vec4 out_FragColor;
void main()
{
  out_FragColor= color;
//  out_FragColor= vec4( 1.0, 1.0, 0.0, 1.0) ; // red
}
` + "\x00"

// MakeGeneratorShaderProgram makes new GL shader program for generating the texture defined with def and
// returns its GL ID.
func MakeGeneratorShaderProgram(def mki3d.TexturionDefType) (programId uint32, err error) {
	vertexShader := MakeGeneratorVertexShader(def)
	// fmt.Printf("vertexShader:\n%v\n", vertexShader)                       //// test
	// fmt.Printf("GeneratorFragmentShader:\n%v\n", GeneratorFragmentShader) //// test
	return NewProgram(vertexShader, GeneratorFragmentShader)
}

// hBuffer is an auxiliary buffer used for texture generation
var hBufferId uint32
var hBufferIdExists = false

// frameBuffer is a frame to which a texture image is attached to be drawed on
var frameBufferId uint32
var frameBufferIdExists = false

func GenerateTexture(def mki3d.TexturionDefType) (textureId uint32, err error) {
	renderTextureShaderProgram, err := MakeGeneratorShaderProgram(def)

	if err != nil {
		return 0, err
	}

	// gl.BindFragDataLocation(renderTextureShaderProgram, 0, gl.Str("out_FragColor\x00")) // test

	/* set vertex attributes locations */
	hLocation := gl.GetAttribLocation(renderTextureShaderProgram, gl.Str("h\x00"))
	if hLocation < 0 {
		panic("hLocation=" + strconv.Itoa(int(hLocation)))
	}

	/* set uniform variables locations */
	vLocation := gl.GetUniformLocation(renderTextureShaderProgram, gl.Str("v\x00"))
	if vLocation < 0 {
		panic("vLocation=" + strconv.Itoa(int(vLocation)))
	}

	/* load hBuffer data if needed */
	if hBufferIdExists == false {
		gl.GenBuffers(1, &hBufferId)
		gl.BindBuffer(gl.ARRAY_BUFFER, hBufferId)

		var hIn [texSize + 4]float32
		for i := 0; i < texSize+4; i++ {
			hIn[i] = float32(i - 2)
		}
		gl.BufferData(gl.ARRAY_BUFFER, len(hIn)*4 /* 4 bytes per float32 */, gl.Ptr(&hIn[0]), gl.STATIC_DRAW)
		hBufferIdExists = true
	}

	/* init VAO */
	var renderTextureVAO uint32
	gl.UseProgram(renderTextureShaderProgram)
	gl.GenVertexArrays(1, &renderTextureVAO)
	gl.BindVertexArray(renderTextureVAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, hBufferId)
	gl.EnableVertexAttribArray(uint32(hLocation))
	gl.VertexAttribPointer(uint32(hLocation), 1, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.BindVertexArray(0) // unbind VAO

	gl.GenTextures(1, &textureId)
	/// TO DO: check textureId

	gl.ActiveTexture(gl.TEXTURE0 + 0)

	// set texture type, image and parameters
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, texSize, texSize, 0, /* border */
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	if frameBufferIdExists == false {
		/* create framebuffer object */
		gl.GenFramebuffers(1, &frameBufferId)
		frameBufferIdExists = true
	}

	// remember default FrameBuffer Object
	var defaultFBO int32
	gl.GetIntegerv(gl.FRAMEBUFFER_BINDING, &defaultFBO)

	// remember viewport
	var viewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &viewport[0]) // save viewport parameters

	gl.UseProgram(renderTextureShaderProgram)

	// fmt.Printf("frameBufferId = %v\n", frameBufferId)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBufferId)
	gl.Viewport(0, 0, texSize, texSize)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, textureId, 0)

	// fmt.Printf("hBufferId = %v\n", hBufferId)
	// gl.BindBuffer(gl.ARRAY_BUFFER, hBufferId)
	// gl.EnableVertexAttribArray(uint32(hLocation))

	gl.BindVertexArray(renderTextureVAO)
	for j := 0; j < texSize+4; j++ {
		// gl.VertexAttribPointer(uint32(hLocation), 1, gl.FLOAT, false, 0, gl.PtrOffset(0)) /// in the loop ?
		gl.Uniform1f(vLocation, float32(j-2))
		gl.DrawArrays(gl.POINTS, 0, texSize+4)
	}
	gl.BindVertexArray(0) // unbind
	// gl.DisableVertexAttribArray(uint32(hLocation))

	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(defaultFBO))          // return to default screen FBO
	gl.Viewport(viewport[0], viewport[1], viewport[2], viewport[3]) // restore viewport

	gl.DeleteVertexArrays(1, &renderTextureVAO)
	gl.DeleteProgram(renderTextureShaderProgram) // delete used program

	return textureId, nil
}
