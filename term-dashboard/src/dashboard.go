package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"os"
	"image/png"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Execute a given bash command and return the numeric value from stdout.
// This only works for commands that will return a integer like value.
func cmdToFloat(cmdString string, isInt bool) (result float64){
	var intResult int64
	cmd := exec.Command("bash", "-c", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing command '%s' : %s", cmdString, err)
		return 0.0
	}
	data := string(output)
	data = strings.TrimSuffix(data, "\n")
	if isInt {
		intResult, err = strconv.ParseInt(data, 10, 64)
		result = float64(intResult)
	} else {
		result, err = strconv.ParseFloat(data, 64)
	}
	if err != nil {
		log.Fatalf("Error parsing command output: '%s' : '%s' ", data, err)
		return 0.0
	}
	return result
}

func soilMoisture() (float64){
	cmd_string := "sort -R /home/pi/data/data.csv | head -n 1 | awk -F, '{print $3}'"
	output := cmdToFloat(cmd_string, true)
	return output
}

func uvIndex() (float64){
	cmd_string := "sort -R /home/pi/data/data.csv | head -n 1 | awk -F, '{print $6}'"
	output := cmdToFloat(cmd_string, false)
	return output
}

func ambientTemperature() (float64){
	cmd_string := "sort -R /home/pi/data/data.csv | head -n 1 | awk -F, '{print $2}'"
	output := cmdToFloat(cmd_string, false)
	return output
}

func initData(n int) (data [][]float64) {
	data = make([][]float64, 1)
	data[0] = make([]float64, n)
	return data
}

func rollArray(data []float64) ([]float64) {
	r := len(data) - 1 % len(data)
	data = append(data[r:], data[:r]...)
	return data
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	existingImageFile, err := os.Open("bonzai.png")
	if err != nil {
		log.Fatalf("Failed to open image file: %s", err)
	}
	defer existingImageFile.Close()

	imageData, err := png.Decode(existingImageFile)
	if err != nil {
		log.Fatalf("Failed to decode image data: %s", err)
	}

	imgWidget := widgets.NewImage(nil)
	imgWidget.Image = imageData
	imgWidget.SetRect(0, 0, 50, 10)
	imgWidget.MonochromeInvert = true

	moistureData := initData(40)
	uvData := initData(40)
	tempData := initData(40)

	moisturePlot := widgets.NewPlot()
	moisturePlot.Title = "Soil Moisture"
	moisturePlot.Data = moistureData
	moisturePlot.AxesColor = ui.ColorWhite
	moisturePlot.LineColors[0] = ui.ColorGreen
	moisturePlot.SetRect(0, 10, 50, 20)

	uvPlot := widgets.NewPlot()
	uvPlot.Title = "UV Index"
	uvPlot.Data = uvData
	uvPlot.AxesColor = ui.ColorWhite
	uvPlot.LineColors[0] = ui.ColorGreen
	uvPlot.SetRect(0, 20, 50, 30)

	tempPlot := widgets.NewPlot()
	tempPlot.Title = "Ambient Temperature"
	tempPlot.Data = tempData
	tempPlot.AxesColor = ui.ColorWhite
	tempPlot.LineColors[0] = ui.ColorGreen
	tempPlot.SetRect(0, 30, 50, 40)

	// grid := ui.NewGrid()
	// termWidth, termHeight := ui.TerminalDimensions()
	// grid.SetRect(0, 0, termWidth, termHeight)

	// grid.Set(
	// 	ui.NewRow(1.0/3,
	// 		ui.NewCol(1.0/2, moisturePlot),
	// 		ui.NewCol(1.0/2, imgWidget),
	// 	),
	// 	ui.NewRow(1.0/3,
	// 		ui.NewCol(1.0/2, uvPlot),
	// 	),
	// 	ui.NewRow(1.0/3,
	// 		ui.NewCol(1.0/2, tempPlot),
	// 	),
	// )

	// ui.Render(grid)
	ui.Render(imgWidget, moisturePlot, uvPlot, tempPlot)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:

			moistureData[0] = rollArray(moistureData[0])
			moistureData[0][0] = soilMoisture()
			moisturePlot.Data = moistureData

			uvData[0] = rollArray(uvData[0])
			uvData[0][0] = uvIndex()
			uvPlot.Data = uvData

			tempData[0] = rollArray(tempData[0])
			tempData[0][0] = ambientTemperature()
			tempPlot.Data = tempData

			// ui.Render(grid)
			ui.Render(imgWidget, moisturePlot, uvPlot, tempPlot)
		}
	}

}
