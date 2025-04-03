package datepicker

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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

func updateGrid(grid *fyne.Container, when time.Time, weekStart time.Weekday, updateWhen func(time.Time), updateSelects func(time.Time), strs *dateTimePickerStrings) {
	objs := []fyne.CanvasObject{}

	left := fyne.TextAlignLeading
	bold := fyne.TextStyle{Bold: true}
	// row of weekdays at the top
	for n := weekStart; n < 7; n++ {
		objs = append(objs, widget.NewLabelWithStyle(strs.weekdays[n], left, bold))
	}
	for n := 0; n < int(weekStart); n++ {
		objs = append(objs, widget.NewLabelWithStyle(strs.weekdays[n], left, bold))
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

func (dtp *DateTimePicker) findMonth(month string) int {
	for n := 0; n < len(dtp.strings.months); n++ {
		if dtp.strings.months[n] == month {
			return n + 1
		}
	}
	return 0
}

// custom entry that's a bit wider
type selectEntry struct {
	widget.SelectEntry
}

func newSelectEntry(options []string) *selectEntry {
	e := &selectEntry{}
	e.ExtendBaseWidget(e)
	e.SetOptions(options)
	return e
}

func (e *selectEntry) MinSize() fyne.Size {
	o := e.SelectEntry.MinSize()
	x := widget.NewLabel("").MinSize()
	o.Width += x.Width
	return o
}

type dateTimePickerStrings struct {
	months   []string
	weekdays []string
	now      string
	time     string
}

type DateTimePicker struct {
	widget.BaseWidget
	container     *fyne.Container
	updateSelects func(time.Time)
	strings       *dateTimePickerStrings
	OnActioned    func(bool)
	when          time.Time
}

func (dtp *DateTimePicker) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(dtp.container)
}

func (dtp *DateTimePicker) initStrings() {
	dtp.strings = &dateTimePickerStrings{}
	dtp.strings.months = []string{
		lang.X("datepicker.month.january", "January"),
		lang.X("datepicker.month.february", "February"),
		lang.X("datepicker.month.march", "March"),
		lang.X("datepicker.month.april", "April"),
		lang.X("datepicker.month.may", "May"),
		lang.X("datepicker.month.june", "June"),
		lang.X("datepicker.month.july", "July"),
		lang.X("datepicker.month.august", "August"),
		lang.X("datepicker.month.september", "September"),
		lang.X("datepicker.month.october", "October"),
		lang.X("datepicker.month.november", "November"),
		lang.X("datepicker.month.december", "December"),
	}

	dtp.strings.weekdays = []string{
		lang.X("datepicker.weekday.sun", "Sun"),
		lang.X("datepicker.weekday.mon", "Mon"),
		lang.X("datepicker.weekday.tue", "Tue"),
		lang.X("datepicker.weekday.wed", "Wed"),
		lang.X("datepicker.weekday.thu", "Thu"),
		lang.X("datepicker.weekday.fri", "Fri"),
		lang.X("datepicker.weekday.sat", "Sat"),
	}

	dtp.strings.now = lang.X("datepicker.now", "Now")
	dtp.strings.time = lang.X("datepicker.time", "Time")
}

func NewDatePicker(when time.Time, weekStart time.Weekday, fn func(time.Time, bool)) *DateTimePicker {
	dtp := &DateTimePicker{}
	dtp.ExtendBaseWidget(dtp)
	dtp.initStrings()
	dtp.when = when
	dtp.OnActioned = func(ok bool) {
		fn(dtp.when, ok)
	}

	grid := container.New(layout.NewGridLayoutWithColumns(7))

	updateWhen := func(t time.Time) {
		dtp.when = t
	}

	monthSelect := widget.NewSelect(dtp.strings.months, func(selected string) {
		i := dtp.findMonth(selected)
		if i == 0 {
			return
		}

		dtp.when = timeJumpYearMonth(dtp.when, dtp.when.Year(), i)

		updateGrid(grid, dtp.when, weekStart, updateWhen, dtp.updateSelects, dtp.strings)
	})
	monthSelect.Selected = dtp.strings.months[dtp.when.Month()-1]

	years := []string{}
	// inverted years, most recent on top for easy selection
	for n := dtp.when.Year() + 10; n >= dtp.when.Year()-100; n-- {
		years = append(years, fmt.Sprintf("%d", n))
	}
	yearSelect := widget.NewSelect(years, func(selected string) {
		i, err := strconv.ParseInt(selected, 10, 64)
		if err != nil {
			return
		}

		dtp.when = timeJumpYearMonth(dtp.when, int(i), int(dtp.when.Month()))

		updateGrid(grid, dtp.when, weekStart, updateWhen, dtp.updateSelects, dtp.strings)
	})
	yearSelect.Selected = fmt.Sprintf("%d", dtp.when.Year())

	dtp.updateSelects = func(t time.Time) {
		// directly assign instead of setter methods to avoid multiple updates
		monthSelect.Selected = t.Month().String()
		monthSelect.Refresh()
		yearSelect.Selected = fmt.Sprintf("%d", t.Year())
		yearSelect.Refresh()

		updateGrid(grid, t, weekStart, updateWhen, dtp.updateSelects, dtp.strings)
	}

	prevMonthButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		dtp.when = timeJumpMonth(dtp.when, -1)

		dtp.updateSelects(dtp.when)
	})

	nextMonthButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		dtp.when = timeJumpMonth(dtp.when, 1)

		dtp.updateSelects(dtp.when)
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

	updateGrid(grid, dtp.when, weekStart, updateWhen, dtp.updateSelects, dtp.strings)

	dtp.container = container.New(
		layout.NewBorderLayout(top, nil, nil, nil),
		top,
		grid,
	)

	return dtp
}

func NewDateTimePicker(inWhen time.Time, weekStart time.Weekday, fn func(time.Time, bool)) *DateTimePicker {
	dtp := NewDatePicker(inWhen, weekStart, fn)

	hours := []string{}
	for n := 0; n <= 23; n++ {
		hours = append(hours, fmt.Sprintf("%02d", n))
	}
	hourInput := newSelectEntry(hours)

	minutes := []string{}
	for n := 0; n <= 59; n++ {
		minutes = append(minutes, fmt.Sprintf("%02d", n))
	}
	minuteInput := newSelectEntry(minutes)

	controlButtons := container.New(
		layout.NewHBoxLayout(),
		widget.NewButton("Now", func() {
			dtp.when = time.Now()

			hourInput.SetText(dtp.when.Format("15"))
			minuteInput.SetText(dtp.when.Format("04"))

			dtp.updateSelects(dtp.when)
		}),
	)

	hourInput.SetText(dtp.when.Format("15"))
	hourInput.OnChanged = func(str string) {
		t, err := time.Parse("15", str)
		if err != nil {
			fyne.LogError("invalid hour", err)
			return
		}

		dtp.when = time.Date(
			dtp.when.Year(),
			dtp.when.Month(),
			dtp.when.Day(),
			t.Hour(),
			dtp.when.Minute(),
			0,
			0,
			dtp.when.Location(),
		)
		dtp.OnActioned(true)
	}

	minuteInput.SetText(dtp.when.Format("04"))
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
		dtp.when = time.Date(
			dtp.when.Year(),
			dtp.when.Month(),
			dtp.when.Day(),
			dtp.when.Hour(),
			int(i),
			0,
			0,
			dtp.when.Location(),
		)
		dtp.OnActioned(true)
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

	top := dtp.container.Objects[0]
	grid := dtp.container.Objects[1]

	dtp.container = container.New(
		layout.NewBorderLayout(top, bottom, nil, nil),
		top,
		grid,
		bottom,
	)

	return dtp
}
