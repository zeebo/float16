package float16

import (
	"strconv"
)

//
// the format of a float16:
//
// 0b0000000000000000
//   |________||||__|
//   |         ||   |__> exponent [0, 15]
//   |         ||______> exponent sign bit (1 means negative)
//   |         |_______> value sign bit (1 means negative)
//   |_________________> value [0, 1023]
//
// this is sufficient to represent all decimals like 1.23e10 for exponents in
// [-15, 15], both positive and negative.
//
// we chose this represenation for a few reasons:
//   1. we will not be doing any computation, so accuracy is only important
//      for storage, giving us a greater representable range (consider againt
//      the range of IEEE half floats)
//   2. the value is first in the bit stream so that we can extract it
//      with a single right shift.
//   3. the exponent is last in the bit stream so that we can extract it
//      with a single mask
//   4. since the sign bits are just single bits, we can extract them with
//      a single mask as well.
//
// some considerations about the representation:
//   1. we waste some entropy using 10 bits for the value portion but
//      only use 1000 of the availble 1023 bits, but this amount of waste
//      is small. log_2(1000) is 9.96.
//   2. we actually have 2^6 ways to represent zero, since the sign and
//      exponent bits would not affect the computation. maybe this can be
//      productive somehow.
//
// this code does not allow numbers like form 0.01e-15 yet, but i'm not worried
// about it.

type Float16 uint16

const (
	expSign   = 1 << 4
	valueSign = 1 << 5
)

func FromFloat64(val float64) (x Float16, ok bool) {
	if val == 0 {
		return 0, true
	}
	if val < 0 {
		x = valueSign
		val *= -1
	}

	switch {
	case val >= 1e16:
		return 0, false

	case val >= 1e15:
		return x | round(val*1e-13)<<6 | 15, true
	case val >= 1e14:
		return x | round(val*1e-12)<<6 | 14, true
	case val >= 1e13:
		return x | round(val*1e-11)<<6 | 13, true
	case val >= 1e12:
		return x | round(val*1e-10)<<6 | 12, true
	case val >= 1e11:
		return x | round(val*1e-9)<<6 | 11, true
	case val >= 1e10:
		return x | round(val*1e-8)<<6 | 10, true
	case val >= 1e9:
		return x | round(val*1e-7)<<6 | 9, true
	case val >= 1e8:
		return x | round(val*1e-6)<<6 | 8, true
	case val >= 1e7:
		return x | round(val*1e-5)<<6 | 7, true
	case val >= 1e6:
		return x | round(val*1e-4)<<6 | 6, true
	case val >= 1e5:
		return x | round(val*1e-3)<<6 | 5, true
	case val >= 1e4:
		return x | round(val*1e-2)<<6 | 4, true
	case val >= 1e3:
		return x | round(val*1e-1)<<6 | 3, true
	case val >= 1e2:
		return x | round(val*1e0)<<6 | 2, true
	case val >= 1e1:
		return x | round(val*1e1)<<6 | 1, true
	case val >= 1e0:
		return x | round(val*1e2)<<6 | 0, true
	case val >= 1e-1:
		return x | round(val*1e3)<<6 | 17, true
	case val >= 1e-2:
		return x | round(val*1e4)<<6 | 18, true
	case val >= 1e-3:
		return x | round(val*1e5)<<6 | 19, true
	case val >= 1e-4:
		return x | round(val*1e6)<<6 | 20, true
	case val >= 1e-5:
		return x | round(val*1e7)<<6 | 21, true
	case val >= 1e-6:
		return x | round(val*1e8)<<6 | 22, true
	case val >= 1e-7:
		return x | round(val*1e9)<<6 | 23, true
	case val >= 1e-8:
		return x | round(val*1e10)<<6 | 24, true
	case val >= 1e-9:
		return x | round(val*1e11)<<6 | 25, true
	case val >= 1e-10:
		return x | round(val*1e12)<<6 | 26, true
	case val >= 1e-11:
		return x | round(val*1e13)<<6 | 27, true
	case val >= 1e-12:
		return x | round(val*1e14)<<6 | 28, true
	case val >= 1e-13:
		return x | round(val*1e15)<<6 | 29, true
	case val >= 1e-14:
		return x | round(val*1e16)<<6 | 30, true
	case val >= 1e-15:
		return x | round(val*1e17)<<6 | 31, true

	default:
		return 0, false
	}
}

func round(val float64) Float16 {
	out := Float16(val)
	return out + Float16(2*(val-float64(out)))
}

// spans from -17 to 13 because we have a three digit number stored
var mulTable = [32]float64{
	17 + -17: 1e-17,
	17 + -16: 1e-16,
	17 + -15: 1e-15,
	17 + -14: 1e-14,
	17 + -13: 1e-13,
	17 + -12: 1e-12,
	17 + -11: 1e-11,
	17 + -10: 1e-10,
	17 + -9:  1e-9,
	17 + -8:  1e-8,
	17 + -7:  1e-7,
	17 + -6:  1e-6,
	17 + -5:  1e-5,
	17 + -4:  1e-4,
	17 + -3:  1e-3,
	17 + -2:  1e-2,
	17 + -1:  1e-1,
	17 + 0:   1e0,
	17 + 1:   1e1,
	17 + 2:   1e2,
	17 + 3:   1e3,
	17 + 4:   1e4,
	17 + 5:   1e5,
	17 + 6:   1e6,
	17 + 7:   1e7,
	17 + 8:   1e8,
	17 + 9:   1e9,
	17 + 10:  1e10,
	17 + 11:  1e11,
	17 + 12:  1e12,
	17 + 13:  1e13,
}

func (f Float16) Float64() float64 {
	exp := int(f & 0xf)
	if f&expSign != 0 {
		exp *= -1
	}
	exp += 15

	val := float64(f >> 6)
	if f&valueSign != 0 {
		val *= -1
	}

	return val * mulTable[exp&31]
}

func (f Float16) String() string {
	return strconv.FormatFloat(f.Float64(), 'e', 2, 64)
}
