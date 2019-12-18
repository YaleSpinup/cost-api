# cost-api

This API provides simple restful API access to Amazon's Cost explorer service.

## Endpoints

```
GET /v1/cost/ping
GET /v1/cost/version
GET /v1/cost/metrics

GET /v1/cost/{account}/spaces/{spaceid}[?start=2019-10-01&end=2019-10-30]

GET /v1/cost/{account}/instances/{id}/metrics/{metric}.png[?start=-P1D&end=PT0H&period=300]
GET /v1/cost/{account}/instances/{id}/metrics/{metric}[?start=-P1D&end=PT0H&period=300]
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

### Get cloudwatch metrics widgets for an instance ID

This will get the passed metric for the passed instance ID in a `image/png` graph for the past 1 day by default. It's also
possible to pass the start time, end time and period (in seconds).  Query parameters must follow
the [CloudWatch Metric Widget Structure](https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Metric-Widget-Structure.html).

#### Request

GET /v1/cost/{account}/instances/{id}/metrics/{metric}.png
GET /v1/cost/{account}/instances/{id}/metrics/{metric}.png?start={StartTime}&end={EndTime}&period={Period}

#### Response

![WidgetExample](/img/example_response.png?raw=true)

### Get cloudwatch metrics widgets URL from S3 for an instance ID

This will get the passed metric for the passed instance ID in a `image/png` graph for the past 1 day by default, cache it in S3
and return the URL. URLs are cached in the API for 5 minutes, the images should be purged from the S3 cache on a schedule. It's also
possible to pass the start time, end time and period (in seconds).  Query parameters must follow the [CloudWatch Metric Widget Structure](https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Metric-Widget-Structure.html).

#### Request

GET /v1/cost/{account}/instances/{id}/metrics/{metric}
GET /v1/cost/{account}/instances/{id}/metrics/{metric}?start={StartTime}&end={EndTime}&period={Period}

#### Response

```json
{
    "ImageURL": "https://s3.amazonaws.com/sometestbucket/abc123_kLbi1SNQlKqMOmpaaJHAQZ3a-acutp5-tc6J0="
}
```

## Caching
Caching data (using go-cache) from AWS Cost Explorer configurable via config.json: CacheExpireTime and CachePurgeTime.  The cache can also be purged via daemon restart. 

## Authentication

Authentication is accomplished via a pre-shared key.  This is done via the `X-Auth-Token` header.

## Author

E Camden Fisher <camden.fisher@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright (c) 2019 Yale University
