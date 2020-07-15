# sms-bizApp

GoLang - Business Application (bizApp) Backend 

For documentation and a live demo see ###


## Setup
### Backend WebServer (GoLang/Aero)
To run the backend webserver simply execute:
```shell
go run main.go
```
OR <br><br>
Use Golang to build/execute the backend by adding a debug configuration from the "Go Build" item in the debug configuration editor.


- The configuration file for the Go webserver is located in ./config.json.

### Database (MongoDB)
- The configuration file for the database is located in /config/config.json.
- The application uses MongDB for the database.

#### Docker Configuration
```bash
version: '3.1'
services:
  mongo:
    image: mongo
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: mongoadmin
      MONGO_INITDB_ROOT_PASSWORD: password
    ports:
      - 27017:27017
    volumes:
      - ./data:/data/db
 ```

### Redis (Session Cache)

#### Docker Configuration
```bash
redis:
  container_name: redis
  image: redis
  ports:
    - "6379:6379"
  volumes:
    - ./data:/data
  restart: always
```
- The data volume for Redis is ./data which resides in the same folder as the docker-compose.yml file. 
- Please edit as necessary.