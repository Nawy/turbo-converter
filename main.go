package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
)

func getMillis(endTime int64) float64 {
	return float64(endTime) * 0.001 * 0.001
}

func convertImage(input, output string, quality int) {
	fImg1, _ := imaging.Open(input)

	m := imaging.Resize(fImg1, 300, 0, imaging.Lanczos)
	m = imaging.Fill(m, 200, 200, imaging.Center, imaging.Lanczos)
	m = imaging.Sharpen(m, 1.2)
	m = imaging.AdjustBrightness(m, 3)
	m = imaging.AdjustContrast(m, 1)
	m = imaging.AdjustGamma(m, 1.0)

	toimg, _ := os.Create(output)
	defer toimg.Close()

	//imaging.Save(m, output)
	jpeg.Encode(toimg, m, &jpeg.Options{Quality: quality})
}

func printHelp() {
	fmt.Println("-i\tinput file")
	fmt.Println("-o\toutput file")
	fmt.Println("-q\tquality 1..100")
}

func isParamExists(key string, args []string) bool {
	var size int = len(args)
	for i := 0; i < size; i++ {
		if args[i] == key {
			return true
		}
	}
	return false
}

func getParam(key string, args []string) string {
	var size int = len(args)
	for i := 0; i < size; i++ {
		if args[i] == key {
			if i+1 < size {
				return args[i+1]
			}
			break
		}
	}
	return ""
}

func main() {

	args := os.Args[1:]
	qualityStr := getParam("-q", args)
	inStr := getParam("-i", args)
	outStr := getParam("-o", args)

	if isParamExists("-h", args) {
		printHelp()
		return
	}

	quality := 70
	if qualityStr != "" {
		quality, _ = strconv.Atoi(qualityStr)
		if quality > 100 || quality < 1 {
			fmt.Println("quallity should be in 1..100")
			return
		}
	}

	if inStr == "" {
		fmt.Println("Input file not found!")
		return
	}

	if outStr == "" {
		fmt.Println("Output file not found!")
		return
	}

	var results float64
	var startTime, endTime int64

	for i := 0; i < 3; i++ {
		startTime = time.Now().UnixNano()
		convertImage(inStr, outStr, quality)
		endTime = time.Now().UnixNano() - startTime
		deltaTime := getMillis(endTime)
		fmt.Printf("#%d = %f ms\n", i, deltaTime)
		results += deltaTime
	}

	fmt.Printf("Total = %f\n", results/3.0)
}
