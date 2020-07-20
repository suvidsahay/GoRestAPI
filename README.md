# Rest API Written in Golang

## Setting up a development environment

```bash
# Get the code
git clone https://github.com/suvidsahay/GoRestAPI Factly
cd Factly
```

## Configuration

### Postgres database setup
Assuming that Postgres(v12 or later) is installed:
* Create role:
  ```bash
  sudo -u postgres createuser -P factly     # prompts for password
  ```
* Create database:
  ```bash
  sudo -u postgres createdb -O factly factly
  ```
  
### Configuration
  Create .env file based in sample.env and change the password provided while creating the factly user
  
 ## Starting the web server
 
 The web server can be started as shown below(assuming Go 1.13 is installed. By default it listens for
 HTTP connections on port 5000, so point your client at
 `localhost:5000`.
 
 ```bash
 go run main.go
```

## Consuming the REST API endpoints

* GET /users
```bash
curl -XGET localhost:5000/users
```

* POST /user
```bash
curl -XPOST -d '{"name":"Something"}' -H 'Content-Type: application/json' localhost:5000/user
```

* PUT /user/{id}
```bash
curl -XPUT -d '{"name":"New Name"}' -H 'Content-Type: application/json' localhost:5000/user/1
```

* DELETE /user/{id}
```bash
curl -XDELETE localhost:5000/user/1
```

## Unit Testing using Go testing package 

Run the following to command to test your application
```bash
go test -v
```
