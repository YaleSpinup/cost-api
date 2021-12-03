# cost-api

This API provides simple restful API access to Amazon's Cost explorer and cloudwatch metrics service.

## Endpoints

```text
GET /v1/cost/ping
GET /v1/cost/version
GET /v1/cost/metrics

GET /v1/cost/{account}/spaces/{spaceid}[?start=2019-10-01&end=2019-10-30][&groupBy=SERVICE]

POST /v1/cost/{account}/spaces/{spaceid}/budgets
GET /v1/cost/{account}/spaces/{spaceid}/budgets
DELETE /v1/cost/{account}/spaces/{spaceid}/budgets/{budget}

GET /v1/cost/{account}/spaces/{space}/instances/{id}/optimizer

GET /v1/inventory/{account}/spaces/{spaceid}

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

GET /v1/cost/{account}/spaces/{spaceid}

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

GET /v1/cost/{account}/spaces/{spaceid}?start=2021-04-01&end=2021-05-31

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

GET /v1/cost/{account}/spaces/{spaceid}?start=2021-04-01&end=2021-05-31&groupby=INSTANCE_TYPE_FAMILY

Supported default 'groupby' values are AZ, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, PURCHASE_TYPE, SERVICE, USAGE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, LEGAL_ENTITY_NAME, DEPLOYMENT_OPTION, DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, REGION, BILLING_ENTITY, RESERVATION_ID, SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM. In addition, the custom RESOURCE_NAME 'groupby' is supported using the Name tag.

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

## Compute Optimizer recommendations

### Get recommendations for an instance id

GET /v1/cost/{account}/spaces/{space}/instances/{id}/optimizer

#### Empty response

Empty response is usually a result of not enough data for a recommendation.

```json
[]
```

#### Response for optimized instance

```json
[
    {
        "AccountId": "1234567890",
        "CurrentInstanceType": "t3a.medium",
        "Finding": "OPTIMIZED",
        "FindingReasonCodes": [],
        "InstanceArn": "arn:aws:ec2:us-east-1:1234567890:instance/i-1122334455",
        "InstanceName": "best.bobs.edu",
        "LastRefreshTimestamp": "2021-06-16T18:21:42.595Z",
        "LookBackPeriodInDays": 14,
        "RecommendationOptions": [
            {
                "InstanceType": "t3a.medium",
                "PerformanceRisk": 1,
                "PlatformDifferences": [],
                "ProjectedUtilizationMetrics": [
                    {
                        "Name": "CPU",
                        "Statistic": "MAXIMUM",
                        "Value": 82.72781814545695
                    }
                ],
                "Rank": 1
            }
        ],
        "RecommendationSources": [
            {
                "RecommendationSourceArn": "arn:aws:ec2:us-east-1:1234567890:instance/i-1122334455",
                "RecommendationSourceType": "Ec2Instance"
            }
        ],
        "UtilizationMetrics": [
            {
                "Name": "CPU",
                "Statistic": "MAXIMUM",
                "Value": 82.72781814545695
            },
            {
                "Name": "EBS_READ_OPS_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 1127.2433333333333
            },
            {
                "Name": "EBS_WRITE_OPS_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 74.5
            },
            {
                "Name": "EBS_READ_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 29981835.9375
            },
            {
                "Name": "EBS_WRITE_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 1845286.4583333333
            },
            {
                "Name": "NETWORK_IN_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 10330.019722222223
            },
            {
                "Name": "NETWORK_OUT_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 507890.63972222223
            },
            {
                "Name": "NETWORK_PACKETS_IN_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 39.420111111111105
            },
            {
                "Name": "NETWORK_PACKETS_OUT_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 328.95061111111113
            }
        ]
    }
]
```

#### Response with recommendations

```json
[
    {
        "AccountId": "1234567890",
        "CurrentInstanceType": "m4.xlarge",
        "Finding": "UNDER_PROVISIONED",
        "FindingReasonCodes": [
            "CPUOverprovisioned",
            "EBSIOPSOverprovisioned",
            "EBSThroughputUnderprovisioned"
        ],
        "InstanceArn": "arn:aws:ec2:us-east-1:1234567890:instance/i-0987654321",
        "InstanceName": "knowledge.bobs.edu",
        "LastRefreshTimestamp": "2021-06-16T18:53:25.669Z",
        "LookBackPeriodInDays": 14,
        "RecommendationOptions": [
            {
                "InstanceType": "t3.xlarge",
                "PerformanceRisk": 3,
                "PlatformDifferences": [
                    "NetworkInterface",
                    "Hypervisor",
                    "StorageInterface"
                ],
                "ProjectedUtilizationMetrics": [
                    {
                        "Name": "CPU",
                        "Statistic": "MAXIMUM",
                        "Value": 44.64812085482684
                    }
                ],
                "Rank": 1
            },
            {
                "InstanceType": "m5.xlarge",
                "PerformanceRisk": 1,
                "PlatformDifferences": [
                    "NetworkInterface",
                    "Hypervisor",
                    "StorageInterface"
                ],
                "ProjectedUtilizationMetrics": [
                    {
                        "Name": "CPU",
                        "Statistic": "MAXIMUM",
                        "Value": 44.64812085482684
                    }
                ],
                "Rank": 2
            },
            {
                "InstanceType": "m4.xlarge",
                "PerformanceRisk": 1,
                "PlatformDifferences": [],
                "ProjectedUtilizationMetrics": [
                    {
                        "Name": "CPU",
                        "Statistic": "MAXIMUM",
                        "Value": 55.5084745762712
                    }
                ],
                "Rank": 3
            }
        ],
        "RecommendationSources": [
            {
                "RecommendationSourceArn": "arn:aws:ec2:us-east-1:1234567890:instance/i-0987654321",
                "RecommendationSourceType": "Ec2Instance"
            }
        ],
        "UtilizationMetrics": [
            {
                "Name": "CPU",
                "Statistic": "MAXIMUM",
                "Value": 55.5084745762712
            },
            {
                "Name": "EBS_READ_OPS_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 753.56
            },
            {
                "Name": "EBS_WRITE_OPS_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 917.98
            },
            {
                "Name": "EBS_READ_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 89428489.58333333
            },
            {
                "Name": "EBS_WRITE_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 83148230.79666667
            },
            {
                "Name": "NETWORK_IN_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 414224.09777777776
            },
            {
                "Name": "NETWORK_OUT_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 29260.118055555555
            },
            {
                "Name": "NETWORK_PACKETS_IN_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 177.9728888888889
            },
            {
                "Name": "NETWORK_PACKETS_OUT_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 35.34688888888889
            }
        ]
    }
]
```

## Inventory Usage

The inventory endpoint returns resources belonging to a space by tag.  It uses the resourcegroupstaggingapi and also parses the ARN to
determine the resource type information.  Since the `resource` prefix is inconsistent, it shouldn't be relied upon for categorization, but the
`service` should be accurate.

### Request

GET /v1/inventory/{account}/spaces/{spaceid}

### Response

```json
[
    {
        "name": "spintst-000a16-TestFS",
        "arn": "arn:aws:elasticfilesystem:us-east-1:1234567890:file-system/fs-aaaabbbb11",
        "partition": "aws",
        "service": "elasticfilesystem",
        "region": "us-east-1",
        "account_id": "1234567890",
        "resource": "file-system/fs-aaaabbbb11"
    },
    {
        "name": "",
        "arn": "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/testTargetGroup-HTTP80/abcdefg12",
        "partition": "aws",
        "service": "elasticloadbalancing",
        "region": "us-east-1",
        "account_id": "1234567890",
        "resource": "targetgroup/testTargetGroup-HTTP80/abcdefg12"
    },
    {
        "name": "spintst-000028",
        "arn": "arn:aws:logs:us-east-1:1234567890:log-group:spintst-000028",
        "partition": "aws",
        "service": "logs",
        "region": "us-east-1",
        "account_id": "1234567890",
        "resource": "log-group:spintst-000028"
    },
    {
        "name": "spintst-000b67-webServiceTest",
        "arn": "arn:aws:secretsmanager:us-east-1:1234567890:secret:spinup/sstst/spintst-000028/spintst-000b67-webServiceTest-api-cred-ibdIk7",
        "partition": "aws",
        "service": "secretsmanager",
        "region": "us-east-1",
        "account_id": "1234567890",
        "resource": "secret:spinup/sstst/spintst-000028/spintst-000b67-webServiceTest-api-cred-ibdIk7"
    }
]
```


## Metrics Usage

### Get Cloudwatch metrics widgets URL from S3 for an instance ID

This will get the passed metric(s) for the passed instance ID or container cluster/service in a `image/png` graph for the past 1 day by default, cache it in S3
and return the URL. URLs are cached in the API for 5 minutes, the images should be purged from the S3 cache on a schedule. It's also
possible to pass the height, width, start time, end time and period (e. `300s` for 300 seconds, `5m` for 5 minutes).  Query parameters must follow
the [CloudWatch Metric Widget Structure](https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Metric-Widget-Structure.html).

### Documentation on Cloudwatch metrics

#### Get a list of metrics per AWS service

```bash
aws --region us-east-1 cloudwatch list-metrics --namespace AWS/RDS |grep MetricName |sort| uniq
```

#### Helpful references

* [Viewing Metrics with CloudWatch](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/viewing_metrics_with_cloudwatch.html)
* [CloudWatch Metrics Widget Structure](https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Metric-Widget-Structure.html)
* [CloudWatch Metrics Developer Guide](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/cloudwatch-metrics.html)
* [AWS Aurora Monitoring Overview](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/MonitoringOverview.html)

#### Example metrics request

```json
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

```text
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
