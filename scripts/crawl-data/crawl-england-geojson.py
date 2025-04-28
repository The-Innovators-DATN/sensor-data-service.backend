# -------------  crawler_all_basins.py  -------------
import requests, json, geopandas as gpd
from bs4 import BeautifulSoup
from edn_format import loads as edn_loads
from shapely.geometry import shape
import psycopg2
import os
from dotenv import find_dotenv, load_dotenv


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
        return None

def map_id_by_condition(conn, table_name, column_name, condition):
    """
    Maps IDs by finding rows in the specified table with a WHERE condition.

    Args:
        conn: The PostgreSQL connection object.
        table_name (str): The name of the table to query.
        column_name (str): The column to retrieve.
        condition (str): The WHERE condition for filtering rows.

    Returns:
        list: A list of IDs matching the condition.
    """
    try:
        with conn.cursor() as cursor:
            query = f"SELECT {column_name} FROM {table_name} WHERE {condition};"
            cursor.execute(query)
            results = cursor.fetchall()
            return [row[0] for row in results]
    except Exception as e:
        print(f"Error querying database: {e}")
        return []

# ── helpers (sIame as before, trimmed) ──────────────────────────
def edn_to_native(obj):
    from edn_format.immutable_list import ImmutableList
    from edn_format.immutable_dict import ImmutableDict

    if isinstance(obj, (ImmutableList, list)):
        return [edn_to_native(x) for x in obj if not (isinstance(x, str) and x.startswith("~"))]
    if isinstance(obj, (ImmutableDict, dict)):
        return {str(k).lstrip("~:"): edn_to_native(v) for k, v in obj.items()}
    if isinstance(obj, str):
        try:
            return float(obj)
        except ValueError:
            return obj
    return obj


def is_lonlat(p): return isinstance(p, list) and len(p) == 2 and all(isinstance(c, (int, float)) for c in p)


def unwrap_coords(coords):
    """
    Recursively unwraps ^6-encoded EDN structures but preserves [lon, lat] as-is.
    """
    if isinstance(coords, list):
        # unwrap if ^6 format: ["^6", [...]]
        if len(coords) == 2 and isinstance(coords[0], str) and coords[0].startswith("^"):
            return unwrap_coords(coords[1])
        if is_lonlat_pair(coords):
            return coords
        return [unwrap_coords(c) for c in coords]
    return coords

def is_lonlat_pair(x):
    return (
        isinstance(x, list)
        and len(x) == 2
        and all(isinstance(c, (float, int)) for c in x)
    )

def extract_named_keys(lst):
    out = {}
    i = 0
    while i < len(lst) - 1:
        k, v = lst[i], lst[i + 1]
        if k == "^:":
            out["name"] = v
        elif k == "^;":
            out["page-url"] = v
        elif k == "^<":
            out["catchment"] = v
        i += 2
    return out

def truncate_coords(coords):
    """
    Recursively trims any coordinate triplet [lon, lat, extra] → [lon, lat]
    """
    if isinstance(coords, list):
        if len(coords) == 3 and all(isinstance(c, (int, float)) for c in coords):
            return coords[:2]
        elif is_lonlat_pair(coords):
            return coords
        else:
            return [truncate_coords(c) for c in coords]
    return coords

# ── main ───────────────────────────────────────────────────────
def crawl_england_basins(out_path="river_basins.geojson"):
    url = "https://environment.data.gov.uk/catchment-planning"
    soup = BeautifulSoup(requests.get(url, timeout=30).content, "html.parser")
    edn_root = edn_loads(soup.select_one("div.cljc-component")["data-init"])

    # 1️⃣ locate the results list
    # 🧠 locate ~:results
    results_idx = edn_root.index("~:results")
    raw_results = edn_root[results_idx + 1]

    # 🔧 unwrap ImmutableList → native list
    raw_results_native = edn_to_native(raw_results)

    # 🎯 unwrap ^6 wrapper if needed
    if isinstance(raw_results_native, list) and raw_results_native[0] == "^6":
        basins_list = raw_results_native[1]
    elif isinstance(raw_results_native, list) and isinstance(raw_results_native[0], list):
        basins_list = raw_results_native
    else:
        raise RuntimeError("~:results has unknown structure")

    print(f"✅ Found {len(basins_list)} river basins.")
    with open("basins_list.json", "w") as f:
        json.dump(basins_list, f, indent=2)
    print("✅ Saved raw basins list to basins_list.json")

    features = []
    for i, raw in enumerate(basins_list):

        # Try both formats
        raw_native = edn_to_native(raw)

        # unwrap '^ ' nếu còn dính bên ngoài
        if isinstance(raw_native, list) and raw_native and raw_native[0] == "^ ":
            raw_native = raw_native[1:]

        # 🧠 Nếu là dạng positional (Anglian)
        if isinstance(raw_native, list) and len(raw_native) == 3 and isinstance(raw_native[2], list):
            raw_dict = {
                "name": raw_native[0],
                "page-url": raw_native[1],
                "catchment": raw_native[2]
            }

        # 🧠 Nếu là dạng named keys (Dee, Humber,...)
        elif isinstance(raw_native, list) and "^:" in raw_native:
            raw_dict = extract_named_keys(raw_native)

        else:
            print(f"⚠️ Skipped basin {i}: unknown structure →", raw_native[:5])
            continue
        # #for each basin write raw_dict to file
        # with open(f"basin_{i}.json", "w") as f:
        #     json.dump(raw_dict, f, indent=2)
        # print(f"Saved raw basin {i} to basin_{i}.json")
        name = raw_dict.get("name")
        page_url = "https://environment.data.gov.uk" + raw_dict.get("page-url", "")
        basin_id = page_url.rsplit("/", 1)[-1]

        catch_blk = raw_dict.get("catchment", {})
        if isinstance(catch_blk, list) and catch_blk and catch_blk[0] == "^ ":
            catch_blk = catch_blk[1:]  # unwrap inner

        if isinstance(catch_blk, list):
            catch_dict = {catch_blk[i]: catch_blk[i+1] for i in range(0, len(catch_blk) - 1, 2)}
        else:
            print(f"Invalid catchment block for {name}")
            continue

        gtype = catch_dict.get("^4", "MultiPolygon")
        raw_coords = catch_dict.get("^5", [])
        with open(f"basin_coords_{i}.json", "w") as f:
            json.dump(raw_coords, f, indent=2)
        print(f"Saved raw coords for {name} to basin_coords_{i}.json")

        try:
            unwrapped = unwrap_coords(edn_to_native(raw_coords))

            coords = truncate_coords(unwrapped)
            with open(f"basin_coords_unwrapped_{i}.json", "w") as f:
                json.dump(coords, f, indent=2)
            print(f"Saved unwrapped coords for {name} to basin_coords_unwrapped_{i}.json")

            # coords = force_multipolygon_shape(coords)
            # coords = close_rings(coords)
            geom = shape({"type": gtype, "coordinates": coords})
            print(f"Parsed geometry for {name}")
        except Exception as e:
            print(f"Error parsing geometry for {name}: {e}")
            continue

        if not basin_id.isdigit():
            print(f"Skipped basin {name} due to invalid ID: {basin_id}")
            continue

        features.append({
            "geometry": geom,
            "id": int(basin_id),
            "name": name,
            "url": page_url
        })

    if not features:
        raise RuntimeError("No valid basin features extracted.")

    print("✅ First feature keys:", features[0].keys())

    gdf = gpd.GeoDataFrame(features, geometry="geometry", crs="EPSG:4326")
    gdf.to_file(out_path, driver="GeoJSON")
    print(f"Saved {len(gdf)} basins → {out_path}")

def crawl_england_catchements(out_path="england_catchments.geojson", river_basin_id=4):  # Default is Humber
    url = f"https://environment.data.gov.uk/catchment-planning/RiverBasinDistrict/{river_basin_id}"
    soup = BeautifulSoup(requests.get(url, timeout=30).content, "html.parser")
    edn_root = edn_loads(soup.select_one("div.cljc-component")["data-init"])

    # Convert to native Python structure
    native_edn_root = edn_to_native(edn_root)

    # Extract catchment data
    catchment_data = native_edn_root[8][1]

    features = []
    for i, raw in enumerate(catchment_data):
        if isinstance(raw, list) and raw and raw[0] == "^ ":
            raw = raw[1:]  # Unwrap
        if not raw:
            continue

        # 1️⃣ Kiểm tra: dạng "Simple" hay "Key-Value"
        if isinstance(raw, list) and isinstance(raw[0], str) and not raw[0].startswith("^"):
            # Dạng 1: Đơn giản (name, url, geometry)
            try:
                name = raw[0]
                page_url = "https://environment.data.gov.uk" + raw[1]
                catchment_id = page_url.rsplit("/", 1)[-1]
                catch_blk = raw[2]

            except (IndexError, ValueError) as e:
                print(f"Invalid simple catchment structure: {e}")
                continue
        elif isinstance(raw, list) and all(isinstance(x, (str, list)) for x in raw):
            # with open(f"catchment_{i}.json", "w") as f:
            #     json.dump(raw, f, indent=2)
            try:

                if raw[0] == "^@":
                    # Start new catchment object
                    name = raw[1]
                    current = {"name": name}
                    page_url = "https://environment.data.gov.uk" + raw[4]
                    catchment_id = page_url.rsplit("/", 1)[-1]
                    catch_blk = raw[6]

                with open(f"catchment_coords_{i}.json", "w") as f:
                    json.dump(catch_blk, f, indent=2)

                print(f"Name: {name}, page-url: {page_url}, id: {catchment_id}")
            except Exception as e:
                print(f"Invalid key-value catchment structure: {e}")
                continue

        if isinstance(catch_blk, list):
            if catch_blk[0] == "^ ":
                catch_blk = catch_blk[1:]
            if len(catch_blk) % 2 == 0:
                catch_blk = {catch_blk[j]: catch_blk[j+1] for j in range(0, len(catch_blk), 2)}
            else:

                raise ValueError(f"Catchment block does not have key-value pairs: {len(catch_blk)}")

        gtype = catch_blk.get("^:", "MultiPolygon")
        raw_coords = catch_blk.get("^;", [])

        try:
            unwrapped = unwrap_coords(edn_to_native(raw_coords))
            coords = truncate_coords(unwrapped)

            geom = shape({"type": gtype, "coordinates": coords})
            print(f"Parsed geometry for {name}")
        except Exception as e:
            print(f"Error parsing geometry for {name}: {e}")
            continue

        if not catchment_id.isdigit():
            print(f"Skipped catchment {name} due to invalid ID: {catchment_id}")
            continue
        
        conn = postgresql_connect()
        if conn:
            # Check if the catchment ID exists in the database
            existing_ids = map_id_by_condition(conn, "catchment", "id", f"name = '{name}'")
            if existing_ids:
                catchment_id = int(existing_ids[0])
                print(f"Catchment ID {existing_ids} already exists in the database.")
            else:
                catchment_id = 0
                print(f"Catchment ID {catchment_id} does not exist in the database.")
            conn.close()
        features.append({
            "geometry": geom,
            "id": int(catchment_id),
            "name": name,
            "url": page_url
        })
    if not features:
        raise RuntimeError("No valid basin features extracted.")

    gdf = gpd.GeoDataFrame(features, geometry="geometry", crs="EPSG:4326")
    gdf.to_file(out_path, driver="GeoJSON")
    print(f"Saved {len(gdf)} basins → {out_path}")



if __name__ == "__main__":
    # crawl_england_basins()
    crawl_england_catchements()
