package check

import (
	"net"
	"strings"
)

func ChkMAC() bool {

	EvidenceOfSandbox := make([]net.HardwareAddr, 0)

	badMacAddresses := [...]string{`00:0C:29`, `00:1C:14`, `00:50:56`, `00:05:69`, `08:00:27`}

	NICs, _ := net.Interfaces()
	for _, NIC := range NICs {
		for _, badMacAddress := range badMacAddresses {
			if strings.Contains(strings.ToLower(NIC.HardwareAddr.String()), strings.ToLower(badMacAddress)) {
				EvidenceOfSandbox = append(EvidenceOfSandbox, NIC.HardwareAddr)
			}
		}
	}

	if len(EvidenceOfSandbox) == 0 {
		return true
	} else {
		return false
	}

}
