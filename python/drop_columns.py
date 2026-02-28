import os
import pandas as pd

# Input and output paths
inp = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered_env_and_nutri_real.csv")
out = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered_env_and_nutri_real_cleaned.csv")

# Read CSV
df = pd.read_csv(inp)

# Columns to remove
cols_to_drop = [
    "environmental_score_grade",
    "nutriscore_grade",
    "carbon-footprint_100g",
    "completeness",
    "countries_en"
]

# Drop columns safely (won’t crash if one is missing)
df = df.drop(columns=cols_to_drop, errors="ignore")

# Save cleaned CSV
df.to_csv(out, index=False)

print(f"Cleaned file saved to: {out}")
