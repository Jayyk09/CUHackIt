import pandas as pd
import os
import csv

src = os.path.expanduser("~/Downloads/en.openfoodfacts.org.products.csv")
out = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered_en_countries.csv")

# Keep list after:
# - removing macro columns
# - keeping only _en versions for fields that have them
# - still keeping environmental + nutri + images + completeness + countries_en
keep = [
  "product_name",
  "environmental_score_score","environmental_score_grade",
  "nutriscore_score","nutriscore_grade",
  "labels_en",
  "allergens_en","traces_en",
  "carbon-footprint_100g",
  "image_url","image_small_url",
  "completeness",
  "countries_en",
]

# Define "English-speaking countries" you want to keep.
# (Tweak this list to match your definition.)
ENGLISH_COUNTRIES = {
  "United States",
  "United Kingdom",
  "Canada",
  "Ireland",
  "Australia",
  "New Zealand",
  "South Africa",
  # optional common additions:
  "Singapore",
  "India",
}

if os.path.exists(out):
    os.remove(out)

chunksize = 200_000
first = True

for chunk in pd.read_csv(
    src,
    sep="\t",
    engine="c",
    quoting=csv.QUOTE_NONE,
    on_bad_lines="skip",
    low_memory=False,
    chunksize=chunksize,
):
    # keep only columns that exist (OFF exports vary)
    keep_existing = [c for c in keep if c in chunk.columns]
    chunk = chunk[keep_existing]

    # Filter to English-speaking countries.
    # countries_en can contain multiple countries separated by commas.
    if "countries_en" in chunk.columns:
        ce = chunk["countries_en"].fillna("").astype(str)

        # keep row if ANY country in the list appears in countries_en
        pattern = "|".join(map(lambda s: s.replace(" ", r"\s+"), ENGLISH_COUNTRIES))
        chunk = chunk[ce.str.contains(pattern, case=False, na=False, regex=True)]

    chunk.to_csv(out, mode="w" if first else "a", index=False, header=first)
    first = False

print("done:", out)
