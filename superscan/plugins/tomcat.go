package plugins

import "github.com/zsdevX/DarkEye/superscan/dic"

func tomcatCheck(plg *Plugins) {
	//爆破manager
	plg.tmp.tls = plg.Web.Tls
	plg.tmp.urlPath = "/manager/html"
	plg.TargetProtocol = "tomcat"
	_401AuthCheck(plg, dic.DIC_USERNAME_TOMCAT, dic.DIC_PASSWORD_TOMCAT)
}
