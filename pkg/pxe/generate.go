package pxe

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed ipxe.tmpl
var ipxetmpl string

type values struct {
	URL    string
	Serial string
}

func GeneratePXE(serial, url string) ([]byte, error) {
	tmpl, err := template.New("pxe").Parse(ipxetmpl)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}
	val := values{URL: url, Serial: serial}
	w := bytes.NewBuffer(nil)
	if err := tmpl.Execute(w, val); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}
	return w.Bytes(), nil
}
