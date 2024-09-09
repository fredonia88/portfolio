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

# generate keys in temporary file and delete
echo "Generating rsa keys..."
TMP_PRIVATE_KEY=$(mktemp)
chmod 600 "$TMP_PRIVATE_KEY"
ssh-keygen -t rsa -b 2048 -N '' -C "" -f "$TMP_PRIVATE_KEY" <<< y
PRIVATE_KEY=$(cat "$TMP_PRIVATE_KEY")
PUBLIC_KEY=$(ssh-keygen -y -f "$TMP_PRIVATE_KEY")
SECRET_PAYLOAD=$(jq -n --arg private "$PRIVATE_KEY" --arg public "$PUBLIC_KEY" \
    '{private_key: $private, public_key: $public}')

rm -f "$TMP_PRIVATE_KEY"

# create or update secret
if [[ "$build" == "create" ]]; then
    echo "Creating the secret..."
    aws secretsmanager create-secret \
        --name "fred-portfolio-ec2-git-keys" \
        --secret-string "$SECRET_PAYLOAD" \
        --profile personal
elif [[ "$build" == "update" ]]; then
    echo "Updating the secret..."
    aws secretsmanager put-secret-value \
        --secret-id "fred-portfolio-ec2-git-keys" \
        --secret-string "$SECRET_PAYLOAD" \
        --profile personal
fi
