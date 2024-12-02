#!/bin/bash

# print error message and exit
exit_with_error() {
    echo "$1"
    exit 1
}

# check if 'env' argument is provided
if ! [[ "$@" =~ "--env=" ]]; then
    exit_with_error "Error: '--env=prod|dev' argument is required."
fi

# parse the arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --env=*)
            env="${1#*=}"
            shift
            ;;
        *)
            exit_with_error "Error: Invalid argument: $1. Valid options are: '--env=prod|dev'"
            ;;
    esac
done

# validate the value of env
if [[ "$env" != "prod" && "$env" != "dev" ]]; then
    exit_with_error "Error: Invalid 'env' value. Allowed values: prod|dev"
fi

export DJANGO_ENV=$env

if [[ "$env" = "dev" ]]; then
    export DJANGO_DEBUG=True
    set -a # ensure variables in .env are automatically exported
    source portfolio/.env
    set +a
else
    SECRET_VALUE=$(aws secretsmanager get-secret-value --secret-id fred-portfolio-django --query SecretString --output text)
    export DJANGO_SECRET_KEY=$(echo $SECRET_VALUE | jq -r .DJANGO_SECRET_KEY)
    export DJANGO_EMAIL_HOST=$(echo $SECRET_VALUE | jq -r .DJANGO_EMAIL_HOST)
    export DJANGO_EMAIL_PORT=$(echo $SECRET_VALUE | jq -r .DJANGO_EMAIL_PORT)
    export DJANGO_EMAIL_USE_TLS=$(echo $SECRET_VALUE | jq -r .DJANGO_EMAIL_USE_TLS)
    export DJANGO_EMAIL_HOST_USER=$(echo $SECRET_VALUE | jq -r .DJANGO_EMAIL_HOST_USER)
    export DJANGO_EMAIL_HOST_PASSWORD=$(echo $SECRET_VALUE | jq -r .DJANGO_EMAIL_HOST_PASSWORD)
    export DJANGO_EMAIL_RECIPIENT=$(echo $SECRET_VALUE | jq -r .DJANGO_EMAIL_RECIPIENT)
    export DJANGO_RECAPTCHA_PUBLIC_KEY=$(echo $SECRET_VALUE | jq -r .DJANGO_RECAPTCHA_PUBLIC_KEY)
    export DJANGO_RECAPTCHA_PRIVATE_KEY=$(echo $SECRET_VALUE | jq -r .DJANGO_RECAPTCHA_PRIVATE_KEY)
    export POSTGRES_DB=$(echo $SECRET_VALUE | jq -r .POSTGRES_DB)
    export POSTGRES_USER=$(echo $SECRET_VALUE | jq -r .POSTGRES_USER)
    export POSTGRES_PASSWORD=$(echo $SECRET_VALUE | jq -r .POSTGRES_PASSWORD)
    export DJANGO_DEBUG=False # overwrite in case this is set to True in .env file
    unset SECRET_VALUE
fi