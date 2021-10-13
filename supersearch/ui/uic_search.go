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
	ActionZoomEye  *widgets.QAction
	ActionAbout    *widgets.QAction
	ActionXray     *widgets.QAction
	ActionPath     *widgets.QAction
	Centralwidget  *widgets.QWidget
	GridLayout     *widgets.QGridLayout
	TextEditBanner *widgets.QTextEdit
	Menubar        *widgets.QMenuBar
	MenuOpen       *widgets.QMenu
	MenuAbout      *widgets.QMenu
	Statusbar      *widgets.QStatusBar
}

func NewMainWindow(p widgets.QWidget_ITF) *MainWindow {
	var par *widgets.QWidget
	if p != nil {
		par = p.QWidget_PTR()
	}
	file := core.NewQFile2(":/ui/search.ui")
	file.Open(core.QIODevice__ReadOnly)
	w := &MainWindow{QMainWindow: widgets.NewQMainWindowFromPointer(uitools.NewQUiLoader(nil).Load(file, par).Pointer())}
	file.Close()
	w.setupUI()
	w.init()
	return w
}
func (w *MainWindow) setupUI() {
	w.Centralwidget = widgets.NewQWidgetFromPointer(w.FindChild("centralwidget", core.Qt__FindChildrenRecursively).Pointer())
	w.TextEditBanner = widgets.NewQTextEditFromPointer(w.FindChild("textEditBanner", core.Qt__FindChildrenRecursively).Pointer())
	w.MenuOpen = widgets.NewQMenuFromPointer(w.FindChild("menuOpen", core.Qt__FindChildrenRecursively).Pointer())
	w.MenuAbout = widgets.NewQMenuFromPointer(w.FindChild("menuAbout", core.Qt__FindChildrenRecursively).Pointer())
	w.ActionZoomEye = widgets.NewQActionFromPointer(w.FindChild("actionZoomEye", core.Qt__FindChildrenRecursively).Pointer())
	w.ActionAbout = widgets.NewQActionFromPointer(w.FindChild("actionAbout", core.Qt__FindChildrenRecursively).Pointer())
	w.ActionXray = widgets.NewQActionFromPointer(w.FindChild("actionXray", core.Qt__FindChildrenRecursively).Pointer())
	w.ActionPath = widgets.NewQActionFromPointer(w.FindChild("actionPath", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout", core.Qt__FindChildrenRecursively).Pointer())
	w.Menubar = widgets.NewQMenuBarFromPointer(w.FindChild("menubar", core.Qt__FindChildrenRecursively).Pointer())
	w.Statusbar = widgets.NewQStatusBarFromPointer(w.FindChild("statusbar", core.Qt__FindChildrenRecursively).Pointer())
}
