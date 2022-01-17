package scan

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type nb struct {
	Hostname string
	Address  []string
	Username string
	Domain   []string
}

func (s *discovery) nb(parent context.Context, ip string) (interface{}, error) {
	c := net.Dialer{Timeout: time.Duration(s.timeout) * time.Microsecond}
	ctx, _ := context.WithCancel(parent)
	conn, err := c.DialContext(ctx, "udp", net.JoinHostPort(ip, "137"))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	this := ProbeNetbios{
		timeout: s.timeout,
		socket:  conn,
		nb:      NetbiosInfo{},
	}
	if err = this.SendStatusRequest(); err != nil {
		return nil, err
	}
	if err = this.ProcessReplies(); err != nil {
		return nil, err
	}
	result := nb{
		Hostname: this.trimName(string(this.nb.statusReply.HostName[:])),
		Address:  make([]string, 0),
		Domain:   make([]string, 0),
	}

	if this.nb.nameReply.Header.RecordType == 0x20 {
		for _, a := range this.nb.nameReply.Addresses {
			result.Address = append(result.Address, net.IP(a.Address[:]).String())
		}
	}

	username := this.trimName(string(this.nb.statusReply.UserName[:]))
	if len(username) > 0 {
		result.Username = username
	}
	for _, rName := range this.nb.statusReply.Names {

		tName := this.trimName(string(rName.Name[:]))
		if tName == result.Hostname {
			continue
		}
		if rName.Flag&0x0800 != 0 {
			continue
		}
		result.Domain = append(result.Domain, tName)
	}
	return &result, nil
}

type NetbiosInfo struct {
	statusRecv  time.Time
	nameSent    time.Time
	nameRecv    time.Time
	statusReply NetbiosReplyStatus
	nameReply   NetbiosReplyStatus
}

type ProbeNetbios struct {
	timeout int
	socket  net.Conn
	nb      NetbiosInfo
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
	Address [4]byte
}

type NetbiosReplyStatus struct {
	Header    NetbiosReplyHeader
	HostName  [15]byte
	UserName  [15]byte
	Names     []NetbiosReplyName
	Addresses []NetbiosReplyAddress
	HWAddr    string
}

func (this *ProbeNetbios) ProcessReplies() error {
	buff := make([]byte, 1500)
	packet := 2
	for packet > 0 {
		packet--
		_ = this.socket.SetDeadline(time.Now().Add(time.Duration(this.timeout) * time.Millisecond))
		rLen, err := this.socket.Read(buff)
		if err != nil {
			return err
		}
		reply := this.ParseReply(buff[0 : rLen-1])
		if len(reply.Names) == 0 && len(reply.Addresses) == 0 {
			continue
		}
		if reply.Header.RecordType == 0x21 {
			this.nb.statusReply = reply
			this.nb.statusRecv = time.Now()
			nTime := time.Time{}
			if this.nb.nameSent == nTime {
				this.nb.nameSent = time.Now()
				name := this.trimName(string(this.nb.statusReply.HostName[:]))
				if err = this.SendNameRequest(name); err != nil {
					return err
				}
			}
		}
		if reply.Header.RecordType == 0x20 {
			this.nb.nameReply = reply
			this.nb.nameRecv = time.Now()
			return nil
		}
	}
	return fmt.Errorf("ProcessReplies no response")
}

func (this *ProbeNetbios) SendRequest(req []byte) error {
	_ = this.socket.SetDeadline(time.Now().Add(time.Duration(this.timeout) * time.Millisecond))
	if _, err := this.socket.Write(req); err != nil {
		return err
	}
	return nil
}

func (this *ProbeNetbios) SendStatusRequest() error {
	return this.SendRequest(this.CreateStatusRequest())
}

func (this *ProbeNetbios) SendNameRequest(name string) error {
	return this.SendRequest(this.CreateNameRequest(name))
}

func (this *ProbeNetbios) EncodeNetbiosName(name [16]byte) [32]byte {
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

func (this *ProbeNetbios) DecodeNetbiosName(name [32]byte) [16]byte {
	decoded := [16]byte{}

	for i := 0; i < 16; i++ {
		if name[(i*2)+0] == 'C' && name[(i*2)+1] == 'A' {
			decoded[i] = 0
		} else {
			decoded[i] = ((name[(i*2)+0] * 16) - 0x41) + (name[(i*2)+1] - 0x41)
		}
	}
	return decoded
}

func (this *ProbeNetbios) ParseReply(buff []byte) NetbiosReplyStatus {

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
			binary.Read(temp, binary.BigEndian, &name)
			resp.Names = append(resp.Names, name)

			if name.Type == 0x20 {
				resp.HostName = name.Name
			}

			if name.Type == 0x03 {
				resp.UserName = name.Name
			}
		}

		var hwbytes [6]uint8
		_ = binary.Read(temp, binary.BigEndian, &hwbytes)
		resp.HWAddr = fmt.Sprintf("%.2x:%.2x:%.2x:%.2x:%.2x:%.2x",
			hwbytes[0], hwbytes[1], hwbytes[2], hwbytes[3], hwbytes[4], hwbytes[5],
		)
		return resp
	}

	// Addresses
	if resp.Header.RecordType == 0x20 {
		var ridx uint16
		for ridx = 0; ridx < (resp.Header.RecordLength / 6); ridx++ {
			addr := NetbiosReplyAddress{}
			_ = binary.Read(temp, binary.BigEndian, &addr)
			resp.Addresses = append(resp.Addresses, addr)
		}
	}

	return resp
}

func (this *ProbeNetbios) CreateStatusRequest() []byte {
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

func (this *ProbeNetbios) CreateNameRequest(name string) []byte {
	nbytes := [16]byte{}
	copy(nbytes[0:15], strings.ToUpper(name)[:])

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

	encoded := this.EncodeNetbiosName(nbytes)
	copy(req[13:45], encoded[0:32])
	return req
}

func (this *ProbeNetbios) trimName(name string) string {
	return strings.TrimSpace(strings.Replace(name, "\x00", "", -1))
}
