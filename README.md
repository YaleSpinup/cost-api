# cost-api

This API provides simple restful API access to Amazon's Cost explorer and cloudwatch metrics service.

## Endpoints

```
GET /v1/cost/ping
GET /v1/cost/version
GET /v1/cost/metrics

GET /v1/cost/{account}/spaces/{spaceid}[?start=2019-10-01&end=2019-10-30]
GET /v1/cost/{account}/spaces/{spaceid}/{resourcename}[?start=2019-10-01&end=2019-10-30]

GET /v1/cost/{account}/instances/{id}/metrics/graph.png?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]
GET /v1/cost/{account}/instances/{id}/metrics/graph?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]

GET /v1/metrics/{account}/instances/{id}/graph?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]
GET /v1/metrics/{account}/clusters/{cluster}/services/{service}/graph?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]
GET /v1/metrics/{account}/buckets/{bucket}/graph?metric={BucketSizeBytes|NumberOfObjects}
GET /v1/metrics/{account}/rds/{type}/{id}/graph?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]
```

## How it works

Costs are Filtered - the keys/values are resource tags

```json
{
  Filter: {
    And: [
      {
        Tags: {
          Key: "Name",
          Values: ["spinup-000cba.spinup.yale.edu"]
        }
      },
      {
        Tags: {
          Key: "spinup:spaceid",
          Values: ["spinup-0002a2"]
        }
      },
      {
        Or: [{
            Tags: {
              Key: "yale:org",
              Values: ["ss"]
            }
          },{
            Tags: {
              Key: "spinup:org",
              Values: ["ss"]
            }
          }]
      }
  }
}
```

## Usage

### Get the cost and usage for a space ID, using tags

By default, this will get the month to date costs for a space id (based on the `spinup:spaceid` tag).

#### Request

```
GET /v1/cost/{account}/spaces/{spaceid}
```

#### Response

```json
[
    {
        "Estimated": true,
        "Groups": [],
        "TimePeriod": {
            "End": "2020-01-15",
            "Start": "2020-01-01"
        },
        "Total": {
            "BlendedCost": {
                "Amount": "0",
                "Unit": "USD"
            },
            "UnblendedCost": {
                "Amount": "0",
                "Unit": "USD"
            },
            "UsageQuantity": {
                "Amount": "0",
                "Unit": "N/A"
            }
        }
    }
]
```

### Get the cost and usage for a resourcename within a space ID, using tags

By default, this will get the month to date costs for a resource name with a space id

tags

- spinup:spaceid
- Name

### Get cloudwatch metrics widgets URL from S3 for an instance ID

This will get the passed metric(s) for the passed instance ID or container cluster/service in a `image/png` graph for the past 1 day by default, cache it in S3
and return the URL. URLs are cached in the API for 5 minutes, the images should be purged from the S3 cache on a schedule. It's also
possible to pass the height, width, start time, end time and period (e. `300s` for 300 seconds, `5m` for 5 minutes).  Query parameters must follow
the [CloudWatch Metric Widget Structure](https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Metric-Widget-Structure.html).

### Documentation on cloudwatch metrics

```
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/viewing_metrics_with_cloudwatch.html

Get you a list of metrics per AWS service
$ aws --region us-east-1 cloudwatch list-metrics --namespace AWS/RDS |grep MetricName |sort| uniq

GetMetricWidget gets a metric widget image for an instance id
https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Metric-Widget-Structure.html
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/viewing_metrics_with_cloudwatch.html
https://docs.aws.amazon.com/AmazonECS/latest/developerguide/cloudwatch-metrics.html
https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/MonitoringOverview.html

Example metrics request
{
  "metrics": [
    [ "AWS/ECS", "CPUUtilization", "ClusterName", "spinup-000393", "ServiceName", "spinup-0010a3-testsvc" ]
  ],
  "stat": "Average"
  "period": 300,
  "start": "-P1D",
  "end": "PT0H"
}
```

#### Request

```
GET /v1/metrics/{account}/instances/{id}/graph?metric={metric1}[&metric={metric2}&....]
GET /v1/metrics/{account}/instances/{id}/graph?metric={metric1}[&metric={metric2}&start={start}&end={end}&period={period}]

GET /v1/metrics/{account}/clusters/{cluster}/services/{service}/graph?metric={metric1}[&metric={metric2}&....]
GET /v1/metrics/{account}/clusters/{cluster}/services/{service}/graph?metric={metric1}[&metric={metric2}&start={start}&end={end}&period={period}]
```

#### Response

```json
{
    "ImageURL": "https://s3.amazonaws.com/sometestbucket/aabbccddeeff-Y3_yCKckBrkUNt3Lh4LzXBFeLXBY5IP1oUED4hyY0cdKneYelKv-xlV7K2F_d0ccwp677A=="
}
```

with an image like this

![WidgetExample](/img/example_response.png?raw=true)

## Image Caching

When image urls are returned for metrics graph data, they are cached in the image cache.  The default implementation of this cache is an S3 bucket where the URLs are returned in the response (and cached in the data cache).

## Data Caching

AWS Cost Explorer data and metrics graph image url is cached (using go-cache).  The cache TTLs are configurable via config.json: CacheExpireTime and CachePurgeTime.  The cache can also be purged via daemon restart.

## Authentication

Authentication is accomplished via a pre-shared key (hashed string).  This is done via the `X-Auth-Token` header.

## Author

E Camden Fisher <camden.fisher@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright (c) 2019 Yale University
