# Go CV REST API


This repository houses the backend service for a full-stack CV application, designed to serve as a personal website. Developed using Go 1.21, this **REST API** incorporates various technologies to deliver a robust and scalable solution.
The frontend app can be found at [cv-frontend-vuejs](https://github.com/aalug/cv-frontend-vuejs)

## App built with Go 1.21

### The app uses:
- Postgres
- Docker
- [Gin](https://github.com/gin-gonic/gin)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://github.com/kyleconroy/sqlc)
- [testify](https://github.com/stretchr/testify)
- [Viper](https://github.com/spf13/viper)
- [gin cors](https://github.com/gin-contrib/cors)
- [gin-swagger](https://github.com/swaggo/gin-swagger)
<hr>

## Getting started
1. Clone the repository
2. Go to the project's root directory
3. Rename `app.env.example` to `app.env` and replace the values
4. Install [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
5. Run in your terminal:
`docker-compose up` to run the containers
6. Now everything should be ready and server running on `SERVER_ADDRESS` specified in `app.env`
<hr>

## Testing
1. Run the containers (`docker-compose up`)
2. Run in your terminal:
    - `make test` to run all tests

   or
    - `make test_coverage p={PATH}` - to get the coverage in the HTML format - where `{PATH}` is the path to the target directory for which you want to generate test coverage. The `{PATH}` should be replaced with the actual path you want to use. For example `./internal/api`

   or
    - use standard `go test` commands (e.g. `go test -v ./internal/api`)
<hr>

## API Endpoint
All endpoints are available to test at http://localhost:8080/swagger/index.html after running the containers.


### GET `/api/v1/cv-profiles/{id}`

This endpoint is used to get the details of a CV profile with a provided ID.

#### Parameters

- `id` (integer, required): The ID of the CV profile. This parameter is included in the path of the request.

#### Responses

- `200 OK`: The request was successful and the response body contains the CV profile details.
- `400 Invalid ID`: The provided ID is invalid.
- `404 CV profile with given ID does not exist`: There is no CV profile with the provided ID. 
- `500 Any other server-side error`: There was a server-side error while processing the request.

#### Produces

The endpoint produces responses in the `application/json` format.

### GET `/api/v1/projects/skill/{id}/{skill}`

This endpoint is used to list projects for a CV profile with a provided ID and skill.

#### Parameters

- `id` (integer, required): The ID of the CV profile. This parameter is included in the path of the request.
- `skill` (string, required): The name of the skill. This parameter is included in the path of the request.
- `page` (integer, required): The page number. This parameter is included in the query of the request.
- `page_size` (integer, required): The page size. This parameter is included in the query of the request.

#### Responses

- `200 OK`: The request was successful and the response body contains a list of projects.
- `400 Invalid ID, skill name, page or page size`: The provided ID, skill name, page or page size is invalid. 
- `404 CV profile with given ID or skill with given name does not exist`: There is no CV profile with the provided ID or no skill with the provided name.
- `500 Any other server-side error`: There was a server-side error while processing the request.

#### Produces

The endpoint produces responses in the `application/json` format.

### GET `/api/v1/projects/{id}`

This endpoint is used to list projects for a CV profile with a provided ID.

#### Parameters

- `id` (integer, required): The ID of the CV profile. This parameter is included in the path of the request.
- `page` (integer, required): The page number. This parameter is included in the query of the request.
- `page_size` (integer, required): The page size. This parameter is included in the query of the request.

#### Responses

- `200 OK`: The request was successful and the response body contains a list of projects. 
- `400 Invalid ID, page or page size`: The provided ID, page or page size is invalid.
- `404 CV profile with given ID does not exist`: There is no CV profile with the provided ID.
- `500 Any other server-side error`: There was a server-side error while processing the request.

#### Produces

The endpoint produces responses in the `application/json` format.

### GET `/api/v1/skills/{id}`

This endpoint is used to list skills for a CV profile with a provided ID.

#### Parameters

- `id` (integer, required): The ID of the CV profile. This parameter is included in the path of the request.

#### Responses

- `200 OK`: The request was successful and the response body contains a list of skills. 
- `400 Invalid ID`: The provided ID is invalid.
- `404 CV profile with given ID does not exist`: There is no CV profile with the provided ID.
- `500 Any other server-side error`: There was a server-side error while processing the request.

#### Produces

The endpoint produces responses in the `application/json` format.
