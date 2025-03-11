import requests
from bs4 import BeautifulSoup
import csv

BASE_URL = "https://vi.wikipedia.org"

def fetch_soup(url):
    """Fetch the URL and return a BeautifulSoup object if successful."""
    response = requests.get(url)
    if response.status_code == 200:
        return BeautifulSoup(response.content, "html.parser")
    else:
        print(f"Failed to retrieve {url}. Status code: {response.status_code}")
        return None

def extract_category_links(soup):
    """
    Given a BeautifulSoup object, extract all links from within the div
    with class "mw-category mw-category-columns" (i.e. within each group).
    Returns a list of dictionaries with 'title' and 'url'.
    """
    links = []
    parent_div = soup.find("div", class_="mw-category mw-category-columns")
    if not parent_div:
        return links

    group_divs = parent_div.find_all("div", class_="mw-category-group")
    for group in group_divs:
        for a in group.find_all("a"):
            href = a.get("href")
            title = a.get("title")
            if href and title:
                full_url = BASE_URL + href
                links.append({"title": title.strip(), "url": full_url})
    return links

def get_first_paragraph_text(url):
    """
    Request the URL, locate the <div> with class "mw-content-ltr mw-parser-output",
    then remove citation tags (<sup class="reference">) and extract the normalized text 
    of its first <p> tag.
    """
    soup = fetch_soup(url)
    if soup:
        content_div = soup.find("div", class_="mw-content-ltr mw-parser-output")
        if content_div:
            first_p = content_div.find("p")
            if first_p:
                # Remove citation elements (e.g., <sup class="reference">...</sup>)
                for sup in first_p.find_all("sup", class_="reference"):
                    sup.decompose()
                # Use get_text with a separator to ensure inline elements are separated by a space.
                text = first_p.get_text(separator=" ", strip=True)
                # Normalize whitespace by joining split words.
                return " ".join(text.split())
    return ""

def normalize_parent(title):
    """
    Remove the prefix 'Thể loại:Sông tại ' from the given title if it exists.
    """
    prefix_province = "Thể loại:Sông tại "
    prefix_river_basin = "Thể loại:"
    if title.startswith(prefix_river_basin):
        return title[len(prefix_river_basin):].strip()
    return title

def process_category_page(url, parent=""):
    """
    Processes a category page at the given URL.
    It extracts all water body links and their description (first paragraph text)
    and records the parent category as 'catchment'.
    If a linked page itself contains category groups, it is processed recursively.
    Returns a list of dictionaries with keys: 'catchment', 'water_body', and 'description'.
    """
    soup = fetch_soup(url)
    if not soup:
        return []
    
    results = []
    links = extract_category_links(soup)
    for link in links:
        water_body = link["title"]
        description = get_first_paragraph_text(link["url"])
        normalized_parent = normalize_parent(parent) if parent else ""
        if description:
            results.append({
            "catchment": normalized_parent,
            "water_body": water_body,
            "description": description
            })
            print({
            "river_basin": normalized_parent,
            "water_body": water_body,
            "description": description
            })
        # Check if the linked page is itself a category page
        child_soup = fetch_soup(link["url"])
        if child_soup and child_soup.find("div", class_="mw-category mw-category-columns"):
            # Recursively process child category page; pass the current water body as parent.
            results.extend(process_category_page(link["url"], parent=water_body))
    return results

# bentre_wb= "https://vi.wikipedia.org/wiki/Th%E1%BB%83_lo%E1%BA%A1i:S%C3%B4ng_t%E1%BA%A1i_B%E1%BA%BFn_Tre"
# province_wb = "https://vi.wikipedia.org/wiki/Th%E1%BB%83_lo%E1%BA%A1i:S%C3%B4ng_Vi%E1%BB%87t_Nam_theo_t%E1%BB%89nh_th%C3%A0nh"
river_basin_water_body = "https://vi.wikipedia.org/wiki/Th%E1%BB%83_lo%E1%BA%A1i:H%E1%BB%87_th%E1%BB%91ng_s%C3%B4ng_Vi%E1%BB%87t_Nam"
all_results = process_category_page(river_basin_water_body)

csv_file = "../../dataset/river_basin_water_bodies.csv"
with open(csv_file, "w", newline="", encoding="utf-8") as f:
    fieldnames = ["catchment", "water_body", "description"]
    writer = csv.DictWriter(f, fieldnames=fieldnames)
    writer.writeheader()
    for row in all_results:
        writer.writerow(row)

print(f"Data written to {csv_file}")
