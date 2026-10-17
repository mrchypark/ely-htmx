"""Offline tests for the readiness gate; simulator behavior requires macOS CI."""
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

SCRIPT = Path(__file__).with_name('prepare-ios-simulator.sh').resolve()


class ReadinessTests(unittest.TestCase):
    def run_gate(self, devices, boot_exit=0):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / 'xcodebuild').write_text('#!/bin/sh\necho "Xcode mock"\n')
            (root / 'xcrun').write_text('''#!/bin/sh
if [ "$2" = "help" ]; then
  exit 0
elif [ "$1 $2" = "simctl list" ]; then
  cat "$TEST_DEVICES"
else
  test "$1 $2 $3 $4" = "simctl bootstatus test-udid -b" || exit 99
  echo "mock readiness"
  exit "$TEST_BOOT_EXIT"
fi
''')
            for command in ('xcodebuild', 'xcrun'):
                (root / command).chmod(0o755)
            (root / 'devices.json').write_text(json.dumps(devices))
            output = root / 'output'
            env = dict(os.environ, PATH=f'{root}:{os.environ["PATH"]}',
                       TEST_DEVICES=str(root / 'devices.json'),
                       TEST_BOOT_EXIT=str(boot_exit), GITHUB_OUTPUT=str(output),
                       ELY_IOS_DIAGNOSTICS_DIR=str(root / 'diagnostics'))
            result = subprocess.run(['bash', str(SCRIPT)], env=env, capture_output=True, text=True)
            return result.returncode, output.read_text() if output.exists() else ''

    def test_ready(self):
        self.assertEqual(self.run_gate({'devices': {'runtime': [{'name': 'iPhone', 'udid': 'test-udid'}]}}),
                         (0, 'device_id=test-udid\n'))

    def test_empty_inventory(self):
        code, output = self.run_gate({'devices': {}})
        self.assertNotEqual(code, 0)
        self.assertEqual(output, '')

    def test_failed_boot_does_not_publish_device(self):
        code, output = self.run_gate({'devices': {'runtime': [{'name': 'iPhone', 'udid': 'test-udid'}]}}, 7)
        self.assertEqual(code, 7)
        self.assertEqual(output, '')


if __name__ == '__main__':
    unittest.main()
