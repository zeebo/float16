package float16

import "math"

func FromFloat64(val float64) (x Float16, ok bool) {
	if val == 0 {
		return 0, true
	}
	if val < 0 {
		x = valueSign
		val *= -1
	}

	vb := math.Float64bits(val)

	switch {
	case vb >= 0x4341c37937e08000 /* 1e+16 */ :
		return 0, false

	case vb >= 0x430c6bf526340000 /* 1e+15 */ :
		return x | round64(val*1e-13)<<6 | 15, true
	case vb >= 0x42d6bcc41e900000 /* 1e+14 */ :
		return x | round64(val*1e-12)<<6 | 14, true
	case vb >= 0x42a2309ce5400000 /* 1e+13 */ :
		return x | round64(val*1e-11)<<6 | 13, true
	case vb >= 0x426d1a94a2000000 /* 1e+12 */ :
		return x | round64(val*1e-10)<<6 | 12, true
	case vb >= 0x42374876e8000000 /* 1e+11 */ :
		return x | round64(val*1e-9)<<6 | 11, true
	case vb >= 0x4202a05f20000000 /* 1e+10 */ :
		return x | round64(val*1e-8)<<6 | 10, true
	case vb >= 0x41cdcd6500000000 /* 1e+09 */ :
		return x | round64(val*1e-7)<<6 | 9, true
	case vb >= 0x4197d78400000000 /* 1e+08 */ :
		return x | round64(val*1e-6)<<6 | 8, true
	case vb >= 0x416312d000000000 /* 1e+07 */ :
		return x | round64(val*1e-5)<<6 | 7, true
	case vb >= 0x412e848000000000 /* 1e+06 */ :
		return x | round64(val*1e-4)<<6 | 6, true
	case vb >= 0x40f86a0000000000 /* 1e+05 */ :
		return x | round64(val*1e-3)<<6 | 5, true
	case vb >= 0x40c3880000000000 /* 1e+04 */ :
		return x | round64(val*1e-2)<<6 | 4, true
	case vb >= 0x408f400000000000 /* 1e+03 */ :
		return x | round64(val*1e-1)<<6 | 3, true
	case vb >= 0x4059000000000000 /* 1e+02 */ :
		return x | round64(val*1e0)<<6 | 2, true
	case vb >= 0x4024000000000000 /* 1e+01 */ :
		return x | round64(val*1e1)<<6 | 1, true
	case vb >= 0x3ff0000000000000 /* 1e+00 */ :
		return x | round64(val*1e2)<<6 | 0, true
	case vb >= 0x3fb999999999999a /* 1e-01 */ :
		return x | round64(val*1e3)<<6 | 17, true
	case vb >= 0x3f847ae147ae147b /* 1e-02 */ :
		return x | round64(val*1e4)<<6 | 18, true
	case vb >= 0x3f50624dd2f1a9fc /* 1e-03 */ :
		return x | round64(val*1e5)<<6 | 19, true
	case vb >= 0x3f1a36e2eb1c432d /* 1e-04 */ :
		return x | round64(val*1e6)<<6 | 20, true
	case vb >= 0x3ee4f8b588e368f1 /* 1e-05 */ :
		return x | round64(val*1e7)<<6 | 21, true
	case vb >= 0x3eb0c6f7a0b5ed8d /* 1e-06 */ :
		return x | round64(val*1e8)<<6 | 22, true
	case vb >= 0x3e7ad7f29abcaf48 /* 1e-07 */ :
		return x | round64(val*1e9)<<6 | 23, true
	case vb >= 0x3e45798ee2308c3a /* 1e-08 */ :
		return x | round64(val*1e10)<<6 | 24, true
	case vb >= 0x3e112e0be826d695 /* 1e-09 */ :
		return x | round64(val*1e11)<<6 | 25, true
	case vb >= 0x3ddb7cdfd9d7bdbb /* 1e-10 */ :
		return x | round64(val*1e12)<<6 | 26, true
	case vb >= 0x3da5fd7fe1796495 /* 1e-11 */ :
		return x | round64(val*1e13)<<6 | 27, true
	case vb >= 0x3d719799812dea11 /* 1e-12 */ :
		return x | round64(val*1e14)<<6 | 28, true
	case vb >= 0x3d3c25c268497682 /* 1e-13 */ :
		return x | round64(val*1e15)<<6 | 29, true
	case vb >= 0x3d06849b86a12b9b /* 1e-14 */ :
		return x | round64(val*1e16)<<6 | 30, true
	case vb >= 0x3cd203af9ee75616 /* 1e-15 */ :
		return x | round64(val*1e17)<<6 | 31, true

	default:
		return 0, false
	}
}
