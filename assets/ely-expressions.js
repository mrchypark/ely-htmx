// Arithmetic is parsed as data. Caller expressions never execute JavaScript.
export function evaluateExpression(source, variables = {}) {
  const text = String(source);
  if (text.length > 1024) throw new Error('Expression is too long.');
  const tokens = text.match(/\d+(?:\.\d*)?|\.\d+|[A-Za-z_][A-Za-z0-9_]*|[()+\-*/%]|\S/g) || [];
  if (!tokens.length || tokens.length > 256) throw new Error('Enter a bounded arithmetic expression.');
  let cursor = 0;
  const number = (depth = 0) => {
    if (depth > 32) throw new Error('Expression nesting is too deep.');
    const token = tokens[cursor++];
    if (token === '+' || token === '-') return (token === '-' ? -1 : 1) * number(depth + 1);
    if (token === '(') {
      const value = sum(depth + 1);
      if (tokens[cursor++] !== ')') throw new Error('Missing closing parenthesis.');
      return value;
    }
    if (/^(?:\d+(?:\.\d*)?|\.\d+)$/.test(token || '')) return Number(token);
    if (/^[A-Za-z_][A-Za-z0-9_]*$/.test(token || '') && Object.hasOwn(variables, token)) {
      const value = Number(variables[token]);
      if (Number.isFinite(value)) return value;
    }
    throw new Error(`Unknown number or variable: ${token || 'end of expression'}`);
  };
  const product = depth => {
    let value = number(depth);
    while (['*', '/', '%'].includes(tokens[cursor])) {
      const operation = tokens[cursor++], right = number(depth);
      if ((operation === '/' || operation === '%') && right === 0) throw new Error('Division by zero.');
      value = operation === '*' ? value * right : operation === '/' ? value / right : value % right;
    }
    return value;
  };
  const sum = depth => {
    let value = product(depth);
    while (['+', '-'].includes(tokens[cursor])) {
      const operation = tokens[cursor++], right = product(depth);
      value += operation === '+' ? right : -right;
    }
    return value;
  };
  const result = sum(0);
  if (cursor !== tokens.length || !Number.isFinite(result)) throw new Error('Invalid arithmetic expression.');
  return result;
}

export function validateCron(source) {
  const fields = String(source).trim().split(/\s+/);
  if (fields.length !== 5) throw new Error('Cron requires minute, hour, day, month and weekday.');
  const bounds = [[0, 59], [0, 23], [1, 31], [1, 12], [0, 7]];
  fields.forEach((field, index) => {
    const [min, max] = bounds[index];
    for (const part of field.split(',')) {
      const match = /^(\*|\d+(?:-\d+)?)(?:\/(\d+))?$/.exec(part);
      if (!match || (match[2] && (!Number(match[2]) || Number(match[2]) > max - min + 1))) throw new Error(`Invalid cron field ${index + 1}.`);
      if (match[1] === '*') continue;
      const range = match[1].split('-').map(Number);
      if (range.some(value => value < min || value > max) || (range.length === 2 && range[0] > range[1])) throw new Error(`Cron field ${index + 1} is outside ${min}–${max}.`);
    }
  });
  return fields.join(' ');
}
