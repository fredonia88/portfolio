# add python:
sudo yum update
sudo yum install python3.11-pip -y
sudo ln -sf /usr/bin/python3.11 /usr/bin/python
sudo ln -sf /usr/bin/pip3.11 /usr/bin/pip

# install git:
sudo dnf install git-all -y

# add docker -- could not get repo to install docker compose, so installing manually (https://docs.docker.com/compose/install/linux/#install-the-plugin-manually):
sudo yum install docker -y
DOCKER_CONFIG=${DOCKER_CONFIG:-/usr/local/lib/docker/cli-plugins}
sudo mkdir -p $DOCKER_CONFIG
sudo curl -SL https://github.com/docker/compose/releases/download/v2.23.1/docker-compose-linux-x86_64 -o $DOCKER_CONFIG/docker-compose
sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
sudo systemctl start docker
docker compose version

# install cron
sudo yum install cronie -y
sudo systemctl start crond.service

# create django user:
sudo useradd -m django

# add django to docker group:
sudo usermod -aG docker django

# allow django to run sudo without a password:
echo "django ALL=(ALL) NOPASSWD:ALL" | sudo tee -a /etc/sudoers

# assume django user:
sudo -u django -H bash << EOF
cd ~
# create directories:
mkdir repos
mkdir venvs

# create virtual environment:
cd venvs
python -m venv portfolio
cd ~

# pull git keys and write them
SECRET=\$(aws secretsmanager get-secret-value --secret-id fred-portfolio-ec2-git-keys --query 'SecretString' --output text)
PRIVATE_KEY=\$(echo \$SECRET | jq -r '.private_key')
PUBLIC_KEY=\$(echo \$SECRET | jq -r '.public_key')
mkdir -p /home/django/.ssh
echo "\$PRIVATE_KEY" > /home/django/.ssh/id_rsa
chmod 600 /home/django/.ssh/id_rsa
echo "\$PUBLIC_KEY" > /home/django/.ssh/id_rsa.pub
chmod 644 /home/django/.ssh/id_rsa.pub

# create ssh config file
echo -e "Host * \n    StrictHostKeyChecking no" > ~/.ssh/config
chmod 644 ~/.ssh/config

# clone repo:
cd repos
git clone git@github.com:fredonia88/portfolio.git

# source venv, cd into repo and install the project
cd portfolio
source ~/venvs/portfolio/bin/activate
python -m pip install . -I --no-cache

# set secrets
source get_secrets.sh --env=prod

# docker compose build and up
docker compose build
docker compose up
EOF
