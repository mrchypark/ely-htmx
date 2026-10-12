package uiexamples

import (
	"fmt"

	"github.com/mrchypark/ely-htmx"
)

// advancedInteractionCases gives the gallery one caller-owned instance for every interaction primary.
func advancedInteractionCases() []AdvancedCase {
	stories := map[string][]string{
		"buttons":     {"Button", "Sizes", "With icons", "Icon button", "ShortcutHint", "ButtonGroup", "SplitButton", "ToggleButton", "ToggleGroup", "SegmentedControl", "FloatingActionButton", "ConfirmButton", "CopyButton", "ActionBar", "BulkActionBar", "ActionSheet", "QuickActions"},
		"navigation":  {"Tabs", "EditorTabs", "Breadcrumb", "Pagination", "Steps", "Sidebar Navigation", "NavigationMenu", "BackForwardNavigation", "Anchor", "JumpTo", "CommandPalette", "QuickOpen", "QuickSwitcher", "SearchPalette"},
		"menus":       {"Menu", "CheckboxMenuItem", "DropdownMenu", "ContextMenu", "MenuBar Menu", "MegaMenu", "SearchableMenu", "PieMenu"},
		"overlays":    {"Popover", "HoverCard", "Dialog", "Lightbox", "Spotlight", "FloatingToolbar", "Callout", "Tooltip"},
		"feedback":    {"Toast", "Notification", "Alert", "Banner", "Callout", "InlineMessage", "EmptyState", "ErrorView", "ErrorBoundary", "Result", "ConfirmationCard", "SavingIndicator", "Countdown"},
		"motion":      {"Spinner", "Skeleton", "Shimmer", "LoadingButton", "LoadMore", "ProgressBar", "LoadingOverlay", "LazyLoad", "Suspense", "RefreshIndicator", "UploadList", "Transition", "Stagger", "AnimatePresence", "Flip", "Reorder", "Marquee", "Ripple", "Glow", "Blink", "ParticleBackground"},
		"interaction": {"Draggable", "Resizable", "LongPress", "Hotkeys", "Undo", "VirtualScroller", "Autofocus"},
	}
	var cases []AdvancedCase
	for chapter, names := range stories {
		for index, name := range names {
			id := fmt.Sprintf("audit-%s-%d", chapter, index+1)
			props := map[string]any{
				"Title": name + " for Northwind review", "Label": "Use " + name, "ActionLabel": "Continue", "Content": "Caller-controlled content for " + id, "Description": "Caller-controlled detail for " + id,
				"Name": "field_" + fmt.Sprintf("%d", index+1), "Value": fmt.Sprintf("%d", index+2), "Placeholder": "Find a caller-owned item", "Required": index%5 == 0, "Disabled": index%13 == 0,
				"Action": "save", "Icon": "✦", "Alt": name + " preview", "Src": "/assets/sample.svg", "Href": "#" + id,
				"Max": 100, "Loading": false, "Autofocus": name == "Autofocus", "DropLabel": "Drop caller item here", "UndoLabel": "Undo", "RedoLabel": "Redo", "BackLabel": "Back", "ForwardLabel": "Forward",
				"Items": []map[string]any{{"Value": "first-" + id, "Label": "First caller item", "Title": "First caller item", "Content": "First independent caller value", "Selected": true, "Href": "#first-" + id}, {"Value": "second-" + id, "Label": "Second caller item", "Title": "Second caller item", "Content": "Second independent caller value", "Href": "#second-" + id}},
			}
			if name == "Breadcrumb" {
				props["Items"] = []map[string]any{{"Label": "Research", "Href": "#research"}, {"Label": "Northwind review"}}
			}
			if name == "ParticleBackground" {
				props["Animation"] = `{"v":"5.7.4","fr":30,"ip":0,"op":30,"w":80,"h":80,"nm":"Caller shapes","ddd":0,"assets":[],"layers":[{"ddd":0,"ind":1,"ty":4,"nm":"Circle","sr":1,"ks":{"o":{"a":0,"k":100},"r":{"a":0,"k":0},"p":{"a":1,"k":[{"t":0,"s":[20,40,0],"e":[60,40,0]},{"t":30,"s":[60,40,0]}]},"a":{"a":0,"k":[0,0,0]},"s":{"a":0,"k":[100,100,100]}},"shapes":[{"ty":"el","p":{"a":0,"k":[0,0]},"s":{"a":0,"k":[28,28]},"nm":"Circle"},{"ty":"fl","c":{"a":0,"k":[0.35,0.4,0.8,1]},"o":{"a":0,"k":100},"nm":"Fill"}],"ip":0,"op":30,"st":0,"bm":0}]}`
				props["Autoplay"] = true
				props["Loop"] = true
			}
			cases = append(cases, AdvancedCase{Family: "interaction", Component: ui.Component{Chapter: chapter, Name: name, ID: id, Props: props}})
		}
	}
	// Exercise the Lottie alias through the same primary with a caller-owned animation.
	cases = append(cases, AdvancedCase{Family: "interaction", Component: ui.Component{Chapter: "motion", Name: "Lottie Player", ID: "audit-motion-lottie", Props: map[string]any{"Title": "Caller Lottie animation", "Label": "Caller animation", "Animation": `{"v":"5.7.4","fr":30,"ip":0,"op":2,"w":8,"h":8,"layers":[]}`, "Autoplay": true, "Loop": true}}})
	return cases
}
