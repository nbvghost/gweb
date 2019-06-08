package tool

import "strconv"

func PrecFloat64(f float64,prec int) float64  {
	v, _ := strconv.ParseFloat(strconv.FormatFloat(f,'f',prec,64), 64)
	return v
}
func PrecFloat32(f float32,prec int) float32  {
	v, _ := strconv.ParseFloat(strconv.FormatFloat(float64(f),'f',prec,32), 32)
	return float32(v)
}
