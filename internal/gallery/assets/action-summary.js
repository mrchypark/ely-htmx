export function actionSummary(detail) {
 return {
  component: detail.component, action: detail.action,
  value: /\/(PasswordInput|PasswordField|APITokenList|ApiKeyManager)$/.test(detail.component) ? '[redacted]' : detail.value,
  values: detail.values?.map(([name,value])=>[name,/password|secret|authorization|api[_-]?key|(?:api|access|refresh)[_-]?token/i.test(name) || (detail.component?.startsWith('account/') && /^(token|key)$/i.test(name)) ? '[redacted]' : value instanceof File ? value.name : value]),
 };
}
