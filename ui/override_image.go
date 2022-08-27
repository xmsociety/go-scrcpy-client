package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"go-scrcpy-client/scrcpy"
	"image"
	"log"
)

type OverRideImageWidget struct {
	widget.BaseWidget
	Image  *canvas.Image
	Client *scrcpy.Client
}

func (o *OverRideImageWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(o.Image)
}

func NewOverRideImageWidget(img image.Image, client *scrcpy.Client) *OverRideImageWidget {
	w := &OverRideImageWidget{Image: canvas.NewImageFromImage(img), Client: client}
	w.ExtendBaseWidget(w)
	return w
}

// Tapped 左键点击
func (o *OverRideImageWidget) Tapped(e *fyne.PointEvent) {
	fmt.Println("Tapped")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if o.Client.Control.ControlConn == nil {
		log.Println("o.Client.Control.ControlConn is nil")
		return
	}
	o.Client.Control.Touch(int(e.Position.X), int(e.Position.Y), scrcpy.ActionDown)
	o.Client.Control.Touch(int(e.Position.X), int(e.Position.Y), scrcpy.ActionUp)
}

// DoubleTapped 左键双击
func (o *OverRideImageWidget) DoubleTapped(e *fyne.PointEvent) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if o.Client.Control.ControlConn == nil {
		log.Println("o.Client.Control.ControlConn is nil")
		return
	}
	fmt.Println("DoubleTapped")
}

// TappedSecondary 右键点击
func (o *OverRideImageWidget) TappedSecondary(e *fyne.PointEvent) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if o.Client.Control.ControlConn == nil {
		log.Println("o.Client.Control.ControlConn is nil")
		return
	}
	fmt.Println("TappedSecondary")
}
