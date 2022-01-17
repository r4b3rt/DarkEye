package scan

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type discovery struct {
	timeout        int
	disco          IdType
	pingCommand    string
	pingShell      string
	pingMatch      string
	withPrivileged bool
	host           []string

	logger *logrus.Logger
}

func NewDiscovery(timeout int, disco IdType) (Scan, error) {
	s := &discovery{
		timeout: timeout,
		logger:  logrus.New(),
	}
	s.disco = disco
	if s.disco == DiscoPing {
		switch runtime.GOOS {
		case "linux":
			s.pingCommand = "ping -c 1 -w 1"
			s.pingMatch = "ttl="
			s.pingShell = "sh -c "
		case "windows":
			s.pingCommand = "ping -n 1 -w 1"
			s.pingMatch = "TTL="
			s.pingShell = "CMD /c"
		case "darwin":
			s.pingCommand = "ping -c 1 -W 1"
			s.pingMatch = ", 0.0%"
			s.pingShell = "sh -c "
		default:
			return nil, fmt.Errorf("unsupport arch %v", runtime.GOOS)
		}
		if s.pingWithPrivileged(context.Background(), "127.0.0.1") == nil {
			s.withPrivileged = true
		}
	}
	return s, nil
}

func (s *discovery) Setup(args ...interface{}) {
	for k, v := range args {
		switch v.(type) {
		case *logrus.Logger:
			s.logger = v.(*logrus.Logger)
		case []string:
			switch k {
			case 1: //args[1] host
				if x, ok := v.([]string); ok {
					s.host = x
				}
			}
		}
	}
}

func (s *discovery) Start(ctx context.Context, ip, port string) (interface{}, error) {
	switch s.disco {
	case DiscoTcp:
		c, err := dail(ctx, "tcp", net.JoinHostPort(ip, port), s.timeout)
		if err != nil {
			return nil, err
		}
		defer c.Close()
		return fmt.Sprintf("%s %s open", ip, port), nil
	case DiscoPing:
		if s.ping(ctx, ip) {
			return fmt.Sprintf("%s alive", ip), nil
		}
	case DiscoHttp:
		return s.http(ctx, ip, port)
	case DiscoNb:
		return s.nb(ctx, ip)
	default:
		return nil, fmt.Errorf("not support check %v", s.disco)
	}
	return nil, nil
}

func (s *discovery) ping(ctx context.Context, ip string) bool {
	if s.withPrivileged {
		return s.pingWithPrivileged(ctx, ip) == nil
	}
	cmd := strings.Split(s.pingShell, " ")
	c := exec.CommandContext(ctx, cmd[0], cmd[1], strings.Join([]string{s.pingCommand, ip}, " "))
	b, _ := c.Output()
	if b != nil {
		return bytes.Contains(b, []byte(s.pingMatch))
	}
	return false
}

func (s *discovery) pingWithPrivileged(ctx context.Context, ip string) error {
	data := []byte{8, 0, 247, 255, 0, 0, 0, 0}
	c, err := dail(ctx, "ip4:icmp", ip, s.timeout)
	if err != nil {
		return err
	}
	defer c.Close()
	if _, err = c.Write(data); err != nil {
		return err
	}
	_ = c.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(s.timeout)))
	r := make([]byte, 1024)
	_, err = c.Read(r)
	if err != nil {
		return err
	}
	return nil
}

func (s *discovery) Identify(_ context.Context, _, _ string) bool {
	return true
}

func (s *discovery) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *discovery) Output() interface{} {
	return nil
}
