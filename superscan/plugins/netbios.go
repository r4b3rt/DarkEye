package plugins

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func init() {
	supportPlugin["netbios"] = "netbios"
	preCheckFuncs[NetBiosPre] = nbCheck
}

func nbCheck(plg *Plugins) {
	plg.TargetPort = "137"
	plg.DescCallback(fmt.Sprintf("Cracking %s %s:%s",
		"netbios", plg.TargetIp, plg.TargetPort))

	plg.RateWait(plg.RateLimiter) //爆破限制
	addr := fmt.Sprintf("%s:%s", plg.TargetIp, plg.TargetPort)
	socket, err := net.Dial("udp", addr)
	if err != nil {
		return
	}
	defer socket.Close()
	if !nbSendRequest(socket, addr, nbCreateStatusRequest()) {
		return
	}
	_ = socket.SetDeadline(time.Now().Add(time.Duration(plg.TimeOut) * time.Millisecond))
	buff := make([]byte, 1500)
	rLen, err := socket.Read(buff)
	if err != nil {
		return
	}
	reply := nbParseReply(buff[0 : rLen-1])
	if len(reply.Names) == 0 && len(reply.Addresses) == 0 {
		return
	}
	// Netbios status
	if reply.Header.RecordType != 0x21 {
		return
	}
	cracked := Account{}
	plg.TargetProtocol = "netbios"
	plg.PortOpened = true
	cracked.HostName = strings.TrimSpace(strings.Replace(string(reply.HostName[:]), "\x00", "", -1))
	cracked.UserName = strings.TrimSpace(strings.Replace(string(reply.UserName[:]), "\x00", "", -1))
	defer func() {
		plg.Lock()
		plg.Cracked = append(plg.Cracked, cracked)
		plg.Unlock()
	}()
	//Netbios name
	if !nbSendRequest(socket, addr, nbCreateNameRequest(cracked.HostName)) {
		return
	}
	//Netbios NAME paring
	_ = socket.SetDeadline(time.Now().Add(time.Duration(plg.TimeOut) * time.Millisecond))
	buff = make([]byte, 1500)
	rLen, err = socket.Read(buff)
	if err != nil {
		return
	}
	reply = nbParseReply(buff[0 : rLen-1])
	if len(reply.Names) == 0 && len(reply.Addresses) == 0 {
		return
	}
	cracked.Ip = make([]string, 0)
	for _, ip := range reply.Addresses {
		cracked.Ip = append(cracked.Ip, net.IPv4(ip.Address[0], ip.Address[1], ip.Address[2], ip.Address[3]).String())
	}
}

func nbSendRequest(socket net.Conn, addr string, req []byte) (send bool) {
	_, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}
	_, err = socket.Write(req)
	if err != nil {
		return
	}
	send = true
	return
}

func nbParseReply(buff []byte) NetbiosReplyStatus {
	resp := NetbiosReplyStatus{}
	temp := bytes.NewBuffer(buff)

	_ = binary.Read(temp, binary.BigEndian, &resp.Header)
	if resp.Header.QuestionCount != 0 {
		return resp
	}
	if resp.Header.AnswerCount == 0 {
		return resp
	}

	// Names
	if resp.Header.RecordType == 0x21 {
		var rcnt uint8
		var ridx uint8
		_ = binary.Read(temp, binary.BigEndian, &rcnt)

		for ridx = 0; ridx < rcnt; ridx++ {
			name := NetbiosReplyName{}
			_ = binary.Read(temp, binary.BigEndian, &name)
			resp.Names = append(resp.Names, name)
			switch name.Type {
			case 0x20:
				resp.HostName = name.Name
			case 0x03:
				resp.UserName = name.Name
			case 0x1e:
				resp.WorkGroup = name.Name
			default:
				//color.Red(strings.TrimSpace(strings.Replace(string(name.Name[:]), "\x00", "", -1)))
			}
		}
		return resp
	}

	// Addresses
	if resp.Header.RecordType == 0x20 {
		var rIdx uint16
		for rIdx = 0; rIdx < (resp.Header.RecordLength / 6); rIdx++ {
			addr := NetbiosReplyAddress{}
			_ = binary.Read(temp, binary.BigEndian, &addr)
			resp.Addresses = append(resp.Addresses, addr)
		}
	}

	return resp
}

func nbCreateStatusRequest() []byte {
	return []byte{
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x20, 0x43, 0x4b, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x00, 0x00, 0x21, 0x00, 0x01,
	}
}

func nbCreateNameRequest(name string) []byte {
	nBytes := [16]byte{}
	copy(nBytes[0:15], []byte(strings.ToUpper(name)[:]))

	req := []byte{
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x20,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x00, 0x00, 0x20, 0x00, 0x01,
	}

	encoded := encodeNetbiosName(nBytes)
	copy(req[13:45], encoded[0:32])
	return req
}

func encodeNetbiosName(name [16]byte) [32]byte {
	encoded := [32]byte{}

	for i := 0; i < 16; i++ {
		if name[i] == 0 {
			encoded[(i*2)+0] = 'C'
			encoded[(i*2)+1] = 'A'
		} else {
			encoded[(i*2)+0] = byte((name[i] / 16) + 0x41)
			encoded[(i*2)+1] = byte((name[i] % 16) + 0x41)
		}
	}

	return encoded
}

type NetbiosReplyHeader struct {
	XID             uint16
	Flags           uint16
	QuestionCount   uint16
	AnswerCount     uint16
	AuthCount       uint16
	AdditionalCount uint16
	QuestionName    [34]byte
	RecordType      uint16
	RecordClass     uint16
	RecordTTL       uint32
	RecordLength    uint16
}

type NetbiosReplyName struct {
	Name [15]byte
	Type uint8
	Flag uint16
}

type NetbiosReplyAddress struct {
	Flag    uint16
	Address [4]uint8
}

type NetbiosReplyStatus struct {
	Header    NetbiosReplyHeader
	HostName  [15]byte
	UserName  [15]byte
	WorkGroup [15]byte
	Names     []NetbiosReplyName
	Addresses []NetbiosReplyAddress
}
