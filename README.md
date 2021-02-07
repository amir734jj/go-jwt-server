# go-jwt-server
Implements bare-bone JWT login/register/logout using:

- goji
- gorm
- jwt-go

```shell
POST /account/register
POST /account/login
GET  /account/logout
```

Auth can enforced using a simple middleware.

