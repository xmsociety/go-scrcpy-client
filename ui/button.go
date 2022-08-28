package ui

import "fyne.io/fyne/v2/widget"

func NewButton(label string, tapped func()) *widget.Button {
	button := &widget.Button{
		Text:       label,
		OnTapped:   tapped,
		Importance: widget.HighImportance,
	}

	button.ExtendBaseWidget(button)
	return button
}
