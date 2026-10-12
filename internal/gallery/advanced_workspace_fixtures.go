package uiexamples

import (
	"fmt"
	"strings"

	"github.com/mrchypark/ely-htmx"
)

// advancedWorkspaceCases gives the gallery one caller-owned case for each primary story.
func advancedWorkspaceCases() []AdvancedCase {
	stories := map[string][]string{
		"editor":    {"CodeEditor", "WhitespaceRenderer", "InlineDiff", "CompletionMenu", "CodeActionMenu", "FindWidget", "SearchPanel", "ProblemsPanel", "DebugConsole", "SettingsEditor", "WelcomePage", "CursorPosition"},
		"terminal":  {"Terminal", "CommandBlock", "AnsiText", "LogViewer", "ProcessList"},
		"git":       {"DiffViewer", "ThreeWayMerge", "ChangesList", "CommitGraph", "BranchSelector", "BlameView", "PullRequestCard"},
		"debug":     {"DebugToolbar", "BreakpointList", "VariablesPanel", "MemoryViewer", "PerformanceProfiler", "NetworkInspector"},
		"documents": {"RichTextEditor", "MarkdownEditor", "BlockEditor", "Block types", "TodoItem", "MarkdownRenderer", "DocumentViewer", "OfficePreview", "EpubReader", "Footnote", "FocusMode", "PrintPreview", "PageTree", "PageCover", "Backlinks", "TemplatePicker", "TrashBin", "VersionHistory"},
		"collab":    {"PresenceAvatars", "CommentThread", "CommentBubble", "CommentSidebar", "Annotation", "SuggestionMode", "ShareDialog", "ActivityLog"},
		"files":     {"FileIcon", "FileItem", "FileList", "FileGrid", "FileExplorer", "FolderTree", "FileContextMenu", "RenameInline", "DragDropFiles", "FilePreview", "RecentFiles", "Favorites", "FileSearch", "DuplicateFinder", "ArchiveViewer", "HexViewer", "FileOperationProgress", "TransferQueue", "DownloadManager", "StorageUsage"},
		"canvas":    {"InfiniteCanvas", "ToolPalette", "LayerPanel", "ColorPanel", "ConstraintsEditor", "NodeEditor", "FlowChart", "MindMap", "Whiteboard"},
		"devtools":  {"ConnectionManager", "SchemaTree", "QueryEditor", "TableStructureEditor", "RedisKeyBrowser", "DocumentViewer", "ApiRequestBuilder", "WebSocketConsole", "GraphQLExplorer", "RegexTester", "ContainerList"},
	}
	var cases []AdvancedCase
	for _, chapter := range []string{"editor", "terminal", "git", "debug", "documents", "collab", "files", "canvas", "devtools"} {
		names := stories[chapter]
		for index, name := range names {
			id := fmt.Sprintf("audit-%s-%d", chapter, index+1)
			value := "<caller value " + id + ">"
			if chapter == "documents" && (name == "MarkdownEditor" || name == "MarkdownRenderer") {
				value = "# Caller markdown\n\nSafe preview"
				if name == "MarkdownRenderer" {
					value += "\n\n## Caller details" + strings.Repeat("\n\nCaller reading content for the scroll container.", 20)
				}
			}
			props := map[string]any{
				"Title": name + " caller data", "Description": "Instance-specific workspace state", "Name": "payload", "Value": value, "Placeholder": "Enter a value", "Required": index%5 == 0, "Disabled": false,
				"Action": "operate", "ActionLabel": "Request host operation", "SubmitLabel": "Send to host", "Status": "Awaiting host", "Filter": "", "Zoom": "100", "CanvasText": name + " canvas", "Language": "go", "Cursor": "Ln 2, Col 8", "Query": "caller", "Command": "caller-command --safe", "Markdown": "# Caller markdown\n\n<escaped content>", "Content": "Caller rich text content", "Ours": "caller ours", "Theirs": "caller theirs", "Merged": "caller merged", "Method": "POST", "Endpoint": "/caller/resource", "Accept": ".txt,.md",
				"Lines": []map[string]any{{"Number": 1, "Text": "caller line one"}, {"Number": 2, "Text": "caller line two"}}, "Hunks": []map[string]any{{"Lines": []map[string]any{{"Type": "remove", "Prefix": "-", "Text": "old caller value"}, {"Type": "add", "Prefix": "+", "Text": "new caller value"}}}},
				"Rows": []map[string]any{{"Author": "Caller", "Revision": "r1", "Line": 1, "Text": "caller row", "Label": "Caller row", "Value": "caller-row", "Detail": "caller detail"}}, "Variables": []map[string]any{{"Name": "callerValue", "Value": "42", "Type": "int"}}, "Entries": []map[string]any{{"Text": "caller terminal output"}}, "History": []string{"caller-command --safe"},
				"Nodes": []map[string]any{{"ID": "n1", "X": 80, "Y": 80, "Label": "start"}, {"ID": "n2", "X": 260, "Y": 120, "Label": "finish"}}, "Edges": []map[string]any{{"FromX": 80, "FromY": 80, "ToX": 260, "ToY": 120}}, "Shapes": []map[string]any{{"ID": "shape-a", "X": 80, "Y": 80, "Width": 120, "Height": 64, "Label": "Caller shape"}},
				"Items":   []map[string]any{{"Value": id + "-one", "Label": "First " + name, "Detail": "caller item"}, {"Value": id + "-two", "Label": "Second " + name, "Detail": "independent item"}},
				"Options": []map[string]any{{"Value": "one", "Label": "One"}, {"Value": "two", "Label": "Two"}}, "Selected": "one",
			}
			switch chapter + "/" + name {
			case "editor/CompletionMenu":
				props["Value"] = "pri"
				props["Candidates"] = []map[string]any{{"ID": "complete-print", "Label": "printCaller", "InsertText": "printCaller()", "Detail": "caller function"}, {"ID": "complete-read", "Label": "readCaller", "InsertText": "readCaller()", "Detail": "caller reader"}}
			case "editor/CodeActionMenu":
				props["Path"], props["Value"] = "src/worker.go", "callerSource"
				props["CodeActions"] = []map[string]any{{"ID": "action-extract", "Label": "Extract caller function", "Command": "refactor.extract", "Detail": "Extract selected expression"}}
			case "editor/ProblemsPanel":
				props["Problems"] = []map[string]any{{"ID": "diag-1", "Severity": "warning", "Message": "Caller value is unused", "File": "src/worker.go", "Line": 27, "Column": 8}}
				props["Output"] = "Caller diagnostic output"
			case "editor/WhitespaceRenderer":
				props["Value"] = "a b\tc\n"
			case "editor/DebugConsole":
				props["Expression"], props["Session"], props["Output"] = "request.attempts", "session-12", "Caller debugger output"
			case "editor/SettingsEditor":
				props["Settings"] = []map[string]any{{"Key": "fontSize", "Label": "Font size", "Type": "number", "Value": 14, "Min": 10, "Max": 32, "Required": true}, {"Key": "language", "Label": "Language", "Value": "go"}}
			case "editor/WelcomePage":
				props["Projects"] = []map[string]any{{"Name": "Caller workspace", "URL": "#caller-project", "Detail": "Local project"}}
				props["Tasks"] = []map[string]any{{"ID": "task-check", "Label": "Check caller project"}}
			case "editor/CursorPosition":
				props["Value"], props["Lines"], props["Name"], props["Required"] = "first\nsecond", nil, "code", true
			case "git/CommitGraph":
				props["Nodes"] = []map[string]any{{"ID": "commit-parent", "X": 80, "Y": 80, "Label": "Caller parent"}, {"ID": "commit-child", "X": 260, "Y": 120, "Label": "Caller child", "Parents": []string{"commit-parent"}}}
				props["Edges"] = nil
			case "git/PullRequestCard":
				props["Number"], props["Subject"], props["PullRequestID"], props["Author"], props["Status"] = 42, "Caller review", "pr-42", "Caller reviewer", "Open"
				props["SourceBranch"], props["TargetBranch"], props["URL"] = "feature/worker", "main", "#caller-pr"
				props["Reviewers"] = []map[string]any{{"Name": "Reviewer one", "Status": "Pending"}}
				props["Comments"] = []map[string]any{{"ID": "comment-8", "Author": "Caller reviewer", "Time": "10:30", "Status": "Open", "Text": "Explain caller change", "Quote": "Caller passage"}}
			case "documents/Block types":
				props["Content"], props["Name"], props["Required"] = "Caller body", "content", true
				props["Blocks"] = []map[string]any{{"ID": "block-heading", "Type": "h2", "Label": "Insert heading", "Text": "Caller inserted heading"}, {"ID": "block-code", "Type": "pre", "Label": "Insert code", "Text": "callerCode()"}}
			case "documents/Footnote":
				props["Text"], props["Term"], props["Definition"], props["Number"] = "Caller passage", "CallerTerm", "Caller footnote definition", 3
			case "documents/FocusMode":
				props["Content"], props["Description"] = "Caller focused paragraph", "Reader tools"
			case "documents/PrintPreview":
				props["DocumentID"], props["Paper"], props["Orientation"] = "doc-21", "A4", "portrait"
				props["Pages"] = []map[string]any{{"Title": "Caller first page", "Content": "Caller page one"}, {"Title": "Caller second page", "Content": "Caller page two"}}
			case "documents/PageTree":
				props["Items"] = []map[string]any{{"Value": "handbook", "Label": "Caller handbook", "Children": []map[string]any{{"Value": "handbook/lift", "Label": "Caller lift page"}}}}
			case "documents/PageCover":
				props["ImageURL"], props["Alt"], props["Icon"], props["Content"] = "/assets/sample.svg", "Caller cover image", "C", "Caller cover description"
				props["Properties"] = []map[string]any{{"Name": "Status", "Value": "Draft"}}
			case "documents/Backlinks":
				props["Links"] = []map[string]any{{"URL": "#caller-source", "Title": "Caller source page", "Excerpt": "Caller linked passage"}}
			case "documents/TemplatePicker":
				props["Templates"] = []map[string]any{{"ID": "tpl-brief", "Name": "Caller brief", "Description": "Caller reusable document", "Content": "Caller template content"}}
			case "collab/PresenceAvatars":
				props["People"] = []map[string]any{{"ID": "user-7", "Name": "Caller writer", "Initials": "CW", "Status": "Editing", "Selection": "Caller selected phrase"}}
			case "collab/CommentBubble", "collab/CommentSidebar", "collab/Annotation":
				props["Comments"] = []map[string]any{{"ID": "comment-8", "Author": "Caller reviewer", "Time": "10:30", "Status": "Open", "Text": "Explain caller change", "Quote": "Caller passage"}}
			case "collab/SuggestionMode":
				props["Suggestions"] = []map[string]any{{"ID": "suggestion-9", "Author": "Caller reviewer", "Status": "Pending", "Original": "Old caller text", "Replacement": "New caller text"}}
			case "collab/ActivityLog":
				props["Activities"] = []map[string]any{{"Author": "Caller writer", "Action": "edited", "Time": "10:40", "DateTime": "2026-10-11T10:40:00+09:00", "Target": "Caller notes", "URL": "#caller-notes", "Detail": "Caller timeline detail"}}
			case "editor/CodeEditor", "documents/RichTextEditor", "documents/BlockEditor", "collab/CommentThread", "devtools/QueryEditor":
				props["Required"] = true
				if chapter == "collab" {
					props["Comments"] = []map[string]any{{"ID": "comment-8", "Author": "Caller reviewer", "Time": "10:30", "Status": "Open", "Text": "Explain caller change", "Quote": "Caller passage"}}
				}
				if chapter == "devtools" {
					props["Name"] = "query"
				}
			case "devtools/ApiRequestBuilder", "devtools/ConnectionManager", "devtools/WebSocketConsole", "collab/ShareDialog":
				props["MethodOptions"] = []map[string]any{{"Value": "GET", "Label": "GET"}, {"Value": "POST", "Label": "POST"}, {"Value": "PATCH", "Label": "PATCH"}}
			case "files/FileIcon", "files/FileItem":
				props["Name"], props["Path"], props["Extension"], props["Icon"], props["Size"], props["Modified"] = "brief.md", "docs/brief.md", "md", "MD", "4 KB", "2026-10-11"
			case "files/FileGrid", "files/RecentFiles", "files/Favorites":
				props["Files"] = []map[string]any{{"Name": "brief.md", "Path": "docs/brief.md", "Extension": "md", "Icon": "MD", "Size": "4 KB", "Modified": "2026-10-11", "Favorite": true}, {"Name": "budget.csv", "Path": "data/budget.csv", "Extension": "csv", "Icon": "CSV", "Size": "12 KB", "Modified": "2026-10-10", "Favorite": false}}
			case "files/FileContextMenu":
				props["Path"], props["Actions"] = "docs/brief.md", []map[string]any{{"Action": "rename-file", "Label": "Rename with host"}, {"Action": "move-file", "Label": "Move with host"}}
			case "files/DuplicateFinder":
				props["DuplicateGroups"] = []map[string]any{{"Hash": "sha256:caller-duplicate", "Files": []map[string]any{{"Path": "docs/brief-copy.md", "Size": "4 KB"}, {"Path": "archive/brief.md", "Size": "4 KB"}}}}
			case "files/FileExplorer", "files/FileList", "files/FolderTree":
				props["Items"] = []map[string]any{{"Value": "docs", "Label": "Documents", "Children": []map[string]any{{"Value": "docs/brief.md", "Label": "brief.md", "Detail": "4 KB"}}}}
			case "files/RenameInline":
				props["Path"], props["Name"], props["Value"], props["Required"] = "docs/brief.md", "name", "brief.md", true
			case "files/FileSearch", "editor/SearchPanel":
				props["Query"] = ""
				props["Items"] = []map[string]any{{"Path": "src/worker.go", "Detail": "worker result"}, {"Path": "docs/brief.md", "Detail": "brief result"}}
			case "files/FilePreview":
				props["Content"] = "Read-only caller file <script>content</script>"
			case "files/HexViewer", "debug/MemoryViewer":
				props["Bytes"], props["Offset"] = []int{0, 65, 127, 255}, 16
			case "files/FileOperationProgress":
				props["Progress"], props["Total"], props["Label"], props["Status"] = 37, 80, "Copy caller files", "Running"
			case "files/TransferQueue", "files/DownloadManager":
				props["Transfers"] = []map[string]any{{"ID": "transfer-7", "Name": "report.pdf", "Progress": 32, "Total": 80, "Status": "Paused"}}
			case "files/StorageUsage":
				props["Used"], props["Capacity"], props["Unit"], props["Label"] = 240, 500, "GB", "Caller storage"
				props["Categories"] = []map[string]any{{"Name": "Documents", "Size": 180}, {"Name": "Images", "Size": 60}}
			case "devtools/SchemaTree":
				props["Nodes"] = []map[string]any{{"Name": "public", "Path": "public", "Kind": "schema", "Children": []map[string]any{{"Name": "projects", "Path": "public.projects", "Kind": "table", "Children": []map[string]any{{"Name": "id", "Path": "public.projects.id", "Kind": "uuid"}}}}}}
			case "devtools/TableStructureEditor":
				props["Action"] = "update-table-structure"
				props["Columns"] = []map[string]any{{"Name": "project_id", "Type": "uuid", "Nullable": false, "TypeOptions": []map[string]any{{"Value": "uuid", "Label": "UUID"}, {"Value": "text", "Label": "Text"}}}, {"Name": "notes", "Type": "text", "Nullable": true, "TypeOptions": []map[string]any{{"Value": "uuid", "Label": "UUID"}, {"Value": "text", "Label": "Text"}}}}
			case "devtools/RedisKeyBrowser":
				props["Keys"], props["Preview"] = []map[string]any{{"Key": "session:42", "Type": "string", "Preview": "caller-session-token"}, {"Key": "project:7", "Type": "hash", "Preview": "{\"name\":\"Caller project\"}"}}, "Select a caller key"
			case "devtools/DocumentViewer":
				props["Fields"], props["JSON"] = []map[string]any{{"Name": "_id", "Type": "string", "Value": "project-7"}, {"Name": "active", "Type": "boolean", "Value": "true"}}, "{\"_id\":\"project-7\",\"active\":true}"
			case "devtools/ContainerList":
				props["Containers"] = []map[string]any{{"ID": "a1b2c3d4", "Image": "caller/api:1.2", "Status": "running", "Ports": "8080:8080"}}
			case "devtools/GraphQLExplorer":
				props["Action"], props["Name"], props["Query"], props["VariablesJSON"], props["Required"] = "graphql-query", "query", "query Project($id: ID!) { project(id: $id) { name } }", "{\"id\":\"project-7\"}", true
			case "devtools/RegexTester":
				props["Action"], props["Pattern"], props["Text"], props["Required"] = "evaluate-regex", "caller-[0-9]+", "caller-42\nother", true
			case "debug/DebugToolbar":
				props["Session"], props["Status"] = "session-12", "Paused"
			case "debug/BreakpointList":
				props["Breakpoints"] = []map[string]any{{"ID": "bp-4", "File": "src/worker.go", "Line": 27, "Enabled": true, "Condition": "count > 3"}}
				props["Stack"] = []map[string]any{{"ID": "frame-2", "Function": "HandleRequest", "File": "src/worker.go", "Line": 27}}
			case "debug/VariablesPanel":
				props["Variables"] = []map[string]any{{"Name": "request", "Type": "Request", "Value": "{…}", "Children": []map[string]any{{"Name": "attempts", "Type": "int", "Value": 4}}}}
				props["Expression"] = "request.attempts"
			case "debug/PerformanceProfiler":
				props["Duration"] = 200
				props["Samples"] = []map[string]any{{"Name": "ParseRequest", "Duration": 40}, {"Name": "Render", "Duration": 120}}
			case "debug/NetworkInspector":
				props["Requests"] = []map[string]any{{"Method": "POST", "URL": "/caller/notes", "Status": 202, "Duration": 18, "Headers": []map[string]any{{"Name": "Content-Type", "Value": "application/json"}}, "Body": "{\"accepted\":true}"}}
			case "terminal/ProcessList":
				props["Processes"] = []map[string]any{{"PID": 4242, "Name": "caller-worker", "CPU": 12.5, "Memory": "48 MB"}}
			case "terminal/AnsiText", "terminal/LogViewer":
				props["Query"] = ""
				props["Entries"] = []map[string]any{{"Time": "10:32", "Level": "info", "Text": "caller ready"}, {"Time": "10:33", "Level": "error", "Text": "caller retry"}}
			case "git/BranchSelector":
				props["Selected"] = "main"
				props["Branches"] = []map[string]any{{"Value": "main", "Label": "main"}, {"Value": "feature/worker", "Label": "feature/worker"}}
			case "git/ChangesList":
				props["Changes"] = []map[string]any{{"Path": "src/worker.go", "Status": "Modified"}, {"Path": "docs/brief.md", "Status": "Added"}}
			case "documents/TodoItem":
				props["Todos"] = []map[string]any{{"ID": "todo-1", "Text": "Review caller brief", "Checked": true}, {"ID": "todo-2", "Text": "Publish brief", "Checked": false}}
			case "documents/DocumentViewer", "documents/OfficePreview", "documents/EpubReader":
				props["Chapters"] = []map[string]any{{"Title": "Overview", "Content": "Caller overview"}, {"Title": "Details", "Content": "Caller document details"}}
				if name == "OfficePreview" {
					props["Columns"], props["Cells"] = []string{"Project", "Budget"}, [][]string{{"Caller", "1200"}}
				}
			case "documents/VersionHistory":
				props["Revisions"] = []map[string]any{{"ID": "rev-1", "Author": "Reviewer", "Time": "09:00", "Content": "Original caller revision"}, {"ID": "rev-2", "Author": "Writer", "Time": "10:00", "Content": "Updated caller revision"}}
			case "documents/TrashBin":
				props["Deleted"] = []map[string]any{{"ID": "deleted-8", "Path": "docs/retired.md", "DeletedAt": "2026-10-09"}}
			case "canvas/ToolPalette":
				props["Tool"] = "select"
				props["Tools"] = []map[string]any{{"Value": "select", "Label": "Select"}, {"Value": "rectangle", "Label": "Rectangle"}}
			case "canvas/LayerPanel":
				props["Layers"] = []map[string]any{{"ID": "shape-a", "Name": "Title layer", "Visible": true}}
			case "canvas/ColorPanel":
				props["Color"], props["Opacity"] = "#336699", 85
			case "canvas/ConstraintsEditor":
				props["Width"], props["Height"], props["Gap"], props["AutoLayout"] = 320, 180, 16, true
			case "canvas/NodeEditor":
				props["Nodes"] = []map[string]any{{"ID": "node-source", "X": 40, "Y": 70, "Width": 130, "Height": 60, "Label": "Read caller"}, {"ID": "node-target", "X": 300, "Y": 70, "Width": 140, "Height": 60, "Label": "Transform caller"}}
				props["Edges"] = []map[string]any{{"FromX": 170, "FromY": 100, "ToX": 300, "ToY": 100}}
			}
			cases = append(cases, AdvancedCase{Family: "workspace", Component: ui.Component{Chapter: chapter, Name: name, ID: id, Props: props}})
		}
	}
	return cases
}
