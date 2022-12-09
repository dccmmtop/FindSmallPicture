package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Comdex/imgo"
	"github.com/fogleman/gg"
)

type Img struct {
	w   int
	h   int
	rgb [][]int
}

func newImg(path string) *Img {

	img, err := imgo.DecodeImage(path) // 获取 图片 image.Image 对象
	if err != nil {
		fmt.Println(err)
	}

	height := imgo.GetImageHeight(img) // 获取 图片 高度[height]
	width := imgo.GetImageWidth(img)   // 获取 图片 宽度[width]
	imgMatrix := imgo.MustRead(path)   // 读取图片RGBA值
	color := make([][]int, height)
	for i := range color {
		color[i] = make([]int, width)
	}

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			color[h][w] = (int(imgMatrix[h][w][0]))*1000 + 1000*int(imgMatrix[h][w][1]) + 1000*int(imgMatrix[h][w][2])
		}
	}
	imgRgb := Img{
		w:   width,
		h:   height,
		rgb: color,
	}
	return &imgRgb
}

func (bigImg *Img) include(smallImg *Img, offsetX int, offsetY int) (error, int, int) {
	same := true
	for h := 0; h < bigImg.h; h++ {
		for w := 0; w < bigImg.w; w++ {
			same = true
			for y := 0; y < smallImg.h; y += offsetY {
				if h+y >= bigImg.h {
					same = false
					break
				}
				for x := 0; x < smallImg.w; x += offsetX {
					if w+x >= bigImg.w {
						same = false
						break
					}
					if bigImg.rgb[h+y][w+x] != smallImg.rgb[y][x] {
						same = false
					}
				}
			}
			if same {
				return nil, w, h
			}
		}
	}
	return fmt.Errorf("not found"), 0, 0
}
func main() {
	bigImgPath := flag.String("big", "", "大图路径")
	smallImgPath := flag.String("small", "", "小图路径")
	draw := flag.Bool("draw", false, "是否需要画出目标图片的位置，默认保存在当前目录下的 draw.png")
	offsetX := flag.Int("offsetX", 10, "比对像素时X轴的步长,默认每隔10个像素点比对一次")
	offsetY := flag.Int("offsetY", 10, "比对像素时Y轴的步长,默认每隔10个像素点比对一次")
	flag.Parse()

	if *bigImgPath == "" || *smallImgPath == "" {
		fmt.Println("图片不能不能为空")
		os.Exit(0)
	}

	smallImg := newImg(*smallImgPath)
	if *offsetX >= smallImg.w || *offsetY >= smallImg.h {
		fmt.Println("X轴偏移量不能大于小图的宽，Y 轴的偏移量不能大于小图的高")
		os.Exit(0)
	}
	bigImg := newImg(*bigImgPath)
	err, sx, sy := bigImg.include(smallImg, *offsetX, *offsetY)
	if err != nil {
		fmt.Println("not found")
		return
	}
	if *draw {

		im, _ := gg.LoadImage(*bigImgPath)
		dc := gg.NewContextForImage(im)
		dc.SetHexColor("#dc0315")
		dc.SetLineWidth(2)
		dc.DrawRoundedRectangle(float64(sx), float64(sy), float64(smallImg.w), float64(smallImg.h), 0)
		dc.Stroke()
		dc.SavePNG("./draw.png")
	}
	fmt.Printf("%d %d\n", sx+smallImg.w/2, sy+smallImg.h/2)
}
