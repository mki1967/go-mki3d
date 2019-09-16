package texture

import (
	"github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/go-mki3d/glmki3d"
	"math"
	"strconv"
)

const texSize = 256


func MakeGeneratorVertexShader( def mki3d.TexturionDefType ) string {
	return ""+
	"const float PI = " +strconv.FormatFloat(math.Pi, 'f', -1, 64)  +";\n"+
	"const int texSize= "+strconv.Itoa(texSize)+";\n"+
	"float G(float x,float y);\n"+
	"float B(float x,float y);\n"+
	"float A(float x,float y);\n"+
	"float R(float x,float y){ return  "+def.R+"; }\n"+
	"float G(float x,float y){ return  "+def.G+"; }\n"+
	"float B(float x,float y){ return  "+def.B+"; }\n"+
	"float A(float x,float y){ return  "+def.A+"; }\n"+
	// "attribute float h;\n"+
	"layout (location = 0) in float h;\n"+
	"uniform float v;\n"+
	// "varying vec4 color;\n"+
	"out vec4 color;\n"+
	"void main()\n"+
	"{\n"+
	"  float  args[6];\n"+
	"  float h=h-float(texSize)/2.0;\n"+
	"  float v=v-float(texSize)/2.0;\n"+
	"  float x= 2.0*h/float(texSize); \n"+
	"  float y= 2.0*v/float(texSize); \n"+
	"  color= vec4( R(x,y), G(x,y), B(x,y), A(x,y) );\n"+
	"  gl_Position = vec4( x, y, 0.0, 1.0 );\n"+ /// w=0.5 for perspective division
	"  gl_PointSize=1.0;\n"+ /// test it
	"}\n"+
	"\x00";
}

var GeneratorFragmentShader = `
#version 330
in vec4 color;
out vec4 out_FragColor;
void main()
{
  out_FragColor= color;
}
` + "\x00"

// MakeGeneratorShaderProgram makes new GL shader program for generating the texture defined with def and
// returns its GL ID.
func MakeGeneratorShaderProgram(def mki3d.TexturionDefType) (program uint32, err error) {
	return glmki3d.NewProgram( MakeGeneratorVertexShader( def ), GeneratorFragmentShader )
}

