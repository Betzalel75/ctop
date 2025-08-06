package widgets

import (
	"fmt"

	"github.com/Betzalel75/ctop/connector"
	"github.com/Betzalel75/ctop/cwidgets/compact"
	"github.com/Betzalel75/ctop/dtop/manager"
	"github.com/Betzalel75/ctop/dtop/resource"
	"github.com/Betzalel75/ctop/dtop/utils"
	api "github.com/fsouza/go-dockerclient"
	ui "github.com/gizak/termui"
)

type AllContainers struct {
	*ui.Block
	Grid           *compact.CompactGrid
	Header         *CTopHeader
	connectorSuper *connector.ConnectorSuper
	dockerManager  *manager.DockerResourceManager
	currentMode    string // "menu", "containers", "images", "volumes", "publish"
	menuItems      []MenuItem
	selectedIndex  int
	statusMsg      string
	publishData    *PublishData
	resources      []resource.ResourceItem // Liste des ressources actuelles
	resourceType   string                  // Type de ressource actuel

	// NOUVEAU: Variables de pagination
	currentPage     int
	itemsPerPage    int
	maxDisplayItems int
}

type MenuItem struct {
	Title       string
	Description string
}

type PublishData struct {
	Registry string
	Username string
	Tag      string
	Step     int // 0=registry, 1=username, 2=tag, 3=select_images, 4=confirm
}

func NewAllContainers(grid *compact.CompactGrid,
	header *CTopHeader, cSuper *connector.ConnectorSuper) *AllContainers {
	block := ui.NewBlock()
	block.BorderLabel = "All"
	block.Border = true

	return &AllContainers{
		Block:          block,
		Grid:           grid,
		Header:         header,
		connectorSuper: cSuper,
		currentMode:    "menu",
		menuItems: []MenuItem{
			{"🐳 Publisher", ""},
			{"📦 Delete Containers", ""},
			{"💿 Delete Images", ""},
			{"💾 Delete Volumes", ""},
		},
		selectedIndex: 0,
		resources:     []resource.ResourceItem{},
		// NOUVEAU: Initialisation pagination
		currentPage:     0,
		itemsPerPage:    10, // Nombre d'items par page
		maxDisplayItems: 10, // Nombre d'items visibles à l'écran
	}
}

// NOUVEAU: Méthodes de pagination
func (a *AllContainers) getTotalPages() int {
	if len(a.resources) == 0 {
		return 0
	}
	return (len(a.resources) + a.itemsPerPage - 1) / a.itemsPerPage
}

func (a *AllContainers) getCurrentPageItems() []resource.ResourceItem {
	start := a.currentPage * a.itemsPerPage
	end := start + a.itemsPerPage
	
	if start >= len(a.resources) {
		return []resource.ResourceItem{}
	}
	
	if end > len(a.resources) {
		end = len(a.resources)
	}
	
	return a.resources[start:end]
}

func (a *AllContainers) adjustSelectedIndexForPage() {
	pageItems := a.getCurrentPageItems()
	if a.selectedIndex >= len(pageItems) {
		a.selectedIndex = len(pageItems) - 1
	}
	if a.selectedIndex < 0 {
		a.selectedIndex = 0
	}
}

func (a *AllContainers) getDockerClient() (*api.Client, error) {
	conn, err := a.connectorSuper.Get()
	if err != nil {
		return nil, err
	}

	dockerConn, ok := conn.(*connector.Docker)
	if !ok {
		return nil, fmt.Errorf("connector is not a Docker connector")
	}

	return dockerConn.GetClient(), nil
}

// Initialiser le DockerResourceManager de manière paresseuse
func (a *AllContainers) getDockerManager() (*manager.DockerResourceManager, error) {
	if a.dockerManager == nil {
		client, err := a.getDockerClient()
		if err != nil {
			return nil, err
		}
		a.dockerManager = manager.NewDockerResourceManager(client)
	}
	return a.dockerManager, nil
}

func (a *AllContainers) Buffer() ui.Buffer {
	buf := a.Block.Buffer()

	// Calculer les positions internes
	innerY := a.Y + 1
	innerWidth := a.Width - 2
	innerX := a.X + 1

	// Rendre le contenu selon le mode actuel
	switch a.currentMode {
	case "menu":
		a.renderMenu(buf, innerX, innerY, innerWidth)
	case "containers", "images", "volumes":
		a.renderResourceList(buf, innerX, innerY, innerWidth)
	case "publish":
		a.renderPublishFlow(buf, innerX, innerY, innerWidth)
	}

	// Afficher le message de statut
	if a.statusMsg != "" {
		statusY := a.Y + a.Height - 2
		for i, ch := range a.statusMsg {
			if i >= innerWidth {
				break
			}
			buf.Set(innerX+i, statusY, ui.Cell{
				Ch: ch,
				Fg: ui.ThemeAttr("status.warn"),
				Bg: ui.ColorDefault,
			})
		}
	}

	return buf
}

func (a *AllContainers) renderMenu(buf ui.Buffer, x, y, width int) {
	title := "Docker Management Menu"
	for i, ch := range title {
		if i >= width {
			break
		}
		buf.Set(x+i, y, ui.Cell{
			Ch: ch,
			Fg: ui.ThemeAttr("header.fg"),
			Bg: ui.ColorDefault,
		})
	}

	for i, item := range a.menuItems {
		itemY := y + 2 + i
		if itemY >= a.Y+a.Height-1 {
			break
		}

		// Marquer l'élément sélectionné
		prefix := "  "
		fg := ui.ThemeAttr("par.text.fg")
		if i == a.selectedIndex {
			prefix = "> "
			fg = ui.ThemeAttr("status.ok")
		}

		text := prefix + item.Title
		for j, ch := range text {
			if j >= width {
				break
			}
			buf.Set(x+j, itemY, ui.Cell{
				Ch: ch,
				Fg: fg,
				Bg: ui.ColorDefault,
			})
		}

		// Description sur la ligne suivante
		if i == a.selectedIndex && itemY+1 < a.Y+a.Height-1 {
			desc := item.Description
			for j, ch := range desc {
				if j >= width {
					break
				}
				buf.Set(x+j, itemY+1, ui.Cell{
					Ch: ch,
					Fg: ui.ThemeAttr("par.text.fg"),
					Bg: ui.ColorDefault,
				})
			}
		}
	}
}


func (a *AllContainers) renderResourceList(buf ui.Buffer, x, y, width int) {
	totalPages := a.getTotalPages()
	currentPageItems := a.getCurrentPageItems()
	
	// Titre avec info pagination
	title := fmt.Sprintf("Manage %s (Page %d/%d - %d/%d items)", 
		utils.Capitalize(a.currentMode), 
		a.currentPage+1, 
		totalPages, 
		len(currentPageItems),
		len(a.resources))
	
	for i, ch := range title {
		if i >= width {
			break
		}
		buf.Set(x+i, y, ui.Cell{
			Ch: ch,
			Fg: ui.ThemeAttr("header.fg"),
			Bg: ui.ColorDefault,
		})
	}

	// Instructions mises à jour
	instructions := "↑↓: navigate, PgUp/PgDn: pages, Space: select, d: delete, r: refresh, q: back"
	instructY := y + 1
	for i, ch := range instructions {
		if i >= width {
			break
		}
		buf.Set(x+i, instructY, ui.Cell{
			Ch: ch,
			Fg: ui.ThemeAttr("par.text.fg"),
			Bg: ui.ColorDefault,
		})
	}

	// Afficher la liste des ressources de la page courante
	startY := y + 3
	a.maxDisplayItems = (a.Height - 8) / 2 // 2 lignes par item + marges
	
	displayItems := currentPageItems
	if len(displayItems) > a.maxDisplayItems {
		displayItems = displayItems[:a.maxDisplayItems]
	}
	
	for i, res := range displayItems {
		itemY := startY + i*2
		if itemY >= a.Y+a.Height-4 {
			break
		}

		// Checkbox et titre
		checkbox := "[ ]"
		if res.Selected {
			checkbox = "[x]"
		}
		
		// Marquer l'item sélectionné avec > et couleur différente
		prefix := "  "
		fg := ui.ThemeAttr("par.text.fg")
		bg := ui.ColorDefault
		
		if i == a.selectedIndex {
			prefix = "> "
			fg = ui.ThemeAttr("status.ok")
			bg = ui.ThemeAttr("par.text.fg")
		}

		titleLine := fmt.Sprintf("%s%s %s", prefix, checkbox, res.Title)
		
		// Rendre la ligne avec background si sélectionnée
		for j := range width {
			ch := ' '
			if j < len(titleLine) {
				ch = rune(titleLine[j])
			}
			buf.Set(x+j, itemY, ui.Cell{
				Ch: ch,
				Fg: fg,
				Bg: bg,
			})
		}

		// Description sur la ligne suivante
		if itemY+1 < a.Y+a.Height-4 {
			desc := "    " + res.Desc
			descFg := ui.ThemeAttr("par.text.fg")
			descBg := ui.ColorDefault
			
			if i == a.selectedIndex {
				descBg = ui.ThemeAttr("par.text.fg")
				descFg = ui.ThemeAttr("par.text.hi")
			}
			
			for j := range width {
				ch := ' '
				if j < len(desc) {
					ch = rune(desc[j])
				}
				buf.Set(x+j, itemY+1, ui.Cell{
					Ch: ch,
					Fg: descFg,
					Bg: descBg,
				})
			}
		}
	}

	// Afficher les informations de pagination en bas
	paginationY := a.Y + a.Height - 4
	
	// Compteur de sélection
	selectedCount := 0
	for _, res := range a.resources { // Compter sur TOUTES les ressources
		if res.Selected {
			selectedCount++
		}
	}
	
	if selectedCount > 0 {
		counter := fmt.Sprintf("Selected: %d/%d items", selectedCount, len(a.resources))
		for i, ch := range counter {
			if i >= width {
				break
			}
			buf.Set(x+i, paginationY, ui.Cell{
				Ch: ch,
				Fg: ui.ThemeAttr("status.ok"),
				Bg: ui.ColorDefault,
			})
		}
	}
	
	// Info de navigation
	if totalPages > 1 {
		navInfo := fmt.Sprintf("Page %d/%d - Use PgUp/PgDn to navigate", a.currentPage+1, totalPages)
		navY := paginationY + 1
		for i, ch := range navInfo {
			if i >= width {
				break
			}
			buf.Set(x+i, navY, ui.Cell{
				Ch: ch,
				Fg: ui.ThemeAttr("par.text.fg"),
				Bg: ui.ColorDefault,
			})
		}
	}
}


func (a *AllContainers) renderPublishFlow(buf ui.Buffer, x, y, width int) {
	title := "🐳 Publish Docker Images"
	for i, ch := range title {
		if i >= width {
			break
		}
		buf.Set(x+i, y, ui.Cell{
			Ch: ch,
			Fg: ui.ThemeAttr("header.fg"),
			Bg: ui.ColorDefault,
		})
	}

	if a.publishData != nil {
		currentY := y + 2
		switch a.publishData.Step {
		case 0:
			prompt := fmt.Sprintf("Registry [%s]: %s", "docker.io", a.publishData.Registry)
			a.renderTextLine(buf, x, currentY, width, prompt)
		case 1:
			prompt := fmt.Sprintf("Username: %s", a.publishData.Username)
			a.renderTextLine(buf, x, currentY, width, prompt)
		case 2:
			prompt := fmt.Sprintf("Tag [%s]: %s", "latest", a.publishData.Tag)
			a.renderTextLine(buf, x, currentY, width, prompt)
		case 3:
			// Afficher la liste des images pour sélection
			a.renderResourceList(buf, x, y, width)
		}
	}
}

func (a *AllContainers) renderTextLine(buf ui.Buffer, x, y, width int, text string) {
	for i, ch := range text {
		if i >= width {
			break
		}
		buf.Set(x+i, y, ui.Cell{
			Ch: ch,
			Fg: ui.ThemeAttr("par.text.fg"),
			Bg: ui.ColorDefault,
		})
	}
}

func (a *AllContainers) Align() {
	a.Width = ui.TermWidth() / 2
	a.Height = ui.TermHeight() - 1
}

// Méthodes de navigation
func (a *AllContainers) HandleKey(key string) bool {
	switch a.currentMode {
	case "menu":
		return a.handleMenuKey(key)
	case "containers", "images", "volumes":
		return a.handleResourceKey(key)
	case "publish":
		return a.handlePublishKey(key)
	}
	return false
}

func (a *AllContainers) handleMenuKey(key string) bool {
	switch key {
	case "up", "k":
		if a.selectedIndex > 0 {
			a.selectedIndex--
			a.statusMsg = fmt.Sprintf("Selected: %s", a.menuItems[a.selectedIndex].Title)
		}
		return true
	case "down", "j":
		if a.selectedIndex < len(a.menuItems)-1 {
			a.selectedIndex++
			a.statusMsg = fmt.Sprintf("Selected: %s", a.menuItems[a.selectedIndex].Title)
		}
		return true
	case "enter":
		return a.selectMenuItem()
	}
	return false
}


func (a *AllContainers) handleResourceKey(key string) bool {
	currentPageItems := a.getCurrentPageItems()
	
	switch key {
	case "q", "esc":
		a.currentMode = "menu"
		a.statusMsg = ""
		a.selectedIndex = 0
		a.currentPage = 0
		return true
		
	case "up", "k":
		if a.selectedIndex > 0 {
			a.selectedIndex--
		} else if a.currentPage > 0 {
			// Aller à la page précédente et se positionner en bas
			a.currentPage--
			newPageItems := a.getCurrentPageItems()
			a.selectedIndex = min(len(newPageItems)-1, a.maxDisplayItems-1)
		}
		return true
		
	case "down", "j":
		maxIndex := min(len(currentPageItems)-1, a.maxDisplayItems-1)
		if a.selectedIndex < maxIndex {
			a.selectedIndex++
		} else if a.currentPage < a.getTotalPages()-1 {
			// Aller à la page suivante et se positionner en haut
			a.currentPage++
			a.selectedIndex = 0
		}
		return true
		
	case "pgup":
		if a.currentPage > 0 {
			a.currentPage--
			a.adjustSelectedIndexForPage()
		}
		return true
		
	case "pgdown":
		if a.currentPage < a.getTotalPages()-1 {
			a.currentPage++
			a.adjustSelectedIndexForPage()
		}
		return true
		
	case "space":
		// Toggle selection du current item
		if a.selectedIndex >= 0 && a.selectedIndex < len(currentPageItems) {
			realIndex := a.currentPage*a.itemsPerPage + a.selectedIndex
			if realIndex < len(a.resources) {
				a.resources[realIndex].Selected = !a.resources[realIndex].Selected
				selectedCount := 0
				for _, res := range a.resources {
					if res.Selected {
						selectedCount++
					}
				}
				a.statusMsg = fmt.Sprintf("%d items selected. Press 'd' to delete.", selectedCount)
			}
		}
		return true
		
	case "d":
		// Delete selected items
		var toDelete []resource.ResourceItem
		for _, res := range a.resources {
			if res.Selected {
				toDelete = append(toDelete, res)
			}
		}
		if len(toDelete) > 0 {
			a.statusMsg = fmt.Sprintf("Deleting %d items...", len(toDelete))
			go a.deleteSelectedItems(toDelete)
		} else {
			a.statusMsg = "No items selected. Press space to select items."
		}
		return true
		
	case "r": // Refresh manual
		a.statusMsg = "Refreshing..."
		a.currentPage = 0
		a.selectedIndex = 0
		a.refreshCurrentList()
		return true
	}
	return false
}


func (a *AllContainers) handlePublishKey(key string) bool {
	switch key {
	case "q", "esc":
		a.currentMode = "menu"
		a.publishData = nil
		a.statusMsg = ""
		a.selectedIndex = 0
		return true
	case "enter":
		if a.publishData != nil && a.publishData.Step == 3 {
			// Confirmer la sélection d'images
			selectedCount := 0
			for _, res := range a.resources {
				if res.Selected {
					selectedCount++
				}
			}
			if selectedCount > 0 {
				a.statusMsg = fmt.Sprintf("Ready to publish %d images. Press 'p' to confirm.", selectedCount)
				a.publishData.Step = 4
			} else {
				a.statusMsg = "No images selected. Press space to select images."
			}
		}
		return true
	case "space":
		if a.publishData != nil && a.publishData.Step == 3 {
			// Toggle selection pour publish
			if a.selectedIndex >= 0 && a.selectedIndex < len(a.resources) {
				a.resources[a.selectedIndex].Selected = !a.resources[a.selectedIndex].Selected
				selectedCount := 0
				for _, res := range a.resources {
					if res.Selected {
						selectedCount++
					}
				}
				a.statusMsg = fmt.Sprintf("%d images selected for publishing.", selectedCount)
			}
		}
		return true
	case "p":
		if a.publishData != nil && a.publishData.Step == 4 {
			// Lancer la publication
			go a.publishSelectedImages()
			a.statusMsg = "Publishing images..."
		}
		return true
	case "up", "k":
		if a.publishData != nil && a.publishData.Step == 3 && a.selectedIndex > 0 {
			a.selectedIndex--
		}
		return true
	case "down", "j":
		if a.publishData != nil && a.publishData.Step == 3 && a.selectedIndex < len(a.resources)-1 {
			a.selectedIndex++
		}
		return true
	}
	return false
}

func (a *AllContainers) selectMenuItem() bool {
	if a.selectedIndex >= 0 && a.selectedIndex < len(a.menuItems) {
		item := a.menuItems[a.selectedIndex]
		a.statusMsg = fmt.Sprintf("Loading %s...", item.Title)

		switch item.Title {
		case "🐳 Publisher":
			a.startPublishFlow()
		case "📦 Delete Containers":
			a.currentMode = "containers"
			a.resourceType = "containers"
			go a.loadContainers()
		case "💿 Delete Images":
			a.currentMode = "images"
			a.resourceType = "images"
			go a.loadImages()
		case "💾 Delete Volumes":
			a.currentMode = "volumes"
			a.resourceType = "volumes"
			go a.loadVolumes()
		}
		a.selectedIndex = 0 // Reset selection pour les listes
		return true
	}
	return false
}

func (a *AllContainers) startPublishFlow() {
	a.currentMode = "publish"
	a.publishData = &PublishData{
		Registry: "docker.io",
		Username: "",
		Tag:      "latest",
		Step:     3, // Aller directement à la sélection d'images
	}
	a.resourceType = "images"
	go a.loadImagesForPublish()
}

func (a *AllContainers) loadContainers() {
	dm, err := a.getDockerManager()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error getting Docker client: %v", err)
		return
	}

	items, err := dm.LoadContainers()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error loading containers: %v", err)
		return
	}

	// Convertir les items en ResourceItem
	a.resources = make([]resource.ResourceItem, 0, len(items))
	for _, item := range items {
		if resItem, ok := item.(resource.ResourceItem); ok {
			a.resources = append(a.resources, resItem)
		}
	}
	a.statusMsg = fmt.Sprintf("Loaded %d containers", len(a.resources))
}

func (a *AllContainers) loadImages() {
	dm, err := a.getDockerManager()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error getting Docker client: %v", err)
		return
	}

	items, err := dm.LoadImages()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error loading images: %v", err)
		return
	}

	// Convertir les items en ResourceItem
	a.resources = make([]resource.ResourceItem, 0, len(items))
	for _, item := range items {
		if resItem, ok := item.(resource.ResourceItem); ok {
			a.resources = append(a.resources, resItem)
		}
	}
	a.statusMsg = fmt.Sprintf("Loaded %d images", len(a.resources))
}

func (a *AllContainers) loadVolumes() {
	dm, err := a.getDockerManager()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error getting Docker client: %v", err)
		return
	}

	items, err := dm.LoadVolumes()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error loading volumes: %v", err)
		return
	}

	// Convertir les items en ResourceItem
	a.resources = make([]resource.ResourceItem, 0, len(items))
	for _, item := range items {
		if resItem, ok := item.(resource.ResourceItem); ok {
			a.resources = append(a.resources, resItem)
		}
	}
	a.statusMsg = fmt.Sprintf("Loaded %d volumes", len(a.resources))
}

func (a *AllContainers) loadImagesForPublish() {
	dm, err := a.getDockerManager()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error getting Docker client: %v", err)
		return
	}

	items, err := dm.LoadImages()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error loading images: %v", err)
		return
	}

	// Filtrer les images <none>:<none> pour publish
	a.resources = make([]resource.ResourceItem, 0)
	for _, item := range items {
		if resItem, ok := item.(resource.ResourceItem); ok {
			if resItem.Title != "<none>:<none>" {
				a.resources = append(a.resources, resItem)
			}
		}
	}
	a.statusMsg = fmt.Sprintf("Loaded %d images for publishing", len(a.resources))
}

func (a *AllContainers) deleteSelectedItems(toDelete []resource.ResourceItem) {
	dm, err := a.getDockerManager()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error getting Docker manager: %v", err)
		return
	}

	var errs []error
	switch a.resourceType {
	case "containers":
		errs = dm.DeleteContainers(toDelete)
	case "images":
		errs = dm.DeleteImages(toDelete)
	case "volumes":
		errs = dm.DeleteVolumes(toDelete)
	}

	if len(errs) > 0 {
		a.statusMsg = fmt.Sprintf("Error deleting: %v", errs[0])
	} else {
		a.statusMsg = fmt.Sprintf("Successfully deleted %d items", len(toDelete))

		// Recharger la liste après suppression
		switch a.resourceType {
		case "containers":
			a.loadContainers()
		case "images":
			a.loadImages()
		case "volumes":
			a.loadVolumes()
		}
	}
}

func (a *AllContainers) publishSelectedImages() {
	dm, err := a.getDockerManager()
	if err != nil {
		a.statusMsg = fmt.Sprintf("Error getting Docker manager: %v", err)
		return
	}

	var toPublish []resource.ResourceItem
	for _, res := range a.resources {
		if res.Selected {
			toPublish = append(toPublish, res)
		}
	}

	if len(toPublish) == 0 {
		a.statusMsg = "No images selected for publishing"
		return
	}

	errs := dm.PublishImages(toPublish, a.publishData.Registry, a.publishData.Username, a.publishData.Tag)

	if len(errs) > 0 {
		a.statusMsg = fmt.Sprintf("Error publishing: %v", errs[0])
	} else {
		a.statusMsg = fmt.Sprintf("Successfully published %d images", len(toPublish))
	}
}

func (a *AllContainers) refreshCurrentList() {
	switch a.currentMode {
	case "containers":
		go a.loadContainers()
	case "images":
		go a.loadImages()
	case "volumes":
		go a.loadVolumes()
	}
}
