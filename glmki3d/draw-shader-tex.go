package glmki3d

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	// "github.com/mki1967/go-mki3d/glmki3d"
	// "github.com/mki1967/go-mki3d/mki3d"
)

var vertexShaderTex = `
#version 330
/* attributes */
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texAttr;
/* uniforms */
uniform mat4 model; 
uniform mat4 view;
uniform mat4 projection;
uniform vec3 light;
uniform float ambient; 
 
/* output to fragment shader */
out vec3 texUVS; // (u,v,shade)
void main() {
    /* compute shaded color */
    vec4 modelNormal=model*vec4(normal, 0.0);
    float shade= abs( dot( modelNormal.xyz, light ) ); // diffuse factor
    texUVS.xy=texAttr;
    texUVS.z=(ambient+(1.0-ambient)*shade); // shading scaling factor
    /* compute projected position */
    gl_Position = projection*view*model*vec4(position, 1.0);
}
` + "\x00"

// fragment shader - the same for segments and triangles
var fragmentShaderTex = `
#version 330
/* input from vertex shader */
in vec3 texUVS;
/* uniforms */
uniform sampler2D texSampler;
/* fragment color output */
out vec4 outputColor;
void main() {
    outputColor = vec4(texUVS.z*texture2D(texSampler, texUVS.xy).rgb, 1.0) ;
}
` + "\x00"

// ShaderTex structure for mki3d shader for drawing textured triangles
// with references to attributes and uniform locations.
type ShaderTex struct {
	// program Id
	ProgramId uint32
	// locations of attributes
	PositionAttr uint32
	NormalAttr   uint32
	TexAttr      uint32
	// locations of uniforms ( why int32 instead of uint32 ? )
	ProjectionUni int32
	ViewUni       int32
	ModelUni      int32
	LightUni      int32
	AmbientUni    int32
	TexSamplerUni int32
}

// MakeShaderTex compiles  mki3d shader for drawing textured triangles and
// returns pointer to its newly created  ShaderTex structure
func MakeShaderTex() (shaderPtr *ShaderTex, err error) {
	program, err := NewProgram(vertexShaderTex, fragmentShaderTex)
	if err != nil {
		return nil, err
	}

	var shader ShaderTex

	// set ProgramId
	shader.ProgramId = program

	// gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00")) // test

	// set attributes
	shader.PositionAttr = uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	shader.NormalAttr = uint32(gl.GetAttribLocation(program, gl.Str("normal\x00")))
	shader.TexAttr = uint32(gl.GetAttribLocation(program, gl.Str("texAttr\x00")))

	// set uniforms
	shader.ProjectionUni = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	shader.ViewUni = gl.GetUniformLocation(program, gl.Str("view\x00"))
	shader.ModelUni = gl.GetUniformLocation(program, gl.Str("model\x00"))
	shader.LightUni = gl.GetUniformLocation(program, gl.Str("light\x00"))
	shader.AmbientUni = gl.GetUniformLocation(program, gl.Str("ambient\x00"))
	shader.TexSamplerUni = gl.GetUniformLocation(program, gl.Str("texSampler\x00"))
	return &shader, nil
}
