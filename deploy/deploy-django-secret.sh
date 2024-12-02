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

ENV_FILE="portfolio/.env"

# check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    exit_with_error "Error: .env file not found!"
fi

# convert .env file to JSON string
SECRET_PAYLOAD=$(jq -n --rawfile env $ENV_FILE '
  $env | split("\n") | map(select(length > 0)) |
  map(
    . as $line |
    {
      ( $line | split("=")[0] ): ( $line | split("=")[1:] | join("=") | gsub("^\"|\"$"; "") )
    }
  ) |
  add
')

# create or update secret
if [[ "$build" == "create" ]]; then
    echo "Creating the secret..."
    aws secretsmanager create-secret \
        --name "fred-portfolio-django" \
        --secret-string "$SECRET_PAYLOAD"
elif [[ "$build" == "update" ]]; then
    echo "Updating the secret..."
    aws secretsmanager put-secret-value \
        --secret-id "fred-portfolio-django" \
        --secret-string "$SECRET_PAYLOAD"
fi
