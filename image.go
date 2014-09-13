package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

const (
	RetinaWidth  = 72
	RetinaHeight = 72
)

type Image struct {
	name string
}

func (img *Image) Name() (name string) {
	return img.name
}

func (img *Image) Capture() (err error) {
	img.name = fmt.Sprintf("/tmp/gozo-image-%d-%d.png", os.Getuid(), os.Getpid())

	// TODO: サブシェルで Ctrl+C しても終わらないな
	err = exec.Command("screencapture", "-i", img.name).Run()
	return
}

func (img *Image) GetProperty(kind string) (result int, err error) {
	b, err := exec.Command("sips", "-g", kind, img.name).Output()
	if err != nil {
		return
	}

	re, _ := regexp.Compile(kind + ": (\\d+)")
	m := re.FindSubmatch(b)

	result, err = strconv.Atoi(string(m[1]))
	return
}

func (img *Image) PixelWidth() (width int, err error) {
	return img.GetProperty("pixelWidth")
}

func (img *Image) DpiWidth() (width int, err error) {
	return img.GetProperty("dpiWidth")
}

func (img *Image) DpiHeight() (height int, err error) {
	return img.GetProperty("dpiHeight")
}

func (img *Image) RemoveProfile() (err error) {
	err = exec.Command(
		"sips",
		"-d",
		"profile",
		"--deleteColorManagementProperties",
		img.name,
	).Run()
	return
}

func (img *Image) DownScale() (err error) {
	width, err := img.PixelWidth()
	if err != nil {
		return
	}
	err = exec.Command("sips", "--resampleWidth", strconv.Itoa(width/2), img.name).Run()
	return
}

func (img *Image) IsRetinaSize() (result bool, err error) {
	w, err := img.DpiWidth()
	if err != nil {
		return
	}

	h, err := img.DpiHeight()
	if err != nil {
		return
	}

	result = w > RetinaWidth && h > RetinaHeight
	return
}

func (img *Image) ToPng() (err error) {
	err = exec.Command("sips", "-s", "format", "png", img.name).Run()
	return
}

func (img *Image) Exists() bool {
	_, err := os.Stat(img.Name())
	return !os.IsNotExist(err)
}

func (img *Image) Remove() (err error) {
	return os.Remove(img.Name())
}

func CaptureImage() (img *Image, err error) {
	img = new(Image)

	err = img.Capture()
	if err != nil {
		return
	}

	if !img.Exists() {
		img = nil
		return
	}

	retina, err := img.IsRetinaSize()
	if err != nil {
		return
	}

	if retina {
		if err = img.DownScale(); err != nil {
			return
		}
	}

	err = img.RemoveProfile()
	if err != nil {
		return
	}
	return
}
