#!/bin/sh

set -e

ROLE_NAME=$1

if [ "$ROLE_NAME" = "" ]; then
  ROLE_NAME="Admin"
fi

ROLE_ARN=$(aws iam get-role --role-name $ROLE_NAME | jq -r ".Role.Arn")
CREDENTIALS="$(aws sts assume-role --role-arn $ROLE_ARN --role-session-name otel --duration-seconds=3600)"
ACCESS_KEY_ID="$(echo $CREDENTIALS | jq -r ".Credentials.AccessKeyId")"
SECRET_ACCESS_KEY="$(echo $CREDENTIALS | jq -r ".Credentials.SecretAccessKey")"
SESSION_TOKEN="$(echo $CREDENTIALS | jq -r ".Credentials.SessionToken")"

#docker run --rm \
#  -p 4317:4317 \
#  -p 55680:55680 \
#  -p 8889:8888 \
#  -e AWS_REGION=ap-southeast-2 \
#  -e "AWS_ACCESS_KEY_ID=$ACCESS_KEY_ID" \
#  -e "AWS_SECRET_ACCESS_KEY=$SECRET_ACCESS_KEY" \
#  -e "AWS_SESSION_TOKEN=$SESSION_TOKEN" \
#  -v "${PWD}/aws-otel-collector-config.yaml":/otel-local-config.yaml \
#  --name awscollector \
#  public.ecr.aws/aws-observability/aws-otel-collector --config otel-local-config.yaml
#

docker run --rm \
  -p 13133:13133 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 55678-55680:55678-55680 \
  -p 8888:8888 \
  -p 9411:9411 \
  -p 16686:16686 \
  -e AWS_REGION=ap-southeast-2 \
  -e "AWS_ACCESS_KEY_ID=$ACCESS_KEY_ID" \
  -e "AWS_SECRET_ACCESS_KEY=$SECRET_ACCESS_KEY" \
  -e "AWS_SESSION_TOKEN=$SESSION_TOKEN" \
  --name otelcollector \
  otel/opentelemetry-collector # --config otel-local-config.yaml
#  -v "${PWD}/aws-otel-collector-config.yaml":/otel-local-config.yaml \

