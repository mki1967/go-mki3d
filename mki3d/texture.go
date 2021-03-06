package mki3d

/* data structures for textured triangles */

//Vector2dType is 2D vector in MKI3D - used for UV texture coordinates.
type Vector2dType [2]float32

//TriangleUVType is a sequence of UV coordinates of endpoints of a textured triangle.
type TriangleUVType [3]Vector2dType

// TexturionDefType is a Texturion definition of a textue.
// See: https://mki1967.github.io/texturion/
type TexturionDefType struct {
	Label string `json:"label"`
	R     string `json:"R"`
	G     string `json:"G"`
	B     string `json:"B"`
	A     string `json:"A"`
}

// TexturedTriangleType is a triangle with its UV endpoint's texture coordinates.
type TexturedTriangleType struct {
	Triangle   TriangleType   `json:"triangle"`
	TriangleUV TriangleUVType `json:"triangleUV"`
}

// TexturedTrianglesType is a sequence of TexturedTriangleType
type TexturedTrianglesType []TexturedTriangleType

// Get the array of triangles of type TrianglesType from TexturedTrianglesType
func (textured TexturedTrianglesType) GetTriangles() TrianglesType {
	triangles := make([]TriangleType, 0, len(textured))
	for _, texTriangle := range textured {
		triangles = append(triangles, texTriangle.Triangle)
	}
	return TrianglesType(triangles)
}

// Gets array which is a sequence of enpoints' UV coordinates of textured triangles
func (texTriangles TexturedTrianglesType) GetUVArrays() []float32 {
	data := make([]float32, 0, 6*len(texTriangles)) // each triangleUV has 3*2 coordinates
	for _, texTriangle := range texTriangles {
		for j := 0; j < 3; j++ {
			data = append(data, texTriangle.TriangleUV[j][0:2]...)
		}
	}
	return data
}

// TextureElementType is a texture definition with the sequence of triangles textured with this texture.
type TextureElementType struct {
	Def               TexturionDefType      `json:"def"`
	TexturedTriangles TexturedTrianglesType `json:"texturedTriangles"`
}

// TextureElementsType is a sequence of TextureElementType
type TextureElementsType []TextureElementType

// TextureType is a set of textures with triangles textured with the textures.
type TextureType struct {
	Elements TextureElementsType `json:"elements"`
	Index    int                 `json:"index"`
}
