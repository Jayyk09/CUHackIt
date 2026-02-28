import os
import pandas as pd

# Input and output paths
inp = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered_env_and_nutri_real.csv")
out = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered_env_and_nutri_real_no_empty_scores.csv")

# Read CSV
df = pd.read_csv(inp)

# Convert score columns to numeric (forces invalid/blank to NaN)
df["environmental_score_score"] = pd.to_numeric(
    df["environmental_score_score"], errors="coerce"
)

df["nutriscore_score"] = pd.to_numeric(
    df["nutriscore_score"], errors="coerce"
)

# Drop rows where either score is missing
df_filtered = df.dropna(subset=[
    "environmental_score_score",
    "nutriscore_score"
])

# Save result
df_filtered.to_csv(out, index=False)

print("Original rows:", len(df))
print("Remaining rows:", len(df_filtered))
print("Saved to:", out)
