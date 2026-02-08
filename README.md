# Portfolio

## How to run locally

Install:

 - VS Code
 - Python 3.14 (see deploy/user-data.sh in repository)
 - Git (see deploy/user-data.sh in repository)
 - Docker (see deploy/user-data.sh in repository)
 - Homebrew (for running locally)
 
Once you've installed the above, install the project. First create a python VE. With the VE activated, run: 

`pip install . -I --no-cache`

Next you need to configure the local database. The project can run using sqlite, but postgres is preferred. To install postgres run: 

 - Install homebrew: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
 - Install postgres (will start on boot): `brew install postgresql@15`
 - Start postgres: `brew services start postgresql@15`
 - Check postgres is running: `brew services list`
 - Stop postgres: `brew services stop postgresql@15`
 - Start postgres shell: `psql postgres`

Next, create a `_dev.env` file in portfolio/ to hold the secrets. You will need the following secrets to run the project: 

 - DJANGO_SECRET_KEY
 - DJANGO_DEBUG
 - DJANGO_EMAIL_HOST
 - DJANGO_EMAIL_PORT
 - DJANGO_EMAIL_USE_TLS
 - DJANGO_EMAIL_HOST_USER
 - DJANGO_EMAIL_HOST_PASSWORD
 - DJANGO_EMAIL_RECIPIENT
 - DJANGO_RECAPTCHA_PUBLIC_KEY
 - DJANGO_RECAPTCHA_PRIVATE_KEY
 - BLS_REGISTRATION_KEY
 - POSTGRES_DB
 - POSTGRES_USER
 - POSTGRES_PASSWORD
 - POSTGRES_SERVER
 - POSTGRES_PORT

 DJANGO_DEBUG and POSTGRES vars should be configured with local creds when using dev.

Next, add the secrets in the `_dev.env` file as env vars: 
 - `source get_secrets.sh --env=dev`

Lastly, run the ETL scripts: 
 - `python portfolio/scripts/cpi_index_by_year.py`
 - `python portfolio/scripts/median_income_by_age_constant_dollars.py`

With the database running and the secrets in place, you should be able to run the django project: 

 - `python manage.py makemigrations web`
 - `python manage.py migrate web`
 - `python manage.py runserver`
 
 ## How to run locally in docker
 
With secrets in `_staging.env` configured, run:
 - `source get_secrets.sh --env=dev`

Then run the below to build the image and compose up: 

 - `docker compose build --no-cache`
 - `docker compose up`

## How to deploy to AWS

Comment out the mounted aws credentials in docker-compose.yml:
- ~/.aws:/root/.aws:ro

Run `bash deploy/deploy-django-secrets.sh --build=create` to parse the `_staging.env` file and add its contents to the AWS secret. This secret will be used in production.

Pull down the secrets you just added to AWS and set as env vars: 
 - `source get_secrets.sh --env=prod`

Run all the scripts in the deploy folder to create the S3 bucket, secrets for the git keys for the EC2 instance (add the secret values manually), the project's secrets (which you should've already done), the Route53 hosted zone for the domains and finally the script to deploy the project's stack (EC2, Load Balancer, Listeners, Security Groups, Target Group, SSL Certificates, etc).

Use this command to create a bash shell in a container: `docker exec -i -t <mycontainer> /bin/bash`