# Secret Auth

Authorise requests by secret cookie

## Configuration

```yaml
# Static configuration

experimental:
  plugins:
    example:
      moduleName: github.com/vslinko/secret-auth
      version: v0.1.0
```

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - secret-auth

  services:
    service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000

  middlewares:
    secret-auth:
      plugin:
        secret-auth:
          cookieName: myCookie
          secretKey: mySecretKey
```
