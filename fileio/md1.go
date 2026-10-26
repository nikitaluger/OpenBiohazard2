package fileio

// .md1 - 3D model file

import (
	"fmt"
	"io"
)

type MD1Header struct {
	SectionLengthBytes uint32
	Unknown            uint32
	NumObj             uint32
}

type MD1ObjectHeader struct {
	TrianglesHeader MD1TrianglesHeader
	QuadsHeader     MD1QuadsHeader
}

type MD1TrianglesHeader struct {
	VertexOffset        uint32
	VertexCount         uint32
	NormalOffset        uint32
	NormalCount         uint32
	TriangleIndexOffset uint32
	TriangleIndexCount  uint32
	TextureOffset       uint32
}

type MD1QuadsHeader struct {
	VertexOffset    uint32
	VertexCount     uint32
	NormalOffset    uint32
	NormalCount     uint32
	QuadIndexOffset uint32
	QuadIndexCount  uint32
	TextureOffset   uint32
}

type MD1Object struct {
	TriangleVertices []MD1Vertex
	TriangleNormals  []MD1Vertex
	TriangleIndices  []MD1TriangleIndex
	TriangleTextures []MD1TriangleTexture
	QuadVertices     []MD1Vertex
	QuadNormals      []MD1Vertex
	QuadIndices      []MD1QuadIndex
	QuadTextures     []MD1QuadTexture
}

type MD1Vertex struct {
	X    int16
	Y    int16
	Z    int16
	Zero uint16
}

type MD1TriangleIndex struct {
	IndexNormal0 uint16
	IndexVertex0 uint16
	IndexNormal1 uint16
	IndexVertex1 uint16
	IndexNormal2 uint16
	IndexVertex2 uint16
}

type MD1TriangleTexture struct {
	U0     uint8 // UV coordinates for vertex 0
	V0     uint8
	ClutId uint16 // Texture clut id, bits 0-5
	U1     uint8  // UV coordinates for vertex 1
	V1     uint8
	Page   uint16 // Texture page
	U2     uint8  // UV coordinates for vertex 2
	V2     uint8
	Zero   uint16
}

type MD1QuadIndex struct {
	IndexNormal0 uint16
	IndexVertex0 uint16
	IndexNormal1 uint16
	IndexVertex1 uint16
	IndexNormal2 uint16
	IndexVertex2 uint16
	IndexNormal3 uint16
	IndexVertex3 uint16
}

type MD1QuadTexture struct {
	U0     uint8 // UV coordinates for vertex 0
	V0     uint8
	ClutId uint16 // Texture clut id, bits 0-5
	U1     uint8  // UV coordinates for vertex 1
	V1     uint8
	Page   uint16 // Texture page
	U2     uint8  // UV coordinates for vertex 2
	V2     uint8
	Zero1  uint16
	U3     uint8 // UV coordinates for vertex 3
	V3     uint8
	Zero2  uint16
}

type MD1Output struct {
	Components []MD1Object
	NumBytes   int64
}

// validateMD1Offset ensures an offset is within the file and returns the
// number of bytes remaining from that offset.
func validateMD1Offset(offset int64, fileLength int64) (int64, error) {
	if offset < 0 || offset >= fileLength {
		return 0, fmt.Errorf("MD1 offset out of range: %d", offset)
	}
	return fileLength - offset, nil
}

// validateMD1Count ensures a slice allocation is both sane and fits in the
// remaining bytes of the file.
func validateMD1Count(count uint32, elemSize int, remaining int64) error {
	const maxCount = 1 << 20 // generous but safe upper bound
	if count > maxCount {
		return fmt.Errorf("MD1 count too large: %d", count)
	}
	if int64(count)*int64(elemSize) > remaining {
		return fmt.Errorf("MD1 count exceeds remaining file data: %d", count)
	}
	return nil
}

func LoadMD1Stream(r io.ReaderAt, fileLength int64) (*MD1Output, error) {
	fileReader := io.NewSectionReader(r, int64(0), fileLength)
	fileStreamReader := NewStreamReader(fileReader)

	// Read header
	md1Header := MD1Header{}
	if err := fileStreamReader.ReadData(&md1Header); err != nil {
		return nil, err
	}

	// Guard against corrupt or hostile headers causing huge allocations.
	if md1Header.NumObj > 1000 {
		return nil, fmt.Errorf("MD1 NumObj too large: %d", md1Header.NumObj)
	}

	// Read header offsets
	modelObjectHeaders := make([]MD1ObjectHeader, int(md1Header.NumObj)/2)
	if err := fileStreamReader.ReadData(&modelObjectHeaders); err != nil {
		return nil, err
	}

	totalBytes := int64(0)

	// Offsets are after model header, which is 12 bytes
	beginOffset := int64(12)
	objects := make([]MD1Object, len(modelObjectHeaders))
	for i := 0; i < len(modelObjectHeaders); i++ {
		modelObjectHeader := modelObjectHeaders[i]

		// Triangle data
		triVertexOffset := beginOffset + int64(modelObjectHeader.TrianglesHeader.VertexOffset)
		remaining, err := validateMD1Offset(triVertexOffset, fileLength)
		if err != nil {
			return nil, err
		}
		if err := validateMD1Count(modelObjectHeader.TrianglesHeader.VertexCount, 8, remaining); err != nil {
			return nil, err
		}
		fileStreamReader.SetPosition(triVertexOffset)
		triangleVertices := make([]MD1Vertex, modelObjectHeader.TrianglesHeader.VertexCount)
		if err := fileStreamReader.ReadData(&triangleVertices); err != nil {
			return nil, err
		}
		totalBytes = maxInt64(totalBytes, triVertexOffset+int64(modelObjectHeader.TrianglesHeader.VertexCount)*8)

		triNormalOffset := beginOffset + int64(modelObjectHeader.TrianglesHeader.NormalOffset)
		remaining, err = validateMD1Offset(triNormalOffset, fileLength)
		if err != nil {
			return nil, err
		}
		if err := validateMD1Count(modelObjectHeader.TrianglesHeader.NormalCount, 8, remaining); err != nil {
			return nil, err
		}
		fileStreamReader.SetPosition(triNormalOffset)
		triangleNormals := make([]MD1Vertex, modelObjectHeader.TrianglesHeader.NormalCount)
		if err := fileStreamReader.ReadData(&triangleNormals); err != nil {
			return nil, err
		}
		totalBytes = maxInt64(totalBytes, triNormalOffset+int64(modelObjectHeader.TrianglesHeader.NormalCount)*8)

		triIndexOffset := beginOffset + int64(modelObjectHeader.TrianglesHeader.TriangleIndexOffset)
		remaining, err = validateMD1Offset(triIndexOffset, fileLength)
		if err != nil {
			return nil, err
		}
		if err := validateMD1Count(modelObjectHeader.TrianglesHeader.TriangleIndexCount, 12, remaining); err != nil {
			return nil, err
		}
		fileStreamReader.SetPosition(triIndexOffset)
		triangleIndices := make([]MD1TriangleIndex, modelObjectHeader.TrianglesHeader.TriangleIndexCount)
		if err := fileStreamReader.ReadData(&triangleIndices); err != nil {
			return nil, err
		}
		totalBytes = maxInt64(totalBytes, triIndexOffset+int64(modelObjectHeader.TrianglesHeader.TriangleIndexCount)*12)

		triTextureOffset := beginOffset + int64(modelObjectHeader.TrianglesHeader.TextureOffset)
		remaining, err = validateMD1Offset(triTextureOffset, fileLength)
		if err != nil {
			return nil, err
		}
		if err := validateMD1Count(modelObjectHeader.TrianglesHeader.TriangleIndexCount, 12, remaining); err != nil {
			return nil, err
		}
		fileStreamReader.SetPosition(triTextureOffset)
		triangleTextures := make([]MD1TriangleTexture, modelObjectHeader.TrianglesHeader.TriangleIndexCount)
		if err := fileStreamReader.ReadData(&triangleTextures); err != nil {
			return nil, err
		}
		totalBytes = maxInt64(totalBytes, triTextureOffset+int64(modelObjectHeader.TrianglesHeader.TriangleIndexCount)*12)

		// Quad data
		quadVertexOffset := beginOffset + int64(modelObjectHeader.QuadsHeader.VertexOffset)
		remaining, err = validateMD1Offset(quadVertexOffset, fileLength)
		if err != nil {
			return nil, err
		}
		if err := validateMD1Count(modelObjectHeader.QuadsHeader.VertexCount, 8, remaining); err != nil {
			return nil, err
		}
		fileStreamReader.SetPosition(quadVertexOffset)
		quadVertices := make([]MD1Vertex, modelObjectHeader.QuadsHeader.VertexCount)
		if err := fileStreamReader.ReadData(&quadVertices); err != nil {
			return nil, err
		}
		totalBytes = maxInt64(totalBytes, quadVertexOffset+int64(modelObjectHeader.QuadsHeader.VertexCount)*8)

		quadNormalOffset := beginOffset + int64(modelObjectHeader.QuadsHeader.NormalOffset)
		remaining, err = validateMD1Offset(quadNormalOffset, fileLength)
		if err != nil {
			return nil, err
		}
		if err := validateMD1Count(modelObjectHeader.QuadsHeader.NormalCount, 8, remaining); err != nil {
			return nil, err
		}
		fileStreamReader.SetPosition(quadNormalOffset)
		quadNormals := make([]MD1Vertex, modelObjectHeader.QuadsHeader.NormalCount)
		if err := fileStreamReader.ReadData(&quadNormals); err != nil {
			return nil, err
		}
		totalBytes = maxInt64(totalBytes, quadNormalOffset+int64(modelObjectHeader.QuadsHeader.NormalCount)*8)

		// A quad has 2 triangles
		quadIndexOffset := beginOffset + int64(modelObjectHeader.QuadsHeader.QuadIndexOffset)
		remaining, err = validateMD1Offset(quadIndexOffset, fileLength)
		if err != nil {
			return nil, err
		}
		if err := validateMD1Count(modelObjectHeader.QuadsHeader.QuadIndexCount, 32, remaining); err != nil {
			return nil, err
		}
		fileStreamReader.SetPosition(quadIndexOffset)
		quadIndices := make([]MD1QuadIndex, modelObjectHeader.QuadsHeader.QuadIndexCount)
		if err := fileStreamReader.ReadData(&quadIndices); err != nil {
			return nil, err
		}
		totalBytes = maxInt64(totalBytes, quadIndexOffset+int64(modelObjectHeader.QuadsHeader.QuadIndexCount)*32)

		quadTextureOffset := beginOffset + int64(modelObjectHeader.QuadsHeader.TextureOffset)
		remaining, err = validateMD1Offset(quadTextureOffset, fileLength)
		if err != nil {
			return nil, err
		}
		if err := validateMD1Count(modelObjectHeader.QuadsHeader.QuadIndexCount, 16, remaining); err != nil {
			return nil, err
		}
		fileStreamReader.SetPosition(quadTextureOffset)
		quadTextures := make([]MD1QuadTexture, modelObjectHeader.QuadsHeader.QuadIndexCount)
		if err := fileStreamReader.ReadData(&quadTextures); err != nil {
			return nil, err
		}
		totalBytes = maxInt64(totalBytes, quadTextureOffset+int64(modelObjectHeader.QuadsHeader.QuadIndexCount)*16)

		objects[i] = MD1Object{
			TriangleVertices: triangleVertices,
			TriangleNormals:  triangleNormals,
			TriangleIndices:  triangleIndices,
			TriangleTextures: triangleTextures,
			QuadVertices:     quadVertices,
			QuadNormals:      quadNormals,
			QuadIndices:      quadIndices,
			QuadTextures:     quadTextures,
		}
	}

	md1Output := &MD1Output{
		Components: objects,
		NumBytes:   totalBytes,
	}
	return md1Output, nil
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
