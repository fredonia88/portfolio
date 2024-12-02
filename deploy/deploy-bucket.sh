#!/bin/bash

# print error message and exit
exit_with_error() {
    echo "$1"
    exit 1
}

echo "Creating the bucket 'fred-portfolio' in us-east-1..."
aws s3api create-bucket \
    --bucket fred-portfolio \
    --region us-east-1
