package glm

import (
	"fmt"

	m32 "github.com/chewxy/math32"
)

type Vec3 [3]float32
type Vec4 [4]float32
type Mat4 [4 * 4]float32

var mat4zero Mat4 = Mat4{}

var mat4id Mat4 = Mat4{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}

func Identity() Mat4 {
	return mat4id
}

func (m Mat4) Translate(offset Vec3) Mat4 {
	m[12] += offset[0]
	m[13] += offset[1]
	m[14] += offset[2]
	return m
}

func Rad(deg float32) float32 {
	return deg / 180 * m32.Pi
}

func RotationX(r float32) (m Mat4) {
	m = mat4id
	m[i4(1, 1)] = m32.Cos(r)
	m[i4(1, 2)] = -m32.Sin(r)
	m[i4(2, 1)] = m32.Sin(r)
	m[i4(2, 2)] = m32.Cos(r)
	return
}

func RotationY(r float32) (m Mat4) {
	m = mat4id
	m[i4(0, 0)] = m32.Cos(r)
	m[i4(0, 3)] = m32.Sin(r)
	m[i4(3, 0)] = -m32.Sin(r)
	m[i4(3, 3)] = m32.Cos(r)
	return
}

func RotationZ(r float32) (m Mat4) {
	m = mat4id
	m[i4(0, 0)] = m32.Cos(r)
	m[i4(0, 1)] = -m32.Sin(r)
	m[i4(1, 0)] = m32.Sin(r)
	m[i4(1, 1)] = m32.Cos(r)
	return m
}

func j4(r, c int) int {
	return c + 4*r
}

func i4(r, c int) int {
	return 4*c + r
}

func (m Mat4) Times(o Mat4) (n Mat4) {
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			for i := 0; i < 4; i++ {
				n[i4(r, c)] += m[i4(r, i)] * o[i4(i, c)]
			}
		}
	}
	return
}

func (m Mat4) Mulv(v Vec4) (n Vec4) {
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			n[r] += m[i4(r, c)] * v[c]
		}
	}
	return
}

func (m *Mat4) Ptr() *float32 {
	return &m[0]
}

func Perspective(fov float32, aspect float32, near float32, far float32) (n Mat4) {
	fmt.Println("perspective", fov, aspect, near, far)
	heightViewport2 := near * m32.Tan(fov/2)
	widthViewport2 := heightViewport2 * aspect

	n[i4(0, 0)] = near / widthViewport2
	n[i4(1, 1)] = near / heightViewport2
	n[i4(2, 2)] = -(far + near) / (far - near)
	n[i4(3, 2)] = -1
	n[i4(2, 3)] = -2 * near * far / (far - near)
	return
}

func Perspect(fovy, aspect, near, far float32) Mat4 {
	nmf, f := near-far, 1./m32.Tan(fovy/2)

	return Mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (near + far) / nmf, -1,
		0, 0, 2. * far * near / nmf, 0,
	}
}
