package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
)

type __zoomeye struct{}

func (*__zoomeye) init() {}

type ZoomEye struct {
	*__zoomeye
	*widgets.QDialog
	GridLayout          *widgets.QGridLayout
	RadioButtonPassword *widgets.QRadioButton
	Key                 *widgets.QLineEdit
	TabWidget           *widgets.QTabWidget
	TabParams           *widgets.QWidget
	LabelNumber         *widgets.QLabel
	LabelType           *widgets.QLabel
	ComboBoxType        *widgets.QComboBox
	SpinBoxNumber       *widgets.QSpinBox
	LabelCondition      *widgets.QLabel
	PlainTextEditSearch *widgets.QPlainTextEdit
	LineEditFacet       *widgets.QLineEdit
	Label               *widgets.QLabel
	ToolButton          *widgets.QToolButton
	TabDomain           *widgets.QWidget
	CheckBoxDomain      *widgets.QCheckBox
	LineEditDomain      *widgets.QLineEdit
	LabelDomain         *widgets.QLabel
	TextEditLog         *widgets.QTextEdit
	RadioButtonApiKey   *widgets.QRadioButton
	PushButtonStart     *widgets.QPushButton
	PushButtonStop      *widgets.QPushButton
}

func NewZoomEye(p widgets.QWidget_ITF) *ZoomEye {
	var par *widgets.QWidget
	if p != nil {
		par = p.QWidget_PTR()
	}
	file := core.NewQFile2(":/ui/zoomeye.ui")
	file.Open(core.QIODevice__ReadOnly)
	w := &ZoomEye{QDialog: widgets.NewQDialogFromPointer(uitools.NewQUiLoader(nil).Load(file, par).Pointer())}
	file.Close()
	w.setupUI()
	w.init()
	return w
}
func (w *ZoomEye) setupUI() {
	w.SpinBoxNumber = widgets.NewQSpinBoxFromPointer(w.FindChild("spinBoxNumber", core.Qt__FindChildrenRecursively).Pointer())
	w.TabDomain = widgets.NewQWidgetFromPointer(w.FindChild("tabDomain", core.Qt__FindChildrenRecursively).Pointer())
	w.CheckBoxDomain = widgets.NewQCheckBoxFromPointer(w.FindChild("checkBoxDomain", core.Qt__FindChildrenRecursively).Pointer())
	w.RadioButtonApiKey = widgets.NewQRadioButtonFromPointer(w.FindChild("radioButtonApiKey", core.Qt__FindChildrenRecursively).Pointer())
	w.PushButtonStart = widgets.NewQPushButtonFromPointer(w.FindChild("pushButtonStart", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout", core.Qt__FindChildrenRecursively).Pointer())
	w.TabWidget = widgets.NewQTabWidgetFromPointer(w.FindChild("tabWidget", core.Qt__FindChildrenRecursively).Pointer())
	w.LineEditFacet = widgets.NewQLineEditFromPointer(w.FindChild("lineEditFacet", core.Qt__FindChildrenRecursively).Pointer())
	w.ToolButton = widgets.NewQToolButtonFromPointer(w.FindChild("toolButton", core.Qt__FindChildrenRecursively).Pointer())
	w.LineEditDomain = widgets.NewQLineEditFromPointer(w.FindChild("lineEditDomain", core.Qt__FindChildrenRecursively).Pointer())
	w.RadioButtonPassword = widgets.NewQRadioButtonFromPointer(w.FindChild("radioButtonPassword", core.Qt__FindChildrenRecursively).Pointer())
	w.LabelNumber = widgets.NewQLabelFromPointer(w.FindChild("labelNumber", core.Qt__FindChildrenRecursively).Pointer())
	w.LabelCondition = widgets.NewQLabelFromPointer(w.FindChild("labelCondition", core.Qt__FindChildrenRecursively).Pointer())
	w.LabelDomain = widgets.NewQLabelFromPointer(w.FindChild("labelDomain", core.Qt__FindChildrenRecursively).Pointer())
	w.LabelType = widgets.NewQLabelFromPointer(w.FindChild("labelType", core.Qt__FindChildrenRecursively).Pointer())
	w.ComboBoxType = widgets.NewQComboBoxFromPointer(w.FindChild("comboBoxType", core.Qt__FindChildrenRecursively).Pointer())
	w.PlainTextEditSearch = widgets.NewQPlainTextEditFromPointer(w.FindChild("PlainTextEditSearch", core.Qt__FindChildrenRecursively).Pointer())
	w.Label = widgets.NewQLabelFromPointer(w.FindChild("label", core.Qt__FindChildrenRecursively).Pointer())
	w.TextEditLog = widgets.NewQTextEditFromPointer(w.FindChild("textEditLog", core.Qt__FindChildrenRecursively).Pointer())
	w.PushButtonStop = widgets.NewQPushButtonFromPointer(w.FindChild("pushButtonStop", core.Qt__FindChildrenRecursively).Pointer())
	w.Key = widgets.NewQLineEditFromPointer(w.FindChild("Key", core.Qt__FindChildrenRecursively).Pointer())
	w.TabParams = widgets.NewQWidgetFromPointer(w.FindChild("tabParams", core.Qt__FindChildrenRecursively).Pointer())
}
