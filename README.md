# Portfolio

## How to run locally

This project runs on an AWS EC2 micro instance and uses AWS's standard AMI which is based on CentOS. To run CentOS as a VM, install [Oracle VM VirtualBox](https://www.virtualbox.org/wiki/Downloads), download the latest [CentOS ISO Image](https://www.centos.org/download/) (use x86_64 architecture) and create the VM in VirtualBox.

Once the VM is running, you'll need to provision its environment with the necessary dev tools. Install:

 - VS Code
 - Python 3.11 (see deploy/user-data.sh in repository)
 - Git (see deploy/user-data.sh in repository)
 - Docker (see deploy/user-data.sh in repository)
 
Once your environment is configured, install the project. First create a python VE. With the VE activated, run: 

`python -m pip install . -I --no-cache`

Create two .env files to hold the secrets for dev and prod environments `_dev.env` and `_prod.env`. The entries in these files can be the same, except for the prod database password. You will need the following secrets to run the project: 

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
 - POSTGRES_DB
 - POSTGRES_USER
 - POSTGRES_PASSWORD

 Next you need to configure the local database. The project can run using sqlite, but postgres is preferred. To install postgres run: 

 - `sudo yum install -y https://download.postgresql.org/pub/repos/yum/15/redhat/rhel-9-x86_64/pgdg-redhat-repo-latest.noarch.rpm` (pull down the postgres repo)
 - `sudo dnf -qy module disable postgresql` (disables active version of postgres)
 - `sudo yum install -y postgresql15 postgresql15-server` (install postgres and and postgres server)
 - `sudo /usr/pgsql-15/bin/postgresql-15-setup initdb` (init the db)
 - `sudo systemctl enable postgresql-15` (enable to start on boot)
 - `sudo systemctl start postgresql-15` (start the service)

Lastly, add the secrets in the `_dev.env` file as env vars: 
 - `source get_secrets.sh --env=dev`

With the database running and the secrets in place, you should be able to run the django project: 

 - `python manage.py makemigrations web`
 - `python manage.py migrate web`
 - `python manage.py runserver`

 
 ## How to run locally in docker
 
With `_prod.env` configured, run `bash deploy/deploy-django-secrets.sh --build=create` to parse the `_prod.env` file and add its contents to an AWS secret. This secret will be used in production.

Pull down the secrets you just added to AWS and set as env vars: 
 - `source get_secrets.sh --env=prod` 

Lastly, run the below to build the image and compose up: 

 - `docker compose build --no-cache`
 - `docker compose up`

## How to deploy to AWS

Run all the scripts in the deploy folder to create the S3 bucket, secrets for the git keys for the EC2 instance (add the secret values manually), the project's secrets (which you should've already done), the Route53 hosted zone for the domains and finally the script to deploy the project's stack (EC2, Load Balancer, Listeners, Security Groups, Target Group, SSL Certificates, etc).