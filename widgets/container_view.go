package widgets

import (
	"github.com/Betzalel75/ctop/connector"
	"github.com/Betzalel75/ctop/cwidgets/compact"
	ui "github.com/gizak/termui"
)

type ContainerView struct {
	*ui.Block
	RunningWidget  *RunningContainers
	AllWidget      *AllContainers
	activeWidget   string // "running" ou "all"
	Grid           *compact.CompactGrid
	Header         *CTopHeader
	connectorSuper *connector.ConnectorSuper
}

func NewContainerView(grid *compact.CompactGrid,
	header *CTopHeader, cSuper *connector.ConnectorSuper) *ContainerView {
	cv := &ContainerView{
		Block:          ui.NewBlock(),
		Grid:           grid,
		Header:         header,
		activeWidget:   "running", // Par défaut sur "running"
		connectorSuper: cSuper,
	}

	cv.Border = false
	cv.RunningWidget = NewRunningContainers(grid, header)
	cv.AllWidget = NewAllContainers(grid, header, cSuper)

	return cv
}

func (cv *ContainerView) Buffer() ui.Buffer {
	buf := cv.Block.Buffer()

	// Calculer les dimensions pour chaque widget (50% chacun)
	halfWidth := cv.Width / 2

	// Positionner le widget Running à gauche
	cv.RunningWidget.X = cv.X
	cv.RunningWidget.Y = cv.Y
	cv.RunningWidget.Width = halfWidth
	cv.RunningWidget.Height = cv.Height

	// Positionner le widget All à droite
	cv.AllWidget.X = cv.X + halfWidth
	cv.AllWidget.Y = cv.Y
	cv.AllWidget.Width = cv.Width - halfWidth
	cv.AllWidget.Height = cv.Height
	
	if cv.activeWidget == "running" {
		cv.RunningWidget.BorderFg = ui.ThemeAttr("status.ok") // Vert pour actif
		cv.RunningWidget.BorderLabel = "Running [ACTIVE]"
		cv.AllWidget.BorderFg = ui.ThemeAttr("border.fg") // Couleur normale
		cv.AllWidget.BorderLabel = "All"
	} else {
		cv.RunningWidget.BorderFg = ui.ThemeAttr("border.fg") // Couleur normale
		cv.RunningWidget.BorderLabel = "Running"
		cv.AllWidget.BorderFg = ui.ThemeAttr("status.ok") // Vert pour actif
		cv.AllWidget.BorderLabel = "All [ACTIVE]"
	}
	
	// Rendre les deux widgets
	buf.Merge(cv.RunningWidget.Buffer())
	buf.Merge(cv.AllWidget.Buffer())

	return buf
}

func (cv *ContainerView) Align() {
	cv.Width = ui.TermWidth()
	cv.Height = ui.TermHeight() - 1 // -1 pour la status line
}

func (cv *ContainerView) SwitchToRunning() {
	cv.activeWidget = "running"
}

func (cv *ContainerView) SwitchToAll() {
	cv.activeWidget = "all"
}

func (cv *ContainerView) GetActiveWidget() string {
	return cv.activeWidget
}

func (cv *ContainerView) IsRunningActive() bool {
	return cv.activeWidget == "running"
}
