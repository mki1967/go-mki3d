package glmki3d

import (
	// "fmt"
	"errors"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/mki1967/go-mki3d/mki3d"
)

// DataShaderTex is a binding between data for texture elements and a shader for textured triangles
type DataShaderTex struct {
	ShaderPtr *ShaderTex // pointer to the GL shader program structure
	// VAO       uint32           // GL Vertex Array Object // each GLDataTexEl has its own VAO
	DataElements []*GLDataTexEl   // silce of GL data structures for texture elements
	UniPtr       *GLUni           // pointer to GL uniform parameters structure
	Mki3dPtr     *mki3d.Mki3dType // pointer to original Mki3dType data
}

// MakeDataShaderTex either returns a pointer to a newly created DataShaderTex or an error.
// The parameters should be pointers to existing and initiated objects.
// If mPtr.Texture!=nil then the function creates DataElements of the returned DataShaderTex
// If mPtr.Texture == nil then the function returns (nil, nil)
func MakeDataShaderTex(sPtr *ShaderTex, uPtr *GLUni, mPtr *mki3d.Mki3dType) (dsPtr *DataShaderTex, err error) {
	if mPtr == nil {
		return nil, errors.New("mPtr == nil // type *Mki3dType ")
	}
	if mPtr.Texture == nil { // there are no textures
		return nil, nil
	}
	if sPtr == nil {
		return nil, errors.New("sPtr == nil // type *ShaderTr ")
	}
	if uPtr == nil {
		return nil, errors.New("uPtr == nil // type *GLUni ")
	}

	dataElements := make([]*GLDataTexEl, 0, len(mPtr.Texture.Elements))
	//// TO DO

	for _, texEl := range mPtr.Texture.Elements {
		dataElement, err := MakeGLDataTexEl(&texEl, sPtr)
		if err != nil {
			return nil, err
		}
		dataElements = append(dataElements, dataElement)

	}

	ds := DataShaderTex{ShaderPtr: sPtr, DataElements: dataElements, UniPtr: uPtr, Mki3dPtr: mPtr}

	return &ds, nil
}

// UniTexSamplerUniToShader stes the sampler to the texture unit
func (ds *DataShaderTex) UniTexSamplerUniToShader(texUnit int32) (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.Uniform1i(ds.ShaderPtr.TexSamplerUni, texUnit)
	return nil
}

// UniLightToShader sets  light uniform parameters from ds.UniPtr to ds.ShaderPtr  (both must be not nil and previously initiated)
func (ds *DataShaderTex) UniLightToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.Uniform3fv(ds.ShaderPtr.LightUni, 1, &(ds.UniPtr.LightUni[0]))
	gl.Uniform1f(ds.ShaderPtr.AmbientUni, ds.UniPtr.AmbientUni)

	return nil
}

// UniModelToShader sets uniform parameter from ds.UniPtr to ds.ShaderPtr
func (ds *DataShaderTex) UniModelToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.UniformMatrix4fv(ds.ShaderPtr.ModelUni, 1, false, &(ds.UniPtr.ModelUni[0]))

	return nil
}

// UniViewToShader sets uniform parameter from ds.UniPtr to ds.ShaderPtr
func (ds *DataShaderTex) UniViewToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.UniformMatrix4fv(ds.ShaderPtr.ViewUni, 1, false, &(ds.UniPtr.ViewUni[0]))

	return nil
}

// UniProjectionToShader sets uniform parameter from ds.UniPtr to ds.ShaderPtr
func (ds *DataShaderTex) UniProjectionToShader() (err error) {
	if ds.ShaderPtr == nil {
		return errors.New("ds.ShaderPtr == nil // type *ShaderTr")
	}
	if ds.UniPtr == nil {
		return errors.New("ds.UniPtr == nil // type *GLUni")
	}

	gl.UseProgram(ds.ShaderPtr.ProgramId)
	gl.UniformMatrix4fv(ds.ShaderPtr.ProjectionUni, 1, false, &(ds.UniPtr.ProjectionUni[0]))

	return nil
}

// InitStage initiates stage parameters in ds.ShaderPtr assuming that ds is a stage
func (ds *DataShaderTex) InitStage() (err error) {
	if ds.Mki3dPtr == nil {
		return errors.New("ds.Mki3dPtr == nil // type *Mki3dType")
	}

	err = ds.UniProjectionToShader() // set projection
	if err != nil {
		return err
	}

	err = ds.UniViewToShader() // set view
	if err != nil {
		return err
	}

	err = ds.UniLightToShader() // set light - for triangles only
	if err != nil {
		return err
	}

	return nil

}

// Draw a model (textured)
func (ds *DataShaderTex) DrawModel() {
	ds.UniModelToShader()
	ds.UniTexSamplerUniToShader(0) // use zero as default texture unit
	gl.UseProgram(ds.ShaderPtr.ProgramId)
	for _, texEl := range ds.DataElements {
		if texEl.VertexCount != 0 {
			gl.BindVertexArray(texEl.VAO)
			gl.ActiveTexture(gl.TEXTURE0 + 0) // use zero as default texture unit
			gl.BindTexture(gl.TEXTURE_2D, texEl.Texture)
			// fmt.Printf("texEl = %v\n", texEl );
			gl.DrawArrays(gl.TRIANGLES, 0, texEl.VertexCount)
			gl.BindVertexArray(0)
		}
	}
}

// Draw a stage (textured).
func (ds *DataShaderTex) DrawStage() {
	ds.InitStage()
	ds.DrawModel()
}
