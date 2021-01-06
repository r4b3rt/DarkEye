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
	Centralwidget              *widgets.QWidget
	GridLayout_9               *widgets.QGridLayout
	Main_tab                   *widgets.QTabWidget
	Fofa_tab                   *widgets.QWidget
	GridLayout_3               *widgets.QGridLayout
	Zoomeye_radioButton        *widgets.QRadioButton
	Fofa_log                   *widgets.QTextEdit
	Frame                      *widgets.QFrame
	GridLayout_2               *widgets.QGridLayout
	Zoomeye_asset_label        *widgets.QLabel
	Zoomeuye_search            *widgets.QLineEdit
	Zoomeye_asset_page_label   *widgets.QLabel
	Zoomeuye_page              *widgets.QLineEdit
	Zoomeye_asset_page_label_2 *widgets.QLabel
	Zoomeye_assetkey_label     *widgets.QLabel
	Zoomeye_apikey             *widgets.QLineEdit
	Fofa_radioButton           *widgets.QRadioButton
	Frame_2                    *widgets.QFrame
	GridLayout                 *widgets.QGridLayout
	Fofa_asset_label           *widgets.QLabel
	Fofa_asset_ip              *widgets.QLineEdit
	Fofa_interval_label        *widgets.QLabel
	Fofa_interval              *widgets.QLineEdit
	Fofa_session_label         *widgets.QLabel
	Fofa_session               *widgets.QLineEdit
	Fofa_start                 *widgets.QPushButton
	Fofa_stop                  *widgets.QPushButton
	Fofa_clear                 *widgets.QPushButton
	St_tab                     *widgets.QWidget
	GridLayout_6               *widgets.QGridLayout
	Frame_3                    *widgets.QFrame
	GridLayout_4               *widgets.QGridLayout
	St_domain_brute            *widgets.QLineEdit
	St_domain_brute_rate_label *widgets.QLabel
	St_domain_brute_rate       *widgets.QLineEdit
	St_domain                  *widgets.QLineEdit
	St_domain_brute_mode_label *widgets.QLabel
	St_domain_label            *widgets.QLabel
	St_dns_label               *widgets.QLabel
	St_dns                     *widgets.QLineEdit
	St_domain_brute_mode       *widgets.QCheckBox
	St_domain_brute_label      *widgets.QLabel
	St_key_label               *widgets.QLabel
	St_apikeylist              *widgets.QComboBox
	St_log                     *widgets.QTextEdit
	St_start                   *widgets.QPushButton
	St_stop                    *widgets.QPushButton
	St_clear                   *widgets.QPushButton
	Sensitive_tab              *widgets.QWidget
	GridLayout_5               *widgets.QGridLayout
	Spider_stop                *widgets.QPushButton
	Spider_start               *widgets.QPushButton
	Spider_log                 *widgets.QTextEdit
	Frame_4                    *widgets.QFrame
	GridLayout_7               *widgets.QGridLayout
	Spider_url_label           *widgets.QLabel
	Spider_import_urls         *widgets.QPushButton
	Spider_node                *widgets.QLabel
	Spider_node_url            *widgets.QLineEdit
	Spider_deps_label          *widgets.QLabel
	Spider_deps                *widgets.QLineEdit
	Spider_cookie_label        *widgets.QLabel
	Spider_cookie              *widgets.QLineEdit
	Spider_url                 *widgets.QLineEdit
	Frame_5                    *widgets.QFrame
	GridLayout_8               *widgets.QGridLayout
	Search_enable              *widgets.QCheckBox
	Label                      *widgets.QLabel
	Search_key_label           *widgets.QLabel
	Search_key                 *widgets.QLineEdit
	Search_query_label         *widgets.QLabel
	Search_query               *widgets.QLineEdit
	Statusbar                  *widgets.QStatusBar
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
	w.Frame = widgets.NewQFrameFromPointer(w.FindChild("frame", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeuye_search = widgets.NewQLineEditFromPointer(w.FindChild("zoomeuye_search", core.Qt__FindChildrenRecursively).Pointer())
	w.Frame_2 = widgets.NewQFrameFromPointer(w.FindChild("frame_2", core.Qt__FindChildrenRecursively).Pointer())
	w.Sensitive_tab = widgets.NewQWidgetFromPointer(w.FindChild("sensitive_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout_7 = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout_7", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_session_label = widgets.NewQLabelFromPointer(w.FindChild("fofa_session_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_session = widgets.NewQLineEditFromPointer(w.FindChild("fofa_session", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_start = widgets.NewQPushButtonFromPointer(w.FindChild("fofa_start", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_rate_label = widgets.NewQLabelFromPointer(w.FindChild("st_domain_brute_rate_label", core.Qt__FindChildrenRecursively).Pointer())
	w.St_start = widgets.NewQPushButtonFromPointer(w.FindChild("st_start", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_stop = widgets.NewQPushButtonFromPointer(w.FindChild("spider_stop", core.Qt__FindChildrenRecursively).Pointer())
	w.Label = widgets.NewQLabelFromPointer(w.FindChild("label", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout_8 = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout_8", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_interval = widgets.NewQLineEditFromPointer(w.FindChild("fofa_interval", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute = widgets.NewQLineEditFromPointer(w.FindChild("st_domain_brute", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_log = widgets.NewQTextEditFromPointer(w.FindChild("spider_log", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_url_label = widgets.NewQLabelFromPointer(w.FindChild("spider_url_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_deps_label = widgets.NewQLabelFromPointer(w.FindChild("spider_deps_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_cookie = widgets.NewQLineEditFromPointer(w.FindChild("spider_cookie", core.Qt__FindChildrenRecursively).Pointer())
	w.Main_tab = widgets.NewQTabWidgetFromPointer(w.FindChild("main_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout_6 = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout_6", core.Qt__FindChildrenRecursively).Pointer())
	w.St_dns = widgets.NewQLineEditFromPointer(w.FindChild("st_dns", core.Qt__FindChildrenRecursively).Pointer())
	w.Frame_4 = widgets.NewQFrameFromPointer(w.FindChild("frame_4", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_url = widgets.NewQLineEditFromPointer(w.FindChild("spider_url", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout_3 = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout_3", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_clear = widgets.NewQPushButtonFromPointer(w.FindChild("fofa_clear", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout_4 = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout_4", core.Qt__FindChildrenRecursively).Pointer())
	w.St_key_label = widgets.NewQLabelFromPointer(w.FindChild("st_key_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_apikey = widgets.NewQLineEditFromPointer(w.FindChild("zoomeye_apikey", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_asset_label = widgets.NewQLabelFromPointer(w.FindChild("zoomeye_asset_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_interval_label = widgets.NewQLabelFromPointer(w.FindChild("fofa_interval_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_deps = widgets.NewQLineEditFromPointer(w.FindChild("spider_deps", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_cookie_label = widgets.NewQLabelFromPointer(w.FindChild("spider_cookie_label", core.Qt__FindChildrenRecursively).Pointer())
	w.St_apikeylist = widgets.NewQComboBoxFromPointer(w.FindChild("st_apikeylist", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_key_label = widgets.NewQLabelFromPointer(w.FindChild("search_key_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_key = widgets.NewQLineEditFromPointer(w.FindChild("search_key", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_radioButton = widgets.NewQRadioButtonFromPointer(w.FindChild("zoomeye_radioButton", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_asset_page_label = widgets.NewQLabelFromPointer(w.FindChild("zoomeye_asset_page_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_asset_page_label_2 = widgets.NewQLabelFromPointer(w.FindChild("zoomeye_asset_page_label_2", core.Qt__FindChildrenRecursively).Pointer())
	w.St_tab = widgets.NewQWidgetFromPointer(w.FindChild("st_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_rate = widgets.NewQLineEditFromPointer(w.FindChild("st_domain_brute_rate", core.Qt__FindChildrenRecursively).Pointer())
	w.St_stop = widgets.NewQPushButtonFromPointer(w.FindChild("st_stop", core.Qt__FindChildrenRecursively).Pointer())
	w.St_clear = widgets.NewQPushButtonFromPointer(w.FindChild("st_clear", core.Qt__FindChildrenRecursively).Pointer())
	w.Statusbar = widgets.NewQStatusBarFromPointer(w.FindChild("statusbar", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_log = widgets.NewQTextEditFromPointer(w.FindChild("fofa_log", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout_2 = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout_2", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_assetkey_label = widgets.NewQLabelFromPointer(w.FindChild("zoomeye_assetkey_label", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain = widgets.NewQLineEditFromPointer(w.FindChild("st_domain", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_mode_label = widgets.NewQLabelFromPointer(w.FindChild("st_domain_brute_mode_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_node = widgets.NewQLabelFromPointer(w.FindChild("spider_node", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_enable = widgets.NewQCheckBoxFromPointer(w.FindChild("search_enable", core.Qt__FindChildrenRecursively).Pointer())
	w.Frame_3 = widgets.NewQFrameFromPointer(w.FindChild("frame_3", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_label = widgets.NewQLabelFromPointer(w.FindChild("st_domain_brute_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_import_urls = widgets.NewQPushButtonFromPointer(w.FindChild("spider_import_urls", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_node_url = widgets.NewQLineEditFromPointer(w.FindChild("spider_node_url", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_mode = widgets.NewQCheckBoxFromPointer(w.FindChild("st_domain_brute_mode", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_start = widgets.NewQPushButtonFromPointer(w.FindChild("spider_start", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_tab = widgets.NewQWidgetFromPointer(w.FindChild("fofa_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.St_dns_label = widgets.NewQLabelFromPointer(w.FindChild("st_dns_label", core.Qt__FindChildrenRecursively).Pointer())
	w.St_log = widgets.NewQTextEditFromPointer(w.FindChild("st_log", core.Qt__FindChildrenRecursively).Pointer())
	w.Centralwidget = widgets.NewQWidgetFromPointer(w.FindChild("centralwidget", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeuye_page = widgets.NewQLineEditFromPointer(w.FindChild("zoomeuye_page", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_stop = widgets.NewQPushButtonFromPointer(w.FindChild("fofa_stop", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_label = widgets.NewQLabelFromPointer(w.FindChild("st_domain_label", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout_5 = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout_5", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_asset_label = widgets.NewQLabelFromPointer(w.FindChild("fofa_asset_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_asset_ip = widgets.NewQLineEditFromPointer(w.FindChild("fofa_asset_ip", core.Qt__FindChildrenRecursively).Pointer())
	w.Frame_5 = widgets.NewQFrameFromPointer(w.FindChild("frame_5", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout_9 = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout_9", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_radioButton = widgets.NewQRadioButtonFromPointer(w.FindChild("fofa_radioButton", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_query_label = widgets.NewQLabelFromPointer(w.FindChild("search_query_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_query = widgets.NewQLineEditFromPointer(w.FindChild("search_query", core.Qt__FindChildrenRecursively).Pointer())
}
