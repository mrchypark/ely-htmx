package uiexamples

import (
	"fmt"

	"github.com/mrchypark/ely-htmx"
)

// advancedFoundationCases gives the gallery one caller-owned case for every Foundation primary.
func advancedFoundationCases() []AdvancedCase {
	stories := map[string][]string{
		"primitives": {"Box", "Text", "Icon", "IconSet", "Svg", "Image", "Canvas", "Divider", "Spacer", "Portal", "Overlay", "Slot", "Show", "For", "Fragment", "VisuallyHidden", "FocusRing", "FocusScope", "ClickOutside", "HoverArea", "Pressable", "KeyboardHandler", "Measure", "IntersectionObserver", "Clipboard", "Tooltip Trigger"},
		"typography": {"Heading", "Title", "Paragraph", "Label", "Caption", "Overline", "Code", "Kbd", "KbdCombo", "Blockquote", "Highlight", "Truncate", "LineClamp", "MiddleEllipsis", "EllipsisTooltip", "Link", "ExternalLink", "RelativeTime", "DateTimeText", "NumberText", "CurrencyText", "PercentText", "FileSizeText", "DurationText", "PluralText", "CopyableText", "SelectableText", "Emoji", "Latex", "AnimatedNumber", "Typewriter", "GradientText", "ShimmerText"},
		"layout":     {"Stack", "Flex", "Grid", "SimpleGrid", "Masonry", "Center", "Container", "AspectRatio", "Inset", "Wrap", "Absolute", "Layer", "ScrollArea", "Scrollbar", "ScrollShadow", "ScrollToTop", "StickyHeader", "SplitPane", "ResizeHandle", "Collapsible", "Accordion", "Card", "Panel", "Section", "Group", "Frame", "Well", "Sidebar", "Drawer", "Sheet", "Page", "AppShell", "MasterDetail", "Dock", "DockZone", "FloatingPanel", "Pane", "Workspace", "Viewport"},
	}
	var cases []AdvancedCase
	for chapter, names := range stories {
		for index, name := range names {
			id := fmt.Sprintf("audit-%s-%d", chapter, index+1)
			props := map[string]any{
				"Title": name + " caller content", "Text": "Caller supplied text for " + id + ".", "Label": "Caller use " + name, "CloseLabel": "Caller close", "FieldLabel": "Caller value", "Name": "value", "Value": "<caller value " + id + ">", "Key": "K", "Required": index%5 == 0, "Disabled": index%11 == 0, "Checked": index%2 == 0,
				"Before": "Caller before " + id, "After": "Caller after " + id, "Alt": name + " caller alternative text", "Caption": "Caller supplied caption", "Citation": "Caller source", "DateTime": "2026-10-11T09:30:00+09:00", "Href": "#" + id, "Src": "/assets/sample.svg",
				"Items": []string{"Caller first " + name, "Caller second " + name, "Caller third " + name},
			}
			switch name {
			case "Canvas":
				props["Points"] = "8,80 96,28 184,58 312,14"
			case "NumberText":
				props["Value"] = "12480"
			case "CurrencyText":
				props["Value"], props["Currency"] = "48000", "KRW"
			case "PercentText":
				props["Value"] = "0.742"
			case "FileSizeText":
				props["Value"] = "2516582"
			case "DurationText":
				props["Value"] = "6138"
			case "PluralText":
				props["Value"], props["Unit"] = "1", "source"
			case "RelativeTime", "DateTimeText":
				props["Value"] = "2026-10-11T09:30:00+09:00"
			case "AnimatedNumber":
				props["Value"], props["From"], props["Duration"] = "120", "0", 300
			case "Typewriter":
				props["Value"], props["Duration"] = "Caller typewriter text", 300
			}
			cases = append(cases, AdvancedCase{Family: "foundation", Component: ui.Component{Chapter: chapter, Name: name, ID: id, Props: props}})
		}
	}
	return cases
}
