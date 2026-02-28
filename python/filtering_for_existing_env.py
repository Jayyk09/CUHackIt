import pandas as pd
import os

src = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered_en_countries.csv")
out = os.path.expanduser("~/Work/Projects/CUHackIt/python/filtered_env_real.csv")

df = pd.read_csv(src)

print("Before:", len(df))

def norm(series: pd.Series) -> pd.Series:
    # normalize: NaN -> "", trim, lowercase
    return series.fillna("").astype(str).str.strip().str.lower()

score = norm(df["environmental_score_score"])
grade = norm(df["environmental_score_grade"])

# Treat these as empty/invalid
INVALID = {"", "nan", "none", "null", "unknown"}

# score is valid if it's not invalid (could be numeric-like string)
score_valid = ~score.isin(INVALID)

# grade is valid if it's one of the expected grades (OFF uses a-e)
grade_valid = grade.isin({"a", "b", "c", "d", "e"})

# keep if either score or grade is valid
df2 = df[score_valid | grade_valid].copy()

print("After:", len(df2))

df2.to_csv(out, index=False)
print("Wrote:", out)
