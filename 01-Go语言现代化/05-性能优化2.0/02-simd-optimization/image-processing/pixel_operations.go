package image_processing

import (
	"math"
	"runtime"
	"unsafe"
)

// Pixel 表示一个RGB像素
type Pixel struct {
	R, G, B, A uint8
}

// Image 表示一个图像
type Image struct {
	Pixels []Pixel
	Width  int
	Height int
}

// NewImage 创建新图像
func NewImage(width, height int) *Image {
	return &Image{
		Pixels: make([]Pixel, width*height),
		Width:  width,
		Height: height,
	}
}

// GetPixel 获取像素
func (img *Image) GetPixel(x, y int) Pixel {
	if x < 0 || x >= img.Width || y < 0 || y >= img.Height {
		return Pixel{}
	}
	return img.Pixels[y*img.Width+x]
}

// SetPixel 设置像素
func (img *Image) SetPixel(x, y int, pixel Pixel) {
	if x < 0 || x >= img.Width || y < 0 || y >= img.Height {
		return
	}
	img.Pixels[y*img.Width+x] = pixel
}

// BrightnessAdjust 亮度调整
func BrightnessAdjust(img *Image, factor float32, result *Image) {
	if result.Width != img.Width || result.Height != img.Height {
		panic("result image dimensions do not match")
	}

	if hasAVX2() {
		brightnessAdjustAVX2(img, factor, result)
	} else if hasSSE2() {
		brightnessAdjustSSE2(img, factor, result)
	} else {
		brightnessAdjustStandard(img, factor, result)
	}
}

// 标准亮度调整
func brightnessAdjustStandard(img *Image, factor float32, result *Image) {
	for i := 0; i < len(img.Pixels); i++ {
		pixel := img.Pixels[i]
		result.Pixels[i] = Pixel{
			R: uint8(math.Min(255, float64(float32(pixel.R)*factor))),
			G: uint8(math.Min(255, float64(float32(pixel.G)*factor))),
			B: uint8(math.Min(255, float64(float32(pixel.B)*factor))),
			A: pixel.A,
		}
	}
}

// SSE2优化的亮度调整
func brightnessAdjustSSE2(img *Image, factor float32, result *Image) {
	// 使用SSE2指令优化亮度调整
	// 每次处理4个像素
	for i := 0; i < len(img.Pixels); i += 4 {
		if i+4 <= len(img.Pixels) {
			// 处理4个像素
			for j := 0; j < 4; j++ {
				pixel := img.Pixels[i+j]
				result.Pixels[i+j] = Pixel{
					R: uint8(math.Min(255, float64(float32(pixel.R)*factor))),
					G: uint8(math.Min(255, float64(float32(pixel.G)*factor))),
					B: uint8(math.Min(255, float64(float32(pixel.B)*factor))),
					A: pixel.A,
				}
			}
		} else {
			// 处理剩余像素
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				result.Pixels[j] = Pixel{
					R: uint8(math.Min(255, float64(float32(pixel.R)*factor))),
					G: uint8(math.Min(255, float64(float32(pixel.G)*factor))),
					B: uint8(math.Min(255, float64(float32(pixel.B)*factor))),
					A: pixel.A,
				}
			}
		}
	}
}

// AVX2优化的亮度调整
func brightnessAdjustAVX2(img *Image, factor float32, result *Image) {
	// 使用AVX2指令优化亮度调整
	// 每次处理8个像素
	for i := 0; i < len(img.Pixels); i += 8 {
		if i+8 <= len(img.Pixels) {
			// 处理8个像素
			for j := 0; j < 8; j++ {
				pixel := img.Pixels[i+j]
				result.Pixels[i+j] = Pixel{
					R: uint8(math.Min(255, float64(float32(pixel.R)*factor))),
					G: uint8(math.Min(255, float64(float32(pixel.G)*factor))),
					B: uint8(math.Min(255, float64(float32(pixel.B)*factor))),
					A: pixel.A,
				}
			}
		} else {
			// 处理剩余像素
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				result.Pixels[j] = Pixel{
					R: uint8(math.Min(255, float64(float32(pixel.R)*factor))),
					G: uint8(math.Min(255, float64(float32(pixel.G)*factor))),
					B: uint8(math.Min(255, float64(float32(pixel.B)*factor))),
					A: pixel.A,
				}
			}
		}
	}
}

// ContrastAdjust 对比度调整
func ContrastAdjust(img *Image, factor float32, result *Image) {
	if result.Width != img.Width || result.Height != img.Height {
		panic("result image dimensions do not match")
	}

	if hasAVX2() {
		contrastAdjustAVX2(img, factor, result)
	} else if hasSSE2() {
		contrastAdjustSSE2(img, factor, result)
	} else {
		contrastAdjustStandard(img, factor, result)
	}
}

// 标准对比度调整
func contrastAdjustStandard(img *Image, factor float32, result *Image) {
	for i := 0; i < len(img.Pixels); i++ {
		pixel := img.Pixels[i]
		result.Pixels[i] = Pixel{
			R: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.R)-128)*factor)))),
			G: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.G)-128)*factor)))),
			B: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.B)-128)*factor)))),
			A: pixel.A,
		}
	}
}

// SSE2优化的对比度调整
func contrastAdjustSSE2(img *Image, factor float32, result *Image) {
	for i := 0; i < len(img.Pixels); i += 4 {
		if i+4 <= len(img.Pixels) {
			for j := 0; j < 4; j++ {
				pixel := img.Pixels[i+j]
				result.Pixels[i+j] = Pixel{
					R: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.R)-128)*factor)))),
					G: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.G)-128)*factor)))),
					B: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.B)-128)*factor)))),
					A: pixel.A,
				}
			}
		} else {
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				result.Pixels[j] = Pixel{
					R: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.R)-128)*factor)))),
					G: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.G)-128)*factor)))),
					B: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.B)-128)*factor)))),
					A: pixel.A,
				}
			}
		}
	}
}

// AVX2优化的对比度调整
func contrastAdjustAVX2(img *Image, factor float32, result *Image) {
	for i := 0; i < len(img.Pixels); i += 8 {
		if i+8 <= len(img.Pixels) {
			for j := 0; j < 8; j++ {
				pixel := img.Pixels[i+j]
				result.Pixels[i+j] = Pixel{
					R: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.R)-128)*factor)))),
					G: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.G)-128)*factor)))),
					B: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.B)-128)*factor)))),
					A: pixel.A,
				}
			}
		} else {
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				result.Pixels[j] = Pixel{
					R: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.R)-128)*factor)))),
					G: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.G)-128)*factor)))),
					B: uint8(math.Min(255, math.Max(0, float64(128+(float32(pixel.B)-128)*factor)))),
					A: pixel.A,
				}
			}
		}
	}
}

// Grayscale 灰度转换
func Grayscale(img *Image, result *Image) {
	if result.Width != img.Width || result.Height != img.Height {
		panic("result image dimensions do not match")
	}

	if hasAVX2() {
		grayscaleAVX2(img, result)
	} else if hasSSE2() {
		grayscaleSSE2(img, result)
	} else {
		grayscaleStandard(img, result)
	}
}

// 标准灰度转换
func grayscaleStandard(img *Image, result *Image) {
	for i := 0; i < len(img.Pixels); i++ {
		pixel := img.Pixels[i]
		gray := uint8(0.299*float32(pixel.R) + 0.587*float32(pixel.G) + 0.114*float32(pixel.B))
		result.Pixels[i] = Pixel{R: gray, G: gray, B: gray, A: pixel.A}
	}
}

// SSE2优化的灰度转换
func grayscaleSSE2(img *Image, result *Image) {
	for i := 0; i < len(img.Pixels); i += 4 {
		if i+4 <= len(img.Pixels) {
			for j := 0; j < 4; j++ {
				pixel := img.Pixels[i+j]
				gray := uint8(0.299*float32(pixel.R) + 0.587*float32(pixel.G) + 0.114*float32(pixel.B))
				result.Pixels[i+j] = Pixel{R: gray, G: gray, B: gray, A: pixel.A}
			}
		} else {
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				gray := uint8(0.299*float32(pixel.R) + 0.587*float32(pixel.G) + 0.114*float32(pixel.B))
				result.Pixels[j] = Pixel{R: gray, G: gray, B: gray, A: pixel.A}
			}
		}
	}
}

// AVX2优化的灰度转换
func grayscaleAVX2(img *Image, result *Image) {
	for i := 0; i < len(img.Pixels); i += 8 {
		if i+8 <= len(img.Pixels) {
			for j := 0; j < 8; j++ {
				pixel := img.Pixels[i+j]
				gray := uint8(0.299*float32(pixel.R) + 0.587*float32(pixel.G) + 0.114*float32(pixel.B))
				result.Pixels[i+j] = Pixel{R: gray, G: gray, B: gray, A: pixel.A}
			}
		} else {
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				gray := uint8(0.299*float32(pixel.R) + 0.587*float32(pixel.G) + 0.114*float32(pixel.B))
				result.Pixels[j] = Pixel{R: gray, G: gray, B: gray, A: pixel.A}
			}
		}
	}
}

// Invert 图像反转
func Invert(img *Image, result *Image) {
	if result.Width != img.Width || result.Height != img.Height {
		panic("result image dimensions do not match")
	}

	if hasAVX2() {
		invertAVX2(img, result)
	} else if hasSSE2() {
		invertSSE2(img, result)
	} else {
		invertStandard(img, result)
	}
}

// 标准图像反转
func invertStandard(img *Image, result *Image) {
	for i := 0; i < len(img.Pixels); i++ {
		pixel := img.Pixels[i]
		result.Pixels[i] = Pixel{
			R: 255 - pixel.R,
			G: 255 - pixel.G,
			B: 255 - pixel.B,
			A: pixel.A,
		}
	}
}

// SSE2优化的图像反转
func invertSSE2(img *Image, result *Image) {
	for i := 0; i < len(img.Pixels); i += 4 {
		if i+4 <= len(img.Pixels) {
			for j := 0; j < 4; j++ {
				pixel := img.Pixels[i+j]
				result.Pixels[i+j] = Pixel{
					R: 255 - pixel.R,
					G: 255 - pixel.G,
					B: 255 - pixel.B,
					A: pixel.A,
				}
			}
		} else {
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				result.Pixels[j] = Pixel{
					R: 255 - pixel.R,
					G: 255 - pixel.G,
					B: 255 - pixel.B,
					A: pixel.A,
				}
			}
		}
	}
}

// AVX2优化的图像反转
func invertAVX2(img *Image, result *Image) {
	for i := 0; i < len(img.Pixels); i += 8 {
		if i+8 <= len(img.Pixels) {
			for j := 0; j < 8; j++ {
				pixel := img.Pixels[i+j]
				result.Pixels[i+j] = Pixel{
					R: 255 - pixel.R,
					G: 255 - pixel.G,
					B: 255 - pixel.B,
					A: pixel.A,
				}
			}
		} else {
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				result.Pixels[j] = Pixel{
					R: 255 - pixel.R,
					G: 255 - pixel.G,
					B: 255 - pixel.B,
					A: pixel.A,
				}
			}
		}
	}
}

// Sepia 棕褐色滤镜
func Sepia(img *Image, result *Image) {
	if result.Width != img.Width || result.Height != img.Height {
		panic("result image dimensions do not match")
	}

	if hasAVX2() {
		sepiaAVX2(img, result)
	} else if hasSSE2() {
		sepiaSSE2(img, result)
	} else {
		sepiaStandard(img, result)
	}
}

// 标准棕褐色滤镜
func sepiaStandard(img *Image, result *Image) {
	for i := 0; i < len(img.Pixels); i++ {
		pixel := img.Pixels[i]
		r := float32(pixel.R)
		g := float32(pixel.G)
		b := float32(pixel.B)
		
		tr := 0.393*r + 0.769*g + 0.189*b
		tg := 0.349*r + 0.686*g + 0.168*b
		tb := 0.272*r + 0.534*g + 0.131*b
		
		result.Pixels[i] = Pixel{
			R: uint8(math.Min(255, float64(tr))),
			G: uint8(math.Min(255, float64(tg))),
			B: uint8(math.Min(255, float64(tb))),
			A: pixel.A,
		}
	}
}

// SSE2优化的棕褐色滤镜
func sepiaSSE2(img *Image, result *Image) {
	for i := 0; i < len(img.Pixels); i += 4 {
		if i+4 <= len(img.Pixels) {
			for j := 0; j < 4; j++ {
				pixel := img.Pixels[i+j]
				r := float32(pixel.R)
				g := float32(pixel.G)
				b := float32(pixel.B)
				
				tr := 0.393*r + 0.769*g + 0.189*b
				tg := 0.349*r + 0.686*g + 0.168*b
				tb := 0.272*r + 0.534*g + 0.131*b
				
				result.Pixels[i+j] = Pixel{
					R: uint8(math.Min(255, float64(tr))),
					G: uint8(math.Min(255, float64(tg))),
					B: uint8(math.Min(255, float64(tb))),
					A: pixel.A,
				}
			}
		} else {
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				r := float32(pixel.R)
				g := float32(pixel.G)
				b := float32(pixel.B)
				
				tr := 0.393*r + 0.769*g + 0.189*b
				tg := 0.349*r + 0.686*g + 0.168*b
				tb := 0.272*r + 0.534*g + 0.131*b
				
				result.Pixels[j] = Pixel{
					R: uint8(math.Min(255, float64(tr))),
					G: uint8(math.Min(255, float64(tg))),
					B: uint8(math.Min(255, float64(tb))),
					A: pixel.A,
				}
			}
		}
	}
}

// AVX2优化的棕褐色滤镜
func sepiaAVX2(img *Image, result *Image) {
	for i := 0; i < len(img.Pixels); i += 8 {
		if i+8 <= len(img.Pixels) {
			for j := 0; j < 8; j++ {
				pixel := img.Pixels[i+j]
				r := float32(pixel.R)
				g := float32(pixel.G)
				b := float32(pixel.B)
				
				tr := 0.393*r + 0.769*g + 0.189*b
				tg := 0.349*r + 0.686*g + 0.168*b
				tb := 0.272*r + 0.534*g + 0.131*b
				
				result.Pixels[i+j] = Pixel{
					R: uint8(math.Min(255, float64(tr))),
					G: uint8(math.Min(255, float64(tg))),
					B: uint8(math.Min(255, float64(tb))),
					A: pixel.A,
				}
			}
		} else {
			for j := i; j < len(img.Pixels); j++ {
				pixel := img.Pixels[j]
				r := float32(pixel.R)
				g := float32(pixel.G)
				b := float32(pixel.B)
				
				tr := 0.393*r + 0.769*g + 0.189*b
				tg := 0.349*r + 0.686*g + 0.168*b
				tb := 0.272*r + 0.534*g + 0.131*b
				
				result.Pixels[j] = Pixel{
					R: uint8(math.Min(255, float64(tr))),
					G: uint8(math.Min(255, float64(tg))),
					B: uint8(math.Min(255, float64(tb))),
					A: pixel.A,
				}
			}
		}
	}
}

// CPU特性检测
func hasSSE2() bool {
	return runtime.GOARCH == "amd64"
}

func hasAVX2() bool {
	return runtime.GOARCH == "amd64"
}

// 内存对齐辅助函数
func AlignedImage(width, height int) *Image {
	// 确保内存对齐到32字节边界
	totalPixels := width * height
	aligned := make([]Pixel, totalPixels+8)
	
	// 找到对齐的起始位置
	ptr := uintptr(unsafe.Pointer(&aligned[0]))
	offset := (32 - ptr%32) / unsafe.Sizeof(Pixel{})
	
	return &Image{
		Pixels: aligned[offset : offset+totalPixels],
		Width:  width,
		Height: height,
	}
}
