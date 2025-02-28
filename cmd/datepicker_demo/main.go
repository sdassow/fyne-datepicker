package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	datepicker "github.com/sdassow/fyne-datepicker"
)

func main() {
	a := app.New()
	w := a.NewWindow("Demo")

	CtrlAltS := &desktop.CustomShortcut{fyne.KeyS, fyne.KeyModifierControl | fyne.KeyModifierAlt}
	w.Canvas().AddShortcut(CtrlAltS, func(_ fyne.Shortcut) {
		makeScreenshot(w)
	})

	dateInput := widget.NewEntry()
	now := time.Now()
	dateInput.SetPlaceHolder(now.Format("2006/01/02"))
	dateInput.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		when, err := time.Parse("2006/01/02 15:04", dateInput.Text)
		if err != nil {
			when = time.Now()
		}

		if dateInput.Text != "" {
			t, err := time.Parse("2006/01/02", dateInput.Text)
			if err == nil {
				when = t
			}
		}

		datepicker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				dateInput.SetText(when.Format("2006/01/02"))
				fmt.Printf("new Date: %s\n", when.Format("2006/01/02"))
			} else {
				fmt.Printf("old Date: %s\n", when.Format("2006/01/02"))
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
	datetimeInput.SetPlaceHolder(now.Format("2006/01/02 15:04"))
	datetimeInput.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		inWhen, err := time.Parse("2006/01/02 15:04", dateInput.Text)
		if err != nil {
			inWhen = time.Now()
		}

		if datetimeInput.Text != "" {
			t, err := time.Parse("2006/01/02 15:04", datetimeInput.Text)
			if err == nil {
				inWhen = t
			}
		}

		picker := datepicker.NewDateTimePicker(inWhen, time.Monday, func(when time.Time, ok bool) {
			if ok {
				datetimeInput.SetText(when.Format("2006/01/02 15:04"))
				fmt.Printf("new DateTime: %s\n", when.Format("2006/01/02 15:04"))
			} else {
				datetimeInput.SetText(inWhen.Format("2006/01/02 15:04"))
				fmt.Printf("old DateTime: %s\n", inWhen.Format("2006/01/02 15:04"))
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
