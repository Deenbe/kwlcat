#!/bin/sh

set -e

ROLE_NAME=$1

if [ "$ROLE_NAME" = "" ]; then
  ROLE_NAME="Admin"
fi

ROLE_ARN=$(aws iam get-role --role-name $ROLE_NAME | jq -r ".Role.Arn")
CREDENTIALS="$(aws sts assume-role --role-arn $ROLE_ARN --role-session-name otel --duration-seconds=3600)"
export ACCESS_KEY_ID="$(echo $CREDENTIALS | jq -r ".Credentials.AccessKeyId")"
export SECRET_ACCESS_KEY="$(echo $CREDENTIALS | jq -r ".Credentials.SecretAccessKey")"
export SESSION_TOKEN="$(echo $CREDENTIALS | jq -r ".Credentials.SessionToken")"

docker-compose up
