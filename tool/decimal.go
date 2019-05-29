package tool

import "strconv"

func PrecFloat(f float64,prec int) float64  {
	v, _ := strconv.ParseFloat(strconv.FormatFloat(f,'f',prec,64), 64)
	return v
}
