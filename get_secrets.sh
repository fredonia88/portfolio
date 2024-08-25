#!/bin/bash

# print error message and exit
exit_with_error() {
    echo "$1"
    exit 1
}

SECRET_VALUE=$(aws secretsmanager get-secret-value --secret-id fred-portfolio-django --query SecretString --output text --profile personal)
export DJANGO_SECRET_KEY=$(echo $SECRET_VALUE | jq -r .secret_key)
export DJANGO_EMAIL_HOST=$(echo $SECRET_VALUE | jq -r .email_host)
export DJANGO_EMAIL_PORT=$(echo $SECRET_VALUE | jq -r .email_port)
export DJANGO_EMAIL_USE_TLS=$(echo $SECRET_VALUE | jq -r .email_user_tls)
export DJANGO_EMAIL_HOST_USER=$(echo $SECRET_VALUE | jq -r .email_host_user)
export DJANGO_EMAIL_HOST_PASSWORD=$(echo $SECRET_VALUE | jq -r .email_host_password)
export DJANGO_EMAIL_RECIPIENT=$(echo $SECRET_VALUE | jq -r .email_recipient)
export DJANGO_RECAPTCHA_PUBLIC_KEY=$(echo $SECRET_VALUE | jq -r .recaptcha_public_key)
export DJANGO_RECAPTCHA_PRIVATE_KEY=$(echo $SECRET_VALUE | jq -r .recaptcha_private_key)
unset SECRET_VALUE