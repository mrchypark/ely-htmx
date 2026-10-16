#!/usr/bin/env bash
# Source this pure validator without building the site.
valid_pages_base() {
  local value="${1:-}"
  [[ "$value" =~ ^/([A-Za-z0-9_.-]+/)*$ ]] || return 1
  [[ "$value" != *'/./'* && "$value" != *'/../'* ]]
}
