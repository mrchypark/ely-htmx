package uiexamples

import (
	"strconv"
	"strings"

	"github.com/mrchypark/ely-htmx"
)

// advancedNativeCases supplies one caller-owned instance for every native primary.
func advancedNativeCases() []AdvancedCase {
	const stories = `shell|TitleBar
shell|ColumnShell
shell|WindowControls
shell|WindowDragRegion
shell|WindowResizeBorder
shell|NativeMenu
shell|Toolbar
shell|StatusBar
shell|ActivityBar
shell|TabBar
shell|SplashScreen
shell|AboutDialog
shell|TrayIcon
shell|DockBadge
shell|SystemNotification
shell|MultiWindow Manager
shell|WindowSwitcher
shell|QuickLauncher
shell|MiniWindow
shell|UpdateBanner
shell|UpdateDialog
shell|CrashReporter
shell|OfflineIndicator
shell|ZoomControl
generative|ImageGenerationPanel
generative|GenerationQueue
generative|GenerationGrid
generative|VariationPicker
generative|InpaintCanvas
generative|AudioGenerationPlayer
generative|TTSVoicePicker
generative|VideoGenerationTimeline
generative|ModelCard
generative|ModelDownloadManager
generative|ModelStatus
generative|FineTuneJobCard
generative|DatasetViewer
generative|EmbeddingVisualizer
generative|PromptPlayground
generative|PromptDiff
media|ImageViewer
media|ImageThumbnail
media|ImageCropper
media|ImageUpload
media|ImageAnnotator
media|VideoPlayer
media|VideoTimeline
media|VideoThumbnailStrip
media|SubtitleEditor
media|AudioPlayer
media|AudioWaveform
media|AudioSpectrum
media|VolumeControl
media|PlaybackSpeedControl
media|MediaControls
media|Playlist
media|ScreenRecorderControls
media|CameraPreview
media|MicLevelMeter
media|DeviceSelector
media|3DModelViewer
theme|Surfaces
theme|Ink
theme|Signals
theme|Series
theme|Syntax
theme|Type
theme|Shape
theme|Motion
theme|ThemeEditor
theme|ThemeImporter
theme|Platform Style
i18n|I18nProvider
i18n|LocaleNumber
i18n|RTL layout
i18n|SkipLink
i18n|ColorBlindSafe Palette
misc|Clock
misc|Stopwatch
misc|Calculator
misc|UnitConverter
misc|QRCodeScanner
misc|Captcha
misc|BarcodeDisplay
misc|Poll
misc|Quiz
misc|ChangeLog Viewer
misc|LicenseViewer
misc|Terms
misc|ExportDialog
misc|WebView
misc|IframeEmbed
misc|PrintButton
tooling|InspectorOverlay
tooling|FPSMeter
tooling|RenderCounter
tooling|EventLogger
tooling|DesignTokenViewer
tooling|IconBrowser
tooling|A11yChecker
tooling|ScreenshotTest Harness
tooling|ComponentGallery
tooling|PlaygroundPanel
rendering|RenderSurface`
	lines := strings.Split(stories, "\n")
	cases := make([]AdvancedCase, 0, len(lines))
	for index, line := range lines {
		part := strings.SplitN(line, "|", 2)
		props := map[string]any{
			"Title":        part[1] + " audit instance",
			"Description":  "Caller supplied content for isolation and escaped-value checks.",
			"Label":        "Caller value",
			"Name":         "native_value",
			"Value":        "value-" + strings.ReplaceAll(strings.ToLower(part[1]), " ", "-"),
			"Action":       "apply",
			"ActionLabel":  "Apply caller data",
			"Min":          0,
			"Max":          100,
			"URL":          "/assets/sample.svg",
			"Alt":          "Caller supplied media " + strconv.Itoa(index+1),
			"TimeZone":     "Asia/Seoul",
			"Rates":        "{\"first\":1,\"second\":2}",
			"Vertices":     "[[-80,-60,0],[80,-60,0],[0,80,0],[0,0,120]]",
			"Faces":        "[[0,1,2],[0,1,3],[1,2,3],[2,0,3]]",
			"CorrectValue": "first",
			"BeforeURL":    "/assets/sample-before.svg",
			"BeforeAlt":    "Caller before image",
			"AfterURL":     "/assets/sample.svg",
			"AfterAlt":     "Caller after image",
			"Compare":      50,
			"Samples": []map[string]any{
				{"X": 5, "Y": 12, "Width": 8, "Height": 18},
				{"X": 20, "Y": 4, "Width": 8, "Height": 34},
				{"X": 35, "Y": 16, "Width": 8, "Height": 14},
			},
			"Frames": []map[string]any{
				{"Time": 0, "URL": "/assets/sample.svg", "Alt": "Caller frame zero"},
				{"Time": 1, "URL": "/assets/sample.svg", "Alt": "Caller frame one"},
			},
			"Options": []map[string]any{
				{"Value": "first", "Label": "First caller option", "Selected": true, "Percent": 60},
				{"Value": "second", "Label": "Second caller option", "Percent": 40},
			},
			"Items": []map[string]any{
				{"Value": "first", "Label": "First caller item", "Content": "Caller-owned row", "URL": "/assets/sample.svg", "Alt": "Caller first image", "Progress": 25},
				{"Value": "second", "Label": "Second caller item", "Content": "Independent row", "URL": "/assets/sample.svg", "Alt": "Caller second image", "Progress": 75},
			},
		}
		if index%17 == 0 {
			props["Disabled"] = true
		}
		if part[0] == "i18n" && part[1] == "RTL layout" {
			props["RTL"] = true
		}
		if part[1] == "VideoPlayer" || part[1] == "VideoTimeline" || part[1] == "VideoThumbnailStrip" || part[1] == "VideoGenerationTimeline" {
			props["URL"] = "/assets/audit-video.mp4"
		} else if strings.Contains(part[1], "Audio") || part[1] == "MediaControls" || part[1] == "VolumeControl" || part[1] == "PlaybackSpeedControl" {
			props["URL"] = "/assets/audit-tone.wav"
		}
		cases = append(cases, AdvancedCase{Component: ui.Component{Chapter: part[0], Name: part[1], ID: "audit-" + part[0] + "-" + strconv.Itoa(index+1), Props: props}, Family: "native"})
	}
	return cases
}
