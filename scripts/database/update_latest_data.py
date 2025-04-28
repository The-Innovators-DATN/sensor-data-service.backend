from shapely.geometry import shape
import psycopg2
import os
from dotenv import find_dotenv, load_dotenv
from clickhouse_driver import Client


# Load environment variables from ../../config/.env
load_dotenv(find_dotenv("../../config/.env"))


def postgresql_connect():
    """
    Connects to a PostgreSQL database and provides a function to map IDs by finding rows with a WHERE condition.
    """
    try:
        # Update these parameters with your PostgreSQL credentials
        conn = psycopg2.connect(
            dbname=os.environ.get("DATABASE_NAME"),
            user=os.environ.get("DATABASE_USER"),
            password=os.environ.get("DATABASE_PASSWORD"),
            host=os.environ.get("DATABASE_HOST"),
            port=os.environ.get("DATABASE_PORT"),
        )
        print("Connected to PostgreSQL database.")
        return conn
    except Exception as e:
        print(f"Failed to connect to PostgreSQL: {e}")
        return 
    
def clickhouse_connect():
    """
    Connects to a ClickHouse database and provides a function to map IDs by finding rows with a WHERE condition.
    """
    try:
        # Update these parameters with your ClickHouse credentials
        conn = Client(
            host=os.environ.get("CLICKHOUSE_HOST", "localhost"),
            user=os.environ.get("CLICKHOUSE_USER", "default"),
            password=os.environ.get("CLICKHOUSE_PASSWORD", ""),
            port=int(os.environ.get("CLICKHOUSE_PORT", 9000)),
            database=os.environ.get("CLICKHOUSE_DATABASE", "default"),
        )
        print("Connected to ClickHouse database.")
        return conn
    except Exception as e:
        print(f"Failed to connect to ClickHouse: {e}")
        return None
    
def update_latest_data():
    """
    Updates the latest data in the PostgreSQL database by fetching it from the ClickHouse database.
    """
    pg_conn = postgresql_connect()
    ch_conn = clickhouse_connect()

    if not pg_conn or not ch_conn:
        print("Failed to connect to one of the databases.")
        return

    try:
        # for each station parameter, get the latest data from ClickHouse
        with pg_conn.cursor() as pg_cursor:
            # Fetch the latest data from ClickHouse
            pg_cursor.execute("SELECT station_id, parameter_id FROM station_parameter")
            station_parameters = pg_cursor.fetchall()
            ## For eact station_parameter, get the latest data from ClickHouse then update again st PostgreSQL to colun last_receiv_at, last_value
            for station_id, parameter_id in station_parameters:
                # Fetch the latest data from ClickHouse
                ch_query = f"""
                    SELECT value, datetime
                    FROM {os.environ.get("CLICKHOUSE_NAME")}.messages_sharded
                    WHERE station_id = '{station_id}' AND metric_id = '{parameter_id}'
                    ORDER BY datetime DESC
                    LIMIT 1
                """
                latest_data = ch_conn.execute(ch_query)
                
                if latest_data:
                    value, received_at = latest_data[0]
                    # Update the PostgreSQL database with the latest data
                    print(f"Latest data for station {station_id}, parameter {parameter_id}: {value}, received at: {received_at}")
                    # print(f"Type of all data: {type(value)}")
                    # print(f"Type of last_rec: {type(received_at)}") 
                    try:
                        # Update the PostgreSQL database with the latest data
                        pg_update_query = """
                            UPDATE station_parameter
                            SET last_value = %s, last_receiv_at = %s
                            WHERE station_id = %s AND parameter_id = %s
                        """
                        pg_cursor.execute(pg_update_query, (value, received_at, station_id, parameter_id))
                        pg_conn.commit()
                    except Exception as e:
                        print(f"Failed to update PostgreSQL: {e}")
    except Exception as e:
        print(f"Error updating latest data: {e}")
    finally:
        pg_conn.close()
        ch_conn.disconnect()

if __name__ == "__main__":
    update_latest_data()