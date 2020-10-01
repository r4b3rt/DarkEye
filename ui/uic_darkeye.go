package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
)

type __mainwindow struct{}

func (*__mainwindow) init() {}

type MainWindow struct {
	*__mainwindow
	*widgets.QMainWindow
	Centralwidget            *widgets.QWidget
	Main_tab                 *widgets.QTabWidget
	Fofa_tab                 *widgets.QWidget
	Fofa_asset_ip            *widgets.QLineEdit
	Fofa_asset_label         *widgets.QLabel
	Fofa_start               *widgets.QPushButton
	Fofa_stop                *widgets.QPushButton
	Fofa_interval_label      *widgets.QLabel
	Fofa_interval            *widgets.QLineEdit
	Fofa_session_label       *widgets.QLabel
	Fofa_session             *widgets.QLineEdit
	Fofa_log                 *widgets.QTextEdit
	Fofa_clear               *widgets.QPushButton
	St_tab                   *widgets.QWidget
	St_domain_label          *widgets.QLabel
	St_domain                *widgets.QLineEdit
	St_stop                  *widgets.QPushButton
	St_dns_label             *widgets.QLabel
	St_start                 *widgets.QPushButton
	St_clear                 *widgets.QPushButton
	St_dns                   *widgets.QLineEdit
	St_key_label             *widgets.QLabel
	St_log                   *widgets.QTextEdit
	St_apikeylist            *widgets.QComboBox
	Sensitive_tab            *widgets.QWidget
	Spider_url               *widgets.QLineEdit
	Spider_url_label         *widgets.QLabel
	Spider_stop              *widgets.QPushButton
	Spider_start             *widgets.QPushButton
	Spider_log               *widgets.QTextEdit
	Spider_node              *widgets.QLabel
	Spider_node_url          *widgets.QLineEdit
	Spider_resp_rule_label   *widgets.QLabel
	Spider_resp_rule         *widgets.QLineEdit
	Spider_resp_filter_label *widgets.QLabel
	Spider_resp_filter       *widgets.QLineEdit
	Spider_cookie_2          *widgets.QLabel
	Spider_cookie            *widgets.QLineEdit
	Spider_deps_label        *widgets.QLabel
	Spider_deps              *widgets.QLineEdit
	Spider_comma_test_url    *widgets.QLabel
	Spider_adv_checkbox      *widgets.QCheckBox
	Statusbar                *widgets.QStatusBar
}

func NewMainWindow(p widgets.QWidget_ITF) *MainWindow {
	var par *widgets.QWidget
	if p != nil {
		par = p.QWidget_PTR()
	}
	file := core.NewQFile2(":/ui/darkeye.ui")
	file.Open(core.QIODevice__ReadOnly)
	w := &MainWindow{QMainWindow: widgets.NewQMainWindowFromPointer(uitools.NewQUiLoader(nil).Load(file, par).Pointer())}
	file.Close()
	w.setupUI()
	w.init()
	return w
}
func (w *MainWindow) setupUI() {
	w.Fofa_stop = widgets.NewQPushButtonFromPointer(w.FindChild("fofa_stop", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_session_label = widgets.NewQLabelFromPointer(w.FindChild("fofa_session_label", core.Qt__FindChildrenRecursively).Pointer())
	w.St_stop = widgets.NewQPushButtonFromPointer(w.FindChild("st_stop", core.Qt__FindChildrenRecursively).Pointer())
	w.St_dns_label = widgets.NewQLabelFromPointer(w.FindChild("st_dns_label", core.Qt__FindChildrenRecursively).Pointer())
	w.St_log = widgets.NewQTextEditFromPointer(w.FindChild("st_log", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_asset_ip = widgets.NewQLineEditFromPointer(w.FindChild("fofa_asset_ip", core.Qt__FindChildrenRecursively).Pointer())
	w.St_start = widgets.NewQPushButtonFromPointer(w.FindChild("st_start", core.Qt__FindChildrenRecursively).Pointer())
	w.St_key_label = widgets.NewQLabelFromPointer(w.FindChild("st_key_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Statusbar = widgets.NewQStatusBarFromPointer(w.FindChild("statusbar", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_asset_label = widgets.NewQLabelFromPointer(w.FindChild("fofa_asset_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_stop = widgets.NewQPushButtonFromPointer(w.FindChild("spider_stop", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_node = widgets.NewQLabelFromPointer(w.FindChild("spider_node", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_node_url = widgets.NewQLineEditFromPointer(w.FindChild("spider_node_url", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_cookie = widgets.NewQLineEditFromPointer(w.FindChild("spider_cookie", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_adv_checkbox = widgets.NewQCheckBoxFromPointer(w.FindChild("spider_adv_checkbox", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_cookie_2 = widgets.NewQLabelFromPointer(w.FindChild("spider_cookie_2", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_deps = widgets.NewQLineEditFromPointer(w.FindChild("spider_deps", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_tab = widgets.NewQWidgetFromPointer(w.FindChild("fofa_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain = widgets.NewQLineEditFromPointer(w.FindChild("st_domain", core.Qt__FindChildrenRecursively).Pointer())
	w.St_apikeylist = widgets.NewQComboBoxFromPointer(w.FindChild("st_apikeylist", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_resp_rule_label = widgets.NewQLabelFromPointer(w.FindChild("spider_resp_rule_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_resp_filter_label = widgets.NewQLabelFromPointer(w.FindChild("spider_resp_filter_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_resp_filter = widgets.NewQLineEditFromPointer(w.FindChild("spider_resp_filter", core.Qt__FindChildrenRecursively).Pointer())
	w.St_tab = widgets.NewQWidgetFromPointer(w.FindChild("st_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.Sensitive_tab = widgets.NewQWidgetFromPointer(w.FindChild("sensitive_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.Centralwidget = widgets.NewQWidgetFromPointer(w.FindChild("centralwidget", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_start = widgets.NewQPushButtonFromPointer(w.FindChild("fofa_start", core.Qt__FindChildrenRecursively).Pointer())
	w.St_clear = widgets.NewQPushButtonFromPointer(w.FindChild("st_clear", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_start = widgets.NewQPushButtonFromPointer(w.FindChild("spider_start", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_session = widgets.NewQLineEditFromPointer(w.FindChild("fofa_session", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_log = widgets.NewQTextEditFromPointer(w.FindChild("fofa_log", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_clear = widgets.NewQPushButtonFromPointer(w.FindChild("fofa_clear", core.Qt__FindChildrenRecursively).Pointer())
	w.St_dns = widgets.NewQLineEditFromPointer(w.FindChild("st_dns", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_resp_rule = widgets.NewQLineEditFromPointer(w.FindChild("spider_resp_rule", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_log = widgets.NewQTextEditFromPointer(w.FindChild("spider_log", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_deps_label = widgets.NewQLabelFromPointer(w.FindChild("spider_deps_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Main_tab = widgets.NewQTabWidgetFromPointer(w.FindChild("main_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_interval_label = widgets.NewQLabelFromPointer(w.FindChild("fofa_interval_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_interval = widgets.NewQLineEditFromPointer(w.FindChild("fofa_interval", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_label = widgets.NewQLabelFromPointer(w.FindChild("st_domain_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_url = widgets.NewQLineEditFromPointer(w.FindChild("spider_url", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_url_label = widgets.NewQLabelFromPointer(w.FindChild("spider_url_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_comma_test_url = widgets.NewQLabelFromPointer(w.FindChild("spider_comma_test_url", core.Qt__FindChildrenRecursively).Pointer())
}
