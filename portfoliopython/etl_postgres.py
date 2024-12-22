from __future__ import annotations
import polars as pl

class Postgres:
    """Class to read and write from postgres db."""

    def __init__(
        self,
        username: str,
        password: str,
        server: str,
        port: int,
        database: str
    ):
        self.uri = f'postgresql://{username}:{password}@{server}:{port}/{database}'

    def read_table(self, query: str) -> pl.DataFrame:
        """Run a query against postgres db into a polars df.
        
        :param query: The query to run.
        :type query: str
        :returns: A polars dataframe.
        :rtype: pl.DataFrame
        """
        return pl.read_database_uri(query=query, uri=self.uri)

    def write_df(self, df: pl.DataFrame, table_name: str) -> None:
        """Write a polars df to postgres db.
        
        :param df: The dataframe to write.
        :type df: pl.DataFrame
        :param table_name: The table name to write to.
        :type table_name: str
        :returns: None.
        :rtype: None
        """
        df.write_database(table_name=table_name, connection=self.uri, if_table_exists='replace')
        