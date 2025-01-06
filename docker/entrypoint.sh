#!/bin/bash

# collect static files and migrate
echo Collecting static files and migrating db...
python manage.py collectstatic --noinput
python manage.py migrate

# run ETL scripts
echo Running ETL scripts...
python scripts/cpi_index_by_year.py
python scripts/median_income_by_age_constant_dollars.py

# start gunicorn server
echo Starting gunicorn server...
gunicorn portfolio.wsgi:application --bind 0.0.0.0:8000