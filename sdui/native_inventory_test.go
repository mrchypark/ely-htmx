package sdui

import (
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"sort"
	"testing"
)

func TestNativeSupportInventoryMatchesPlatformCapabilities(t *testing.T) {
	data, err := os.ReadFile("../native/support-matrix.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		ProtocolVersion string                                        `json:"protocolVersion"`
		DefaultStatus   string                                        `json:"defaultStatus"`
		Components      []struct{ Type, IOS, Android, Limits string } `json:"components"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	if inventory.ProtocolVersion != NativeVersion || inventory.DefaultStatus != "unsupported" {
		t.Fatal("inventory must fail closed for this protocol version")
	}
	var names []string
	for _, c := range inventory.Components {
		if c.IOS == "" || c.Android == "" || c.Limits == "" {
			t.Fatalf("missing implementation or limits for %s", c.Type)
		}
		names = append(names, c.Type)
	}
	want := NativeCapabilities()
	sort.Strings(want)
	sort.Strings(names)
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("inventory %v does not match Go capabilities %v", names, want)
	}
	for _, source := range []struct{ path, pattern string }{
		{"../native/ios/Sources/ElyNativeCore/Protocol.swift", `public static let capabilities = \[([^\]]+)\]`},
		{"../native/android/protocol/src/main/kotlin/dev/ely/nativeprofile/Protocol.kt", `val CAPABILITIES = listOf\(([^)]+)\)`},
	} {
		content, err := os.ReadFile(source.path)
		if err != nil {
			t.Fatal(err)
		}
		match := regexp.MustCompile(source.pattern).FindSubmatch(content)
		if len(match) != 2 {
			t.Fatalf("cannot locate capabilities in %s", source.path)
		}
		var got []string
		for _, name := range regexp.MustCompile(`"([a-z]+)"`).FindAllSubmatch(match[1], -1) {
			got = append(got, string(name[1]))
		}
		sort.Strings(got)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s capabilities %v differ from inventory %v", source.path, got, want)
		}
	}
}
