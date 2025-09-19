package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Calculator struct {
	display     *canvas.Text
	firstNumber float64
	operation   string
	startNew    bool
}

func newCalculator() *Calculator {
	calc := &Calculator{
		display:   canvas.NewText("0", theme.ForegroundColor()),
		operation: "",
		startNew:  true,
	}

	calc.display.Alignment = fyne.TextAlignTrailing
	calc.display.TextSize = 32

	return calc
}

func (c *Calculator) setDisplay(text string) {
	c.display.Text = text
	c.display.Refresh()
}

func (c *Calculator) typeDigit(digit string) {
	if c.startNew {
		c.setDisplay(digit)
		c.startNew = false
	} else {
		current := c.display.Text
		if current == "0" {
			c.setDisplay(digit)
		} else {
			c.setDisplay(current + digit)
		}
	}
}

func (c *Calculator) typeDot() {
	if c.startNew {
		c.setDisplay("0.")
		c.startNew = false
		return
	}
	if !strings.Contains(c.display.Text, ".") {
		c.setDisplay(c.display.Text + ".")
	}
}

func (c *Calculator) setOperation(op string) {
	if c.operation != "" {
		c.calculate()
	}
	val, _ := strconv.ParseFloat(c.display.Text, 64)
	c.firstNumber = val
	c.operation = op
	c.startNew = true
}

func (c *Calculator) calculate() {
	second, _ := strconv.ParseFloat(c.display.Text, 64)
	var result float64

	switch c.operation {
	case "+":
		result = c.firstNumber + second
	case "-":
		result = c.firstNumber - second
	case "*":
		result = c.firstNumber * second
	case "/":
		if second != 0 {
			result = c.firstNumber / second
		} else {
			c.setDisplay("Erro")
			c.reset()
			return
		}
	default:
		result = second
	}

	c.setDisplay(fmt.Sprintf("%g", result))
	c.operation = ""
	c.startNew = true
}

func (c *Calculator) reset() {
	c.setDisplay("0")
	c.firstNumber = 0
	c.operation = ""
	c.startNew = true
}

func (c *Calculator) makeButton(label string, tapped func()) *widget.Button {
	btn := widget.NewButton(label, tapped)
	btn.Importance = widget.HighImportance
	return btn
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Calculadora em Go")
	myWindow.Resize(fyne.NewSize(320, 460))

	calc := newCalculator()

	buttons := container.NewGridWithColumns(4,
		calc.makeButton("C", calc.reset),
		calc.makeButton("±", func() {
			if val, err := strconv.ParseFloat(calc.display.Text, 64); err == nil {
				calc.setDisplay(fmt.Sprintf("%g", -val))
			}
		}),
		calc.makeButton("%", func() {
			if val, err := strconv.ParseFloat(calc.display.Text, 64); err == nil {
				calc.setDisplay(fmt.Sprintf("%g", val/100))
			}
		}),
		calc.makeButton("/", func() { calc.setOperation("/") }),

		calc.makeButton("7", func() { calc.typeDigit("7") }),
		calc.makeButton("8", func() { calc.typeDigit("8") }),
		calc.makeButton("9", func() { calc.typeDigit("9") }),
		calc.makeButton("*", func() { calc.setOperation("*") }),

		calc.makeButton("4", func() { calc.typeDigit("4") }),
		calc.makeButton("5", func() { calc.typeDigit("5") }),
		calc.makeButton("6", func() { calc.typeDigit("6") }),
		calc.makeButton("-", func() { calc.setOperation("-") }),

		calc.makeButton("1", func() { calc.typeDigit("1") }),
		calc.makeButton("2", func() { calc.typeDigit("2") }),
		calc.makeButton("3", func() { calc.typeDigit("3") }),
		calc.makeButton("+", func() { calc.setOperation("+") }),

		calc.makeButton("0", func() { calc.typeDigit("0") }),
		calc.makeButton(".", calc.typeDot),
		calc.makeButton("=", calc.calculate),
	)

	// Centralizar o texto no display
	displayContainer := container.NewHBox()
	displayContainer.Add(calc.display)
	displayContainer.Objects[0].(*canvas.Text).Alignment = fyne.TextAlignTrailing

	content := container.NewBorder(
		container.NewVBox(
			container.NewHBox(fyne.NewContainerWithLayout(layoutCenterRight{}, displayContainer)),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		buttons,
	)

	myWindow.SetContent(content)

	myWindow.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		key := keyEvent.Name

		switch key {
		case fyne.Key0, fyne.Key1, fyne.Key2, fyne.Key3, fyne.Key4,
			fyne.Key5, fyne.Key6, fyne.Key7, fyne.Key8, fyne.Key9:
			digit := string(key[len(key)-1])
			calc.typeDigit(digit)

		case fyne.KeyPeriod, fyne.KeyComma:
			calc.typeDot()

		case fyne.KeyPlus:
			calc.setOperation("+")
		case fyne.KeyMinus:
			calc.setOperation("-")
		case fyne.KeyAsterisk:
			calc.setOperation("*")
		case fyne.KeySlash:
			calc.setOperation("/")

		case fyne.KeyReturn, fyne.KeyEnter:
			calc.calculate()

		case fyne.KeyEscape:
			calc.reset()

		case fyne.KeyBackspace:
			current := calc.display.Text
			if len(current) > 1 {
				calc.setDisplay(current[:len(current)-1])
			} else {
				calc.setDisplay("0")
			}
		}
	})

	myWindow.ShowAndRun()
}

type layoutCenterRight struct{}

func (l layoutCenterRight) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 {
		return
	}
	obj := objects[0]
	obj.Resize(obj.MinSize())
	obj.Move(fyne.NewPos(size.Width-obj.MinSize().Width-8, (size.Height-obj.MinSize().Height)/2))
}

func (l layoutCenterRight) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	return objects[0].MinSize()
}
