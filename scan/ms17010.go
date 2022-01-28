package scan

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

var (
	negotiateProtocolRequest, _  = hex.DecodeString(strings.ReplaceAll("zzzzzz85ff534d4272zzzzzzzz1853czzzzzzzzzzzzzzzzzzzzzzzzzzzzzfffezzzz4zzzzz62zzz25z432z4e4554574f524b2z5z524f4752414d2z312e3zzzz24c414e4d414e312e3zzzz257696e646f77732z666f722z576f726b67726f757z732z332e3161zzz24c4d312e32583z3z32zzz24c414e4d414e322e31zzz24e542z4c4d2z3z2e3132zz","z","0"))
	sessionSetupRequest, _       = hex.DecodeString(strings.ReplaceAll("zzzzzz88ff534d4273zzzzzzzz18z7czzzzzzzzzzzzzzzzzzzzzzzzzzzzzfffezzzz4zzzzdffzz88zzz411zazzzzzzzzzzzzzzz1zzzzzzzzzzzzzzd4zzzzzz4bzzzzzzzzzzzz57zz69zz6ezz64zz6fzz77zz73zz2zzz32zz3zzz3zzz3zzz2zzz32zz31zz39zz35zzzzzz57zz69zz6ezz64zz6fzz77zz73zz2zzz32zz3zzz3zzz3zzz2zzz35zz2ezz3zzzzzzz","z","0"))
	treeConnectRequest, _        = hex.DecodeString(strings.ReplaceAll("zzzzzz6zff534d4275zzzzzzzz18z7czzzzzzzzzzzzzzzzzzzzzzzzzzzzzfffezzz84zzzz4ffzz6zzzz8zzz1zz35zzzz5czz5czz31zz39zz32zz2ezz31zz36zz38zz2ezz31zz37zz35zz2ezz31zz32zz38zz5czz49zz5zzz43zz24zzzzzz3f3f3f3f3fzz","z","0"))
	transNamedPipeRequest, _     = hex.DecodeString(strings.ReplaceAll("zzzzzz4aff534d4225zzzzzzzz18z128zzzzzzzzzzzzzzzzzzzzzzzzzzz88ea3z1z852981zzzzzzzzzffffffffzzzzzzzzzzzzzzzzzzzzzzzz4azzzzzz4azzz2zz23zzzzzzz7zz5c5z495z455czz","z","0"))
	trans2SessionSetupRequest, _ = hex.DecodeString(strings.ReplaceAll("zzzzzz4eff534d4232zzzzzzzz18z7czzzzzzzzzzzzzzzzzzzzzzzzzzzz8fffezzz841zzzfzczzzzzzz1zzzzzzzzzzzzzza6d9a4zzzzzzzczz42zzzzzz4ezzz1zzzezzzdzzzzzzzzzzzzzzzzzzzzzzzzzzzz","z","0"))
)

type ScanStatus string

const (
	statusUnknown    = ScanStatus("?")
	statusVulnerable = ScanStatus("+")
	statusBackdored  = ScanStatus("!")
)

type Target struct {
	IP      string
	Netmask string
}

type Result struct {
	Netmask string
	IP      string
	Text    string
	Error   error
	Status  ScanStatus
}

func (s *discovery)scanMs17010(ip string, result *nb) {
	conn, err := net.DialTimeout("tcp",
		net.JoinHostPort(ip,"445"), time.Millisecond * time.Duration(s.timeout))
	if err != nil {
		return
	}

	conn.SetDeadline(time.Now().Add(time.Second * 10))
	conn.Write(negotiateProtocolRequest)
	reply := make([]byte, 1024)
	if n, err := conn.Read(reply); err != nil || n < 36 {
		return
	}

	if binary.LittleEndian.Uint32(reply[9:13]) != 0 {
		return
	}

	conn.Write(sessionSetupRequest)

	n, err := conn.Read(reply)
	if err != nil || n < 36 {
		return
	}

	if binary.LittleEndian.Uint32(reply[9:13]) != 0 {
		return
	}

	var os string
	sessionSetupResponse := reply[36:n]
	if wordCount := sessionSetupResponse[0]; wordCount != 0 {
		byteCount := binary.LittleEndian.Uint16(sessionSetupResponse[7:9])
		if n != int(byteCount)+45 {
			logrus.Debug(ip, ":invalid session setup AndX response")
		} else {
			for i := 10; i < len(sessionSetupResponse)-1; i++ {
				if sessionSetupResponse[i] == 0 && sessionSetupResponse[i+1] == 0 {
					os = string(sessionSetupResponse[10:i])
					break
				}
			}
		}
	}
	result.Os = strings.Replace(os, "\x00", "", -1)
	userID := reply[32:34]
	treeConnectRequest[32] = userID[0]
	treeConnectRequest[33] = userID[1]
	conn.Write(treeConnectRequest)

	if n, err := conn.Read(reply); err != nil || n < 36 {
		return
	}

	treeID := reply[28:30]
	transNamedPipeRequest[28] = treeID[0]
	transNamedPipeRequest[29] = treeID[1]
	transNamedPipeRequest[32] = userID[0]
	transNamedPipeRequest[33] = userID[1]

	conn.Write(transNamedPipeRequest)
	if n, err := conn.Read(reply); err != nil || n < 36 {
		return
	}

	if reply[9] == 0x05 && reply[10] == 0x02 && reply[11] == 0x00 && reply[12] == 0xc0 {
		result.vulDesc = fmt.Sprintf("Seems vulnerable for MS17-010. Operation System: %s.", result.Os)

		trans2SessionSetupRequest[28] = treeID[0]
		trans2SessionSetupRequest[29] = treeID[1]
		trans2SessionSetupRequest[32] = userID[0]
		trans2SessionSetupRequest[33] = userID[1]

		conn.Write(trans2SessionSetupRequest)

		if n, err := conn.Read(reply); err != nil || n < 36 {
			return
		}

		if reply[34] == 0x51 {
			result.vulDesc += fmt.Sprintf(" Seems to be infected by DoublePulsar.")
		}
	}
}

