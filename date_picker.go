package datepicker

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var datePickerMonths = []string{
	"January",
	"Februrary",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

func daysInMonth(t time.Time, o int) int {
	// get first day of the given month, add a month, and go one day back
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1+o, -1).Day()
}

func firstWeekdayOfMonth(t time.Time) time.Weekday {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Weekday()
}

func lastWeekdayOfMonth(t time.Time) time.Weekday {
	// get first day of the given month, add a month, and go back one day
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, -1).Weekday()
}

func timeJumpMonth(when time.Time, offset int) time.Time {
	if offset == 0 {
		return when
	}

	// adjust day when higher than days in the destination month
	dstDays := daysInMonth(when, offset)
	if when.Day() > dstDays {
		when = time.Date(
			when.Year(),
			when.Month(),
			dstDays,
			when.Hour(),
			when.Minute(),
			when.Second(),
			when.Nanosecond(),
			when.Location(),
		)
	}

	return when.AddDate(0, offset, 0)
}

func timeJumpYearMonth(when time.Time, year int, month int) time.Time {
	dst := time.Date(
		year,
		time.Month(month),
		1,
		when.Hour(),
		when.Minute(),
		when.Second(),
		when.Nanosecond(),
		when.Location(),
	)

	// adjust day when higher than days in the destination month
	dstDays := daysInMonth(dst, 0)
	dstDay := when.Day()
	if dstDay > dstDays {
		dstDay = dstDays
	}

	return time.Date(
		dst.Year(),
		dst.Month(),
		dstDay,
		when.Hour(),
		when.Minute(),
		when.Second(),
		when.Nanosecond(),
		when.Location(),
	)
}

// weeks start on Monday, why is it otherwise called "the weekend"?!
func adjustWeekday(d time.Weekday, weekStart time.Weekday) int {
	if weekStart == 0 {
		return int(d)
	}
	return ((7 - int(weekStart)) + int(d)) % 7
}

// offsets correspond to time.Weekday, hence starting on Monday
var weekdayLabels = []*widget.Label{
	widget.NewLabelWithStyle("Sun", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	widget.NewLabelWithStyle("Mon", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	widget.NewLabelWithStyle("Tue", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	widget.NewLabelWithStyle("Wed", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	widget.NewLabelWithStyle("Thu", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	widget.NewLabelWithStyle("Fri", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	widget.NewLabelWithStyle("Sat", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
}

func updateGrid(grid *fyne.Container, when time.Time, weekStart time.Weekday, updateWhen func(time.Time), updateSelects func(time.Time)) {
	objs := []fyne.CanvasObject{}

	// row of weekdays at the top
	for n := weekStart; n < 7; n++ {
		objs = append(objs, weekdayLabels[n])
	}
	for n := 0; n < int(weekStart); n++ {
		objs = append(objs, weekdayLabels[n])
	}

	firstWeekday := adjustWeekday(firstWeekdayOfMonth(when), weekStart)
	lastWeekday := adjustWeekday(lastWeekdayOfMonth(when), weekStart)
	days := daysInMonth(when, 0)
	daysPrevMonth := daysInMonth(when, -1)

	// empty fields for days of the previous month that cut into the first week
	for n := 1; n <= firstWeekday; n++ {
		day := daysPrevMonth - firstWeekday + n
		button := widget.NewButton(fmt.Sprintf("%d", day), func() {
			when = time.Date(
				when.Year(),
				when.Month(),
				1,
				when.Hour(),
				when.Minute(),
				when.Second(),
				when.Nanosecond(),
				when.Location(),
			).AddDate(0, -1, day-1)

			updateWhen(when)
			updateSelects(when)
		})
		button.Importance = widget.LowImportance

		objs = append(objs, button)
	}

	var buttons []*widget.Button
	for n := 1; n <= days; n++ {
		var button *widget.Button
		var day = n

		button = widget.NewButton(fmt.Sprintf("%02d", n), func() {
			when = time.Date(
				when.Year(),
				when.Month(),
				day,
				when.Hour(),
				when.Minute(),
				when.Second(),
				when.Nanosecond(),
				when.Location(),
			)
			updateWhen(when)

			// reset importance of all buttons
			for _, b := range buttons {
				b.Importance = widget.MediumImportance
			}
			// only highlight selected day
			button.Importance = widget.HighImportance
			grid.Refresh()
		})

		// initially highlight a given day
		if n == when.Day() {
			button.Importance = widget.HighImportance
		}

		buttons = append(buttons, button)

		objs = append(objs, button)
	}

	// empty fields for days after the previous month
	for n := 1; lastWeekday < 7 && n < 7-(lastWeekday%7); n++ {
		day := n
		button := widget.NewButton(fmt.Sprintf("%d", day), func() {
			when = time.Date(
				when.Year(),
				when.Month(),
				1,
				when.Hour(),
				when.Minute(),
				when.Second(),
				when.Nanosecond(),
				when.Location(),
			).AddDate(0, 1, day-1)

			updateWhen(when)
			updateSelects(when)
		})
		button.Importance = widget.LowImportance

		objs = append(objs, button)
	}

	// add up to another empty row to compensate for months with a high first weekday
	for n := len(objs); n < 7*7; n++ {
		objs = append(objs, widget.NewLabel(""))
	}

	grid.Objects = objs
	grid.Refresh()
}

func findMonth(month string) int {
	for n := 0; n < len(datePickerMonths); n++ {
		if datePickerMonths[n] == month {
			return n + 1
		}
	}
	return 0
}

func NewDatePicker(when time.Time, weekStart time.Weekday, fn func(time.Time, bool)) fyne.CanvasObject {
	var updateSelects func(time.Time)

	grid := container.New(layout.NewGridLayoutWithColumns(7))

	updateWhen := func(t time.Time) {
		when = t
	}

	monthSelect := widget.NewSelect(datePickerMonths, func(selected string) {
		i := findMonth(selected)
		if i == 0 {
			return
		}

		when = timeJumpYearMonth(when, when.Year(), i)

		updateGrid(grid, when, weekStart, updateWhen, updateSelects)
	})
	monthSelect.Selected = when.Month().String()

	years := []string{}
	// inverted years, most recent on top for easy selection
	for n := when.Year() + 10; n >= when.Year()-100; n-- {
		years = append(years, fmt.Sprintf("%d", n))
	}
	yearSelect := widget.NewSelect(years, func(selected string) {
		i, err := strconv.ParseInt(selected, 10, 64)
		if err != nil {
			return
		}

		when = timeJumpYearMonth(when, int(i), int(when.Month()))

		updateGrid(grid, when, weekStart, updateWhen, updateSelects)
	})
	yearSelect.Selected = fmt.Sprintf("%d", when.Year())

	updateSelects = func(t time.Time) {
		// directly assign instead of setter methods to avoid multiple updates
		monthSelect.Selected = t.Month().String()
		monthSelect.Refresh()
		yearSelect.Selected = fmt.Sprintf("%d", t.Year())
		yearSelect.Refresh()

		updateGrid(grid, t, weekStart, updateWhen, updateSelects)
	}

	prevMonthButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		when = timeJumpMonth(when, -1)

		updateSelects(when)
	})

	nextMonthButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		when = timeJumpMonth(when, 1)

		updateSelects(when)
	})

	top := container.New(
		// previous and next button left and right
		layout.NewBorderLayout(nil, nil, prevMonthButton, nextMonthButton),
		prevMonthButton,
		nextMonthButton,
		// month and year dropdowns centered in the middle
		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			monthSelect,
			yearSelect,
			layout.NewSpacer(),
		),
	)

	updateGrid(grid, when, weekStart, updateWhen, updateSelects)

	hours := []string{}
	for n := 0; n < 23; n++ {
		hours = append(hours, fmt.Sprintf("%02d", n))
	}
	hourInput := widget.NewSelectEntry(hours)

	minutes := []string{}
	for n := 0; n < 59; n++ {
		minutes = append(minutes, fmt.Sprintf("%02d", n))
	}
	minuteInput := widget.NewSelectEntry(minutes)

	controlButtons := container.New(
		layout.NewHBoxLayout(),
		widget.NewButton("Now", func() {
			when = time.Now()

			hourInput.SetText(when.Format("15"))
			minuteInput.SetText(when.Format("04"))

			updateSelects(when)
		}),
		widget.NewButton("Cancel", func() {
			fn(when, false)
		}),
		widget.NewButton("Ok", func() {
			fn(when, true)
		}),
	)

	hourInput.SetText(when.Format("15"))
	hourInput.OnChanged = func(str string) {
		t, err := time.Parse("15", str)
		if err != nil {
			fyne.LogError("invalid hour", err)
			return
		}

		when = time.Date(
			when.Year(),
			when.Month(),
			when.Day(),
			t.Hour(),
			when.Minute(),
			0,
			0,
			when.Location(),
		)
	}

	minuteInput.SetText(when.Format("04"))
	minuteInput.OnChanged = func(str string) {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			fyne.LogError("failed to parse minte value", err)
			return
		}
		if i < 0 || i > 59 {
			fyne.LogError("minute value out of range", err)
			return
		}
		when = time.Date(
			when.Year(),
			when.Month(),
			when.Day(),
			when.Hour(),
			int(i),
			0,
			0,
			when.Location(),
		)
	}

	timeForm := widget.NewForm(
		widget.NewFormItem("Time",
			container.NewHBox(hourInput, widget.NewLabel(":"), minuteInput),
		),
	)

	bottom := container.New(
		layout.NewBorderLayout(nil, nil, nil, controlButtons),
		controlButtons,
		timeForm,
	)

	return container.New(
		layout.NewBorderLayout(top, bottom, nil, nil),
		top,
		grid,
		bottom,
	)
}
