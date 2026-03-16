remote := "alex@alex-het"
remote_path := "/srv/metricweek"

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
