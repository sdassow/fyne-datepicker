package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/sdassow/fyne-datepicker"
)

func main() {
	a := app.New()
	w := a.NewWindow("Demo")

	dateInput := widget.NewEntry()
	dateInput.SetPlaceHolder("0000/00/00 00:00")
	dateInput.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		var d *dialog.CustomDialog

		when := time.Now()

		if dateInput.Text != "" {
			t, err := time.Parse("2006/01/02 15:04", dateInput.Text)
			if err == nil {
				when = t
			}
		}

		datepicker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				dateInput.SetText(when.Format("2006/01/02 15:04"))
			}
			d.Hide()
		})

		d = dialog.NewCustomWithoutButtons("Choose date and time", datepicker, w)
		d.Show()
	})

	label := widget.NewLabelWithStyle("Demo", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	w.SetContent(container.New(
		layout.NewBorderLayout(label, nil, nil, nil),
		label,
		widget.NewForm(
			widget.NewFormItem("Timestamp", dateInput),
		),
	))

	w.Resize(fyne.Size{
		Width:  640,
		Height: 480,
	})

	w.ShowAndRun()
}
