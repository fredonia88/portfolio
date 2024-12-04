#!/bin/bash

# print error message and exit
exit_with_error() {
    echo "$1"
    exit 1
}

# check if 'build' argument is provided
if ! [[ "$@" =~ "--build=" ]]; then
    exit_with_error "Error: '--build=create|update' argument is required."
fi

# parse the arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --build=*)
            build="${1#*=}"
            shift
            ;;
        *)
            exit_with_error "Error: Invalid argument: $1. Valid option is: '--build=create|update'"
            ;;
    esac
done

# validate the value of build
if [[ "$build" != "create" && "$build" != "update" ]]; then
    exit_with_error "Error: Invalid 'build' value. Allowed values: 'create|update'"
fi

DOMAIN_NAME="loriendream.com"
COMMENT="Hosted zone for $DOMAIN_NAME"
CALLER_REFERENCE=$(date +%s)  # Unique identifier for the request

# create the hosted zone
echo "Creating hosted zone for $DOMAIN_NAME..."
CREATE_RESPONSE=$(aws route53 create-hosted-zone \
    --name "$DOMAIN_NAME" \
    --caller-reference "$CALLER_REFERENCE" \
    --hosted-zone-config Comment="$COMMENT"
)

# get the hosted zone id
HOSTED_ZONE_ID=$(echo "$CREATE_RESPONSE" | jq -r '.HostedZone.Id' | awk -F'/' '{print $3}')

if [ -z "$HOSTED_ZONE_ID" ]; then
  echo "Failed to create hosted zone. Exiting."
  exit_with_error
fi

# retrieve the ns records
echo "Retrieving NS records for the hosted zone..."
NS_RECORDS=$(aws route53 list-resource-record-sets \
    --hosted-zone-id "$HOSTED_ZONE_ID" \
    --query "ResourceRecordSets[?Type == 'NS'].ResourceRecords[*].Value" \
    --output text
)

if [ -z "$NS_RECORDS" ]; then
  echo "No NS records found. Exiting."
  exit 1
fi

# echo ns records and hosted zone id to console
echo -e "NS Records for $DOMAIN_NAME:\n$NS_RECORDS"
echo -e "\nHosted Zone ID: $HOSTED_ZONE_ID"
