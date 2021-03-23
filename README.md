# gastro-backend



To populate config.yml, set up a database that is supported by [GORM](https://gorm.io) (mysql/mariadb/postgres/sqlite3), and follow the following format:

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
