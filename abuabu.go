package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	args := os.Args
	numOfShades := 2

	if len(args) < 2 {
		fmt.Println("require 1 argument")
		os.Exit(1)
	}

	if len(args) == 3 {
		nosString := args[2]
		nos, err := strconv.Atoi(nosString)
		if err != nil {
			fmt.Println("number of shades is invalid")
			os.Exit(1)
		}

		if nos < 2 && nos > 256 {
			fmt.Println("number of shades should between 2 -> 256")
			os.Exit(1)
		}

		numOfShades = nos
	}

	fmt.Println(numOfShades)

	filePath := args[1]

	// open image file
	fileIn, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error open file | %s\n", err.Error())
		os.Exit(1)
	}

	defer fileIn.Close()

	// get the file extension, eg: .jpg, .png ...
	fileExtension := path.Ext(filePath)

	// get base file, eg: when you open file to specific folder like this: /User/john/Documents/john.jpg
	// you will get result: john.jpg
	fileBase := path.Base(filePath)
	fileName := strings.TrimSuffix(fileBase, fileExtension)

	// return directory of filePath
	fileDir := path.Dir(filePath)

	// create new filename with customer prefix
	grayFileName := fmt.Sprintf("%s/%s_abu%s", fileDir, fileName, fileExtension)

	// decode image file to Image
	img, f, err := image.Decode(fileIn)
	if err != nil {
		fmt.Printf("error decode file | %s\n", err.Error())
		os.Exit(1)
	}

	// get rectangle representation of image X = Width, Y = Height
	bounds := img.Bounds()
	point := bounds.Size()

	rect := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{point.X, point.Y},
	}

	imgOut := image.NewRGBA(rect)

	// todo
	//imgOut.Set(point.X, point.Y, color.Black)

	// loop to every pixel
	for x := 0; x < point.X; x++ {

		// X's Y
		for y := 0; y < point.Y; y++ {

			pixel := img.At(x, y)

			// get original red, green, blue and alpha
			r, g, b, a := color.RGBAModel.Convert(pixel).RGBA()

			// http://en.wikipedia.org/wiki/Luma_%28video%29
			red := float64(r) * 0.299
			green := float64(g) * 0.587
			blue := float64(b) * 0.114

			// get average
			average := uint8((red + green + blue) / 3)

			// constract new color based on above calculation
			col := color.RGBA{R: average, G: average, B: average, A: uint8(a)}

			imgOut.Set(point.X, point.Y, col)
		}
	}

	fileOut, err := os.Create(grayFileName)
	if err != nil {
		fmt.Printf("error create file | %s\n", err.Error())
		os.Exit(1)
	}

	defer fileOut.Close()

	encode(f, imgOut, fileOut)
}

func encode(f string, img image.Image, input io.Writer) error {
	switch f {
	case "jpeg":
		return jpeg.Encode(input, img, nil)
	case "png":
		return png.Encode(input, img)
	case "gif":
		return gif.Encode(input, img, nil)
	default:
		return errors.New("unrecognized image format")
	}
}
