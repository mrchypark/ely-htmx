package uiexamples

// Source scenes from examples/gallery/pages/canvas at 49ed1f8. The public
// CanvasScene contract uses lower-case JSON keys; legacy props remain supported.
func workspaceCanvasScene(kind string) map[string]any {
	type M = map[string]any
	shape := func(id, typ string, x, y, w, h int, hue int, text string) M {
		return M{"id": id, "name": id, "type": typ, "x": x, "y": y, "w": w, "h": h, "hue": hue, "text": text}
	}
	link := func(from, to, text string) M { return M{"from": from, "to": to, "text": text} }
	s := M{"width": 620, "height": 300, "view": M{"x": 0, "y": 0, "zoom": 1}}
	basic := []M{shape("card", "rect", 40, 40, 220, 140, 0, ""), shape("badge", "ellipse", 300, 60, 90, 90, 5, ""), shape("label", "text", 60, 210, 160, 32, 0, "Drag me")}
	switch kind {
	case "InfiniteCanvas":
		s["width"], s["height"], s["view"], s["rulers"], s["minimap"] = 820, 460, M{"x": -60, "y": -60, "zoom": .7}, true, true
		s["artboards"] = []M{{"id": "home", "name": "Home · 390", "x": 0, "y": 0, "w": 390, "h": 640}, {"id": "detail", "name": "Detail · 390", "x": 470, "y": 0, "w": 390, "h": 640}}
		shapes := []M{shape("hero", "rect", 24, 32, 342, 180, 0, ""), shape("title", "text", 24, 232, 342, 32, 0, "Studio Atrium"), shape("dot", "ellipse", 24, 290, 56, 56, 5, ""), shape("badge", "polygon", 100, 290, 56, 56, 2, ""), shape("rule", "line", 24, 380, 342, 0, 4, ""), shape("flow", "arrow", 390, 120, 80, 0, 3, ""), shape("sketch", "path", 494, 60, 160, 50, 1, ""), shape("note", "note", 494, 150, 180, 120, 2, "Ask Chloé about the stair")}
		shapes[6]["points"] = []M{{"x": 0, "y": 40}, {"x": 40, "y": 0}, {"x": 80, "y": 30}, {"x": 120, "y": 10}, {"x": 160, "y": 50}}
		s["shapes"] = shapes
		s["guides"] = []M{{"axis": "x", "at": 24}, {"axis": "y", "at": 212}}
	case "ToolPalette", "LayerPanel":
		s["width"], s["height"], s["shapes"], s["selected"] = 560, 360, basic, []string{"card"}
		for _, item := range basic {
			item["name"] = map[string]string{"card": "Card", "badge": "Badge", "label": "Label"}[item["id"].(string)]
		}
		s["assets"] = []M{shape("Button", "rect", 0, 0, 120, 40, 0, ""), shape("Dot", "ellipse", 0, 0, 40, 40, 1, ""), shape("Hexagon", "polygon", 0, 0, 60, 60, 2, ""), shape("Pointer", "arrow", 0, 0, 80, 40, 3, ""), shape("Divider", "line", 0, 0, 120, 0, 4, ""), shape("Heading", "text", 0, 0, 160, 32, 5, "Heading"), shape("Idea", "note", 0, 0, 120, 120, 6, "Idea"), shape("Todo", "note", 0, 0, 120, 120, 7, "Todo")}
	case "ColorPanel":
		s["style"] = M{"color": "#f4efe6", "start": "#e9e4da", "end": "#8a7f6d", "stops": []M{{"At": 0, "Color": "#e9e4da"}, {"At": 100, "Color": "#8a7f6d"}}, "shadows": []M{{"x": 0, "y": 8, "blur": 24, "spread": 0, "color": "#1816132e"}}, "border": M{"width": 1, "color": "#d8d1c4", "dash": "solid", "radii": []int{12, 12, 12, 12}}, "type": M{"family": "Inter", "weight": 500, "size": 20, "lineHeight": 1.3, "tracking": 0, "align": "left"}}
		s["style"].(M)["swatches"] = []M{{"name": "Paper", "color": "#f4efe6"}, {"name": "Sand", "color": "#e9e4da"}, {"name": "Stone", "color": "#b8ad9b"}, {"name": "Moss", "color": "#7d8a6a"}, {"name": "Ink", "color": "#2b2a27"}, {"name": "Clay", "color": "#b86f50"}}
	case "ConstraintsEditor":
		s["rules"] = M{"across": "both", "down": "start", "flow": "row", "gap": 12, "paddingX": 16, "paddingY": 16, "placeX": "start", "placeY": "center", "name": "Card", "exports": []M{{"scale": 1, "format": "png"}, {"scale": 2, "format": "png"}, {"scale": 1, "format": "svg"}}}
	case "NodeEditor":
		port := func(id, kind, name string) M { return M{"id": id, "kind": kind, "name": name} }
		s["view"] = M{"x": 0, "y": 0, "zoom": .8}
		s["nodes"] = []M{{"id": "load", "name": "Load image", "x": 20, "y": 40, "outputs": []M{port("image", "image", "Image")}}, {"id": "prompt", "name": "Prompt", "x": 20, "y": 170, "outputs": []M{port("text", "text", "Text")}}, {"id": "model", "name": "Model", "x": 290, "y": 60, "inputs": []M{port("image", "image", "Image"), port("text", "text", "Prompt"), port("steps", "number", "Steps")}, "outputs": []M{port("image", "image", "Image")}}, {"id": "save", "name": "Save image", "x": 560, "y": 90, "inputs": []M{port("image", "image", "Image")}}}
		s["links"] = []M{{"from": "load", "fromPort": "image", "to": "model", "toPort": "image"}, {"from": "prompt", "fromPort": "text", "to": "model", "toPort": "text"}}
	case "FlowChart":
		s["view"], s["tools"] = M{"x": 0, "y": 0, "zoom": .8}, []string{"select", "connector"}
		s["shapes"] = []M{shape("start", "ellipse", 20, 132, 110, 56, 5, "Start"), shape("draft", "rect", 180, 132, 150, 56, 0, "Write the draft"), shape("check", "diamond", 380, 122, 140, 56, 2, "Reviewed?"), shape("ship", "rect", 580, 132, 120, 56, 0, "Publish"), shape("end", "ellipse", 580, 250, 120, 56, 5, "End"), shape("revise", "rect", 370, 250, 160, 56, 0, "Ask for changes")}
		s["links"] = []M{link("start", "draft", ""), link("draft", "check", ""), link("check", "ship", "Yes"), link("check", "revise", "No"), link("ship", "end", "")}
	case "MindMap":
		z := 588.0 / 1040
		s["view"] = M{"x": -310 / z, "y": -150 / z, "zoom": z}
		s["selected"] = []string{"site"}
		s["topics"] = []M{{"id": "root", "text": "Product launch"}, {"id": "site", "text": "Website", "parent": "root"}, {"id": "copy", "text": "Copy", "parent": "site"}, {"id": "screens", "text": "Screenshots", "parent": "site"}, {"id": "press", "text": "Press", "parent": "root"}, {"id": "kit", "text": "Press kit", "parent": "press"}, {"id": "budget", "text": "Budget", "parent": "root"}, {"id": "team", "text": "Team", "parent": "root"}, {"id": "roles", "text": "Roles", "parent": "team"}}
	case "Whiteboard":
		s["height"] = 340
		s["tools"] = []string{"select", "hand", "note", "rect", "ellipse", "line", "arrow", "connector", "text", "pen", "brush"}
		s["shapes"] = []M{shape("why", "note", 30, 80, 150, 150, 2, "Why do people leave in week two?"), shape("guess", "note", 230, 80, 150, 150, 5, "Setup takes too long"), shape("try", "note", 430, 80, 150, 150, 1, "Try a sample project"), shape("goal", "text", 30, 270, 220, 32, 0, "Week two retention")}
		s["links"] = []M{link("why", "guess", ""), link("guess", "try", "")}
	}
	return s
}
