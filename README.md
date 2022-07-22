# Boilerplate Learn Go
My boilerplate to learn backend in Golang, combining Echo &amp; net/http

## Source
- [Basic to Advance Go (Bahasa Indonesia)](https://dasarpemrogramangolang.novalagung.com/) - [Novalagung](https://github.com/novalagung)
- [JWT Auth](https://developer.vonage.com/blog/20/03/13/using-jwt-for-authentication-in-a-golang-application-dr) 
- [Echo Boilerplate](https://github.com/nixsolutions/golang-echo-boilerplate)
- [Secure Refresh Token Reuse Logic](https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/)

```
Folder
|
├─ cmd
│  └─ main.go
├─ configs
│  └─ config.go
├─ go.mod
├─ internal
│  ├─ handler
│  │  ├─ admin_handler.go
│  │  ├─ auth_handler.go
│  │  ├─ private_handler.go
│  │  └─ public_handler.go
│  ├─ middlewares
│  │  └─ middleware.go
│  ├─ models
│  │  └─ user.go
│  ├─ repository
│  │  └─ user
│  │     └─ user_repository.go
│  ├─ routes
│  │  └─ routes.go
│  └─ validations
│     ├─ auth_validation.go
│     └─ validate.go
├─ pkg
│  ├─ auth
│  │  ├─ auth.go
│  │  ├─ redis.go
│  │  └─ token.go
│  ├─ database
│  │  └─ connection.go
│  └─ utils
│     └─ password.go
└─ README.md

```