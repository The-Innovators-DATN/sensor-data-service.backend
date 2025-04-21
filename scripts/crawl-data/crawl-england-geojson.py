# -------------  crawler_all_basins.py  -------------
import requests, json, geopandas as gpd
from bs4 import BeautifulSoup
from edn_format import loads as edn_loads
from shapely.geometry import shape


# ── helpers (same as before, trimmed) ──────────────────────────
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


if __name__ == "__main__":
    crawl_england_basins()
