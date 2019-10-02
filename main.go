package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	_ "github.com/fogleman/gg"
)

var firstString float64
var stringSpacing float64
var nutY float64
var redByte float64
var width int
var height int
var maxX float64
var maxY float64

func init() {
	width = 200
	height = 400
	maxX = 199.
	maxY = 399.
	firstString = 15.
	stringSpacing = 34.
	redByte = 0.
	nutY = 25.
}
func main() {

	chord := "D:x,x,o,2,3,2"

	dc := getEmptyDiagram()
	if !strings.Contains(chord, ":") {
		fmt.Println("Invalid chord " + chord)
		return
	}
	//chordName := strings.Split(chord, ":")[0]
	drawChord(dc, strings.Split(chord, ":")[1])
	dc.SavePNG("out.png")

}

func drawChord(dc *gg.Context, chord string) {
	fingers := strings.Split(chord, ",")
	minFret := 99
	maxFret := 5

	for stringNumber, finger := range fingers {
		if strings.Contains("xo", finger) {
			continue
		}
		if stringNumber < minFret {
			minFret = stringNumber
		}
		if stringNumber > maxFret {
			maxFret = stringNumber
		}
	}
	if maxFret == 5 {
		minFret = 1
	}
	fretHeight := drawFrets(dc, minFret, maxFret)

	for stringNumber, finger := range fingers {
		switch finger {
		case "o":
			drawOpenString(dc, float64(stringNumber))
		case "x":
			drawMutedString(dc, float64(stringNumber))

		default:
			fretNumber, err := strconv.Atoi(finger)
			if err == nil {
				drawPressedString(dc, float64(stringNumber), fretHeight, fretNumber)
			}
		}
	}
}
func drawFrets(dc *gg.Context, minFret, maxFret int) float64 {
	fretCount := maxFret - minFret + 1
	fretHeight := (float64(height) - nutY - 2) / float64(fretCount)
	for i := 1; i <= fretCount; i++ {
		dc.SetRGB(redByte, 0, 0)
		dc.DrawRectangle(firstString, nutY+(fretHeight*float64(i)), (stringSpacing * 5), 3)
		dc.Fill()
	}
	return fretHeight
}

func drawPressedString(dc *gg.Context, stringNumber, fretHeight float64, fretNumber int) {
	x := firstString + (stringNumber * stringSpacing) + 1
	y := (fretHeight * float64(fretNumber)) - (fretHeight / 2)
	dc.SetRGB(redByte, 0, 0)
	dc.DrawCircle(x, nutY+y, stringSpacing/3)
	dc.Fill()
}

func drawOpenString(dc *gg.Context, stringNumber float64) {
	x := firstString + (stringNumber * stringSpacing) + 1
	dc.SetRGB(redByte, 0, 0)
	dc.DrawCircle(x, nutY-(stringSpacing/3)-1, stringSpacing/3.8)
	dc.Stroke()
}

func drawMutedString(dc *gg.Context, stringNumber float64) {
	x := firstString + (stringNumber * stringSpacing) + 1
	dc.SetRGB(redByte, 0, 0)
	dc.DrawLine(x-(stringSpacing/3.8), 4, x+(stringSpacing/3.8), nutY-4)
	dc.DrawLine(x-(stringSpacing/3.8), nutY-4, x+(stringSpacing/3.8), 4)
	dc.Stroke()
}

func getEmptyDiagram() *gg.Context {

	dc := gg.NewContext(width, height)
	//Background
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(0, 0, maxX+1, maxY)
	dc.Fill()

	//Nut
	dc.SetRGB(redByte, 0, 0)
	dc.DrawRectangle(firstString+1, nutY, (stringSpacing * 5), 4)
	dc.Stroke()

	//Strings

	stringX := firstString

	dc.DrawRectangle(stringX, nutY, 3, maxY)
	dc.Fill()

	for i := 1; i < 6; i++ {
		stringX = stringX + stringSpacing
		dc.DrawRectangle(stringX, nutY, 2, maxY)
		dc.Fill()
	}

	return dc
}
