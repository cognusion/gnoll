package main

import (
	_ "embed"
	"fmt"
	"slices"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//go:embed Icon.png
var iconData []byte

func main() {

	a := app.NewWithID("io.cognusion.gnoll")
	a.SetIcon(&fyne.StaticResource{StaticName: "Icon.png", StaticContent: iconData})

	w := a.NewWindow("Gnoll")
	var history = make([]int64, 0)

	// build the numberlist
	var numberlist = make([]string, 100)
	for i := 0; i < 100; i++ {
		numberlist[i] = fmt.Sprintf("%d", i+1)
	}

	// Label for historical stats
	statsLabel := canvas.NewText("", nil)
	statsLabel.TextSize = 10
	//statsLabel := widget.NewLabel("")
	statsLabel.Alignment = fyne.TextAlignTrailing

	// Label for results
	topLabel := widget.NewLabel("Roll!")
	topLabel.Alignment = fyne.TextAlignCenter

	// Number of dice entry box
	ntext := widget.NewSelect(numberlist, nil)
	ntext.Alignment = fyne.TextAlignTrailing
	ntext.SetSelectedIndex(0) // 1

	// Label with a static "d"
	dLabel := widget.NewLabel("d")

	// Number of faces on the die entry box
	ftext := widget.NewSelect(numberlist, nil)
	ftext.Alignment = fyne.TextAlignLeading
	ftext.SetSelectedIndex(5) // 6

	/*
		// Checkbox
		totalCheck := widget.NewCheck("Total?", nil)
		totalCheck.SetChecked(true)
	*/

	// Button to roll NdF dice
	rollButton := widget.NewButton("Roll!", func() {
		faces := strings.TrimSpace(ftext.Selected)
		num := strings.TrimSpace(ntext.Selected)
		s, sum := RollNdF(num, faces, true)
		topLabel.SetText(s)
		history = append(history, sum)

		statsLabel.Text = fmt.Sprintf("Min: %d Mean: %.2f Median: %d Max: %d", slices.Min(history), mean(history), median(history), slices.Max(history))
		statsLabel.Refresh()
	})
	rollButton.Importance = widget.HighImportance

	clrButton := widget.NewButton("Clear", func() {
		statsLabel.Text = ""
		statsLabel.Refresh()
		topLabel.SetText("Roll!")
		history = make([]int64, 0)
	})
	clrButton.Importance = widget.DangerImportance

	rollBox := container.New(layout.NewCenterLayout(), container.New(layout.NewHBoxLayout(), rollButton, clrButton))

	// centered, horizontal box for holding our entry boxes
	diceBox := container.New(layout.NewCenterLayout(), container.New(layout.NewHBoxLayout(), ntext, dLabel, ftext))
	// vertical box for holding the top label, the diceBox, and our button
	box := container.New(layout.NewVBoxLayout(), statsLabel, topLabel, diceBox, rollBox)
	// centered box to hold the main box
	content := container.New(layout.NewCenterLayout(), box)

	w.SetContent(content)
	w.ShowAndRun()
}

func median(data []int64) int64 {
	dataCopy := make([]int64, len(data))
	copy(dataCopy, data)

	sort.Slice(dataCopy, func(i, j int) bool { return dataCopy[i] < dataCopy[j] })

	var median int64
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	} else {
		median = dataCopy[l/2]
	}

	return median
}

func mean(data []int64) float64 {
	if len(data) == 0 {
		return 0
	}
	var sum float64
	for _, d := range data {
		sum += float64(d)
	}
	return sum / float64(len(data))
}
