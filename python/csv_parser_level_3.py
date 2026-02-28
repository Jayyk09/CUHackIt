import pandas as pd
import os

src = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered_env_real.csv")
out = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered_env_and_nutri_real.csv")

df = pd.read_csv(src)

print("Before:", len(df))

def norm(series: pd.Series) -> pd.Series:
    return series.fillna("").astype(str).str.strip().str.lower()

score = norm(df["nutriscore_score"])
grade = norm(df["nutriscore_grade"])

INVALID = {"", "nan", "none", "null", "unknown"}

score_valid = ~score.isin(INVALID)          # numeric-like is fine, just not invalid
grade_valid = grade.isin({"a","b","c","d","e"})

df2 = df[score_valid | grade_valid].copy()

print("After:", len(df2))

df2.to_csv(out, index=False)
print("Wrote:", out)
