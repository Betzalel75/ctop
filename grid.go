package main

import (
	"github.com/Betzalel75/ctop/config"
	"github.com/Betzalel75/ctop/cwidgets/single"
	ui "github.com/gizak/termui"
)

func ShowConnError(err error) (exit bool) {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	setErr := func(err error) {
		errView.Append(err.Error())
		errView.Append("attempting to reconnect...")
		ui.Render(errView)
	}

	HandleKeys("exit", func() {
		exit = true
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(ui.Event) {
		_, err := cursor.RefreshContainers()
		if err == nil {
			ui.StopLoop()
			return
		}
		setErr(err)
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		errView.Resize()
		ui.Clear()
		ui.Render(errView)
		log.Infof("RESIZE")
	})

	errView.Resize()
	setErr(err)
	ui.Loop()
	return exit
}

func RedrawRows(clr bool) {
	// reinit body rows
	cGrid.Clear()

	// build layout
	y := 1
	if config.GetSwitchVal("enableHeader") {
		header.SetCount(cursor.Len())
		header.SetFilter(config.GetVal("filterStr"))
		y += header.Height()
	}

	cGrid.SetY(y)

	for _, c := range cursor.filtered {
		cGrid.AddRows(c.Widgets)
	}

	if clr {
		ui.Clear()
		log.Debugf("screen cleared")
	}
	if config.GetSwitchVal("enableHeader") {
		ui.Render(header)
	}
	// MODIFIÉ: Utiliser containerView au lieu de runningWidget
	containerView.Align()
	ui.Render(containerView)
}

func SingleView() MenuFn {
	c := cursor.Selected()
	if c == nil {
		return nil
	}

	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	ex := single.NewSingle()
	c.SetUpdater(ex)

	ex.Align()
	ui.Render(ex)

	HandleKeys("up", ex.Up)
	HandleKeys("down", ex.Down)
	ui.Handle("/sys/kbd/", func(ui.Event) { ui.StopLoop() })

	ui.Handle("/timer/1s", func(ui.Event) { ui.Render(ex) })
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ex.SetWidth(ui.TermWidth())
		ex.Align()
		log.Infof("resize: width=%v max-rows=%v", ex.Width, cGrid.MaxRows())
	})

	ui.Loop()
	c.SetUpdater(c.Widgets)
	return nil
}

func RefreshDisplay() error {
	// skip display refresh during scroll
	if !cursor.isScrolling {
		needsClear, err := cursor.RefreshContainers()
		if err != nil {
			return err
		}
		RedrawRows(needsClear)
	}
	return nil
}

// func Display() bool {
// 	var menu MenuFn
// 	var connErr error

// 	cGrid.SetWidth(ui.TermWidth())
// 	ui.DefaultEvtStream.Hook(logEvent)

// 	// initial draw
// 	header.Align()
// 	status.Align()
// 	cursor.RefreshContainers()
// 	RedrawRows(true)

// 	// HandleKeys("up", cursor.Up)
// 	// HandleKeys("down", cursor.Down)

// 	// HandleKeys("pgup", cursor.PgUp)
// 	// HandleKeys("pgdown", cursor.PgDown)

// 	// HANDLERS POUR LE WIDGET RUNNING (par défaut)
// 	HandleKeys("up", func() {
// 		if containerView.IsRunningActive() {
// 			cursor.Up()
// 		}
// 	})
// 	HandleKeys("down", func() {
// 		if containerView.IsRunningActive() {
// 			cursor.Down()
// 		}
// 	})

// 	HandleKeys("pgup", func() {
// 		if containerView.IsRunningActive() {
// 			cursor.PgUp()
// 		}
// 	})
// 	HandleKeys("pgdown", func() {
// 		if containerView.IsRunningActive() {
// 			cursor.PgDown()
// 		}
// 	})

// 	HandleKeys("exit", ui.StopLoop)
// 	HandleKeys("help", func() {
// 		menu = HelpMenu
// 		ui.StopLoop()
// 	})

// 	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
// 		menu = ContainerMenu
// 		ui.StopLoop()
// 	})
// 	ui.Handle("/sys/kbd/<left>", func(ui.Event) {
// 		menu = LogMenu
// 		ui.StopLoop()
// 	})
// 	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
// 		menu = SingleView
// 		ui.StopLoop()
// 	})
// 	ui.Handle("/sys/kbd/l", func(ui.Event) {
// 		menu = LogMenu
// 		ui.StopLoop()
// 	})
// 	ui.Handle("/sys/kbd/e", func(ui.Event) {
// 		menu = ExecShell
// 		ui.StopLoop()
// 	})
// 	ui.Handle("/sys/kbd/w", func(ui.Event) {
// 		menu = OpenInBrowser()
// 	})
// 	ui.Handle("/sys/kbd/o", func(ui.Event) {
// 		menu = SingleView
// 		ui.StopLoop()
// 	})
// 	ui.Handle("/sys/kbd/a", func(ui.Event) {
// 		config.Toggle("allContainers")
// 		connErr = RefreshDisplay()
// 		if connErr != nil {
// 			ui.StopLoop()
// 		}
// 	})
// 	ui.Handle("/sys/kbd/D", func(ui.Event) {
// 		dumpContainer(cursor.Selected())
// 	})
// 	ui.Handle("/sys/kbd/f", func(ui.Event) {
// 		menu = FilterMenu
// 		ui.StopLoop()
// 	})
// 	ui.Handle("/sys/kbd/H", func(ui.Event) {
// 		config.Toggle("enableHeader")
// 		RedrawRows(true)
// 	})
// 	ui.Handle("/sys/kbd/r", func(e ui.Event) {
// 		config.Toggle("sortReversed")
// 	})
// 	ui.Handle("/sys/kbd/s", func(ui.Event) {
// 		menu = SortMenu
// 		ui.StopLoop()
// 	})
// 	ui.Handle("/sys/kbd/c", func(ui.Event) {
// 		menu = ColumnsMenu
// 		ui.StopLoop()
// 	})
// 	ui.Handle("/sys/kbd/S", func(ui.Event) {
// 		path, err := config.Write()
// 		if err == nil {
// 			log.Statusf("wrote config to %s", path)
// 		} else {
// 			log.StatusErr(err)
// 		}
// 		ui.StopLoop()
// 	})

// 	ui.Handle("/timer/1s", func(e ui.Event) {
// 		if log.StatusQueued() {
// 			ui.StopLoop()
// 		}
// 		connErr = RefreshDisplay()
// 		if connErr != nil {
// 			ui.StopLoop()
// 		}
// 	})

// 	// NOUVEAU: Ajouter les contrôles pour basculer entre les widgets
// 	ui.Handle("/sys/kbd/<tab>", func(ui.Event) {
// 		if containerView.IsRunningActive() {
// 			containerView.SwitchToAll()
// 		} else {
// 			containerView.SwitchToRunning()
// 		}
// 		RedrawRows(false)
// 	})

// 	ui.Handle("/sys/kbd/1", func(ui.Event) {
// 		containerView.SwitchToRunning()
// 		RedrawRows(false)
// 	})

// 	ui.Handle("/sys/kbd/2", func(ui.Event) {
// 		containerView.SwitchToAll()
// 		RedrawRows(false)
// 	})

// 	ui.Handle("/sys/kbd/", func(e ui.Event) {
// 		if containerView.GetActiveWidget() == "all" {
// 			key := strings.ReplaceAll(e.Path, "/sys/kbd/", "")
// 			if containerView.AllWidget.HandleKey(key) {
// 				RedrawRows(false)
// 			}
// 		}
// 	})

// 	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
// 		header.Align()
// 		status.Align()
// 		containerView.Align()
// 		// runningWidget.Align()
// 		cursor.ScrollPage()
// 		cGrid.SetWidth(ui.TermWidth())
// 		log.Infof("resize: width=%v max-rows=%v", cGrid.Width, cGrid.MaxRows())
// 		RedrawRows(true)
// 	})

// 	ui.Loop()

// 	if connErr != nil {
// 		return ShowConnError(connErr)
// 	}

// 	if log.StatusQueued() {
// 		for sm := range log.FlushStatus() {
// 			if sm.IsError {
// 				status.ShowErr(sm.Text)
// 			} else {
// 				status.Show(sm.Text)
// 			}
// 		}
// 		return false
// 	}

// 	if menu != nil {
// 		for menu != nil {
// 			menu = menu()
// 		}
// 		return false
// 	}

// 	return true
// }

func Display() bool {
	var menu MenuFn
	var connErr error

	cGrid.SetWidth(ui.TermWidth())
	ui.DefaultEvtStream.Hook(logEvent)

	// initial draw
	header.Align()
	status.Align()
	cursor.RefreshContainers()
	RedrawRows(true)

	// HANDLERS POUR LE WIDGET RUNNING (par défaut)
	HandleKeys("up", func() {
		if containerView.IsRunningActive() {
			cursor.Up()
		}
	})
	HandleKeys("down", func() {
		if containerView.IsRunningActive() {
			cursor.Down()
		}
	})

	HandleKeys("pgup", func() {
		if containerView.IsRunningActive() {
			cursor.PgUp()
		}
	})
	HandleKeys("pgdown", func() {
		if containerView.IsRunningActive() {
			cursor.PgDown()
		}
	})

	HandleKeys("exit", ui.StopLoop)
	HandleKeys("help", func() {
		menu = HelpMenu
		ui.StopLoop()
	})

	// HANDLERS CONDITIONNELS selon le widget actif
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		if containerView.IsRunningActive() {
			menu = ContainerMenu
			ui.StopLoop()
		} else {
			// Déléguer au widget All
			if containerView.AllWidget.HandleKey("enter") {
				RedrawRows(false)
			}
		}
	})

	ui.Handle("/sys/kbd/<left>", func(ui.Event) {
		if containerView.IsRunningActive() {
			menu = LogMenu
			ui.StopLoop()
		}
	})

	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
		if containerView.IsRunningActive() {
			menu = SingleView
			ui.StopLoop()
		}
	})

	ui.Handle("/sys/kbd/l", func(ui.Event) {
		if containerView.IsRunningActive() {
			menu = LogMenu
			ui.StopLoop()
		}
	})

	ui.Handle("/sys/kbd/e", func(ui.Event) {
		if containerView.IsRunningActive() {
			menu = ExecShell
			ui.StopLoop()
		}
	})

	ui.Handle("/sys/kbd/w", func(ui.Event) {
		if containerView.IsRunningActive() {
			menu = OpenInBrowser()
		}
	})

	ui.Handle("/sys/kbd/o", func(ui.Event) {
		if containerView.IsRunningActive() {
			menu = SingleView
			ui.StopLoop()
		}
	})

	ui.Handle("/sys/kbd/a", func(ui.Event) {
		if containerView.IsRunningActive() {
			config.Toggle("allContainers")
			connErr = RefreshDisplay()
			if connErr != nil {
				ui.StopLoop()
			}
		}
	})

	ui.Handle("/sys/kbd/D", func(ui.Event) {
		if containerView.IsRunningActive() {
			dumpContainer(cursor.Selected())
		}
	})

	ui.Handle("/sys/kbd/f", func(ui.Event) {
		if containerView.IsRunningActive() {
			menu = FilterMenu
			ui.StopLoop()
		}
	})

	ui.Handle("/sys/kbd/H", func(ui.Event) {
		config.Toggle("enableHeader")
		RedrawRows(true)
	})

	// ui.Handle("/sys/kbd/r", func(e ui.Event) {
	// 	if containerView.IsRunningActive() {
	// 		config.Toggle("sortReversed")
	// 	}
	// })

	ui.Handle("/sys/kbd/s", func(ui.Event) {
		if containerView.IsRunningActive() {
			menu = SortMenu
			ui.StopLoop()
		}
	})

	ui.Handle("/sys/kbd/c", func(ui.Event) {
		menu = ColumnsMenu
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/S", func(ui.Event) {
		path, err := config.Write()
		if err == nil {
			log.Statusf("wrote config to %s", path)
		} else {
			log.StatusErr(err)
		}
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/r", func(e ui.Event) {
		if containerView.IsRunningActive() {
			config.Toggle("sortReversed")
		} else {
			// Déléguer au widget All pour le refresh
			if containerView.AllWidget.HandleKey("r") {
				RedrawRows(false)
			}
		}
	})

	ui.Handle("/sys/kbd/q", func(e ui.Event) {
		if containerView.IsRunningActive() {
			ui.StopLoop() // Quitter l'application
		} else {
			// Dans le widget All, 'q' doit revenir au menu ou quitter selon le contexte
			if containerView.AllWidget.HandleKey("q") {
				RedrawRows(false)
			} else {
				ui.StopLoop() // Quitter si on est déjà au menu principal
			}
		}
	})
	// Ajouter aussi les handlers pour PageUp/PageDown
	ui.Handle("/sys/kbd/<prior>", func(ui.Event) { // PageUp
		if !containerView.IsRunningActive() {
			if containerView.AllWidget.HandleKey("pgup") {
				RedrawRows(false)
			}
		}
	})

	ui.Handle("/sys/kbd/<next>", func(ui.Event) { // PageDown
		if !containerView.IsRunningActive() {
			if containerView.AllWidget.HandleKey("pgdown") {
				RedrawRows(false)
			}
		}
	})

	// Ajouter handler pour la barre d'espace
	ui.Handle("/sys/kbd/<space>", func(ui.Event) {
		if !containerView.IsRunningActive() {
			if containerView.AllWidget.HandleKey("space") {
				RedrawRows(false)
			}
		}
	})

	// Handler pour la touche 'd' (delete)
	ui.Handle("/sys/kbd/d", func(ui.Event) {
		if !containerView.IsRunningActive() {
			if containerView.AllWidget.HandleKey("d") {
				RedrawRows(false)
			}
		}
	})
	//
	// 

	ui.Handle("/timer/1s", func(e ui.Event) {
		if log.StatusQueued() {
			ui.StopLoop()
		}

		// Ne rafraîchir que si on est sur l'onglet Running
		if containerView.IsRunningActive() {
			connErr = RefreshDisplay()
			if connErr != nil {
				ui.StopLoop()
			}
		}
	})

	// CONTRÔLES DE NAVIGATION ENTRE LES WIDGETS
	ui.Handle("/sys/kbd/<tab>", func(ui.Event) {
		if containerView.IsRunningActive() {
			containerView.SwitchToAll()
		} else {
			containerView.SwitchToRunning()
		}
		RedrawRows(false)
	})

	ui.Handle("/sys/kbd/1", func(ui.Event) {
		containerView.SwitchToRunning()
		RedrawRows(false)
	})

	ui.Handle("/sys/kbd/2", func(ui.Event) {
		containerView.SwitchToAll()
		RedrawRows(false)
	})

	// GESTIONNAIRE SPÉCIFIQUE POUR LE WIDGET ALL
	// Utiliser des touches spécifiques pour éviter les conflits
	ui.Handle("/sys/kbd/k", func(ui.Event) {
		if !containerView.IsRunningActive() {
			if containerView.AllWidget.HandleKey("up") {
				RedrawRows(false)
			}
		} else {
			cursor.Up() // Comportement par défaut pour Running
		}
	})

	ui.Handle("/sys/kbd/j", func(ui.Event) {
		if !containerView.IsRunningActive() {
			if containerView.AllWidget.HandleKey("down") {
				RedrawRows(false)
			}
		} else {
			cursor.Down() // Comportement par défaut pour Running
		}
	})

	// Touches fléchées pour le widget All
	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		if !containerView.IsRunningActive() {
			if containerView.AllWidget.HandleKey("up") {
				RedrawRows(false)
			}
		} else {
			cursor.Up()
		}
	})

	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		if !containerView.IsRunningActive() {
			if containerView.AllWidget.HandleKey("down") {
				RedrawRows(false)
			}
		} else {
			cursor.Down()
		}
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		header.Align()
		status.Align()
		containerView.Align()
		cursor.ScrollPage()
		cGrid.SetWidth(ui.TermWidth())
		log.Infof("resize: width=%v max-rows=%v", cGrid.Width, cGrid.MaxRows())
		RedrawRows(true)
	})

	ui.Loop()

	if connErr != nil {
		return ShowConnError(connErr)
	}

	if log.StatusQueued() {
		for sm := range log.FlushStatus() {
			if sm.IsError {
				status.ShowErr(sm.Text)
			} else {
				status.Show(sm.Text)
			}
		}
		return false
	}

	if menu != nil {
		for menu != nil {
			menu = menu()
		}
		return false
	}

	return true
}
