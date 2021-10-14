package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
)

type __dialogpreference struct{}

func (*__dialogpreference) init() {}

type DialogPreference struct {
	*__dialogpreference
	*widgets.QDialog
	GridLayout        *widgets.QGridLayout
	GroupBox          *widgets.QGroupBox
	GridLayout_2      *widgets.QGridLayout
	PushButtonWorkDir *widgets.QPushButton
	LineEditWorkDir   *widgets.QLineEdit
	LabelWebHook      *widgets.QLabel
	LineEditWebHook   *widgets.QLineEdit
	PushButtonOK      *widgets.QPushButton
}

func NewDialogPreference(p widgets.QWidget_ITF) *DialogPreference {
	var par *widgets.QWidget
	if p != nil {
		par = p.QWidget_PTR()
	}
	file := core.NewQFile2(":/ui/preference.ui")
	file.Open(core.QIODevice__ReadOnly)
	w := &DialogPreference{QDialog: widgets.NewQDialogFromPointer(uitools.NewQUiLoader(nil).Load(file, par).Pointer())}
	file.Close()
	w.setupUI()
	w.init()
	return w
}
func (w *DialogPreference) setupUI() {
	w.LabelWebHook = widgets.NewQLabelFromPointer(w.FindChild("labelWebHook", core.Qt__FindChildrenRecursively).Pointer())
	w.LineEditWebHook = widgets.NewQLineEditFromPointer(w.FindChild("lineEditWebHook", core.Qt__FindChildrenRecursively).Pointer())
	w.PushButtonOK = widgets.NewQPushButtonFromPointer(w.FindChild("pushButtonOK", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout", core.Qt__FindChildrenRecursively).Pointer())
	w.GroupBox = widgets.NewQGroupBoxFromPointer(w.FindChild("groupBox", core.Qt__FindChildrenRecursively).Pointer())
	w.GridLayout_2 = widgets.NewQGridLayoutFromPointer(w.FindChild("gridLayout_2", core.Qt__FindChildrenRecursively).Pointer())
	w.PushButtonWorkDir = widgets.NewQPushButtonFromPointer(w.FindChild("pushButtonWorkDir", core.Qt__FindChildrenRecursively).Pointer())
	w.LineEditWorkDir = widgets.NewQLineEditFromPointer(w.FindChild("lineEditWorkDir", core.Qt__FindChildrenRecursively).Pointer())
}
