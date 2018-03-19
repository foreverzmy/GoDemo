package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
)

// Image struct
type Image struct {
	W int
	H int
}

// Bounds 图片边界
func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.W, i.H)
}

// ColorModel 颜色模型
func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

// At
func (i Image) At(x, y int) color.Color {
	return color.RGBA{uint8(x), uint8(y), 255, 255}
}

// ImageShow 输出图片
func ImageShow(out io.Writer) {
	m := Image{100, 100}
	png.Encode(out, &m)
}

// ImageHandler print png to page
func ImageHandler(w http.ResponseWriter, r *http.Request) {
	ImageShow(w)
}
