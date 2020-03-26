package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

var (
	width int
	path  string
)

func init() {
	flag.IntVar(&width, "w", 80, "width for ascii art")
	flag.StringVar(&path, "p", "test.png", "image file path (png/jpeg support)")
	flag.Parse()
}

func main() {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("open error:", err)
		return
	}
	defer f.Close()

	img, fn, err := image.Decode(f)
	if err != nil {
		fmt.Printf("decode error %s: %v", fn, err)
		return
	}

	fmt.Println(string(conv2Ascii(resize(img, width))))
}

func conv2Ascii(img image.Image, w, h int) []byte {
	table := []byte("MND8OZ$7I?+=~:,.")
	buf := new(bytes.Buffer)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			g := color.GrayModel.Convert(img.At(i, j))
			var pos int
			if c, ok := g.(color.Gray); ok {
				pos = (len(table) - 1) * int(c.Y) / math.MaxUint8
			}
			_ = buf.WriteByte(table[pos])
		}
		_ = buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func resize(in image.Image, newWidth int) (*image.NRGBA, int, int) {
	bounds := in.Bounds()
	ratio := bounds.Dx() / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.Y/ratio, bounds.Max.X/ratio, bounds.Max.Y/ratio))
	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+ratio, i+1 {
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)})
		}
	}
	return out, out.Rect.Size().X, out.Rect.Size().Y
}
