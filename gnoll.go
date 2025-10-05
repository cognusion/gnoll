package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math/big"
	"os"
	"slices"
	"strconv"
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

	diceString := flag.String("dice", "", "If set, returns the result without creating a GUI. e.g. 2d10, 4d6.")
	flag.Parse()

	if diceString != nil && *diceString != "" {
		ds := strings.Split(*diceString, "d")
		if len(ds) > 2 {
			fmt.Printf("Error creating dice '%s': Too many 'd's\n", *diceString)
			os.Exit(1)
		}
		var (
			n string = "1"
			d string
		)
		if len(ds) == 1 {
			d = ds[0]
		} else {
			n = ds[0]
			d = ds[1]
		}

		die, err := NewDieFromString(d)
		if err != nil {
			fmt.Printf("Error creating dice '%s': %s\n", *diceString, err)
			os.Exit(1)
		}
		ni, err := strconv.Atoi(n)
		if err != nil {
			fmt.Printf("Error creating N dice '%s': %s\n", *diceString, err)
			os.Exit(1)
		}

		var (
			t  = big.NewInt(0)
			r  = big.NewInt(0)
			ts string
		)

		for range ni {
			r = die.Roll()
			t = t.Add(t, r)
			if ts == "" {
				ts = r.String()
			} else {
				ts = fmt.Sprintf("%s + %s", ts, r.String())
			}
		}
		fmt.Printf("Roll: %s = %s = %s\n", *diceString, ts, t.String())
		os.Exit(0)
	}

	a := app.NewWithID("io.cognusion.gnoll")
	a.SetIcon(&fyne.StaticResource{StaticName: "Icon.png", StaticContent: iconData})

	w := a.NewWindow("Gnoll")
	var history = make([]int64, 0)

	// build the numberlist
	var numberlist = make([]string, 100)
	for i := range 100 {
		numberlist[i] = fmt.Sprintf("%d", i+1)
	}

	// Label for historical stats
	statsLabel := canvas.NewText("", nil)
	statsLabel.TextSize = 10
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

	slices.Sort(dataCopy)

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
