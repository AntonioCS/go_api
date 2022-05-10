# Go API Project

THIS IS A TEST API PROJECT

THIS IS NOT PRODUCTION LEVEL CODE


## Endpoints

### List all guests
```
GET /api/guestbook/list
```

### Add a guest
```
POST /api/guestbook/{name}
Body
{
    "table": <int>,
    "accompanying_guests": <int>
}	
```
### Update guest
```
PUT /api/guestbook/{name}
Body
{    
    "accompanying_guests": <int>
}	
```

### Delete guest	
```
DELETE /api/guestbook/{name}
```

### List empty seats
```
GET /api/guestbook/emptySeats
```


## Make targets

All the following targets should be executed in the root dir and not in the project dir (the root dir is basically where this file is located)

To get the project up and running execute the following on the target line (in the root dir):
```
    make start
```
Note: This will only start the Mysql container 

To shut down the project run:
```
    make stop
```

To run the api tests please run (after running `make start`):
```
  make go-test-api
```
Note: The test depends on the initial data (in sql/dump.sql)
If I had more time I would have created proper fixtures for the project. 

You can also run 
```
    make go-tests
```
To execute unit tests on `GuestBook.go`

If you want to test on Postman you can run:

```
    make go-run
```

This will start the server on the port specified in `project/.env` (should be port 3000)

Call via POST with the json in the BODY.

Example using Postman:

![Postman example](https://gcdnb.pbrd.co/images/n8O1uv0vtp1H.png)


To create a deploy image please run
```
    make go-deploy
```


To list all available options run:
```
    make
```

## TODOs

- Add error checking
- Create more isolated tests
- Possibly add a way to create fixtures
- Optimize code

