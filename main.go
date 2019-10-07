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
var fontPath string
var watermark string

func init() {
	width = 200
	height = 400
	maxX = 199.
	maxY = 375.
	firstString = 15.
	stringSpacing = 34.
	redByte = 0.
	nutY = 65.
	//fontPath = "fonts/assurant_standard/Assurant-Standard.ttf"
	//fontPath = "fonts/dealerplate_california/dealerplate california.ttf"
	fontPath = "fonts/techna_sans/TechnaSans-Regular.ttf"
	//watermark = "©Johan Idstam"
}
func main() {

	doDrawStringNames := false
	doDrawNoteNames := true
	chords, err := readLines("chords.txt")
	if err != nil {
		panic(err)
	}
	for _, chord := range chords {
		if !strings.Contains(chord, ":") {
			fmt.Println("Invalid chord " + chord)
			return
		}
		chordName := strings.Split(chord, ":")[0]

		//Lefty
		watermark = "left handed"
		dc := getEmptyDiagram(true, doDrawStringNames)
		drawChord(dc, strings.Split(chord, ":")[1], chordName, true, doDrawNoteNames)
		path := "chordImages/" + chordName + "_LH.png"
		dc.SavePNG(path)

		//Righty
		watermark = "right handed"
		dc = getEmptyDiagram(false, doDrawStringNames)
		drawChord(dc, strings.Split(chord, ":")[1], chordName, false, doDrawNoteNames)
		path = "chordImages/" + chordName + "_RH.png"
		dc.SavePNG(path)

	}

}

func drawChord(dc *gg.Context, chord, chordName string, isLeft, doDrawNoteNames bool) {
	fingers := strings.Split(chord, ",")
	minFret := 99
	maxFret := 5
	var err error

	for _, finger := range fingers {
		if strings.Contains("xo", finger) {
			continue
		}
		fingerNumber, err := strconv.Atoi(finger)
		if err != nil {
			panic(err)
		}
		if fingerNumber < minFret {
			minFret = fingerNumber
		}
		if fingerNumber > maxFret {
			maxFret = fingerNumber
		}
	}
	if maxFret == 5 {
		minFret = 1
	}
	fretHeight := drawFrets(dc, minFret, maxFret)
	drawChordName(dc, chordName)
	fretNumber := 0

	for stringNumber, finger := range fingers {
		stringToUse := stringNumber
		if isLeft {
			stringToUse = 5 - stringNumber
		}
		switch finger {
		case "o":
			drawOpenString(dc, float64(stringToUse))
			fretNumber = 0
		case "x":
			drawMutedString(dc, float64(stringToUse))
			fretNumber = -1
		default:
			fretNumber, err = strconv.Atoi(finger)
			if err == nil {
				drawPressedString(dc, float64(stringToUse), fretHeight, fretNumber-minFret+1)
			}
		}
		if doDrawNoteNames {
			if stringNumber == 6 {
				panic(fmt.Errorf("Invalid chord:, %s", chordName))

			}
			drawNoteName(dc, stringToUse, fretNumber, isLeft)
		}
	}
}

func drawFrets(dc *gg.Context, minFret, maxFret int) float64 {
	fretCount := maxFret - minFret + 1
	if fretCount < 5 {
		minFret--
		fretCount++
	}
	fretHeight := (maxY - nutY - 2) / float64(fretCount)
	for i := 0; i <= fretCount; i++ {
		dc.SetRGB(redByte, 0, 0)
		dc.DrawRectangle(firstString, nutY+(fretHeight*float64(i)), (stringSpacing*5)+2, 3)
		dc.Fill()
		if i < fretCount {
			drawFretNumber(dc, i, fretHeight, minFret)
		}
	}

	return fretHeight
}

func drawFretNumber(dc *gg.Context, fretNumber int, fretHeight float64, minFret int) {
	if err := dc.LoadFontFace(fontPath, 15); err != nil {
		panic(err)
	}
	fretNumStr := strconv.Itoa(fretNumber + minFret)
	x := firstString - 10

	dc.DrawString(fretNumStr, x, nutY+(fretHeight*float64(fretNumber))+20)
	dc.Fill()
}
func drawChordName(dc *gg.Context, chordName string) {
	if err := dc.LoadFontFace(fontPath, 40); err != nil {
		panic(err)
	}
	w, _ := dc.MeasureString(chordName)
	midFret := float64(width) / 2
	midLabel := w / 2
	x := midFret - midLabel

	dc.DrawString(chordName, x, (nutY / 2))
	dc.Fill()

	dc.SetRGB(.5, .5, .5)
	if err := dc.LoadFontFace(fontPath, 14); err != nil {
		panic(err)
	}

	chordName = watermark
	w, _ = dc.MeasureString(chordName)
	midFret = float64(width) / 2
	midLabel = w / 2
	x = midFret - midLabel

	dc.DrawString(chordName, x, (nutY/2)+12)
	dc.Fill()
	dc.SetRGB(redByte, 0, 0)
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

func drawNoteName(dc *gg.Context, stringNumber, fret int, isLeft bool) {

	if fret == -1 {
		return
	}
	if err := dc.LoadFontFace(fontPath, 15); err != nil {
		panic(err)
	}
	x := firstString + (float64(stringNumber) * stringSpacing) - 4

	note := Fret[fret][stringNumber]
	if isLeft {
		note = reverseNotes(Fret[fret])[stringNumber]
	}
	dc.DrawString(note, x, maxY+15)
	dc.Fill()
}
func reverseNotes(chord []string) []string {
	ret := []string{}
	for i := len(chord) - 1; i >= 0; i-- {
		ret = append(ret, chord[i])
	}
	return ret
}

func drawStringName(dc *gg.Context, stringNumber int, isLeft bool) {

	if err := dc.LoadFontFace(fontPath, 15); err != nil {
		panic(err)
	}
	x := firstString + (float64(stringNumber) * stringSpacing) - 4
	stringNames := []string{"E", "A", "D", "G", "B", "e"}
	if isLeft {
		stringNames = []string{"e", "B", "G", "D", "A", "E"}
	}
	dc.DrawString(stringNames[stringNumber], x, maxY+15)
	dc.Fill()
}

func getEmptyDiagram(isLeft, doDrawStringNames bool) *gg.Context {

	dc := gg.NewContext(width, height)
	//Background
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	dc.SetRGB(redByte, 0, 0)

	stringX := firstString
	thickE := 4.
	if isLeft {
		dc.DrawRectangle(stringX+(stringSpacing*5)-(thickE/3), nutY, thickE, maxY-nutY)
	} else {
		dc.DrawRectangle(stringX, nutY, thickE, maxY-nutY)
	}
	dc.Fill()

	for i := 0; i < 6; i++ {
		dc.DrawRectangle(stringX, nutY, 2, maxY-nutY)
		dc.Fill()
		stringX = stringX + stringSpacing
		if doDrawStringNames {
			drawStringName(dc, i, isLeft)
		}
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
