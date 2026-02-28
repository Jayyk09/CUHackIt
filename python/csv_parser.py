import pandas as pd
import os
import csv

src = os.path.expanduser("~/Downloads/en.openfoodfacts.org.products.csv")
out = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered.csv")

keep = [
  "product_name",
  "environmental_score_score","environmental_score_grade",
  "nutriscore_score","nutriscore_grade",
  "labels","labels_tags","labels_en",
  "allergens","allergens_en","traces","traces_en",
  "carbon-footprint_100g",
  "image_url","image_small_url",
  "completeness",
  "countries_tags","countries_en",
]

if os.path.exists(out):
  os.remove(out)

chunksize = 200_000
first = True

for chunk in pd.read_csv(
    src,
    sep="\t",
    engine="c",
    quoting=csv.QUOTE_NONE,   # <-- ignore quotes completely
    on_bad_lines="skip",      # <-- skip truly broken rows
    low_memory=False,
    chunksize=chunksize,
):
    keep_existing = [c for c in keep if c in chunk.columns]
    chunk = chunk[keep_existing]
    chunk.to_csv(out, mode="w" if first else "a", index=False, header=first)
    first = False

print("done:", out)
