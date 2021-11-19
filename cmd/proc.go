package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/infastin/gops/pkg"
)

const appID = "com.github.infastin.fyneps"

func main() {
	app := app.NewWithID(appID)
	win := app.NewWindow("FyneView")

	procs, err := ps.Processes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t := widget.NewTable(
		func() (int, int) { return len(procs), 9 },
		func() fyne.CanvasObject {
			return widget.NewLabel("1231512")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", procs[id.Row].Pid()))
			case 1:
				label.SetText(fmt.Sprintf("%d", procs[id.Row].PPid()))
			case 2:
				label.SetText(fmt.Sprintf("%d", procs[id.Row].Gid()))
			case 3:
				label.SetText(fmt.Sprintf("%d", procs[id.Row].Uid()))
			case 4:
				label.SetText(procs[id.Row].Executable())
			case 5:
				label.SetText(fmt.Sprintf("%d", procs[id.Row].NumThreads()))
			case 6:
				label.SetText(fmt.Sprintf("%d", procs[id.Row].StartTime()))
			case 7:
				label.SetText(fmt.Sprintf("%d", procs[id.Row].VirtMem()))
			case 8:
				label.SetText(fmt.Sprintf("%d", procs[id.Row].PhysMem()))
			}
		})

	container := container.NewBorder(nil, nil, nil, nil, t)

	win.SetContent(container)

	win.ShowAndRun()
}
