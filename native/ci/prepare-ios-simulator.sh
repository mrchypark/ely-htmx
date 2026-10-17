#!/usr/bin/env bash
# Fail before UI tests if the selected simulator is not ready. No test retries.
set -euo pipefail
OUTPUT_DIR="${ELY_IOS_DIAGNOSTICS_DIR:-/tmp/ely-ios-diagnostics}"
mkdir -p "$OUTPUT_DIR"
xcodebuild -version > "$OUTPUT_DIR/xcode-version.txt"
xcrun simctl help bootstatus > "$OUTPUT_DIR/bootstatus-help.txt"
xcrun xcresulttool help export diagnostics > "$OUTPUT_DIR/diagnostics-help.txt"
xcrun simctl list devices available -j > "$OUTPUT_DIR/devices.json"
DEVICE_ID=$(python3 - "$OUTPUT_DIR/devices.json" <<'PY'
import json, sys
with open(sys.argv[1]) as source:
    devices = json.load(source)
for group in devices.get('devices', {}).values():
    for device in group:
        if 'iPhone' in device.get('name', '') and device.get('isAvailable', True):
            print(device['udid'])
            sys.exit(0)
sys.exit('No available iPhone simulator found')
PY
)
xcrun simctl bootstatus "$DEVICE_ID" -b 2>&1 | tee "$OUTPUT_DIR/bootstatus.log"
printf '%s\n' "$DEVICE_ID" > "$OUTPUT_DIR/device-id.txt"
if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
  printf 'device_id=%s\n' "$DEVICE_ID" >> "$GITHUB_OUTPUT"
fi
