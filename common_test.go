package float16

import (
	"math"
	"math/rand"
	"testing"
)

func TestFloat16(t *testing.T) {
	for e := -15; e <= 13; e++ {
		for base := 1; base < 1000; base++ {
			val := float64(base) * math.Pow(10, float64(e))

			f16, ok := FromFloat64(val)
			if !ok {
				t.Errorf("val:%v not ok", val)
				continue
			}

			got := f16.Float64()
			if diff := math.Abs(val - got); diff >= 0.001 {
				t.Errorf("val:%v got:%v", val, got)
				t.Errorf("diff:%v", diff)
				t.Errorf("val:%064b", math.Float64bits(val))
				t.Errorf("got:%064b", math.Float64bits(got))
				continue
			}
		}
	}
}

func BenchmarkFromFloat64(b *testing.B) {
	out := make([]float64, 1024)
	for i := range out {
		exp := math.Pow(10, float64(rand.Intn(29)-15))
		out[i] = float64(rand.Intn(1999)-999) * exp
	}

	b.SetBytes(8)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		FromFloat64(out[i&1023])
	}
}

var f64sink float64

func BenchmarkFloat16_Float64(b *testing.B) {
	out := make([]Float16, 1024)
	for i := range out {
		exp := math.Pow(10, float64(rand.Intn(29)-15))
		in := float64(rand.Intn(1999)-999) * exp
		val, ok := FromFloat64(in)
		if !ok {
			b.Fatalf("invalid value generated: %v", in)
		}
		out[i] = val
	}

	b.SetBytes(8)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		f64sink = out[i&1023].Float64()
	}
}
