package uiexamples

// Gallery source: examples/gallery/pages/shell at 49ed1f8. OS effects are host requests.
func nativeShellFixtures(name string, props map[string]any) {
	p := map[string]any{"HideHeader": true}
	tool := func(label, icon, action string) map[string]any {
		return map[string]any{"Label": label, "Icon": icon, "Action": action}
	}
	platforms := []map[string]any{{"Platform": "mac", "Label": "Ely — Projects", "Active": false, "Specimen": true, "Leading": []map[string]any{tool("Sidebar", "panel-left", "toggle-sidebar")}, "Actions": []map[string]any{tool("Search", "search", "search")}}, {"Platform": "windows", "Label": "Ely — Projects", "Active": false, "Specimen": true, "Leading": []map[string]any{tool("Sidebar", "panel-left", "toggle-sidebar")}, "Actions": []map[string]any{tool("Search", "search", "search")}}, {"Platform": "linux", "Label": "Ely — Projects", "Active": false, "Specimen": true, "Leading": []map[string]any{tool("Sidebar", "panel-left", "toggle-sidebar")}, "Actions": []map[string]any{tool("Search", "search", "search")}}}
	switch name {
	case "TitleBar":
		p["PreviewStack"] = true
		for i, caption := range []string{"macOS", "Windows", "Linux"} {
			platforms[i]["SpecimenCaption"] = caption
		}
		platforms[2]["SpecimenNote"] = "TitleBar::new(id) with no platform leaves room for the real macOS lights."
		for k, v := range platforms[0] {
			p[k] = v
		}
		p["PreviewExamples"] = platforms
	case "WindowControls":
		p["Platform"], p["Active"] = "mac", false
		p["PreviewExamples"] = []map[string]any{{"Platform": "mac", "Active": false, "SpecimenCaption": "Mac"}, {"Platform": "windows", "Active": false, "SpecimenCaption": "Windows"}, {"Platform": "linux", "Active": false, "SpecimenCaption": "Linux"}}
	case "ColumnShell":
		p["PreviewStack"] = true
		leading := []map[string]any{tool("Sidebar", "panel-left", "toggle-sidebar"), tool("Back", "arrow-left", "back"), tool("Forward", "arrow-right", "forward")}
		side := func(mail bool) map[string]any {
			items := []map[string]any{{"Value": "new", "Label": "New task", "Icon": "square-pen", "Selected": true}, {"Value": "search", "Label": "Search", "Icon": "search"}, {"Value": "inbox", "Label": "Inbox", "Icon": "inbox"}}
			if mail {
				items = []map[string]any{{"Value": "inbox", "Label": "Inbox", "Icon": "inbox", "Selected": true}, {"Value": "sent", "Label": "Sent", "Icon": "send"}, {"Value": "archive", "Label": "Archive", "Icon": "archive"}}
			}
			return map[string]any{"Width": 240, "Leading": leading, "Items": items, "Sidebar": true}
		}
		two := []map[string]any{side(false), {"Title": "New task", "Centered": true, "Content": "What should we build?", "Actions": []map[string]any{tool("Info", "info", "info"), tool("Details", "panel-right", "toggle-details")}}}
		three := []map[string]any{side(true), {"Width": 220, "Title": "Inbox", "Actions": []map[string]any{tool("Search", "search", "search")}, "Rows": []map[string]any{{"Label": "Quarterly review", "Content": "Numbers, notes and the next three moves."}, {"Label": "Design sync", "Content": "Tokens for the new sidebar."}, {"Label": "Launch plan", "Content": "Dates for the beta."}}}, {"Title": "Quarterly review", "Content": "Revenue rose eight percent. Three moves come next: hire, ship the beta, close the books.", "Actions": []map[string]any{tool("Reply", "reply", "reply"), tool("Archive", "archive", "archive")}}}
		p["Columns"], p["Width"], p["Height"], p["Platform"] = two, 640, 280, "linux"
		p["PreviewExamples"] = []map[string]any{{"Columns": two, "Width": 640, "Height": 280, "Platform": "linux", "SpecimenCaption": "Two columns"}, {"Columns": three, "Width": 760, "Height": 280, "Platform": "linux", "SpecimenCaption": "Three columns"}}
	case "WindowDragRegion":
		p["Label"] = "Double-click to zoom this window"
	case "WindowResizeBorder":
		p["SourceNote"] = "Linux draws its own edges. This macOS window keeps the system's, so the border lays out only."
	case "NativeMenu":
		p["SourceNotes"] = []string{"Look up: Ely and View come from cx.set_menus(menus)", "Right-click the Dock icon: New Window comes from cx.set_dock_menu(items)", "MenuBar draws the same menus inside the window. It lands with Menus, chapter 8."}
	case "Toolbar":
		p["Groups"] = [][]map[string]any{{tool("Undo", "undo-2", "undo"), tool("Redo", "redo-2", "redo")}, {tool("Bold", "bold", "bold"), tool("Italic", "italic", "italic"), tool("Underline", "underline", "underline")}, {tool("Align left", "align-left", "align-left"), tool("Align center", "align-center", "align-center"), tool("Align right", "align-right", "align-right")}}
		p["Actions"] = []map[string]any{tool("Share", "share-2", "share")}
	case "StatusBar":
		p["Specimen"] = true
		p["Items"] = []map[string]any{{"Label": "main", "Icon": "git-branch", "Action": "switch-branch", "Tooltip": "Switch branch"}, {"Label": "Synced", "Icon": "refresh-cw"}}
		p["RightItems"] = []map[string]any{{"Label": "Ln 12, Col 4"}, {"Label": "UTF-8"}, {"Icon": "bell", "Tooltip": "No notifications", "Action": "notifications"}}
	case "ActivityBar":
		p["Specimen"], p["Selected"] = true, "files"
		p["Items"] = []map[string]any{{"Value": "files", "Label": "Explorer", "Icon": "folder"}, {"Value": "search", "Label": "Search", "Icon": "search"}, {"Value": "git", "Label": "Source Control", "Icon": "git-branch", "Badge": 3}, {"Value": "debug", "Label": "Run and Debug", "Icon": "bug"}, {"Value": "extensions", "Label": "Extensions", "Icon": "puzzle"}}
		p["FooterItems"] = []map[string]any{{"Value": "account", "Label": "Account", "Icon": "user"}, {"Value": "settings", "Label": "Settings", "Icon": "settings"}}
	case "TabBar":
		p["Specimen"], p["Selected"], p["LocalAdd"] = true, "roadmap", true
		p["Items"] = []map[string]any{{"Value": "inbox", "Label": "Inbox", "Icon": "inbox"}, {"Value": "roadmap", "Label": "Roadmap", "Icon": "map", "Selected": true}, {"Value": "notes", "Label": "Meeting notes", "Icon": "notebook-pen"}}
	case "SplashScreen":
		p["Label"], p["Tagline"], p["Status"], p["Progress"] = "Ely", "Components for GPUI", "Loading workspace…", 62
	case "AboutDialog":
		p["Label"], p["Value"], p["Content"], p["Copyright"] = "Ely", "0.1.0 (2026.09)", "A component library for GPUI, in light and dark.", "© 2026 Ely contributors"
		p["Items"] = []map[string]any{{"Label": "Website", "URL": "https://git.zacharyzhang.com/ZacharyZhang-NY/Ely-GPUI-Components"}, {"Label": "License", "URL": "https://git.zacharyzhang.com/ZacharyZhang-NY/Ely-GPUI-Components#license"}}
	case "TrayIcon":
		p["Shown"], p["Symbol"], p["LastPick"] = false, "circle.hexagongrid", "nothing yet"
		p["Items"] = []map[string]any{{"Label": "Open gallery", "Value": "0"}, {"Label": "Toggle theme", "Value": "1"}, {"Separator": true}, {"Label": "Quit", "Value": "3"}}
	case "DockBadge":
		p["Value"] = ""
	case "SystemNotification":
		p["SourceNote"] = "Not shown: gpui on the web posts no system notification."
	case "MultiWindow Manager":
		p["Items"] = []map[string]any{}
	case "WindowSwitcher":
		p["Items"] = []map[string]any{}
	case "QuickLauncher":
		p["SourceNote"] = "It is `navigation::QuickLauncher`; the Navigation page, chapter 7, opens one."
	case "MiniWindow":
		p["LauncherOnly"] = true
		p["Items"] = []map[string]any{{"Value": "mini", "Label": "Mini player", "Action": "open-mini-player", "Width": 320, "Height": 96, "Corner": "bottom-right"}, {"Value": "pip", "Label": "Picture in picture", "Action": "open-picture-in-picture", "Width": 384, "Height": 216, "Corner": "top-right"}}
	case "UpdateBanner":
		p["Banners"] = []map[string]any{{"Value": "available", "Version": "0.2.0", "State": "available", "DownloadAction": "download-update", "LaterAction": "later"}, {"Value": "downloading", "Version": "0.2.0", "State": "downloading", "Progress": 42, "LaterAction": "later"}, {"Value": "ready", "Version": "0.2.0", "State": "ready", "RestartAction": "restart-update", "LaterAction": "later"}}
	case "UpdateDialog":
		p["Label"], p["CurrentVersion"], p["Version"], p["State"] = "Ely", "0.1.0", "0.2.0", "available"
		p["DownloadAction"], p["LaterAction"] = "download-update", "later"
		p["Items"] = []map[string]any{{"Label": "Window chrome for macOS, Windows and Linux."}, {"Label": "Docks remember their layout across launches."}, {"Label": "Sticky headers stay inside their scroll area."}, {"Label": "Drawers follow the pointer outside their bounds."}}
	case "CrashReporter":
		p["Label"], p["Content"], p["Details"] = "Ely quit unexpectedly", "The report says where it stopped. Sending it helps fix it.", true
		p["Value"] = "Exception Type:  EXC_BAD_ACCESS (SIGSEGV)\nCrashed Thread:  0  main\n0   ely           0x0000000104a1c3f0 layout::dock::restore + 112\n1   ely           0x0000000104a1b9d4 workspace::apply + 88\n2   ely           0x0000000104a0f2a8 app::open_recent + 404\n3   libdyld.dylib 0x000000018f2d4f28 start + 4"
		p["Action"], p["DismissAction"] = "send-report", "dismiss-report"
	case "OfflineIndicator":
		p["State"], p["ShowChoices"] = "online", true
	case "ZoomControl":
		p["Percent"] = 100
	}
	if banners, ok := p["Banners"].([]map[string]any); ok {
		for _, banner := range banners {
			banner["DownloadAction"], banner["RestartAction"] = "download-update", "restart-update"
		}
	}
	// Clear generic fixture content only; callers never pass through this gallery function.
	for _, key := range []string{"Title", "Description", "Label", "Content", "Value", "Items", "Options", "Action", "ActionLabel", "Progress", "PreviewExamples"} {
		delete(props, key)
	}
	if _, ok := p["PreviewExamples"]; !ok {
		copy := map[string]any{}
		for k, v := range p {
			copy[k] = v
		}
		p["PreviewExamples"] = []map[string]any{copy}
	}
	for _, example := range p["PreviewExamples"].([]map[string]any) {
		example["HideHeader"] = true
	}
	for k, v := range p {
		props[k] = v
	}
}
