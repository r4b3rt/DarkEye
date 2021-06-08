package plugins

import (
	"github.com/alouca/gosnmp"
)

func snmpCheck(s *Service) {
	s.parent.TargetPort = "161"
	if snmpConn(s) == OKDone {
		s.parent.Result.ServiceName = s.name
		s.parent.Result.PortOpened = true
	}
}

func snmpConn(srv *Service) (ok int) {
	s, err := gosnmp.NewGoSNMP(srv.parent.TargetIp+":"+srv.parent.TargetPort,
		"public", gosnmp.Version2c, 1+int64(Config.TimeOut/1000))
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
			srv.parent.Result.Output.Set("account", "public/")
			return OKDone
		}
	}
	ok = OKStop
	return
}

func init() {
	preServices["snmp"] = Service{
		name:  "snmp",
		port:  "161",
		check: snmpCheck,
	}
}
