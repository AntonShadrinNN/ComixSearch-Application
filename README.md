# API Documentation

## Comics search application
This API allows users to search for comics on [xkcd](https://xkcd.com/) based on given keywords. The API provides two endpoints: one for searching comics and another for accessing the API documentation.

### Direct Usage
1. Install go version 1.22 or higher. See [official docs](https://go.dev/doc/install).
2. Install dependencies. Run `go mod tidy`.
3. Run app. See [allowed usage](#getting-started).
4. Use [endpoints](#endpoints).

## Endpoints

### 1. Search Comics
**POST /api/v1/search**

This endpoint searches for comics on xkcd based on the provided keywords.

#### Request
- **URL:** `/api/v1/search`
- **Method:** `POST`
- **Content-Type:** `application/json`

#### Body Parameters
- **keywords** (string, required): The keywords to search for in the comics.
  - Example: `{ "keywords": "earth" }`

#### Query Parameters
- **limit** (integer, optional): The maximum number of comics to return.
  - Example: `/api/v1/search?limit=5`

#### Response
- **Status Code:** `200 OK` if the request is successful.
- **Content-Type:** `application/json`
- **Body:**
  ```json
  {
        "comices": {
          {
            "earth":"http://xkcd/earth"
          },
        },
        "error": "error",
  }
  ```

#### Example Request
```bash
curl -X POST http://localhost:8080/api/v1/search \
     -H "Content-Type: application/json" \
     -d '{"keywords": "earth"}'
```

### 2. API Documentation
**GET /api/v1/docs**

This endpoint provides access to the API documentation via Swagger.

#### Request
- **URL:** `/api/v1/docs`
- **Method:** `GET`

#### Response
- **Status Code:** `200 OK`
- **Content-Type:** `text/html`

This will serve the Swagger UI documentation for the API, where users can interact with the API and see details about each endpoint.

#### Example Request
```bash
curl http://localhost:8080/api/v1/docs
```

## Error Handling
The API uses standard HTTP status codes to indicate the success or failure of an API request. The following status codes may be returned:

- **200 OK:** The request was successful.
- **400 Bad Request:** The request could not be understood or was missing required parameters.
- **500 Internal Server Error:** An error occurred on the server.