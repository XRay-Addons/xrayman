import requests
import json
from bs4 import BeautifulSoup

URL = "https://www.color-hex.com/color-palettes/popular.php"

headers = {
    "User-Agent": "Mozilla/5.0 (compatible; PaletteScraper/1.0)"
}

palettes = []

print("Fetching popular palettes page...")

r = requests.get(URL, headers=headers, timeout=10)
r.raise_for_status()

soup = BeautifulSoup(r.text, "html.parser")

# Находим все <a>, которые содержат палитру
links = soup.find_all("a", href=True)
for a in links:
    divs = a.find_all("div", class_="palettecolordiv")
    if not divs:
        continue

    colors = []
    for d in divs:
        style = d.get("style", "")
        if "background-color:" in style:
            hex_color = style.split("background-color:")[-1].strip().rstrip(";")
            if hex_color.startswith("#") and len(hex_color) == 7:
                colors.append(hex_color.lower())

    if len(colors) >= 3:
        palettes.append(colors)
        print(f"✔ Found palette: {colors}")

# Сохраняем палитры в JSON
with open("palettes.json", "w", encoding="utf-8") as f:
    json.dump(palettes, f, indent=2)

print(f"\n✅ Collected {len(palettes)} popular palettes and saved to palettes.json")
