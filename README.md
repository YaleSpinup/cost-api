# cost-api

This API provides simple restful API access to Amazon's Cost explorer service.

## Endpoints

```
GET /v1/cost/ping
GET /v1/cost/version
GET /v1/cost/metrics

GET /v1/cost/{account}/spaces/{spaceid}
```

## Usage

### Get the cost and usage for a space ID

By default, this will get the month to date costs for a space id (based on the `spinup:spaceid` tag).

#### Request

GET /v1/cost/{account}/spaces/{spaceid}

#### Response

```json
{
    "TBD"
}
```

## Authentication

Authentication is accomplished via a pre-shared key.  This is done via the `X-Auth-Token` header.

## Author

E Camden Fisher <camden.fisher@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright (c) 2019 Yale University
