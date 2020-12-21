package plugins

import (
	"github.com/alouca/gosnmp"
)

func init() {
	preCheckFuncs[SnmpPre] = snmpCheck
	supportPlugin["snmp"] = "snmp"
}

func snmpCheck(plg *Plugins) {
	plg.TargetPort = "161"
	if snmpConn(plg) == OKDone {
		plg.TargetProtocol = "snmp"
		plg.PortOpened = true
	}
}

func snmpConn(plg *Plugins) (ok int) {
	s, err := gosnmp.NewGoSNMP(plg.TargetIp+":"+plg.TargetPort,
		"public", gosnmp.Version2c, 1+int64(plg.TimeOut/1000))
	if err != nil {
		return OKStop
	}
	resp, err := s.Get(".1.3.6.1.2.1.1.5.0")
	if err != nil {
		return OKStop
	}
	for _, v := range resp.Variables {
		switch v.Type {
		case gosnmp.OctetString:
			ck := Account{
				Username: "public",
			}
			ck.Desc = v.Value.(string)
			plg.Lock()
			plg.Cracked = append(plg.Cracked, ck)
			plg.Unlock()
			return OKDone
		}
	}
	ok = OKStop
	return
}
