package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/zsdevX/DarkEye/osint/social"
	"sync"
)

func restFul() {
	ctrl := &MainController{}
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.EnableGzip = true

	beego.Router("/", ctrl)
	beego.Router("/search", ctrl, "post:Search")
	beego.Run(":8080")
}

type MainController struct {
	OsInt *OsInt
	beego.Controller
	L sync.Mutex
}

func (ctrl *MainController) Get() {
	ctrl.TplName = "hello_world.html"
	ctrl.Data["name"] = "/search"
	_ = ctrl.Render()
}

func (ctrl *MainController) Search() {
	req, ok := ctrl.parseReqOK()
	if !ok {
		return
	}
	so, err := social.New(nil, req)
	if err != nil {
		ctrl.http200(-1, err.Error())
		return
	}
	if err = so.Profile(req); err != nil {
		ctrl.http200(-1, err.Error())
		return
	}
	if err = so.Follow(req); err != nil {
		ctrl.http200(-1, err.Error())
		return
	}
	if err = so.Follower(req); err != nil {
		ctrl.http200(-1, err.Error())
		return
	}
	if err = osIntRuntimeOptions.updateGraph(so); err != nil {
		ctrl.http200(-1, err.Error())
		return
	}
	ctrl.http200(0, "查询成功")
}

func (ctrl *MainController) parseReqOK() (*social.Request, bool) {
	req := social.Request{}
	if err := json.Unmarshal(ctrl.Ctx.Input.RequestBody, &req); err != nil {
		ctrl.http200(-1, err.Error())
		return nil, false
	}
	return &req, true
}

func (ctrl *MainController) http200(code int, message string) () {
	ctrl.Data["json"] =  fmt.Sprintf(`{"success":%d,"message":"%s"}`, code, message)
	ctrl.ServeJSON()
}
