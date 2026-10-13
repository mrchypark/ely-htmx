package uiexamples

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNativeMediaSourceFixtures435To450(t *testing.T) {
	fixtures := map[string]map[string]any{}
	for _, c := range advancedNativeCases() {
		if c.Component.Chapter == "media" {
			fixtures[c.Component.Name] = c.Component.Props
		}
	}
	for _, name := range []string{"ImageViewer", "ImageThumbnail", "ImageCropper", "ImageUpload", "ImageAnnotator", "VideoPlayer", "VideoTimeline", "VideoThumbnailStrip", "SubtitleEditor", "AudioPlayer", "AudioWaveform", "AudioSpectrum", "VolumeControl", "PlaybackSpeedControl", "MediaControls", "Playlist"} {
		p := fixtures[name]
		examples, ok := p["PreviewExamples"].([]map[string]any)
		if !ok || len(examples) != 1 {
			t.Fatalf("%s lost source preview", name)
		}
		for key, value := range examples[0] {
			a, _ := json.Marshal(value)
			b, _ := json.Marshal(p[key])
			if string(a) != string(b) {
				t.Fatalf("%s preview diverges at %s", name, key)
			}
		}
	}
	thumbs := fixtures["ImageThumbnail"]["Items"].([]map[string]any)
	if len(thumbs) != 5 || thumbs[2]["Content"] != "HDR" || thumbs[4]["Ratio"] != "1" {
		t.Fatal("source photos and shapes lost")
	}
	if fixtures["VideoPlayer"]["URL"] != "" || fixtures["VideoPlayer"]["Position"] != 34 || fixtures["AudioPlayer"]["URL"] != "" || fixtures["AudioPlayer"]["Duration"] != 222 {
		t.Fatal("source host playback data lost")
	}
	if fixtures["VideoThumbnailStrip"]["TrimStart"] != 12 || fixtures["VideoThumbnailStrip"]["TrimEnd"] != 66 {
		t.Fatal("source trim lost")
	}
	frames := fixtures["VideoThumbnailStrip"]["Frames"].([]map[string]any)
	if len(frames) != 8 || !strings.Contains(frames[2]["URL"].(string), "/atrium.jpg") || !strings.Contains(frames[3]["URL"].(string), "atrium-stair") || !strings.Contains(frames[6]["URL"].(string), "native-camera") {
		t.Fatal("source 3/3/2 frame distribution lost")
	}
	speeds := fixtures["PlaybackSpeedControl"]["SpeedOptions"].([]map[string]any)
	if len(speeds) != 6 || speeds[1]["Value"] != .75 || speeds[3]["Value"] != 1.25 {
		t.Fatal("source speeds lost")
	}
	playlist := fixtures["Playlist"]["Items"].([]map[string]any)
	if playlist[3]["Value"] != "dune-hours" || playlist[4]["Artist"] != "Marram" || fixtures["MediaControls"]["Repeat"] != "all" {
		t.Fatal("source tracks or repeat lost")
	}
	cues := fixtures["SubtitleEditor"]["Items"].([]map[string]any)
	if cues[1]["Start"] != "00:31.000" || cues[2]["End"] != "01:10.000" || cues[1]["Active"] != true || fixtures["SubtitleEditor"]["Position"] != "00:34.000" {
		t.Fatal("source cue timings lost")
	}
}
