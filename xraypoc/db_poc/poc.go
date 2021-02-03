package db_poc

type Poc struct {
	Name string
	Data string
}

var POCS = []Poc{
	{
		Name: "shiro.yml",
		Data: `name: poc-yaml-shiro
rules:
  - method: GET
    path: /
    headers:
      Cookie: rememberMe=1
    expression: |
      response.headers["Set-Cookie"].contains("rememberMe=deleteMe")
detail:
  author: laura_lion
  links:
    - https://github.com/Laura0xiaoshizi/xray_pocs/blob/main/shiro.yml
`},
}
