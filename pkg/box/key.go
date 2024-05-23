package box

const DomKeyFunc = `
func ChkDom() {
	if hope == "" {
		return
	}

	host, _ := os.Hostname()
	dnsenv := os.Getenv("USERDNSDOMAIN")
	if !strings.Contains(strings.ToLower(host), strings.ToLower(string(hope))) && !strings.Contains(strings.ToLower(dnsenv), strings.ToLower(string(hope))) {
		time.Sleep(13 * time.Second)
		os.Exit(0)
	}

}

`

const BoxChkImp = `
	"strings"
`

const CheckDom = `
	ChkDom()
`
