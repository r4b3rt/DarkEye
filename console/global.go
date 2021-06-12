package main

type analysisEntity struct {
	Ip      string `mapstructure:"ip"`
	Port    string `mapstructure:"port"`
	Service string `mapstructure:"service" csv:"serv,omitempty"`
	//web
	HttpUrl    string `mapstructure:"http_url" csv:"url,omitempty"`
	HttpTitle  string `mapstructure:"http_title" csv:"title,omitempty"`
	HttpServer string `mapstructure:"http_server" csv:"server,omitempty"`
	HttpCode   string `mapstructure:"http_code" csv:"code,omitempty"`
	//host
	NetBiosHostname  string `mapstructure:"netbios_hostname" csv:"host,omitempty"`
	NetBiosDomain    string `mapstructure:"netbios_domain" csv:"domain,omitempty"`
	NetBiosUserName  string `mapstructure:"netbios_username" csv:"user,omitempty"`
	NetBiosIpAddress string `mapstructure:"netbios_ip_address" csv:"ip,omitempty"`
	SmbShared        string `mapstructure:"smb_shared" csv:"shared,omitempty"`

	//crack
	Account string `mapstructure:"account" csv:"acct,omitempty"`
	//Finger
	Finger string `mapstructure:"finger"`
	Helper string `mapstructure:"helper" csv:"hlp,omitempty"`
}
