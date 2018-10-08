package screen

import (
	"bytes"
	"github.com/lobocv/itree/ctx"
	"github.com/nsf/termbox-go"
	"strings"
	"math"
)

/*
Screen drawing methods
*/

type ScreenState int

const (
	Directory ScreenState = iota
	Help
)

type Screen struct {
	SearchString  []rune
	state         ScreenState
}

func GetScreen() Screen {
	return Screen{make([]rune, 0, 100), Directory}
}

func (s *Screen) Print(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (s *Screen) PrintDirContents(x, y int, dirlist ctx.DirView) error {
	var levelOffsetX, levelOffsetY int // Draw position offset
	var stretch int                    // Length of line connecting subdirectories
	var maxLineWidth int               // Length of longest item in the directory
	var scrollOffsety int			   // Offset to scroll the visible directory text by

	_, screenHeight := termbox.Size()

	levelOffsetX = x
	levelOffsetY = y

	// Determine the scrolling offset
	scrollOffsety = levelOffsetY
	for _, dir := range dirlist {
		scrollOffsety += dir.FileIdx
	}
	scrollOffsety -= screenHeight - levelOffsetY
	if scrollOffsety < 0 {
		scrollOffsety = 0
	} else {
		pagejump := float64(screenHeight) / 5
		scrollOffsety = int(math.Ceil(float64(scrollOffsety) / pagejump) * pagejump)
	}

	// Recurse through the directory list, drawing a tree structure
	for level, dir := range dirlist {
		maxLineWidth = 0

		for ii, f := range dir.Files {

			// Keep track of the longest length item in the directory
			filename_len := len(f.Name())
			if filename_len > maxLineWidth {
				maxLineWidth = filename_len
			}

			// Determine the color of the currently printing directory item
			var color termbox.Attribute
			if dir.FileIdx == ii && level == len(dirlist)-1 {
				color = termbox.ColorCyan
			} else {
				if _, ok := dir.FilteredFiles[ii]; ok {
					color = termbox.ColorGreen
				} else if f.IsDir() {
					color = termbox.ColorYellow
				} else {
					color = termbox.ColorWhite
				}

			}

			line := bytes.Buffer{}
			if ii == 0 {
				line.WriteString(strings.Repeat("─", stretch))
			}

			switch ii {
			case 0:
				if level > 0 {
					if len(dir.Files) < 2 {
						line.WriteString("─")
					} else {
						line.WriteString("┬─")
					}
				} else {
					line.WriteString("├─")
				}
			case len(dir.Files) - 1:
				line.WriteString("└─")
			default:
				line.WriteString("├─")
			}

			// Create the item label, add / if it is a directory
			line.WriteString(f.Name())
			if f.IsDir() {
				line.WriteString("/")
			}

			// Calculate the draw position
			x := levelOffsetY + ii - scrollOffsety
			y := levelOffsetX
			if ii == 0 {
				// The first item is connected to the parent directory with a line
				// shift the position left to account for this line
				y -= stretch
			}
			s.Print(y, x, color, termbox.ColorDefault, line.String())
		}

		// Determine the length of line we need to draw to connect to the next directory
		if len(dir.Files) > 0 {
			stretch = maxLineWidth - len(dir.Files[dir.FileIdx].Name())
		}

		// Shift the draw position in preparation for the next directory
		levelOffsetY += dir.FileIdx
		levelOffsetX += maxLineWidth + 2

	}

	return nil
}

func (s *Screen) GetState() ScreenState {
	return s.state
}
func (s *Screen) SetState(state ScreenState) {
	s.state = state
}

func (s *Screen) Draw(x, y int, view ctx.DirView) {
	switch s.state {
	case Help:
		s.Print(0, 0, termbox.ColorWhite, termbox.ColorDefault, "itree - An interactive tree application for file system navigation.")
		s.Print(0, 2, termbox.ColorWhite, termbox.ColorDefault, "h - Toggle hidden files and folders.")
		s.Print(0, 3, termbox.ColorWhite, termbox.ColorDefault, "e - Log2 skip up.")
		s.Print(0, 4, termbox.ColorWhite, termbox.ColorDefault, "d - Log2 skip down.")
		s.Print(0, 5, termbox.ColorWhite, termbox.ColorDefault, "c - Toggle position between first and last file.")
	case Directory:
		s.PrintDirContents(x, y, view)
	}

	termbox.Flush()
}

func (d *Screen) ClearScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}
