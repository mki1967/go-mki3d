package mki3d

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Float32 arrays that can be used loaded as input to gl.BufferData
// for triangle shader.
type TriangleArrays struct {
	Positions []float32
	Colors    []float32
	Normals   []float32
}

// Float32 arrays that can be used loaded as input to gl.BufferData
// for segment shader.
type SegmentArrays struct {
	Positions []float32
	Colors    []float32
}

// TriangleArrays and SegmentArrays bundled in one structure.
type BufferData struct {
	TrArrPtr  *TriangleArrays
	SegArrPtr *SegmentArrays
}

// Gets array which is a sequence of enpoints' positions' coordinates
func (triangles TrianglesType) GetPositionArrays() []float32 {
	data := make([]float32, 0, 9*len(triangles)) // each triangle has 3*3 coordinates
	for _, triangle := range triangles {
		for j := 0; j < 3; j++ {
			data = append(data, triangle[j].Position[0:3]...)
		}
	}
	return data
}

// Gets array which is a sequence of enpoints' colors' coordinates
func (triangles TrianglesType) GetColorArrays() []float32 {
	data := make([]float32, 0, 9*len(triangles)) // each triangle has 3*3 coordinates
	for _, triangle := range triangles {
		for j := 0; j < 3; j++ {
			data = append(data, triangle[j].Color[0:3]...)
		}
	}
	return data
}

// Gets array which is a sequence of triangles' normal coordinates repeated for each endpoint
func (triangles TrianglesType) GetNormalArrays() []float32 {
	data := make([]float32, 0, 9*len(triangles)) // each triangle has 3*3 coordinates
	for _, triangle := range triangles {
		// compute normal
		a := mgl32.Vec3(triangle[0].Position)
		b := mgl32.Vec3(triangle[1].Position)
		c := mgl32.Vec3(triangle[2].Position)
		normal := (b.Sub(a)).Cross(c.Sub(a))
		if normal.Dot(normal) > 0 {
			normal = normal.Normalize()
		}
		// append to buffers
		for j := 0; j < 3; j++ {
			data = append(data, normal[0:3]...)
		}
	}
	return data
}

// Gets TriangleArrays from mki3dData.
func (mki3dData *Mki3dType) GetTriangleArrays() *TriangleArrays {
	return &TriangleArrays{
		Positions: mki3dData.Model.Triangles.GetPositionArrays(),
		Colors:    mki3dData.Model.Triangles.GetColorArrays(),
		Normals:   mki3dData.Model.Triangles.GetNormalArrays(),
	}
}

// Gets array which is a sequence of enpoints' positions coordinates
func (segments SegmentsType) GetPositionArrays() []float32 {
	data := make([]float32, 0, 6*len(segments)) // each triangle has 3*3 coordinates
	for _, segment := range segments {
		for j := 0; j < 2; j++ {
			data = append(data, segment[j].Position[0:3]...)
		}
	}
	return data
}

// Gets array which is a sequence of enpoints' colors coordinates
func (segments SegmentsType) GetColorArrays() []float32 {
	data := make([]float32, 0, 6*len(segments)) // each triangle has 3*3 coordinates
	for _, segment := range segments {
		for j := 0; j < 2; j++ {
			data = append(data, segment[j].Color[0:3]...)
		}
	}
	return data
}

// Gets SegmentArrays from mki3dData.
func (mki3dData *Mki3dType) GetSegmentArrays() *SegmentArrays {
	return &SegmentArrays{
		Positions: mki3dData.Model.Segments.GetPositionArrays(),
		Colors:    mki3dData.Model.Segments.GetColorArrays(),
	}
}

// Gets BufferData from mki3dData.
func (mki3dData *Mki3dType) GetBufferData() *BufferData {
	tPtr := mki3dData.GetTriangleArrays()
	sPtr := mki3dData.GetSegmentArrays()

	return &BufferData{TrArrPtr: tPtr, SegArrPtr: sPtr}

}
