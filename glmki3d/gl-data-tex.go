package glmki3d

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	// "github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/go-mki3d/mki3d"
)

// references to the objects defining the shape and parameters of mki3d object

// GLDataTexEl contains references to GL texture and buffers for  shader's input attributes
// for TexturedTriangles of single TexturedElement
type GLDataTexEl struct {
	// texture object
	Texture uint32
	// buffer objects in GL
	// triangles:
	VertexCount int32  // the last argument for gl.DrawArrays
	PositionBuf uint32 // positions of the endpoints
	NormalBuf   uint32 // normals of the endpooints
	TexUVBuf    uint32 // UV coordinates of the endpoints
	// VAO for the TexturedElement
	VAO uint32
}

// Delete the texture and buffers in GL, when they are not needed any more
func (glData *GLDataTexEl) Delete() {
	vbo := []uint32{glData.PositionBuf, glData.NormalBuf, glData.TexUVBuf}
	gl.DeleteBuffers(3, &vbo[0])
	textures := []uint32{glData.Texture}
	gl.DeleteTextures(1, &textures[0])
}

// LoadTriangleBufs loads data from mki3dData to the GL buffers referenced by glData
func (glData *GLDataTexEl) LoadTriangleBufs(texEl *mki3d.TextureElementType) {
	triangles := texEl.TexturedTriangles.GetTriangles()
	glData.VertexCount = int32(3 * len(triangles))
	if glData.VertexCount == 0 {
		return // do not create empty buffers
	}
	dataPos := triangles.GetPositionArrays()
	fmt.Printf("dataPos = %v\n", dataPos) //// test
	dataNor := triangles.GetNormalArrays()
	fmt.Printf("dataNor = %v\n", dataNor) //// test
	dataUV := texEl.TexturedTriangles.GetUVArrays()
	fmt.Printf("dataUV = %v\n", dataUV) //// test
	/* transfer data to the GL memory */
	gl.BindBuffer(gl.ARRAY_BUFFER, glData.PositionBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataPos)*4 /* 4 bytes per float32 */, gl.Ptr(dataPos), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, glData.NormalBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataNor)*4 /* 4 bytes per float32 */, gl.Ptr(dataNor), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, glData.TexUVBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataUV)*4 /* 4 bytes per float32 */, gl.Ptr(dataUV), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0) // unbind
}

// MakeGLBufTr either returns pointer to a new GLBufTr or an error
func MakeGLDataTexEl(texEl *mki3d.TextureElementType, shaderPtr *ShaderTex) (*GLDataTexEl, error) {
	if shaderPtr == nil {
		return nil, errors.New("shaderPtr == nil // type *ShaderTex")
	}

	var glData GLDataTexEl
	var vbo [3]uint32 // 5 is the number of buffers
	gl.GenBuffers(3, &vbo[0])
	// TO DO: test for error ...

	// assign buffer ids from vbo array
	glData.PositionBuf = vbo[0]
	glData.NormalBuf = vbo[1]
	glData.TexUVBuf = vbo[2]

	// load data from mki3dData
	glData.LoadTriangleBufs(texEl)

	//// Make texture

	texture, err := GenerateTexture(texEl.Def)
	if err != nil {
		return nil, err
	}

	glData.Texture = texture

	/// make and init VAO

	gl.UseProgram(shaderPtr.ProgramId)
	gl.GenVertexArrays(1, &(glData.VAO))
	gl.BindVertexArray(glData.VAO)
	// bind vertex positions
	gl.BindBuffer(gl.ARRAY_BUFFER, glData.PositionBuf)
	gl.EnableVertexAttribArray(shaderPtr.PositionAttr)
	gl.VertexAttribPointer(shaderPtr.PositionAttr, 3, gl.FLOAT, false, 0 /* stride */, gl.PtrOffset(0))

	// bind vertex UV
	gl.BindBuffer(gl.ARRAY_BUFFER, glData.TexUVBuf)
	gl.EnableVertexAttribArray(shaderPtr.TexAttr)
	gl.VertexAttribPointer(shaderPtr.TexAttr, 2, gl.FLOAT, false, 0 /* stride */, gl.PtrOffset(0))

	// bind vertex normals
	gl.BindBuffer(gl.ARRAY_BUFFER, glData.NormalBuf)
	gl.EnableVertexAttribArray(shaderPtr.NormalAttr)
	gl.VertexAttribPointer(shaderPtr.NormalAttr, 3, gl.FLOAT, false, 0 /* stride */, gl.PtrOffset(0))

	gl.BindVertexArray(0) // unbind VAO

	return &glData, nil
}
