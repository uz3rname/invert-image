# Image color inversion service

Service for inverting images and retrieving previously inverted images from database.

It is written in Golang `v. 1.16` using `Fiber` framework and `GORM` with `PostgreSQL` as DBMS.
Image manipulations are done via https://github.com/disintegration/imaging library.

### Deployment

* Create `.env` file from `.env.example`. For deployment via `docker-compose` you only need to specify `PORT`, other
variables are used for local development and can be safely ignored.
* Run `docker-compose up -d`

API and web interface will be available at specified port. Swagger is located at `/api/docs`

### Tests

To run API e2e tests, execute the following:

```
go test -v
```
