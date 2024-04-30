# A Proxy Server For Local Develpment

Hosts a proxy for the backend and frontend on `http://localhost:3255`

## Requirements
- nginx

## Setup

Add the following to your nginx configuration file within your http block:

```nginx
http {
    # ...
    include <pwd>/nginx.conf;
    include <pwd>/mime.types;
}
```

Make sure nginx (group `nginx` or `www-data`) has read access to the file.
