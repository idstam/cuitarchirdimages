package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
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
	nutY = 75.
}
func main() {

	chords, err := readLines("chords.txt")
	if err != nil {
		panic(err)
	}
	for _, chord := range chords {
		dc := getEmptyDiagramRH()
		if !strings.Contains(chord, ":") {
			fmt.Println("Invalid chord " + chord)
			return
		}
		chordName := strings.Split(chord, ":")[0]
		drawChordRH(dc, strings.Split(chord, ":")[1], chordName)
		dc.SavePNG("chordImages/" + chordName + "_RH.png")
	}

	chords, err = readLines("chords.txt")
	if err != nil {
		panic(err)
	}
	for _, chord := range chords {
		dc := getEmptyDiagramLH()
		if !strings.Contains(chord, ":") {
			fmt.Println("Invalid chord " + chord)
			return
		}
		chordName := strings.Split(chord, ":")[0]
		drawChordLH(dc, strings.Split(chord, ":")[1], chordName)
		dc.SavePNG("chordImages/" + chordName + "_LH.png")
	}

}

func drawChordLH(dc *gg.Context, chord, chordName string) {
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
	drawChordName(dc, chordName)
	for stringNumber, finger := range fingers {
		switch finger {
		case "o":
			drawOpenString(dc, float64(5-stringNumber))
		case "x":
			drawMutedString(dc, float64(5-stringNumber))

		default:
			fretNumber, err := strconv.Atoi(finger)
			if err == nil {
				drawPressedString(dc, float64(5-stringNumber), fretHeight, fretNumber)
			}
		}
	}
}

func drawChordRH(dc *gg.Context, chord, chordName string) {
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
	drawChordName(dc, chordName)
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

func drawChordName(dc *gg.Context, chordName string) {
	w, _ := dc.MeasureString(chordName)
	midFret := (float64(width) - (firstString * 2)) / 2
	midLabel := w / 2
	x := midFret - midLabel
	//if err := dc.LoadFontFace("fonts/dealerplate_california/dealerplate california.ttf", 40); err != nil {
	if err := dc.LoadFontFace("fonts/assurant_standard/Assurant-Standard.ttf", 40); err != nil {
		panic(err)
	}

	dc.DrawString(chordName, x, (nutY / 2))
	dc.Fill()
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
	dc.DrawLine(x-(stringSpacing/4), nutY-(stringSpacing/2), x+(stringSpacing/4), nutY-4)
	dc.DrawLine(x-(stringSpacing/4), nutY-4, x+(stringSpacing/4), nutY-(stringSpacing/2))
	dc.Stroke()
}

func getEmptyDiagramRH() *gg.Context {

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
func getEmptyDiagramLH() *gg.Context {

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

	dc.DrawRectangle(stringX+(stringSpacing*5)-3, nutY, 3, maxY)
	dc.Fill()

	for i := 0; i < 6; i++ {
		dc.DrawRectangle(stringX, nutY, 2, maxY)
		dc.Fill()
		stringX = stringX + stringSpacing
	}

	return dc
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
//https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
