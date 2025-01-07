from portfoliopython.etl_postgres import Postgres
import requests
from io import BytesIO
import polars as pl
import os

headers = {
    'referer': 'https://www.bls.gov/cpi/research-series/r-cpi-u-rs-home.htm',
    #'sec-ch-ua': '"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"',
    #'sec-ch-ua-mobile': '?1',
    #'sec-ch-ua-platform': '"Android"',
    #'upgrade-insecure-requests': '1',
    'user-agent': 'Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Mobile Safari/537.36'
}
response = requests.get('https://www.bls.gov/cpi/research-series/r-cpi-u-rs-allitems.xlsx', headers=headers)

data = BytesIO(response.content)
df = pl.read_excel(data)

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

postgres = Postgres(
    username=os.getenv('POSTGRES_USER'),
    password=os.getenv('POSTGRES_PASSWORD'),
    server=os.getenv('POSTGRES_SERVER'),
    port=os.getenv('POSTGRES_PORT'),
    database=os.getenv('POSTGRES_DB')
)
postgres.write_df(df, 'cpi_index_by_year')