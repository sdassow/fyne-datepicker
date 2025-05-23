package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/sdassow/fyne-datepicker"
)

func main() {
	a := app.New()
	w := a.NewWindow("Demo")

	CtrlAltS := &desktop.CustomShortcut{fyne.KeyS, fyne.KeyModifierControl | fyne.KeyModifierAlt}
	w.Canvas().AddShortcut(CtrlAltS, func(_ fyne.Shortcut) {
		makeScreenshot(w)
	})

	dateInput := widget.NewEntry()
	dateInput.SetPlaceHolder("0000/00/00")
	dateInput.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		when := time.Now()

		if dateInput.Text != "" {
			t, err := time.Parse("2006/01/02", dateInput.Text)
			if err == nil {
				when = t
			}
		}

		datepicker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				dateInput.SetText(when.Format("2006/01/02"))
			}
		})

		dialog.ShowCustomConfirm(
			"Choose date",
			"Ok",
			"Cancel",
			datepicker,
			datepicker.OnActioned,
			w,
		)
	})

	datetimeInput := widget.NewEntry()
	datetimeInput.SetPlaceHolder("0000/00/00 00:00")
	datetimeInput.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		when := time.Now()

		if datetimeInput.Text != "" {
			t, err := time.Parse("2006/01/02 15:04", datetimeInput.Text)
			if err == nil {
				when = t
			}
		}

		picker := datepicker.NewDateTimePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				datetimeInput.SetText(when.Format("2006/01/02 15:04"))
			}
		})

		dialog.ShowCustomConfirm(
			"Choose date and time",
			"Ok",
			"Cancel",
			picker,
			picker.OnActioned,
			w,
		)
	})

	label := widget.NewLabelWithStyle("Demo", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	w.SetContent(container.New(
		layout.NewBorderLayout(label, nil, nil, nil),
		label,
		widget.NewForm(
			widget.NewFormItem("Date", dateInput),
			widget.NewFormItem("Timestamp", datetimeInput),
		),
	))

	w.Resize(fyne.Size{
		Width:  640,
		Height: 480,
	})

	w.ShowAndRun()
}
