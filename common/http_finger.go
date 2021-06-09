package common

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type httpFingerPrint struct {
	Cms      string   `json:"cms"`
	Method   string   `json:"method"`
	Location string   `json:"location"`
	Keyword  []string `json:"keyword"`
}

type httpFingerSrc struct {
	Fingerprint []httpFingerPrint `json:"fingerprint"`
}

func init() {
	if err := json.Unmarshal([]byte(httpFingers), &httpFinger); err != nil {
		Log("finger.init", err.Error(), FAULT)
	}

	newFp := make([]httpFingerPrint, 0)
	for _, f := range httpFinger.Fingerprint {
		//只留存匹配方法为keyword
		if f.Method == "keyword" {
			newFp = append(newFp, f)
		}
	}
	httpFinger.Fingerprint = newFp
	logrus.Info("web finger: ", len(newFp))
}

//dataSrc: https://raw.githubusercontent.com/EASY233/Finger/main/library/finger.json
var (
	httpFinger = httpFingerSrc{
		Fingerprint: make([]httpFingerPrint, 0),
	}
	httpFingers = `
	{
"fingerprint": [

	{
		"cms": "seeyon",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/seeyon/USER-DATA/IMAGES/LOGIN/login.gif"
		]
	},
	{
		"cms": "seeyon",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/seeyon/common/"
		]
	},
	{
		"cms": "Spring env",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"servletContextInitParams"
		]
	},
	{
		"cms": "微三云管理系统",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"WSY_logo",
			"管理系统 MANAGEMENT SYSTEM"
		]
	},
	{
		"cms": "Spring env",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"logback"
		]
	},
	{
		"cms": "Weblogic",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Error 404--Not Found"
		]
	},
	{
		"cms": "Weblogic",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Error 403--"
		]
	},
	{
		"cms": "Weblogic",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/console/framework/skins/wlsconsole/images/login_WebLogic_branding.png"
		]
	},
	{
		"cms": "Weblogic",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Welcome to Weblogic Application Server"
		]
	},
	{
		"cms": "Weblogic",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"\u003ci\u003eHypertext Transfer Protocol -- HTTP/1.1\u003c/i\u003e"
		]
	},
	{
		"cms": "Sangfor SSL VPN",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/por/login_psw.csp"
		]
	},
	{
		"cms": "Sangfor SSL VPN",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"loginPageSP/loginPrivacy.js"
		]
	},
	{
		"cms": "e-mobile",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"weaver,e-mobile"
		]
	},
	{
		"cms": "ecology",
		"method": "keyword",
		"location": "header",
		"keyword": [
			"ecology_JSessionid"
		]
	},
	{
		"cms": "Shiro",
		"method": "keyword",
		"location": "header",
		"keyword": [
			"rememberMe="
		]
	},
	{
		"cms": "Shiro",
		"method": "keyword",
		"location": "header",
		"keyword": [
			"=deleteMe"
		]
	},
	{
		"cms": "泛微云桥 e-Bridge",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"wx.weaver"
		]
	},
	{
		"cms": "泛微云桥 e-Bridge",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"e-Bridge"
		]
	},
	{
		"cms": "泛微-协同办公OA",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"cloudstore/resource/pc/polyfill/polyfill.min.js"
		]
	},
	{
		"cms": "泛微-协同办公OA",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"wui/theme/ecology8/page/images/login/username_wev8.png"
		]
	},
	{
		"cms": "泛微-协同办公OA",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/wui/index.html#/?logintype=1"
		]
	},
	{
		"cms": "Swagger UI",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Swagger UI"
		]
	},
	{
		"cms": "Ruijie",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"4008 111 000"
		]
	},
	{
		"cms": "Huawei SMC",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Script/SmcScript.js?version="
		]
	},
	{
		"cms": "H3C Router",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/wnm/ssl/web/frame/login.html"
		]
	},
	{
		"cms": "Cisco SSLVPN",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/+CSCOE+/logon.html"
		]
	},
	{
		"cms": "通达OA",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/images/tongda.ico"
		]
	},
	{
		"cms": "通达OA",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Office Anywhere"
		]
	},
	{
		"cms": "通达OA",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"通达OA"
		]
	},
	{
		"cms": "深信服 waf",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"rsa.js",
			"commonFunction.js"
		]
	},
	{
		"cms": "深信服 waf",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Redirect to...",
			"/LogInOut.php"
		]
	},
	{
		"cms": "网御 vpn",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/vpn/common/js/leadsec.js",
			"/vpn/user/common/custom/auth_home.css"
		]
	},
	{
		"cms": "启明星辰天清汉马USG防火墙",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/cgi-bin/webui?op=get_product_model"
		]
	},
	{
		"cms": "蓝凌 OA",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"sys/ui/extend/theme/default/style/icon.css",
			"sys/ui/extend/theme/default/style/profile.css"
		]
	},
	{
		"cms": "深信服上网行为管理系统",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"utccjfaewjb = function(str, key)"
		]
	},
	{
		"cms": "深信服上网行为管理系统",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"document.write(WRFWWCSFBXMIGKRKHXFJ"
		]
	},
	{
		"cms": "深信服应用交付报表系统",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/reportCenter/index.php?cls_mode=cluster_mode_others"
		]
	},
	{
		"cms": "金蝶云星空",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"HTML5/content/themes/kdcss.min.css"
		]
	},
	{
		"cms": "金蝶云星空",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/ClientBin/Kingdee.BOS.XPF.App.xap"
		]
	},
	{
		"cms": "CoreMail",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"coremail/common"
		]
	},
	{
		"cms": "启明星辰天清汉马USG防火墙",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"天清汉马USG"
		]
	},
	{
		"cms": "Jboss",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"jboss.css"
		]
	},
	{
		"cms": "Gitlab",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"assets/gitlab_logo"
		]
	},
	{
		"cms": "宝塔-BT.cn",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"入口校验失败"
		]
	},
	{
		"cms": "宝塔-BT.cn",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"扫码登录更安全",
			"bt.cn",
			"/login"
		]
	},
	{
		"cms": "宝塔-BT.cn",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"站点创建成功",
			"bt.cn"
		]
	},
	{
		"cms": "禅道",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"self.location",
			"Lw=="
		]
	},
	{
		"cms": "禅道",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/theme/default/images/main/zt-logo.png"
		]
	},
	{
		"cms": "禅道",
		"method": "keyword",
		"location": "header",
		"keyword": [
			"zentaosid"
		]
	},
	{
		"cms": "用友软件",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"UFIDA Software CO.LTD all rights reserved"
		]
	},
	{
		"cms": "YONYOU NC",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"uclient.yonyou.com",
			"UClient"
		]
	},
	{
		"cms": "宝塔-BT.cn",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"宝塔Linux面板"
		]
	},
	{
		"cms": "RabbitMQ",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"\u003ctitle\u003eRabbitMQ Management\u003c/title\u003e"
		]
	},
	{
		"cms": "Zabbix",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"zabbix",
			"Zabbix SIA"
		]
	},
	{
		"cms": "联软准入",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"网络准入",
			"leagsoft",
			"redirect"
		]
	},
	{
		"cms": "列目录",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Index of /"
		]
	},
	{
		"cms": "列目录",
		"method": "keyword",
		"location": "body",
		"keyword": [
			" - /\u003c/title\u003e"
		]
	},
	{
		"cms": "浪潮服务器IPMI管理口",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"img/inspur_logo.png",
			"Management System"
		]
	},
	{
		"cms": "RegentApi_v2.0",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"RegentApi_v2.0"
		]
	},
	{
		"cms": "Tomcat默认页面",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"/manager/status",
			"/manager/html"
		]
	},
	{
		"cms": "Discuz!",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Discuz!",
			"Comsenz",
			"cache/"
		]
	},
	{
		"cms": "深信服WEB防篡改管理系统",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"WEB防篡改",
			"cgi-bin/tamper_admin.cgi"
		]
	},
	{
		"cms": "YApi 可视化接口管理平台",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"YApi",
			"id=\"yapi\"",
			"prd",
			"可视化接口管理平台"
		]
	},
	{
		"cms": "WeiPHP",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"weiphp.css",
			"weiphp",
			"Public/static"
		]
	},
	{
		"cms": "Nagios XI",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Nagios XI",
			"nagiosxi",
			"Nagios"
		]
	},
	{
		"cms": "群晖 NAS",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"DiskStation",
			"webman/modules",
			"NAS"
		]
	},
	{
		"cms": "山石网科 防火墙",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"Hillstone",
			"licenseAggrement",
			"GLOBAL_CONFIG.js"
		]
	},
	{
		"cms": "360天堤新一代智慧防火墙",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"360天堤",
			"360",
			"360防火墙"
		]
	},
	{
		"cms": "360网神防火墙系统",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"resources/image/logo_header.png",
			"360",
			"网神防火墙系统"
		]
	},
	{
		"cms": "网神SecGate 3600防火墙",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"网神SecGate",
			"3600防火墙",
			"css/lsec/login.css"
		]
	},
	{
		"cms": "蓝盾防火墙",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"蓝盾",
			"Bluedon",
			"default/js/act/login.js"
		]
	},
	{
		"cms": "LanProxy",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"LanProxy",
			"password",
			"lanproxy-config"
		]
	},
	{
		"cms": "ManageEngine ADManager Plus",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"ADManager",
			"Hashtable.js",
			"ManageEngine"
		]
	},
	{
		"cms": "中新金盾信息安全管理系统",
		"method": "keyword",
		"location": "body",
		"keyword": [
			"中新金盾信息安全管理系统",
			"login",
			"useusbkey"
		]
	}
]
}
`
)
