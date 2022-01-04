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
	Serial     string
	Kernel     string
	Initrd     string
	Installer  string
	InitrdBits string
	Rootfs     string
	IpxeCfg    string
}

func GeneratePXE(serial, url string) ([]byte, error) {
	tmpl, err := template.New("pxe").Parse(ipxetmpl)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}
	val := values{
		Serial:     serial,
		Kernel:     fmt.Sprintf("%s%s", url, "kernel"),
		Initrd:     fmt.Sprintf("%s%s", url, "initrd.img"),
		Installer:  fmt.Sprintf("%s%s", url, "installer.img"),
		InitrdBits: fmt.Sprintf("%s%s", url, "initrd.bits"),
		Rootfs:     fmt.Sprintf("%s%s", url, "rootfs.img"),
		IpxeCfg:    fmt.Sprintf("%s%s", url, "ipxe.efi.cfg"),
	}
	w := bytes.NewBuffer(nil)
	if err := tmpl.Execute(w, val); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}
	return w.Bytes(), nil
}
