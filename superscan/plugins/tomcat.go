package plugins

import "github.com/zsdevX/DarkEye/superscan/dic"

func tomcatCheck(s *Service) {
	s.user = dic.DIC_USERNAME_TOMCAT
	s.pass = dic.DIC_PASSWORD_TOMCAT
	s.name = "tomcat"
	s.parent.Result.ServiceName = s.name
	url := s.parent.Result.Web.Url
	s.parent.Result.Web.Url += "/manager/html"
	s.connect = _401AuthConn
	//爆破manager
	_401AuthCheck(s)
	s.parent.Result.Web.Url = url
}
