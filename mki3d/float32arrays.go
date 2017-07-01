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

// Gets TriangleArrays from mki3dData.
func (mki3dData *Mki3dType) GetTriangleArrays() *TriangleArrays {
	dataPos := make([]float32, 0, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	dataCol := make([]float32, 0, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	dataNor := make([]float32, 0, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	i := 0
	for _, triangle := range mki3dData.Model.Triangles {
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
			dataPos = append(dataPos, triangle[j].Position[0:3]...)
			dataCol = append(dataCol, triangle[j].Color[0:3]...)
			dataNor = append(dataNor, normal[0:3]...)
			i = i + 3
		}
	}
	return &TriangleArrays{Positions: dataPos, Colors: dataCol, Normals: dataNor}
}

// Gets SegmentArrays from mki3dData.
func (mki3dData *Mki3dType) GetSegmentArrays() *SegmentArrays {
	dataPos := make([]float32, 0, 6*len(mki3dData.Model.Segments)) // each segment has 2*3 coordinates
	dataCol := make([]float32, 0, 6*len(mki3dData.Model.Segments)) // each segment has 2*3 coordinates
	i := 0
	for _, segment := range mki3dData.Model.Segments {
		for j := 0; j < 2; j++ {
			dataPos = append(dataPos, segment[j].Position[0:3]...)
			dataCol = append(dataCol, segment[j].Color[0:3]...)
			i = i + 2
		}
	}
	return &SegmentArrays{Positions: dataPos, Colors: dataCol}
}

// Gets BufferData from mki3dData.
func (mki3dData *Mki3dType) GetBufferData() *BufferData {
	tPtr := mki3dData.GetTriangleArrays()
	sPtr := mki3dData.GetSegmentArrays()

	return &BufferData{TrArrPtr: tPtr, SegArrPtr: sPtr}

}
