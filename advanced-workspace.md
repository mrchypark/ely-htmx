# Workspace component API

Use `ui.Component{Chapter, Name, ID, Props}`. The parent resolver owns the outer wrapper and protected ID, Chapter, Kind, Variant and ComponentName. Supply a unique ID per instance. No example content is supplied by the production templates.

Common properties: `Title`, `Description`, `Disabled`, `Required`, `Name`, `Value`, `Action`, `ActionLabel`, and optional `ActionValue`. `ws-action` accepts `Type` (button by default, submit inside editor forms). Caller strings are escaped. Read-only displays remain read-only.

## Editing and host payloads

CodeEditor and the source/query renderers use a contenteditable textbox with a named hidden input, not a textarea. Supply `Value` or `Lines[{Number,Text}]`, plus `Language` and `Cursor`. Lines initialize the visible code when supplied. The hidden field is synchronized on initialization, input and submission; line numbers are excluded. `Name` defaults to code. QueryEditor can use Name=query. Submit emits Action or save-code with current plain code in native `FormData.get(Name)`.

RichTextEditor, BlockEditor and CommentThread use a contenteditable textbox initialized with escaped `Content`. Name defaults to content. Submit emits Action or save-content with current serialized content. Serialization preserves paragraph, line-break, emphasis, list, heading, quote and code tags, removes attributes, unwraps unsupported tags, and drops active script/style/embed elements. It does not trust caller HTML. Both editor forms block empty Required content before the base dispatcher runs, expose an alert and set aria-invalid. Disabled/Readonly editor surfaces cannot submit edits.

MarkdownEditor uses a native named textarea initialized from `Value`. Initial preview uses the shared safe Goldmark GFM helper. Real textarea input emits preview-markdown with `source`; it does not locally parse text or pretend that preview succeeded. The gallery host POSTs source to /components/markdown-preview and replaces preview nodes with the safe MarkdownRenderer response. Production hosts must bind the same action or provide their own renderer. Last rendered preview remains until a host response.

MarkdownRenderer uses a read-only article with safe GFM from `Value`; raw HTML and dangerous URL schemes are disabled by the shared helper. TableOfContents and DocumentOutline aliases generate links from the actual rendered headings. Heading IDs start with the wrapper ID. ReadingProgress measures its own scroll region with a native progress element.

## Domain properties and actions

| Components | Caller properties | Observable local behavior / host request |
| --- | --- | --- |
| FindWidget, SearchPanel, FileSearch | Query; Items with Path/Label/Detail | Native search/filter, visible result count. |
| DiffViewer, InlineDiff | Hunks with Lines(Type,Prefix,Text) | Escaped add/remove/context hunk display; Action requests host work. |
| ThreeWayMerge | Ours, Theirs, Merged | Native resolution radios and editable merged result; resolve-conflict (or Action) submits resolution and merged. |
| CommitGraph, BlameView | Nodes(ID,X,Y,Label), Edges(FromX,FromY,ToX,ToY); Rows(Author,Revision,Line,Text) | SVG graph geometry and selectable blame rows. |
| Terminal, CommandBlock | Entries(Text), Output, Command, History | History fills current command input; execute (or Action) submits command to host. Output only comes from caller data. |
| AnsiText, LogViewer | Entries(Time,Level,Text), Query | Read-only log rows with native filter. ANSI terminal emulation is not implemented. |
| ProcessList | Processes(PID,Name,CPU,Memory) | Domain table; terminate-process requests PID, without removing a row or claiming termination. |
| BranchSelector, ChangesList | Branches(Value,Label), Selected; Changes(Path,Status), Message | Native branch select emits switch-branch with branch; stage-change requests path. CommitInput alias submits current message via commit. |
| DebugToolbar | Session, Status | debug-run/debug-pause/debug-step/debug-stop request session; displayed status remains caller-owned. |
| BreakpointList | Breakpoints(ID,File,Line,Enabled,Condition), Stack(ID,Function,File,Line) | Native checkbox emits toggle-breakpoint with breakpoint/file/line/enabled; stack selection stays local. |
| VariablesPanel | Recursive Variables(Name,Type,Value,Children), Expression | Native recursive disclosure; ConsoleREPL alias emits evaluate-watch with current expression. |
| MemoryViewer, HexViewer | Bytes (integer bytes), Offset, Label | Read-only offset and hexadecimal byte cells. |
| PerformanceProfiler | Duration, Samples(Name,Duration) | Native duration meters using caller totals. |
| NetworkInspector | Requests(Method,URL,Status,Duration,Headers(Name,Value),Body) | Read-only metadata and native response disclosure. No request is sent. |
| TodoItem | Todos(ID,Text,Checked) | Native checkboxes and completed count; base ely:change carries checkbox value/checked. |
| DocumentViewer, OfficePreview, EpubReader | Content/Value, URL, Chapters(Title,Content); Columns, Cells for OfficePreview | Read-only caller text/table or validated link; chapter navigation. Binary document/EPUB parsing requires host-supplied content. |
| VersionHistory, TrashBin | Revisions(ID,Author,Time,Content); Deleted(ID,Path,DeletedAt) | Local revision selection; restore-file requests deleted ID. |
| RenameInline | Path, Name(default name), Value | Native rename form submits original path and edited name via rename; Cancel restores initial name. |
| FileExplorer, FileList, FolderTree | Recursive Items(Value,Label,Detail,Children), Filter | Collapse/expand, local selection and filter. |
| DragDropFiles, FilePreview, ArchiveViewer | Accept; ImageURL/Alt or Content/Value; archive Content | Native multiple file chooser sends actual FormData files via upload (or Action); read-only escaped preview/manifest. Archive extraction requires host processing. |
| FileOperationProgress | Progress, Total, Label, Status | Native caller-driven progress. No operation is started or marked successful. |
| TransferQueue, DownloadManager | Transfers(ID,Name,Progress,Total,Status) | Native progress; resume-transfer/cancel-transfer request transfer ID. |
| StorageUsage | Used, Capacity, Unit, Categories(Name,Size) | Native storage meter and category sizes. |
| FileIcon, FileItem, FileGrid, RecentFiles | File fields Name, Path, Extension, Icon, Size, Modified; FileGrid/RecentFiles take Files | Escaped file icon and metadata cards. Recent file time and path remain caller-owned. |
| FileContextMenu, Favorites, DuplicateFinder | Path and Actions(Action,Label); Files(Favorite); DuplicateGroups(Hash,Files(Path,Size)) | Context buttons send the caller path. Favorite is a native checkbox. Duplicate deletion only requests delete-duplicate with a path; no file is removed until the host responds. |
| Canvas/NodeEditor | Shapes or Nodes(ID,X,Y,Width,Height,Label), Edges(FromX,FromY,ToX,ToY), Zoom | Caller-driven SVG; keyboard selection emits select-shape; pointer movement changes SVG position; zoom changes visible CSSOM transform. |
| ToolPalette, LayerPanel | Tools(Value,Label), Tool; Layers(ID,Name,Visible) | Native tool select and matching shape visibility. |
| ColorPanel, ConstraintsEditor | Color, Opacity; Width, Height, Gap, AutoLayout | Color preview; update-constraints submits fields. ExportPanel alias requests export-selection with current form fields. |
| ApiRequestBuilder, ConnectionManager, WebSocketConsole, ShareDialog | Method, MethodOptions(Value,Label), Endpoint, Value(body) | Native method/endpoint/body form emits Action or request with current fields. Default methods are GET/POST/PUT/PATCH/DELETE/HEAD/OPTIONS; MethodOptions replaces them. Provider-specific protocols, sharing and connection establishment require host logic. |
| SchemaTree, TableStructureEditor | Recursive Nodes(Name,Path,Kind,Children); Columns(Name,Type,Nullable,TypeOptions) | Native recursive disclosure; table column name/type/nullability form submits current fields to update-table-structure or Action. |
| RedisKeyBrowser, devtools DocumentViewer | Keys(Key,Type,Preview), Preview; Fields(Name,Type,Value), JSON | Local key filtering and preview selection emits select-redis-key; deletion requests delete-redis-key. DocumentViewer is read-only escaped key/value and JSON display. |
| ContainerList, GraphQLExplorer, RegexTester | Containers(ID,Image,Status,Ports); Query, VariablesJSON; Pattern, Text | Start/stop requests carry container ID without status changes. GraphQL and regex forms submit current named data to the host; regex evaluation is not performed in the browser. Their result output changes only after a same-action `ely:result` or `ely:error`; objects are safely serialized as text, and a host may instead re-render controlled caller properties. |

## Additional editor, document and collaboration APIs

| Components | Caller properties | Local behavior / host request |
| --- | --- | --- |
| CompletionMenu | Candidates(ID,Label,InsertText,Detail), Value | Native source textarea and candidate search. Selection replaces the selected source range and emits select-completion with candidate, insertText and current source. |
| CodeActionMenu | CodeActions(ID,Label,Command,Detail,Disabled), Path, Value | apply-code-action carries action ID plus command, path and source in FormData. Commands are opaque host identifiers and are never executed locally. |
| ProblemsPanel | Problems(ID,Severity,Message,File,Line,Column), Output | Domain diagnostics and open-location with complete named location fields. |
| WhitespaceRenderer, CursorPosition | Value; Lines, Language, Cursor for source editing | Read-only whitespace shows space/tab/newline markers. CursorPosition uses a real editable selection and updates line/column; current text is submitted through the named code field. |
| DebugConsole, SettingsEditor | Expression, Session, Output; Settings(Key,Label,Type,Value,Required,Min,Max) | Native expression form submits session to Action/evaluate-expression. Native bounded settings fields submit current key/value pairs to Action/save-settings. Outputs stay caller-owned. |
| WelcomePage | Projects(Name,URL,Detail), Tasks(ID,Label) | Real project links and run-task host requests. |
| PullRequestCard | PullRequestID, Number, Subject, Author, Status, SourceBranch, TargetBranch, URL, Reviewers(Name,Status), Comments | Metadata, request-review and close-pull-request requests. ReviewComment alias adds the same comment reply/resolve API below. |
| Block types | Blocks(ID,Type,Label,Text), Content, Name | Selectable insertion appends a safe paragraph/heading/quote/code block into the editable content. Unknown types become paragraphs. Submit returns serialized inserted content. |
| Footnote | Text, Number, Definition, Term | Real local note link with instance-prefixed target. Glossary/TooltipTerm aliases use native details disclosure. |
| FocusMode, PrintPreview | Content, Description; DocumentID, Paper, Orientation, Pages(Title,Content) | Focus toggle hides local chrome and focuses read-only content. Print preview is read-only; Action/print-document submits document, paper, orientation and all page content. No native print dialog is opened automatically. |
| PageTree, PageCover, Backlinks | Recursive Items(Value,Label,Children); ImageURL, Alt, Icon, Content, Properties(Name,Value); Links(URL,Title,Excerpt) | Recursive page disclosure, actual cover image, read-only properties and real backlink links. GraphView alias includes selectable caller SVG Nodes/Edges. |
| TemplatePicker | Templates(ID,Name,Description,Content) | select-template carries template ID and complete caller content in FormData. |
| PresenceAvatars | People(ID,Name,Initials,AvatarURL,Status,Selection) | Accessible avatar/initials and caller presence status; follow-user requests participant ID. No fabricated live activity. |
| CommentBubble, CommentSidebar, Annotation, CommentThread | Comments(ID,Author,Time,Status,Quote,Text), Content for new comment editor | Author/text/status and quotes (marked for Annotation); resolve-comment carries comment ID and form data. reply-comment submits comment ID and current required reply. CommentThread also has the rich content form. |
| SuggestionMode | Suggestions(ID,Author,Status,Original,Replacement) | Native del/ins text; accept-suggestion/reject-suggestion carry ID, original and replacement. Status changes require host results. |
| ActivityLog | Activities(Author,Action,Time,DateTime,Target,URL,Detail) | Read-only ordered timeline with native time and target links. |

Rich editable forms have Bold, Italic and Code controls. Select text in the instance, then apply formatting. These controls use native DOM Range extraction/insertion without execCommand. The current safe serialization is returned through the named hidden field. Active content and attributes are removed; limits of 100 nesting levels and 10000 nodes produce a visible error and block submission rather than silently truncating.

Canvas backgrounds support pointer pan with a real CSSOM transform combined with zoom. Whiteboard adds a Draw mode: native pointer samples create an SVG polyline; draw-stroke emits points, color and width for host persistence. A stroke is bounded to 5000 samples. Pointer cancel removes an unfinished local stroke. Shape movement and keyboard selection remain instance-local. CommitGraph can derive SVG edges from Nodes' Parents IDs; absent parents are ignored. ThreeWayMerge choices now populate the merged textarea from the selected caller version before user edits and host submission.

The base module owns generic action dispatch: forms include native FormData and values, buttons include their value and any containing form data. Custom operations supply additional fields through emit. Selection/edit/filter/canvas changes use instance-local ely:change. Consumers must read current editor data from FormData, not original props or submitter labels. No component performs provider calls, executes arbitrary code, or reports provider success without host results.

## Coverage and verification limits

Nine chapter templates and synthetic caller-data fixtures cover 96 primary names: editor 12, terminal 5, git 7, debug 6, documents 18, collab 8, files 20, canvas 9, devtools 11. Focused Go tests validate rendering, escaping, native domain semantics, aliases and editor form landmarks. Gallery-only browser checks validate current submitted editor/query content, rich serialization, required errors, method selection, real Markdown preview responses, Markdown alias navigation/progress and upgraded domain interactions. The parent runs browser checks on all aliases.

All 96 primaries have semantic renderers and shared controls. Their checks require actual controls and test caller data, current submitted fields, local selection/formatting/insertion, link/disclosure semantics or computed canvas transforms and stroke geometry. The gallery audit checks each registered qualified name using two independent instances.

Remaining domain limits: code editing does not provide language parsing, folding, multicursor editing, Vim keymaps or a language server; completion candidates and code actions are caller/host supplied. Complex table/image/math block editing, mentions and link-format dialogs are not implemented by the basic rich formatter. Canvas has no automatic graph layout, snapping, pressure-sensitive brush, path editing or undo history. Live collaboration, invitation permissions, persistent review/resolve state, printing/export, debug expression evaluation and task execution require host backends. Aliases share the primary API except the explicit variant controls documented above; alias-specific advanced features must not be inferred from catalogue names.
