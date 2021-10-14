package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
)

type __dialogxray struct{}

func (*__dialogxray) init() {}

type DialogXray struct {
	*__dialogxray
	*widgets.QDialog
	GridLayout            *widgets.QGridLayout
	Plugins               *widgets.QLabel
	PushButtonStart       *widgets.QPushButton
	LineEditPocs          *widgets.QLineEdit
	OutputLabel           *widgets.QLabel
	ComboBoxOutput        *widgets.QComboBox
	PushButtonStop        *widgets.QPushButton
	LineEditWorkDir       *widgets.QLineEdit
	TabWidget             *widgets.QTabWidget
	TabActive             *widgets.QWidget
	LabelTargetType       *widgets.QLabel
	LineEditTarget        *widgets.QLineEdit
	ComboBoxTargetType    *widgets.QComboBox
	LabelTarget           *widgets.QLabel
	LabelParamExtra       *widgets.QLabel
	LineEditExtra         *widgets.QLineEdit
	TabPassive            *widgets.QWidget
	LineEditListen        *widgets.QLineEdit
	LabelListen           *widgets.QLabel
	TabLog                *widgets.QWidget
	TextEditLog           *widgets.QTextEdit
	PushButtonPocs        *widgets.QPushButton
	PushButtonXrayWorkDir *widgets.QPushButton
	LineEditPlugins       *widgets.QLineEdit
}

func NewDialogXray(p widgets.QWidget_ITF) *DialogXray {
	var par *widgets.QWidget
	if p != nil {
		par = p.QWidget_PTR()
	}
	file := core.NewQFile2(":/ui/xray.ui")
	file.Open(core.QIODevice__ReadOnly)
	w := &DialogXray{QDialog: widgets.NewQDialogFromPointer(uitools.NewQUiLoader(nil).Load(file, par).Pointer())}
	file.Close()
	w.setupUI()
	w.init()
	return w
}
func (w *DialogXray) setupUI() {
	w.LineEditPlugins = widgets.NewQLineEditFromPointer(w.FindChild("lineEditPlugins", core.Qt__FindChildrenRecursively).Pointer())
	w.OutputLabel = widgets.NewQLabelFromPointer(w.FindChild("outputLabel", core.Qt__FindChildrenRecursively).Pointer())
	w.PushButtonStop = widgets.NewQPushButtonFromPointer(w.FindChild("pushButtonStop", core.Qt__FindChildrenRecursively).Pointer())
	w.LabelParamExtra = widgets.NewQLabelFromPointer(w.FindChild("labelParamExtra", core.Qt__FindChildrenRecursively).Pointer())
	w.TabLog = widgets.NewQWidgetFromPointer(w.FindChild("tabLog", core.Qt__FindChildrenRecursively).Pointer())
	w.PushButtonPocs = widgets.NewQPushButtonFromPointer(w.FindChild("pushButtonPocs", core.Qt__FindChildrenRecursively).Pointer())
	w.PushButtonXrayWorkDir = widgets.NewQPushButtonFromPointer(w.FindChild("pushButtonXrayWorkDir", core.Qt__FindChildrenRecursively).Pointer())
	w.Plugins = widgets.NewQLabelFromPointer(w.FindChild("plugins", core.Qt__FindChildrenRecursively).Pointer())
	w.TabWidget = widgets.NewQTabWidgetFromPointer(w.FindChild("tabWidget", core.Qt__FindChildrenRecursively).Pointer())
	w.TabActive = widgets.NewQWidgetFromPointer(w.FindChild("tabActive", core.Qt__FindChildrenRecursively).Pointer())
	w.TabPassive = widgets.NewQWidgetFromPointer(w.FindChild("tabPassive", core.Qt__FindChildrenRecursively).Pointer())
	w.LineEditExtra = widgets.NewQLineEditFromPointer(w.FindChild("lineEditExtra", core.Qt__FindChildrenRecursively).Pointer())
	w.LineEditListen = widgets.NewQLineEditFromPointer(w.FindChild("lineEditListen", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout", core.Qt__FindChildrenRecursively).Pointer())
	w.LineEditPocs = widgets.NewQLineEditFromPointer(w.FindChild("lineEditPocs", core.Qt__FindChildrenRecursively).Pointer())
	w.LineEditWorkDir = widgets.NewQLineEditFromPointer(w.FindChild("lineEditWorkDir", core.Qt__FindChildrenRecursively).Pointer())
	w.LabelTargetType = widgets.NewQLabelFromPointer(w.FindChild("labelTargetType", core.Qt__FindChildrenRecursively).Pointer())
	w.LineEditTarget = widgets.NewQLineEditFromPointer(w.FindChild("lineEditTarget", core.Qt__FindChildrenRecursively).Pointer())
	w.ComboBoxTargetType = widgets.NewQComboBoxFromPointer(w.FindChild("comboBoxTargetType", core.Qt__FindChildrenRecursively).Pointer())
	w.TextEditLog = widgets.NewQTextEditFromPointer(w.FindChild("textEditLog", core.Qt__FindChildrenRecursively).Pointer())
	w.PushButtonStart = widgets.NewQPushButtonFromPointer(w.FindChild("pushButtonStart", core.Qt__FindChildrenRecursively).Pointer())
	w.ComboBoxOutput = widgets.NewQComboBoxFromPointer(w.FindChild("comboBoxOutput", core.Qt__FindChildrenRecursively).Pointer())
	w.LabelTarget = widgets.NewQLabelFromPointer(w.FindChild("labelTarget", core.Qt__FindChildrenRecursively).Pointer())
	w.LabelListen = widgets.NewQLabelFromPointer(w.FindChild("labelListen", core.Qt__FindChildrenRecursively).Pointer())
}
