remote := "alex@alex-het"
remote_path := "/srv/metricweek"

# Version is canonical in pyproject.toml; all lib versions should stay in sync
version := `grep '^version' libs/python/pyproject.toml | sed 's/version = "//;s/"//'`

# Release all libraries — see individual targets for credential requirements
release: release-js release-python release-go release-rust release-ruby

# Publish to npm (requires: npm login)
release-js:
    cd libs/js && npm publish

# Tag to trigger PyPI publish via GitHub Actions
release-python:
    git tag "libs/python/v{{version}}"
    git push origin "libs/python/v{{version}}"

# Tag Go module for pkg.go.dev
release-go:
    git tag "libs/go/v{{version}}"
    git push origin "libs/go/v{{version}}"

# Publish to crates.io (requires: cargo login)
release-rust:
    cd libs/rust && cargo publish

# Build gem and push to RubyGems (requires: gem signin)
release-ruby:
    #!/usr/bin/env bash
    set -euo pipefail
    cd libs/ruby
    gem build metric_calendar.gemspec
    gem push "MetricCalendar-{{version}}.gem"
    rm -f "MetricCalendar-{{version}}.gem"

# Build JS library and copy browser bundle to www/
build-js:
    cd libs/js && npm run build
    cp libs/js/dist/metric-calendar.iife.js www/metric-calendar.js

# Generate ICS calendar feeds
generate-calendar:
    go run ./cmd/generate-ics/

# Deploy to production
deploy: build-js generate-calendar
    rsync -avz --delete www/ {{remote}}:{{remote_path}}/
    scp Caddyfile {{remote}}:{{remote_path}}/Caddyfile
    ssh {{remote}} "sudo systemctl reload caddy"

# Serve locally for testing
serve:
    cd www && python3 -m http.server 8000
