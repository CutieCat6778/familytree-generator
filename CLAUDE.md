# 🌳 Family Tree Generator Based on Real-Life Data

## Project Overview

The goal of this project is to build a **Family Tree Generator** that creates realistic, synthetic family trees using **real-world demographic, economic, and social data**.  
The generated families should statistically reflect real-life conditions such as life expectancy, birth rates, migration, unemployment, education expenditure, and population trends across different countries.

This project uses publicly available datasets (CSV + JSON) to ensure realism while maintaining privacy by generating **fully synthetic individuals**.

---

## Objectives

- Generate realistic multi-generation family trees
- Base all demographic behavior on real-world statistics
- Support country-specific variations
- Ensure reproducibility and configurability
- Output results in structured CSV (and optional JSON)

---

## Data Sources

The project uses the following datasets:

### Core Demographic Data

- `population.csv`
- `birth_rate.csv`
- `death_rate.csv`
- `life_exp_at_birth.csv`
- `migration_rate.csv`
- `imr.csv` (Infant Mortality Rate)

### Economic & Social Indicators

- `gdp_per_cap.csv`
- `unemployment_rate.csv`
- `youth_unemployment_rate.csv`
- `education_expenditure.csv`
- `labor_force.csv`
- `inflation_rate.csv`

### Health & Lifestyle

- `alcohol.csv`
- `tobacco_use.csv`
- `underweight_u5.csv`

### Identity Data

- `forenames.csv`
- `surnames.csv`
- `countries-code.json`

# How it works?

1. **Data Ingestion**: Load and preprocess the datasets to extract relevant statistics.
2. **Family Tree Generation**:
   - Start with a root individual.
   - Recursively generate parents, siblings, spouses, and children based on statistical probabilities derived from the datasets.
3. **Country-Specific Adjustments**: Modify demographic behaviors based on the selected country.
4. **Output**: Save the generated family tree in CSV format, with an option for
   JSON output.
5. **Configuration**: Allow users to specify parameters such as the number of generations, country, and output format.
6. **Reproducibility**: Use a random seed to ensure that the same input parameters yield the same family tree.
7. **Visualization (Optional)**: Provide a graphical representation of the family tree.

---
