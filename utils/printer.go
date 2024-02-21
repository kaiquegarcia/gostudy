package utils

import "fmt"

var DefaultPrinter = NewPrinterByFunction(fmt.Printf)

type PrinterFunc func(format string, arguments ...interface{}) (int, error)

type Printer interface {
	Printf(format string, arguments ...interface{}) (int, error)
}

func NewPrinterByFunction(f PrinterFunc) Printer {
	return &printerByFunc{f: f}
}

type printerByFunc struct {
	f PrinterFunc
}

func (pf *printerByFunc) Printf(format string, arguments ...interface{}) (int, error) {
	return pf.f(format, arguments...)
}
