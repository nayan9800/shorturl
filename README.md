# ShortUrl 
Shorturl is simple MVC pattern web Application created in **Go** language. 

## Description

The shorturl is url shortner web application.
backend is wriiten in **Go** , and frontend is simple **HTML** with **BootStrap** for styling.  **PostgreSQL** is used  for database.  


## How to Run

To Run using docker use
```shell
docker-compose up
```
to run locally 
```shell
go run main.go
```

## Environment Variables
Application uses several envoirnment variables.
put this envoirment variables in **.env** file

### POSTGRES_USER
postgres database user.
 
### POSTGRES_PASSWORD
postgres password used by postgres user and application to access database

### POSTGRES_DB
deafault database name used by database and application