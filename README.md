# ServerCarte

## What is ServerCarte?

------------

ServerCarte is the core backend for both GastroAdmin dashboard and Tribeca application. Using this server, a restaurant owner can generate and edit content for the restaurant's menu. This includes funcitonality for adding add on choices, condiments, and side order options to each menu item.

## 3rd Party Libraries Used

* [GORM](https://gorm.io) (an abstraction layer for multiple sql based databases) 
* [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) (For JWT authentication)
* [Gin-Gonic](https://gin-gonic.com) (Go Web Framework. A little bit overkill for this project's use case)

## Architecture

This project is based on the principles of [Clean Architecture](https://archive.org/details/CleanArchitecture). It consists of four layers (Domains/Entities, Application Business Rules, Application Interface, Frameworks & Drivers) implementing a separation of concerns between each layer. The following directory structure is should demonstrate this architecture:

```
├── cmd
│   └── migration
├── domain (entities)
│   ├── account
|   |   ├── service.go
|   |   ├── model.go
|   |   └── repo.go
│   ├── menu
|   |   ├── service.go
|   |   ├── model.go
|   |   └── repo.go
│   ├── security
│   └── user
|   |   ├── service.go
|   |   ├── model.go
|   |   └── repo.go
├── internal (framework & drivers)
│   ├── account
│   │   ├── delivery
│   │   │   └── ginHTTP
│   │   └── repository
│   │       └── gorm
│   ├── authentication
│   │   ├── delivery
│   │   │   └── ginHTTP
│   │   └── framework
│   │       └── jwt
│   ├── config
│   ├── delivery
│   │   └── ginHTTP
│   ├── helpers
│   ├── logger
│   ├── menu
│   │   ├── delivery
│   │   │   └── ginHTTP
│   │   └── repository
│   │       └── gorm
│   ├── security
│   │   └── bcrypto
│   ├── store
│   │   └── gormDB
│   └── user
│       ├── delivery
│       │   └── ginHTTP
│       └── repository
│           └── gorm
└── server


```
Throughout this project I used the following template for each entity: 

```
entity
├── model (consistings of structs)
├── repository (interfaces for persistent storage)
└── service (interfaces containing application logic)
```

In `/internal` the actual technology or framework is specified and implements interfaces as specified above. For further references, [zhashkevych/go-clean-architecture](https://github.com/zhashkevych/go-clean-architecture), [err0r500/go-realworld-clean](https://github.com/err0r500/go-realworld-clean).

## API Documentation

-------------

This server provides the following endpoints
```
GET    /api/v1/menus       
GET    /api/v1/sections    
GET    /api/v1/sections/:id
GET    /api/v1/items       
GET    /api/v1/items/:id   
POST   /api/v1/sections    
PATCH  /api/v1/sections/:id
DELETE /api/v1/sections/:id
POST   /api/v1/items       
PATCH  /api/v1/items/:id   
DELETE /api/v1/items/:id   
GET    /user/:id           
PATCH  /user/:id           
DELETE /user/:id           
POST   /login              
GET    /accounts           
POST   /account            
PATCH  /account            
DELETE /account 
```

## Prerequisite

* Latest version of `Go`
* One of five databases (MySQL/ PostGreSQL, SQLite, SQLServer and Clickhouse)

## Getting Started

First need to specify the database that will act as persistent storage. Fortunately [GORM](https://gorm.io) provides a choice between MySQL/ PostGreSQL, SQLite, SQLServer and Clickhouse. Set up the preferred database with the appropriate user privileges and create a database. 

This project also provides a choice in encryption technology to use for JWT as well as password policy settings.

Populate `config.yml` and specify your configuration:

```yaml
database:
  type: <mysql|postgres|sqlite3> (default: mysql)
  host: <host> (default: localhost)
  port: <port> (default: 3306)
  user: <username>
  pass: <password>
  name: <dbName> (default: gastro)
server:
  host: <ip address for REST server> (default: localhost)
  port: <port> (default: 8080)
  read_timeout_seconds: <int> (default: 5)
  write_timeout_seconds: <int> (default: 5)
authentication:
  algorithm: <HS256|RSA|RSA-PSS|ECDSA> (default: HS256)
  expiration_period: <int in minutes> (default: 4320)
  minimum_key_length: <int in byte length> (default: 128)
  secret_key: <randomly generated string of characters at least 32 chars long>
security:
  length: <int> (default: 8)
  mixed_case: <true|false> (default: false)
  alpha_num: <true|false> (default: false)
  special_char: <true|false> (default: false)
  check_previous: <true|false> (default: false)
  ```
  
  _Hint: to generate a secret key run_
  `$ date +%s | sha256sum | base64 | head -c 64 ; echo`
  
  then 
  `$ go run cmd/main.go -c config.yml -p`
