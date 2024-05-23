package gengo

import (
	"bytes"
	"io"
	"os"
	"text/template"

	"github.com/kopp0ut/godropit/pkg/delivery"
)

// Add in specific funcs for new dropper type.
func GenDTypeRemote(rd DtypeRemote) (string, error) {

	tmpl, err := template.New("remoteDropper").Parse(dtypeRemote)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, rd)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil

}

func GenDTypeChild(cd DtypeChild) (string, error) {
	tmpl, err := template.New("childDropper").Parse(dtypeChild)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, cd)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil

}

func (d *Dropper) writeFinalGoTemplate(writer io.Writer) error {
	tmpl, err := template.New("newDropper").Parse(GoDropperMain)
	if err != nil {
		return err
	}
	if d.Shared {
		d.C = delivery.DllImport
		d.ProcAttach = delivery.DllFunc
		d.FuncName = "SystemFunction_032"
		d.Export = "export"

	} else {
		d.FuncName = "doStuff"
		d.Export = ""
		d.C = ""
	}
	err = tmpl.Execute(writer, d)
	if err != nil {
		return err
	}
	return nil
}

func (s *Smuggler) writeFinalTemplate(writer io.Writer) error {
	tmpl, err := template.New("newDropper").Parse(smuggleMain)
	if err != nil {
		return err
	}

	s.FuncName = "doStuff"

	err = tmpl.Execute(writer, s)
	if err != nil {
		return err
	}

	return nil
}

func PrintTemplateStr(Data, tmplstr string) {
	tmpl, err := template.New("StrTemplate").Parse(tmplstr)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, Data)
	if err != nil {
		panic(err)
	}
}
