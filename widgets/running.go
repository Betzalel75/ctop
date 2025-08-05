package widgets

import (
	"github.com/Betzalel75/ctop/cwidgets/compact"
	ui "github.com/gizak/termui"
)

type RunningContainers struct {
	*ui.Block
	Grid   *compact.CompactGrid
	Header *CTopHeader
}

func NewRunningContainers(grid *compact.CompactGrid, header *CTopHeader) *RunningContainers {
	block := ui.NewBlock()
	block.BorderLabel = "Running"
	block.Border = true
	
	return &RunningContainers{
		Block:  block,
		Grid:   grid,
		Header: header,
	}
}

func (r *RunningContainers) Buffer() ui.Buffer {
	buf := r.Block.Buffer()
	
	// Calculer les positions internes
	innerY := r.Y + 1 // +1 pour la bordure
	innerWidth := r.Width - 2 // -2 pour les bordures gauche/droite
	innerX := r.X + 1 // +1 pour la bordure gauche
	
	// Positionner et rendre le header si activé
	if r.Header != nil {
		// Ajuster la position X du header pour qu'il soit dans le cadre
		r.Header.Time.X = innerX + 1
		r.Header.Count.X = innerX + 23
		r.Header.Filter.X = innerX + 39
		r.Header.bg.X = innerX
		r.Header.bg.SetWidth(innerWidth)
		
		buf.Merge(r.Header.Buffer())
		innerY += r.Header.Height()
	}
	
	// Positionner et rendre la grille
	r.Grid.Y = innerY
	r.Grid.SetWidth(innerWidth)
	r.Grid.Align()
	buf.Merge(r.Grid.Buffer())
	
	return buf
}

func (r *RunningContainers) Align() {
	// Méthode pour s'adapter aux changements de taille d'écran
	r.Width = ui.TermWidth()
	r.Height = ui.TermHeight() - 1 // -1 pour la status line
	
	// Réaligner les composants internes
	if r.Header != nil {
		r.Header.Align()
	}
	
	innerWidth := r.Width - 2
	r.Grid.SetWidth(innerWidth)
	r.Grid.Align()
}
