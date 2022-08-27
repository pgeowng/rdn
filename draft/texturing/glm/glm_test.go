package glm

import (
	"fmt"
	"testing"
)

func TestMulv(t *testing.T) {

	cases := []struct {
		mat Mat4
		in  Vec4
		out string
	}{
		{
			Mat4{
				1, 5, 9, 13,
				2, 6, 10, 14,
				3, 7, 11, 15,
				4, 8, 12, 16,
			},
			Vec4{1, 2, 3, 4},
			"[30 70 110 150]",
		},
		{
			Mat4{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				10, 0, 0, 1,
			},
			Vec4{10, 10, 10, 1},
			"[20 10 10 1]",
		},
		{
			Mat4{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				10, 0, 0, 1,
			},
			Vec4{0, 0, -1, 0},
			"[0 0 -1 0]",
		},
	}

	for _, c := range cases {
		t.Run(c.out, func(t *testing.T) {
			res := fmt.Sprint(c.mat.Mulv(c.in))
			if res != c.out {
				t.Fatal(res)
			}
		})
	}

}

func TestTimes(t *testing.T) {
	cases := []struct {
		left   Mat4
		right  Mat4
		result string
	}{
		{
			Mat4{
				1, 2, 3, 4,
				5, 6, 7, 8,
				9, 10, 11, 12,
				13, 14, 15, 16,
			},
			Mat4{
				20, -24, 28, -32,
				-21, 25, -29, 33,
				22, -26, 30, -34,
				23, 27, -31, 35,
			},
			"[-264 -272 -280 -288 272 280 288 296 -280 -288 -296 -304 334 388 442 496]",
		},
	}

	for _, c := range cases {
		t.Run(c.result, func(t *testing.T) {
			res := fmt.Sprint(c.left.Times(c.right))
			if res != c.result {
				t.Fatal(res)
			}
		})
	}
}

func TestPerspective(t *testing.T) {
	cases := []struct {
		fov    float32
		aspect float32
		near   float32
		far    float32
		out    string
	}{
		{
			Deg2Rad(45),
			float32(640) / 480,
			0.1,
			100,
			"[1.81066 0 0 0 0 2.4142134 0 0 0 0 -1.002002 -1 0 0 -0.2002002 0]",
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%f %f %f %f", c.fov, c.aspect, c.near, c.far), func(t *testing.T) {

			res := fmt.Sprint(Perspective(c.fov, c.aspect, c.near, c.far))

			if res != c.out {
				t.Fatal(res)
			}

		})
	}
}
