//go:generate go run gen.go

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

func round64(val float64) Float16 {
	return (Float16(val) + 5) / 10
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
