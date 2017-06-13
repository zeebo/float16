// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"math"
	"strings"
	"text/template"
)

var funcs = template.FuncMap{
	"bits": bits,
}

var tmpl = template.Must(template.New("").Funcs(funcs).Parse(strings.TrimSpace(`
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
	case vb >= {{ bits 1e16 }} /* 1e+16 */ :
		return 0, false
{{ range .Cases }}
	case vb >= {{ bits .Value }} /* {{ printf "%1.0e" .Value }} */:
		return x | round64(val*1e{{ .Mul }})<<6 | {{ .Or }}, true
{{- end }}

	default:
		return 0, false
	}
}
`)))

func bits(val float64) string {
	return fmt.Sprintf("0x%016x", math.Float64bits(val))
}

func main() {
	type Case struct {
		Value float64
		Mul   int
		Or    int
	}

	var cases []Case

	for i := 15; i >= -15; i-- {
		or := i
		if i < 0 {
			or = 16 - i
		}
		cases = append(cases, Case{
			Value: math.Pow(10, float64(i)),
			Mul:   2 - i + 1,
			Or:    or,
		})
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, map[string]interface{}{
		"Cases": cases,
	})
	if err != nil {
		panic(err)
	}

	output, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile("from.go", output, 0644)
}
