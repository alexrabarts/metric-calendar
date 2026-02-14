remote := "alex@alex-het"
remote_path := "/srv/metricweek"

# Deploy to production
deploy:
    rsync -avz --delete www/ {{remote}}:{{remote_path}}/
    scp Caddyfile {{remote}}:{{remote_path}}/Caddyfile
    ssh {{remote}} "sudo systemctl reload caddy"

# Serve locally for testing
serve:
    cd www && python3 -m http.server 8000
