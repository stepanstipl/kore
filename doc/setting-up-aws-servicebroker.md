# Setting up the AWS servicebroker

Kore has built-in support for the [aws-servicebroker](https://github.com/awslabs/aws-servicebroker) and it's able to install it as a service provider.

## Prerequisites

1. Make sure you have the `services` feature gate enabled in the API/UI:

```
KORE_FEATURE_GATES="services=true"
```

Or in Helm:

```
api:
  feature_gates: [services=true]
ui:
  feature_gates: [services=true]
```

1. Make sure you have the `services` feature gate enabled in the CLI:

```
$ kore alpha feature-gates enable services
```

## Set up AWS access

1. Create an IAM user with the following policy:

    ```
    {
        "Version": "2012-10-17",
        "Statement": [
          {
            "Sid": "SsmForSecretBindings",
            "Action": "ssm:PutParameter",
            "Resource": "arn:aws:ssm:<REGION>:<ACCOUNT_ID>:parameter/asb-*",
            "Effect": "Allow"
          },
          {
            "Sid": "AllowCfnToGetTemplates",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::[S3 bucket name]/templates/*",
            "Effect": "Allow"
          },
          {
             "Sid": "CloudFormation",
             "Action": [
                "cloudformation:CreateStack",
                "cloudformation:DeleteStack",
                "cloudformation:DescribeStacks",
                "cloudformation:DescribeStackEvents",
                "cloudformation:UpdateStack",
                "cloudformation:CancelUpdateStack"
             ],
             "Resource": [
                "arn:aws:cloudformation:<REGION>:<ACCOUNT_ID>:stack/aws-service-broker-*/*"
             ],
             "Effect": "Allow"
          },
         {
            "Sid": "ServiceClassPermissions",
            "Action": [
               "athena:*",
               "dynamodb:*",
               "kms:*",
               "elasticache:*",
               "elasticmapreduce:*",
               "kinesis:*",
               "rds:*",
               "redshift:*",
               "route53:*",
               "s3:*",
               "sns:*",
               "sqs:*",
               "ec2:*",
               "iam:*",
               "lambda:*"
            ],
            "Resource": [
               "*"
            ],
            "Effect": "Allow"
         }
       ]
    }
    ```

1. Create a set of access keys

1. Create a secret in Kore with the access key id and secret access key

```
$ kore create secret aws-broker \
    -t kore-admin \
    --type aws-credentials \
    --description "AWS service broker credentials" \
    --from-literal "access_key_id=${AWS_ACCESS_KEY_ID}" \
    --from-literal "access_secret_key=$AWS_SECRET_ACCESS_KEY"
```

## Install the service provider in Kore

    ```
    $ cat <<EOF | kore apply -f -
    ---
    apiVersion: serviceproviders.kore.appvia.io/v1
    kind: ServiceProvider
    metadata:
      name: aws-broker
      namespace: kore
    spec:
      description: AWS Service Broker
      type: aws-servicebroker
      summary: AWS Service Broker
      configuration:
        tableName: [DynamoDB Table name]
        region: eu-west-2
        s3BucketName: [S3 bucket name]
        s3BucketRegion: eu-west-2
        allowEmptyCredentialSchema: true
        defaultPlans:
          - auroramysql-custom
          - aurorapostgresql-custom
          - cognito-custom
          - documentdb-custom
          - elasticache-custom
          - elasticsearch-custom
          - emr-custom
          - mq-custom
          - rdsmariadb-custom
          - rdsmssql-custom
          - rdsmysql-custom
          - rdsoracle-custom
          - rdspostgresql-custom
          - redshift-custom
          - s3-custom
      configurationFrom:
        - name: aws_access_key_id
          secretKeyRef:
            name: aws-broker
            namespace: kore-admin
            key: access_key_id
        - name: aws_secret_access_key
          secretKeyRef:
            name: aws-broker
            namespace: kore-admin
            key: access_secret_key
    EOF
    ```

It will take a couple of minutes until the service provider will be ready to be used.

## How to enable/disable service kinds

```
$ kore alpha patch servicekind s3 spec.enabled [true|false]
```

