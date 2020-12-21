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
	Main_tab                   *widgets.QTabWidget
	Fofa_tab                   *widgets.QWidget
	Fofa_log                   *widgets.QTextEdit
	Splitter_4                 *widgets.QSplitter
	Fofa_start                 *widgets.QPushButton
	Fofa_stop                  *widgets.QPushButton
	Fofa_clear                 *widgets.QPushButton
	Asset_toolBox              *widgets.QToolBox
	Page_5                     *widgets.QWidget
	Widget                     *widgets.QWidget
	HorizontalLayout           *widgets.QHBoxLayout
	Zoomeye_asset_label        *widgets.QLabel
	Zoomeuye_search            *widgets.QLineEdit
	Zoomeye_assetkey_label     *widgets.QLabel
	Zoomeye_key                *widgets.QComboBox
	Page_6                     *widgets.QWidget
	Widget1                    *widgets.QWidget
	HorizontalLayout_2         *widgets.QHBoxLayout
	Fofa_asset_label           *widgets.QLabel
	Fofa_asset_ip              *widgets.QLineEdit
	Fofa_interval_label        *widgets.QLabel
	Fofa_interval              *widgets.QLineEdit
	Fofa_session_label         *widgets.QLabel
	Fofa_session               *widgets.QLineEdit
	Splitter_3                 *widgets.QSplitter
	Zoomeye_radioButton        *widgets.QRadioButton
	Fofa_radioButton           *widgets.QRadioButton
	St_tab                     *widgets.QWidget
	St_log                     *widgets.QTextEdit
	Splitter                   *widgets.QSplitter
	St_start                   *widgets.QPushButton
	St_stop                    *widgets.QPushButton
	St_clear                   *widgets.QPushButton
	Widget2                    *widgets.QWidget
	HorizontalLayout_3         *widgets.QHBoxLayout
	St_domain_label            *widgets.QLabel
	St_domain                  *widgets.QLineEdit
	St_dns_label               *widgets.QLabel
	St_dns                     *widgets.QLineEdit
	St_domain_brute_mode_label *widgets.QLabel
	St_domain_brute_mode       *widgets.QCheckBox
	Widget3                    *widgets.QWidget
	HorizontalLayout_4         *widgets.QHBoxLayout
	St_domain_brute_label      *widgets.QLabel
	St_domain_brute            *widgets.QLineEdit
	St_domain_brute_rate_label *widgets.QLabel
	St_domain_brute_rate       *widgets.QLineEdit
	St_key_label               *widgets.QLabel
	St_apikeylist              *widgets.QComboBox
	Sensitive_tab              *widgets.QWidget
	Spider_log                 *widgets.QTextEdit
	ToolBox_spider             *widgets.QToolBox
	Page                       *widgets.QWidget
	Spider_cookie              *widgets.QLineEdit
	Spider_cookie_label        *widgets.QLabel
	Splitter_2                 *widgets.QSplitter
	Spider_url_label           *widgets.QLabel
	Spider_url                 *widgets.QLineEdit
	Spider_import_urls         *widgets.QPushButton
	Spider_node                *widgets.QLabel
	Spider_node_url            *widgets.QLineEdit
	Spider_deps_label          *widgets.QLabel
	Spider_deps                *widgets.QLineEdit
	Page_2                     *widgets.QWidget
	Search_enable              *widgets.QCheckBox
	Widget4                    *widgets.QWidget
	HorizontalLayout_5         *widgets.QHBoxLayout
	Search_key_label           *widgets.QLabel
	Search_key                 *widgets.QLineEdit
	Search_query_label         *widgets.QLabel
	Search_query               *widgets.QLineEdit
	Splitter_7                 *widgets.QSplitter
	Spider_start               *widgets.QPushButton
	Spider_stop                *widgets.QPushButton
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
	w.Spider_deps_label = widgets.NewQLabelFromPointer(w.FindChild("spider_deps_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_query_label = widgets.NewQLabelFromPointer(w.FindChild("search_query_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_stop = widgets.NewQPushButtonFromPointer(w.FindChild("spider_stop", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_start = widgets.NewQPushButtonFromPointer(w.FindChild("fofa_start", core.Qt__FindChildrenRecursively).Pointer())
	w.St_stop = widgets.NewQPushButtonFromPointer(w.FindChild("st_stop", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_rate_label = widgets.NewQLabelFromPointer(w.FindChild("st_domain_brute_rate_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Page = widgets.NewQWidgetFromPointer(w.FindChild("page", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_deps = widgets.NewQLineEditFromPointer(w.FindChild("spider_deps", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_key = widgets.NewQLineEditFromPointer(w.FindChild("search_key", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_query = widgets.NewQLineEditFromPointer(w.FindChild("search_query", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_stop = widgets.NewQPushButtonFromPointer(w.FindChild("fofa_stop", core.Qt__FindChildrenRecursively).Pointer())
	w.HorizontalLayout = widgets.NewQHBoxLayoutFromPointer(w.FindChild("horizontalLayout", core.Qt__FindChildrenRecursively).Pointer())
	w.Widget1 = widgets.NewQWidgetFromPointer(w.FindChild("widget1", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_log = widgets.NewQTextEditFromPointer(w.FindChild("spider_log", core.Qt__FindChildrenRecursively).Pointer())
	w.HorizontalLayout_4 = widgets.NewQHBoxLayoutFromPointer(w.FindChild("horizontalLayout_4", core.Qt__FindChildrenRecursively).Pointer())
	w.Main_tab = widgets.NewQTabWidgetFromPointer(w.FindChild("main_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.Widget = widgets.NewQWidgetFromPointer(w.FindChild("widget", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_asset_ip = widgets.NewQLineEditFromPointer(w.FindChild("fofa_asset_ip", core.Qt__FindChildrenRecursively).Pointer())
	w.St_start = widgets.NewQPushButtonFromPointer(w.FindChild("st_start", core.Qt__FindChildrenRecursively).Pointer())
	w.St_key_label = widgets.NewQLabelFromPointer(w.FindChild("st_key_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_cookie = widgets.NewQLineEditFromPointer(w.FindChild("spider_cookie", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_key_label = widgets.NewQLabelFromPointer(w.FindChild("search_key_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_asset_label = widgets.NewQLabelFromPointer(w.FindChild("fofa_asset_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_radioButton = widgets.NewQRadioButtonFromPointer(w.FindChild("fofa_radioButton", core.Qt__FindChildrenRecursively).Pointer())
	w.St_clear = widgets.NewQPushButtonFromPointer(w.FindChild("st_clear", core.Qt__FindChildrenRecursively).Pointer())
	w.St_dns_label = widgets.NewQLabelFromPointer(w.FindChild("st_dns_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Search_enable = widgets.NewQCheckBoxFromPointer(w.FindChild("search_enable", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_tab = widgets.NewQWidgetFromPointer(w.FindChild("fofa_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.HorizontalLayout_3 = widgets.NewQHBoxLayoutFromPointer(w.FindChild("horizontalLayout_3", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_rate = widgets.NewQLineEditFromPointer(w.FindChild("st_domain_brute_rate", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_url_label = widgets.NewQLabelFromPointer(w.FindChild("spider_url_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_cookie_label = widgets.NewQLabelFromPointer(w.FindChild("spider_cookie_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_node = widgets.NewQLabelFromPointer(w.FindChild("spider_node", core.Qt__FindChildrenRecursively).Pointer())
	w.Splitter_7 = widgets.NewQSplitterFromPointer(w.FindChild("splitter_7", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_interval = widgets.NewQLineEditFromPointer(w.FindChild("fofa_interval", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_session = widgets.NewQLineEditFromPointer(w.FindChild("fofa_session", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_mode = widgets.NewQCheckBoxFromPointer(w.FindChild("st_domain_brute_mode", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_label = widgets.NewQLabelFromPointer(w.FindChild("st_domain_brute_label", core.Qt__FindChildrenRecursively).Pointer())
	w.St_log = widgets.NewQTextEditFromPointer(w.FindChild("st_log", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute_mode_label = widgets.NewQLabelFromPointer(w.FindChild("st_domain_brute_mode_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Splitter_4 = widgets.NewQSplitterFromPointer(w.FindChild("splitter_4", core.Qt__FindChildrenRecursively).Pointer())
	w.HorizontalLayout_2 = widgets.NewQHBoxLayoutFromPointer(w.FindChild("horizontalLayout_2", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_interval_label = widgets.NewQLabelFromPointer(w.FindChild("fofa_interval_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_session_label = widgets.NewQLabelFromPointer(w.FindChild("fofa_session_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_log = widgets.NewQTextEditFromPointer(w.FindChild("fofa_log", core.Qt__FindChildrenRecursively).Pointer())
	w.St_tab = widgets.NewQWidgetFromPointer(w.FindChild("st_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.HorizontalLayout_5 = widgets.NewQHBoxLayoutFromPointer(w.FindChild("horizontalLayout_5", core.Qt__FindChildrenRecursively).Pointer())
	w.Statusbar = widgets.NewQStatusBarFromPointer(w.FindChild("statusbar", core.Qt__FindChildrenRecursively).Pointer())
	w.Sensitive_tab = widgets.NewQWidgetFromPointer(w.FindChild("sensitive_tab", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_url = widgets.NewQLineEditFromPointer(w.FindChild("spider_url", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_node_url = widgets.NewQLineEditFromPointer(w.FindChild("spider_node_url", core.Qt__FindChildrenRecursively).Pointer())
	w.Page_5 = widgets.NewQWidgetFromPointer(w.FindChild("page_5", core.Qt__FindChildrenRecursively).Pointer())
	w.Page_6 = widgets.NewQWidgetFromPointer(w.FindChild("page_6", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_brute = widgets.NewQLineEditFromPointer(w.FindChild("st_domain_brute", core.Qt__FindChildrenRecursively).Pointer())
	w.St_apikeylist = widgets.NewQComboBoxFromPointer(w.FindChild("st_apikeylist", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_start = widgets.NewQPushButtonFromPointer(w.FindChild("spider_start", core.Qt__FindChildrenRecursively).Pointer())
	w.Splitter_3 = widgets.NewQSplitterFromPointer(w.FindChild("splitter_3", core.Qt__FindChildrenRecursively).Pointer())
	w.Widget2 = widgets.NewQWidgetFromPointer(w.FindChild("widget2", core.Qt__FindChildrenRecursively).Pointer())
	w.ToolBox_spider = widgets.NewQToolBoxFromPointer(w.FindChild("toolBox_spider", core.Qt__FindChildrenRecursively).Pointer())
	w.Spider_import_urls = widgets.NewQPushButtonFromPointer(w.FindChild("spider_import_urls", core.Qt__FindChildrenRecursively).Pointer())
	w.Splitter_2 = widgets.NewQSplitterFromPointer(w.FindChild("splitter_2", core.Qt__FindChildrenRecursively).Pointer())
	w.Page_2 = widgets.NewQWidgetFromPointer(w.FindChild("page_2", core.Qt__FindChildrenRecursively).Pointer())
	w.Fofa_clear = widgets.NewQPushButtonFromPointer(w.FindChild("fofa_clear", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_key = widgets.NewQComboBoxFromPointer(w.FindChild("zoomeye_key", core.Qt__FindChildrenRecursively).Pointer())
	w.Splitter = widgets.NewQSplitterFromPointer(w.FindChild("splitter", core.Qt__FindChildrenRecursively).Pointer())
	w.St_dns = widgets.NewQLineEditFromPointer(w.FindChild("st_dns", core.Qt__FindChildrenRecursively).Pointer())
	w.Widget4 = widgets.NewQWidgetFromPointer(w.FindChild("widget4", core.Qt__FindChildrenRecursively).Pointer())
	w.Asset_toolBox = widgets.NewQToolBoxFromPointer(w.FindChild("asset_toolBox", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_asset_label = widgets.NewQLabelFromPointer(w.FindChild("zoomeye_asset_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeuye_search = widgets.NewQLineEditFromPointer(w.FindChild("zoomeuye_search", core.Qt__FindChildrenRecursively).Pointer())
	w.Widget3 = widgets.NewQWidgetFromPointer(w.FindChild("widget3", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain = widgets.NewQLineEditFromPointer(w.FindChild("st_domain", core.Qt__FindChildrenRecursively).Pointer())
	w.Centralwidget = widgets.NewQWidgetFromPointer(w.FindChild("centralwidget", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_assetkey_label = widgets.NewQLabelFromPointer(w.FindChild("zoomeye_assetkey_label", core.Qt__FindChildrenRecursively).Pointer())
	w.Zoomeye_radioButton = widgets.NewQRadioButtonFromPointer(w.FindChild("zoomeye_radioButton", core.Qt__FindChildrenRecursively).Pointer())
	w.St_domain_label = widgets.NewQLabelFromPointer(w.FindChild("st_domain_label", core.Qt__FindChildrenRecursively).Pointer())
}
