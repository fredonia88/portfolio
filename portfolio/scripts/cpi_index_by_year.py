"""
Process CPI-U-RS data from manually downloaded BLS Excel file.

BLS blocks automated requests, so file must be manually downloaded.
Script reads local Excel file, processes data, writes CSV, uploads to S3, and optionally writes to Postgres.

Download from: https://www.bls.gov/cpi/research-series/r-cpi-u-rs-home.htm
File: "Updated R-CPI-U-RS, All items, 1977-2024 (XLSX)"
Save to: files/r-cpi-u-rs-allitems.xlsx

Usage:
  --process-data: Process Excel, write CSV, upload to S3
  --write-db: Write to Postgres database
"""


from portfoliopython.etl_postgres import Postgres
import polars as pl
import os
import argparse
import boto3
from io import BytesIO


def process_data() -> pl.DataFrame:
    """Read, process, and save CPI data to CSV and S3."""

    df = pl.read_excel('files/r-cpi-u-rs-allitems.xlsx')

    headers = df.row(4)
    df = df[5:]
    df.columns = headers

    df = df.select(['YEAR', 'AVG'])
    df = df.rename({'YEAR': 'year', 'AVG': 'cpi_index_avg'})

    df = df.with_columns(
        pl.col('year').cast(pl.Int64),
        pl.col('cpi_index_avg').cast(pl.Float64)
    )

    df = df.with_columns(
        pl.when(pl.col('year') == 1977)
        .then(100)
        .otherwise(pl.col('cpi_index_avg'))
        .alias('cpi_index_avg')
    )

    print(df)

    # Write CSV locally
    csv_path = 'files/cpi_index_by_year.csv'
    df.write_csv(csv_path)

    # Upload to S3
    s3 = boto3.client('s3')
    s3.upload_file(csv_path, 'fred-portfolio', 'data/cpi_index_by_year.csv')

    return df


def write_to_db() -> None:
    """Read CSV from S3 and write to Postgres."""

    s3 = boto3.client('s3')
    obj = s3.get_object(Bucket='fred-portfolio', Key='data/cpi_index_by_year.csv')
    df = pl.read_csv(BytesIO(obj['Body'].read()))

    postgres = Postgres(
        username=os.getenv('POSTGRES_USER'),
        password=os.getenv('POSTGRES_PASSWORD'),
        server=os.getenv('POSTGRES_SERVER'),
        port=os.getenv('POSTGRES_PORT'),
        database=os.getenv('POSTGRES_DB')
    )
    postgres.write_df(df, 'cpi_index_by_year')


def main() -> None:
    parser = argparse.ArgumentParser(description='Process CPI data')
    parser.add_argument('--process-data', action='store_true', help='Create csv and upload to S3')
    parser.add_argument('--write-db', action='store_true', help='Read from S3 and write to database')
    args = parser.parse_args()

    if args.process_data:
        process_data()
    if args.write_db:
        write_to_db()


if __name__ == '__main__':
    main()