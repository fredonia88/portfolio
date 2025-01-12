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

# remove cloudformation and deploy directories on s3
aws s3 rm --recursive s3://fred-portfolio/cloudformation/
aws s3 rm --recursive s3://fred-portfolio/deploy/

# add cloudformation and deploy directories to s3
aws s3 cp --recursive cloudformation/ s3://fred-portfolio/cloudformation/
aws s3 cp --recursive deploy/ s3://fred-portfolio/deploy/

# create or update stack
if [[ "$build" == "create" ]]; then
    echo "Creating the stack 'fred-portfolio-cfs' in us-east-1..."
    aws cloudformation create-stack \
        --stack-name fred-portfolio-cfs \
        --template-url https://fred-portfolio.s3.amazonaws.com/cloudformation/portfolio-cfs.yml \
        --capabilities CAPABILITY_NAMED_IAM \
        --region us-east-1
        #--parameters ParameterKey=Env,ParameterValue=$profile \
        #    ParameterKey=DB,ParameterValue=$db \
        #    ParameterKey=SecretTemplate,ParameterValue=$secrets \
        #    ParameterKey=ClusterId,ParameterValue=$clusterid \
elif [[ "$build" == "update" ]]; then
    echo "Updating the stack 'fred-portfolio-cfs' in us-east-1..."
    aws cloudformation update-stack \
        --stack-name fred-portfolio-cfs \
        --template-url https://fred-portfolio.s3.amazonaws.com/cloudformation/portfolio-cfs.yml \
        --capabilities CAPABILITY_NAMED_IAM \
        --region us-east-1
        #--parameters ParameterKey=Env,ParameterValue=$profile \
        #    ParameterKey=DB,ParameterValue=$db \
        #    ParameterKey=SecretTemplate,ParameterValue=$secrets \
        #    ParameterKey=ClusterId,ParameterValue=$clusterid \
fi