import assert from 'node:assert/strict';
import {evaluateExpression, validateCron} from './assets/ely-expressions.js';

assert.equal(evaluateExpression('price * quantity + (10 - 4) / 2', {price: 8, quantity: 3}), 27);
assert.equal(evaluateExpression('-(.5 + 2) * -2'), 5);
for (const source of ['1 / 0', '1 % 0', 'unknown + 1', '__proto__', 'process.exit()', '1 +', '1 2', '2 ** 3', '1..2', '('.repeat(40)+'1'+')'.repeat(40)]) assert.throws(() => evaluateExpression(source));
assert.equal(validateCron('*/5 0-23 * 1,6 0-7'), '*/5 0-23 * 1,6 0-7');
for (const source of ['* * * *', '60 * * * *', '* 24 * * *', '* * 0 * *', '* * * 13 *', '* * * * 8', '*/0 * * * *', '10-2 * * * *', 'a * * * *']) assert.throws(() => validateCron(source));
