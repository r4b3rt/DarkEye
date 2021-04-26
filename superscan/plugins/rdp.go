package plugins

/*
#ifndef DNO_RDP_SUPPORT
#include <freerdp/gdi/gdi.h>
#include <freerdp/freerdp.h>

int rdp_connect(char *server, char *port, char *domain, char *login, char *password) {

  int32_t err = 0;
  freerdp *instance = 0;
  wLog *root = WLog_GetRoot();
  WLog_SetStringLogLevel(root, "OFF");

  instance = freerdp_new();
  if (instance == NULL || freerdp_context_new(instance) == FALSE) {
    return -1;
  }

  instance->settings->Username = login;
  instance->settings->Password = password;
  instance->settings->IgnoreCertificate = TRUE;
  instance->settings->AuthenticationOnly = TRUE;
  instance->settings->ServerHostname = server;
  instance->settings->ServerPort = atoi(port);
  instance->settings->Domain = domain;
  freerdp_connect(instance);
  err = freerdp_get_last_error(instance->context);
  return err;
}
#else
int rdp_connect(char *server, char *port, char *domain, char *login, char *password) {
	return -1;
}
#endif
*/
import "C"
import (
	"context"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"unsafe"
)

func rdpCheck(s *Service) {
	s.crack()
}

func RdpConn(_ context.Context, s *Service, user, pass string) (ok int) {
	username := C.CString(user)
	password := C.CString(pass)
	server := C.CString(s.parent.TargetIp)
	port := C.CString(s.parent.TargetPort)
	domain := C.CString("")

	defer func() {
		C.free(unsafe.Pointer(username))
		C.free(unsafe.Pointer(password))
		C.free(unsafe.Pointer(domain))
		C.free(unsafe.Pointer(port))
		C.free(unsafe.Pointer(server))

	}()
	ret := C.rdp_connect(server, port, domain, username, password)
	switch ret {
	case 0:
		// login success
		return OKDone
	case 0x00020009:
	case 0x00020014:
	case 0x00020015:
		// login failure
		return OKNext
	case 0x0002000d:
		return OKNext
	case 0x00020006:
	case 0x00020008:
	case 0x0002000c:
		return OKStop
	default:
		return OKTerm
	}
	return OKTerm
}
func init() {
	services["rdp"] = Service{
		name:    "rdp",
		port:    "3389",
		user:    dic.DIC_USERNAME_RDP,
		pass:    dic.DIC_PASSWORD_RDP,
		check:   rdpCheck,
		connect: RdpConn,
		thread:  1,
	}
}
