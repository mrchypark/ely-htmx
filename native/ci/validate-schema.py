"""Validate real Go fixtures and reject unusable native action shapes."""
import copy
import json
from pathlib import Path
from jsonschema import Draft202012Validator

root = Path(__file__).resolve().parents[2]
validator = Draft202012Validator(json.loads((root / 'native/profile-v1.schema.json').read_text()))
for path in (root / 'native/fixtures').glob('profile-*.json'):
    validator.validate(json.loads(path.read_text()))
    print('Validated', path.name)

ready = json.loads((root / 'native/fixtures/profile-controls-ready.json').read_text())
action = {key: ready[key] for key in ['protocolVersion', 'screenID', 'screenRevision', 'sessionID', 'requestID', 'dataVersion', 'fields']}
action.update(nodeID='profile-form', event='submit', action='save-profile', target='profile-region',
              editRevisions={key: 0 for key in ready['fields']}, capabilities=ready['requiredCapabilities'])
validator.validate(action)
invalid = copy.deepcopy(action)
invalid['editRevisions'] = {}
assert list(validator.iter_errors(invalid)), 'Empty edit revisions must fail schema validation'
for suffix in ['\n', '\r', '\r\n', '\u2028']:
    invalid = copy.deepcopy(action)
    invalid['fields']['company' + suffix] = invalid['fields'].pop('company')
    assert list(validator.iter_errors(invalid)), 'Field names must match the complete grammar'
print('Rejected empty revisions and trailing-line-break field names')
