from portfoliopython.etl_postgres import Postgres
import requests
import json
import polars as pl
import os

"""
Series IDs
https://www.bls.gov/webapps/legacy/cpswktab3.htm
LEU0252881500 -- 16 years and over, use
LEU0252886300 -- 16 to 24 years, use
LEU0252886500 -- 16 to 19 years
LEU0252887100 -- 20 to 24 years
LEU0252887700 -- 25 years and over
LEU0252887900 -- 25 to 54 years
LEU0252888500 -- 25 to 34 years, use
LEU0252889100 -- 35 to 44 years, use
LEU0252889700 -- 45 to 54 years, use
LEU0252890300 -- 55 years and over
LEU0252890900 -- 55 to 64 years, use
LEU0252891500 -- 65 years and over, use
"""

# pull data from the api, create a df from each body and add to the dfs list
start = 1979
end = 2023
dfs = []
for i in range(start, end, 20):
    url = 'https://api.bls.gov/publicAPI/v2/timeseries/data/'
    startyear = str(i)
    endyear = str(i + 19) if i + 19 <= end else str(end)
    print(startyear, endyear)
    data = json.dumps({
        "seriesid":["LEU0252881500", "LEU0252886300", "LEU0252888500", "LEU0252889100", "LEU0252889700", "LEU0252890900", "LEU0252891500"],
        "startyear":startyear,
        "endyear": endyear,
        "catalog":True,#|false,
        "calculations":False,#|false,
        "annualaverage":True,
        "aspects":False,#|false,
        "registrationkey":os.getenv('BLS_REGISTRATION_KEY')
    })
    headers = {'Content-type': 'application/json'}
    response = requests.post(url, data=data, headers=headers)
    json_data = json.loads(response.text)
    
    if response.next:
        raise Exception('Response contains next! Be sure to accomodate this!')

    datasets = json_data['Results']['series']
    for i in range(len(datasets)):
        
        catalog = json_data['Results']['series'][i]['catalog']
        data = json_data['Results']['series'][i]['data']
        
        catalog_df = pl.DataFrame(catalog)
        data_df = pl.DataFrame(data)
        data_df = data_df.filter(pl.col('period') == "Q05")
        catalog_df = catalog_df.select(pl.all().repeat_by(len(data_df)).flatten())
        df = pl.concat([catalog_df, data_df], how='horizontal')
        df = df.drop('footnotes') # column is inconsistent
        dfs.append(df)

# concat all dfs into a single df, add CPI data and convert to constant dollars
final_df = pl.concat(dfs, how='vertical')
final_df = pl.DataFrame(final_df, schema_overrides={"year": pl.Int64, "value": pl.Float64})
final_df = final_df.filter(pl.col('demographic_age') != '16 years and over')

cpi_df = pl.read_csv('/home/fred/Downloads/CPI Index by Year - Sheet1.csv')
cpi_df = cpi_df.rename({'YEAR': 'year', 'AVG': 'cpi_index_avg'})

final_df = final_df.join(cpi_df, 'year', 'left')
base_index = final_df.filter(pl.col('year') == 2023).select('cpi_index_avg').min().item()
final_df = final_df.with_columns(pl.lit(base_index).alias('base_index'))

final_df = final_df.with_columns((pl.col('base_index') / pl.col('cpi_index_avg')).alias('ratio'))
final_df = final_df.with_columns((pl.col('ratio') * pl.col('value') * 52).alias('yearly_value_constant_dollars'))

# ETL dataset into postgres
final_df = final_df.select(['year', 'demographic_age', 'yearly_value_constant_dollars'])
final_df = final_df.sort(['demographic_age', 'year'], descending=[False, False])

postgres = Postgres(
    username=os.getenv('POSTGRES_USER'),
    password=os.getenv('POSTGRES_PASSWORD'),
    server=os.getenv('POSTGRES_SERVER'),
    port=os.getenv('POSTGRES_PORT'),
    database=os.getenv('POSTGRES_DB')
)
postgres.write_df(final_df, 'median_income_by_age_constant_dollars', add_id_column=True)
