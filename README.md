# Portfolio

## How to get up and running

This project runs on an AWS EC2 micro instance and uses AWS's standard AMI which is based on CentOS. To run CentOS as a VM, install Oracle VM VirtualBox, download the latest CentOS ISO Image [here](https://www.centos.org/download/) (use x86_64 architecture) and create the VM in VirtualBox.

Once the VM is running, you'll need to configure it. Install:

 - VS Code
 - Python 3.11 (see deploy/user-data.sh in repository)
 - Git (see deploy/user-data.sh in repository)
 - Docker (see deploy/user-data.sh in repository)
 
Once your environment is configured, you can install the project. First create a python VE, then run: 

`python -m pip install . -I --no-cache`

Create two .env files to hold the secrets for dev and prod environments `_dev.env` and and `_prod.env`. The entries in these files can be the same, except for the database passwords in prod. You will need the following secrets to run the entire project: 

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

 Next you need to configure the local database. The project can run in sqlite, but postgres is preferred.


 Once you have `_prod.env` configured, you can run `deploy/deploy-django-secrets` which will parse the .env file and add them to the AWS in a secret.