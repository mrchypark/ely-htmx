#!/usr/bin/env bash
# Run on a macOS Intel host with Java 21 and an already licensed Android SDK.
# This script never accepts SDK licenses or changes host security permissions.
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
cd "$ROOT"
RESULTS="$ROOT/test-results/android-native"
mkdir -p "$RESULTS"
export ANDROID_HOME="${ANDROID_HOME:-${ANDROID_SDK_ROOT:?An existing Android SDK is required}}"
export PATH="$ANDROID_HOME/cmdline-tools/latest/bin:$ANDROID_HOME/emulator:$ANDROID_HOME/platform-tools:$PATH"
SERIAL=emulator-5554
SERVER_PID=
EMULATOR_PID=
cleanup() {
  result=$?
  set +e
  adb -s "$SERIAL" logcat -d > "$RESULTS/logcat.txt" 2>&1
  adb -s "$SERIAL" pull /sdcard/Android/data/dev.ely.nativeprofile/files/screenshots "$RESULTS/"
  if [ -n "$EMULATOR_PID" ]; then kill "$EMULATOR_PID" 2>/dev/null; fi
  if [ -n "$SERVER_PID" ]; then kill "$SERVER_PID" 2>/dev/null; fi
  exit "$result"
}
trap cleanup EXIT

# stdin is closed: fail rather than accepting new/unexpected license terms.
sdkmanager 'system-images;android-35;google_apis;x86_64' < /dev/null
printf 'no\n' | avdmanager create avd --force --name ely-profile --device pixel_2 --package 'system-images;android-35;google_apis;x86_64'
emulator -accel-check
emulator -avd ely-profile -port 5554 -no-window -no-audio -no-boot-anim -no-snapshot -gpu swiftshader_indirect -accel on > "$RESULTS/emulator.log" 2>&1 &
EMULATOR_PID=$!
for attempt in $(seq 1 180); do
  kill -0 "$EMULATOR_PID" || { cat "$RESULTS/emulator.log"; exit 1; }
  if [ "$(adb -s "$SERIAL" shell getprop sys.boot_completed 2>/dev/null | tr -d '\r')" = 1 ]; then break; fi
  sleep 2
done
test "$(adb -s "$SERIAL" shell getprop sys.boot_completed | tr -d '\r')" = 1
adb -s "$SERIAL" shell input keyevent 82
adb -s "$SERIAL" shell getprop ro.build.fingerprint > "$RESULTS/device.txt"

SERVER_BINARY="${RUNNER_TEMP:-/tmp}/ely-native-emulator-profile"
go build -o "$SERVER_BINARY" ./examples/native-sdui
"$SERVER_BINARY" > "$RESULTS/server.log" 2>&1 &
SERVER_PID=$!
for attempt in $(seq 1 30); do
  curl --fail --silent http://127.0.0.1:8094/ > /dev/null && break
  sleep 1
done
curl --fail --silent http://127.0.0.1:8094/ > /dev/null
cd native/android
ANDROID_SERIAL="$SERIAL" ./gradlew --no-daemon :app:connectedDebugAndroidTest
