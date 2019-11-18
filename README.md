# cost-api

This API provides simple restful API access to Amazon's Cost explorer service.

## Endpoints

```
GET /v1/cost/ping
GET /v1/cost/version
GET /v1/cost/metrics

GET /v1/cost/{account}/spaces/{spaceid}
GET /v1/cost/{account}/spaces/{spaceid}[?StartTime=2019-10-01&EndTime=2019-10-30]
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

## Caching
Caching data (using go-cache) from AWS Cost Explorer configurable via config.json: CacheExpireTime and CachePurgeTIme.  The cache can also be purged via daemon restart. 

## Authentication

Authentication is accomplished via a pre-shared key.  This is done via the `X-Auth-Token` header.

## Author

E Camden Fisher <camden.fisher@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright (c) 2019 Yale University
