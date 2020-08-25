package tview

import (
	"github.com/gdamore/tcell"
)

// Checkbox implements a simple box for boolean values which can be checked and
// unchecked.
//
// See https://github.com/rivo/tview/wiki/Checkbox for an example.
type MultiCheckbox struct {
	*Box

	// Whether or not this box is checked.
	checked uint32

	bits uint

	focusedbit uint

	// The text to be displayed before the input area.
	label string

	// The screen width of the label area. A value of 0 means use the width of
	// the label text.
	labelWidth int

	// The label color.
	labelColor tcell.Color

	// The background color of the input area.
	fieldBackgroundColor tcell.Color

	// The text color of the input area.
	fieldTextColor tcell.Color

	// An optional function which is called when the user changes the checked
	// state of this checkbox.
	changed func(checked uint32)

	// An optional function which is called when the user indicated that they
	// are done entering text. The key which was pressed is provided (tab,
	// shift-tab, or escape).
	done func(tcell.Key)

	// A callback function set by the Form class and called when the user leaves
	// this form item.
	finished func(tcell.Key)
}

// NewCheckbox returns a new input field.
func NewMultiCheckbox() *MultiCheckbox {
	return &MultiCheckbox{
		Box:                  NewBox(),
		labelColor:           Styles.SecondaryTextColor,
		fieldBackgroundColor: Styles.ContrastBackgroundColor,
		fieldTextColor:       Styles.PrimaryTextColor,
	}
}

// SetChecked sets the state of the checkbox.
func (c *MultiCheckbox) SetChecked(checked uint32) *MultiCheckbox {
	c.checked = checked
	return c
}

// IsChecked returns whether or not the box is checked.
func (c *MultiCheckbox) IsChecked() uint32 {
	return c.checked
}

// SetLabel sets the text to be displayed before the input area.
func (c *MultiCheckbox) SetLabel(label string) *MultiCheckbox {
	c.label = label
	return c
}

// GetLabel returns the text to be displayed before the input area.
func (c *MultiCheckbox) GetLabel() string {
	return c.label
}

// SetLabelWidth sets the screen width of the label. A value of 0 will cause the
// primitive to use the width of the label string.
func (c *MultiCheckbox) SetLabelWidth(width int) *MultiCheckbox {
	c.labelWidth = width
	return c
}

// SetLabelColor sets the color of the label.
func (c *MultiCheckbox) SetLabelColor(color tcell.Color) *MultiCheckbox {
	c.labelColor = color
	return c
}

// SetFieldBackgroundColor sets the background color of the input area.
func (c *MultiCheckbox) SetFieldBackgroundColor(color tcell.Color) *MultiCheckbox {
	c.fieldBackgroundColor = color
	return c
}

// SetFieldTextColor sets the text color of the input area.
func (c *MultiCheckbox) SetFieldTextColor(color tcell.Color) *MultiCheckbox {
	c.fieldTextColor = color
	return c
}

// SetFormAttributes sets attributes shared by all form items.
func (c *MultiCheckbox) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) FormItem {
	c.labelWidth = labelWidth
	c.labelColor = labelColor
	c.backgroundColor = bgColor
	c.fieldTextColor = fieldTextColor
	c.fieldBackgroundColor = fieldBgColor
	return c
}

// GetFieldWidth returns this primitive's field width.
func (c *MultiCheckbox) GetFieldWidth() int {
	return int(c.bits)
}

func (c *MultiCheckbox) SetBits(bits int) *MultiCheckbox {
	c.bits = uint(bits)
	return c
}

// SetChangedFunc sets a handler which is called when the checked state of this
// checkbox was changed by the user. The handler function receives the new
// state.
func (c *MultiCheckbox) SetChangedFunc(handler func(checked uint32)) *MultiCheckbox {
	c.changed = handler
	return c
}

// SetDoneFunc sets a handler which is called when the user is done using the
// checkbox. The callback function is provided with the key that was pressed,
// which is one of the following:
//
//   - KeyEscape: Abort text input.
//   - KeyTab: Move to the next field.
//   - KeyBacktab: Move to the previous field.
func (c *MultiCheckbox) SetDoneFunc(handler func(key tcell.Key)) *MultiCheckbox {
	c.done = handler
	return c
}

// SetFinishedFunc sets a callback invoked when the user leaves this form item.
func (c *MultiCheckbox) SetFinishedFunc(handler func(key tcell.Key)) FormItem {
	c.finished = handler
	return c
}

// Draw draws this primitive onto the screen.
func (c *MultiCheckbox) Draw(screen tcell.Screen) {
	c.Box.Draw(screen)

	// Prepare
	x, y, width, height := c.GetInnerRect()
	rightLimit := x + width
	if height < 1 || rightLimit <= x {
		return
	}

	// Draw label.
	if c.labelWidth > 0 {
		labelWidth := c.labelWidth
		if labelWidth > rightLimit-x {
			labelWidth = rightLimit - x
		}
		Print(screen, c.label, x, y, labelWidth, AlignLeft, c.labelColor)
		x += labelWidth
	} else {
		_, drawnWidth := Print(screen, c.label, x, y, rightLimit-x, AlignLeft, c.labelColor)
		x += drawnWidth
	}

	// Draw checkboxs
	var i uint

	for i = 0; i < c.bits; i++ {
		fieldStyle := tcell.StyleDefault.Background(c.fieldBackgroundColor).Foreground(c.fieldTextColor)
		if c.focus.HasFocus() && c.focusedbit == i {
			fieldStyle = fieldStyle.Background(c.fieldTextColor).Foreground(c.fieldBackgroundColor)
		}
		checkedRune := 'X'
		if (c.checked & (0x1 << i)) == 0 {
			checkedRune = ' '
		}
		screen.SetContent(x+int(i), y, checkedRune, nil, fieldStyle)
	}
}

// InputHandler returns the handler for this primitive.
func (c *MultiCheckbox) InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive)) {
	return c.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p Primitive)) {
		// Process key event.
		switch key := event.Key(); key {
		case tcell.KeyRune, tcell.KeyEnter: // Check.
			if key == tcell.KeyRune && event.Rune() != ' ' {
				break
			}
			c.checked ^= (0x1 << c.focusedbit)
			if c.changed != nil {
				c.changed(c.checked)
			}
		case tcell.KeyTab, tcell.KeyBacktab, tcell.KeyEscape, tcell.KeyUp, tcell.KeyDown: // We're done.
			if c.done != nil {
				c.done(key)
			}
			if c.finished != nil {
				c.finished(key)
			}
		case tcell.KeyLeft:
			c.focusedbit--
			if c.focusedbit > c.bits {
				c.focusedbit = c.bits - 1
			}
		case tcell.KeyRight:
			c.focusedbit++
			if c.focusedbit > c.bits {
				c.focusedbit = 0
			}
		}
	})
}

// MouseHandler returns the mouse handler for this primitive.
func (c *MultiCheckbox) MouseHandler() func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
	return c.WrapMouseHandler(func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
		x, y := event.Position()
		_, rectY, _, _ := c.GetInnerRect()
		if !c.InRect(x, y) {
			return false, nil
		}

		// Process mouse event.
		if action == MouseLeftClick && y == rectY {
			setFocus(c)
			c.checked ^= (0x1 << c.focusedbit)
			if c.changed != nil {
				c.changed(c.checked)
			}
			consumed = true
		}

		return
	})
}
