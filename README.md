# cost-api

This API provides simple restful API access to Amazon's Cost explorer and cloudwatch metrics service.

## Endpoints

```
GET /v1/cost/ping
GET /v1/cost/version
GET /v1/cost/metrics

GET /v1/cost/{account}/spaces/{spaceid}[?start=2019-10-01&end=2019-10-30][&groupBy=SERVICE]
GET /v1/cost/{account}/spaces/{spaceid}/{resourcename}[?start=2019-10-01&end=2019-10-30]

POST /v1/cost/{account}/spaces/{spaceid}/budgets
GET /v1/cost/{account}/spaces/{spaceid}/budgets
DELETE /v1/cost/{account}/spaces/{spaceid}/budgets/{budget}

### DEPRECATED ###
GET /v1/cost/{account}/instances/{id}/metrics/graph.png?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]
GET /v1/cost/{account}/instances/{id}/metrics/graph?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]
##################

GET /v1/metrics/{account}/instances/{id}/graph?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]
GET /v1/metrics/{account}/clusters/{cluster}/services/{service}/graph?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]
GET /v1/metrics/{account}/buckets/{bucket}/graph?metric={BucketSizeBytes|NumberOfObjects}
GET /v1/metrics/{account}/rds/{type}/{id}/graph?metric={metric1}[&metric={metric2}&start=-P1D&end=PT0H&period=300]
```

## Cost Usage

### Get the cost and usage for a space ID

By default, this will get the month to date costs for a space id (based on the `spinup:spaceid` tag).  Date ranges and grouping by
different dimensions is supported by passing query parameters.

#### Request month to date costs for a space

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

#### Request costs for a space for a date range

```
GET /v1/cost/{account}/spaces/{spaceid}?start=2021-04-01&end=2021-05-31
```

#### Response

```json
[
    {
        "Estimated": false,
        "Groups": [],
        "TimePeriod": {
            "End": "2021-05-01",
            "Start": "2021-04-01"
        },
        "Total": {
            "BlendedCost": {
                "Amount": "8.1395432009",
                "Unit": "USD"
            },
            "UnblendedCost": {
                "Amount": "8.1395437889",
                "Unit": "USD"
            },
            "UsageQuantity": {
                "Amount": "37095.8855728516",
                "Unit": "N/A"
            }
        }
    },
    {
        "Estimated": true,
        "Groups": [],
        "TimePeriod": {
            "End": "2021-05-31",
            "Start": "2021-05-01"
        },
        "Total": {
            "BlendedCost": {
                "Amount": "3.2187441928",
                "Unit": "USD"
            },
            "UnblendedCost": {
                "Amount": "3.2187483458",
                "Unit": "USD"
            },
            "UsageQuantity": {
                "Amount": "8466.9944576915",
                "Unit": "N/A"
            }
        }
    }
]
```

#### Request costs for a space by date range and grouped by a dimension

```
GET /v1/cost/{account}/spaces/{spaceid}?start=2021-04-01&end=2021-05-31&groupby=INSTANCE_TYPE_FAMILY
```

Valid 'groupby' values are AZ, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, PURCHASE_TYPE, SERVICE, USAGE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, LEGAL_ENTITY_NAME, DEPLOYMENT_OPTION, DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, REGION, BILLING_ENTITY, RESERVATION_ID, SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM.

#### Response

```json
[
    {
        "Estimated": false,
        "Groups": [
            {
                "Keys": [
                    "NoInstanceTypeFamily"
                ],
                "Metrics": {
                    "BlendedCost": {
                        "Amount": "2.3325985365",
                        "Unit": "USD"
                    },
                    "UnblendedCost": {
                        "Amount": "2.3325991245",
                        "Unit": "USD"
                    },
                    "UsageQuantity": {
                        "Amount": "36869.1647378516",
                        "Unit": "N/A"
                    }
                }
            },
            {
                "Keys": [
                    "m5a"
                ],
                "Metrics": {
                    "BlendedCost": {
                        "Amount": "4.008678362",
                        "Unit": "USD"
                    },
                    "UnblendedCost": {
                        "Amount": "4.008678362",
                        "Unit": "USD"
                    },
                    "UsageQuantity": {
                        "Amount": "39.124167",
                        "Unit": "N/A"
                    }
                }
            },
            {
                "Keys": [
                    "t3"
                ],
                "Metrics": {
                    "BlendedCost": {
                        "Amount": "0.0238622176",
                        "Unit": "USD"
                    },
                    "UnblendedCost": {
                        "Amount": "0.0238622176",
                        "Unit": "USD"
                    },
                    "UsageQuantity": {
                        "Amount": "0.573611",
                        "Unit": "N/A"
                    }
                }
            },
            {
                "Keys": [
                    "t3a"
                ],
                "Metrics": {
                    "BlendedCost": {
                        "Amount": "1.7744040848",
                        "Unit": "USD"
                    },
                    "UnblendedCost": {
                        "Amount": "1.7744040848",
                        "Unit": "USD"
                    },
                    "UsageQuantity": {
                        "Amount": "187.023057",
                        "Unit": "N/A"
                    }
                }
            }
        ],
        "TimePeriod": {
            "End": "2021-05-01",
            "Start": "2021-04-01"
        },
        "Total": {}
    },
    {
        "Estimated": true,
        "Groups": [
            {
                "Keys": [
                    "NoInstanceTypeFamily"
                ],
                "Metrics": {
                    "BlendedCost": {
                        "Amount": "3.1292908776",
                        "Unit": "USD"
                    },
                    "UnblendedCost": {
                        "Amount": "3.1292950306",
                        "Unit": "USD"
                    },
                    "UsageQuantity": {
                        "Amount": "8466.2786246915",
                        "Unit": "N/A"
                    }
                }
            },
            {
                "Keys": [
                    "c5"
                ],
                "Metrics": {
                    "BlendedCost": {
                        "Amount": "0.06346661",
                        "Unit": "USD"
                    },
                    "UnblendedCost": {
                        "Amount": "0.06346661",
                        "Unit": "USD"
                    },
                    "UsageQuantity": {
                        "Amount": "0.373333",
                        "Unit": "N/A"
                    }
                }
            },
            {
                "Keys": [
                    "m5a"
                ],
                "Metrics": {
                    "BlendedCost": {
                        "Amount": "0.025561092",
                        "Unit": "USD"
                    },
                    "UnblendedCost": {
                        "Amount": "0.025561092",
                        "Unit": "USD"
                    },
                    "UsageQuantity": {
                        "Amount": "0.297222",
                        "Unit": "N/A"
                    }
                }
            },
            {
                "Keys": [
                    "t3a"
                ],
                "Metrics": {
                    "BlendedCost": {
                        "Amount": "0.0004256132",
                        "Unit": "USD"
                    },
                    "UnblendedCost": {
                        "Amount": "0.0004256132",
                        "Unit": "USD"
                    },
                    "UsageQuantity": {
                        "Amount": "0.045278",
                        "Unit": "N/A"
                    }
                }
            }
        ],
        "TimePeriod": {
            "End": "2021-05-31",
            "Start": "2021-05-01"
        },
        "Total": {}
    }
]
```

### Get the cost and usage for a resource (name) within a space ID

By default, this will get the month to date costs for a resource name with a space id

```
GET /v1/cost/{account}/spaces/{spaceid}/{resourcename}
```

### How it works

Costs are Filtered - the keys/values are resource tags

```json
{
  "Filter": {
    "And": [
      {
        "Tags": {
          "Key": "Name",
          "Values": ["spinup-000cba.spinup.yale.edu"]
        }
      },
      {
        "Tags": {
          "Key": "spinup:spaceid",
          "Values": ["spinup-0002a2"]
        }
      },
      {
        "Or": [
          {
            "Tags": {
              "Key": "yale:org",
              "Values": ["ss"]
            }
          },
          {
            "Tags": {
              "Key": "spinup:org",
              "Values": ["ss"]
            }
          }]
      }
    ]
  }
}
```

## Budget Usage

### Create Budgets Alerts

#### Request

POST /v1/cost/{account}/spaces/{spaceid}/budgets

```json
{
    "Amount": "10",
    "TimeUnit": "MONTHLY",
    "Alerts": [
        {
            "ComparisonOperator": "GREATER_THAN",
            "NotificationType": "FORECASTED",
            "Threshold": 100,
            "ThresholdType": "PERCENTAGE",
            "Addresses": ["some.user@yale.edu", "some.other@yale.edu"]
        }
    ]
}
```

#### Response

```json
{
    "Amount": "10",
    "Name": "spintst-000028-MONTHLY-01",
    "TimeUnit": "MONTHLY",
    "Alerts": [
        {
            "ComparisonOperator": "GREATER_THAN",
            "NotificationState": "",
            "NotificationType": "FORECASTED",
            "Threshold": 100,
            "ThresholdType": "PERCENTAGE",
            "Addresses": ["some.user@yale.edu", "some.other@yale.edu"]
        }
    ]
}
```

### List Budgets Alerts

GET /v1/cost/{account}/spaces/{spaceid}/budgets

### Response

```json
[
    "spintst-000028-MONTHLY-01"
]
```

### GET details about a  Budgets Alert

GET /v1/cost/{account}/spaces/{spaceid}/budgets/{budget}

### Response

```json
{
    "Amount": "10",
    "Name": "spintst-000028-MONTHLY-01",
    "TimeUnit": "MONTHLY",
    "Alerts": [
        {
            "ComparisonOperator": "GREATER_THAN",
            "NotificationState": "",
            "NotificationType": "FORECASTED",
            "Threshold": 100,
            "ThresholdType": "PERCENTAGE",
            "Addresses": ["some.user@yale.edu", "some.other@yale.edu"]
        }
    ]
}
```

### Delete Budgets Alert

DELETE /v1/cost/{account}/spaces/{spaceid}/budgets/{budget}

### Response

```json
"OK"
```

## Metrics Usage

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
