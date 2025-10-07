import json
import csv

def main():
    # --- 1. Open input JSON file ---
    try:
        with open("/developer/familytree/main/data/fornames.json", "r", encoding="utf-8") as input_file:
            j = json.load(input_file)
    except FileNotFoundError:
        print("❌ Could not open input JSON file")
        return
    except json.JSONDecodeError as e:
        print(f"❌ JSON parse error: {e}")
        return

    # --- 2. Open output CSV file ---
    try:
        with open("output.csv", "w", newline="", encoding="utf-8") as output_file:
            writer = csv.writer(output_file, quoting=csv.QUOTE_ALL)

            # --- 3. Write CSV header ---
            writer.writerow(["country", "region", "gender", "rank", "name"])

            # --- 4. Loop through countries ---
            for country_code, regions in j.items():
                if not isinstance(regions, list):
                    continue  # Skip if not list

                for region_block in regions:
                    if not isinstance(region_block, dict):
                        continue  # Skip if not dict

                    region = region_block.get("region", "Unknown")

                    names = region_block.get("names")
                    if not isinstance(names, list):
                        continue

                    for name_entry in names:
                        if not isinstance(name_entry, dict):
                            continue

                        gender = name_entry.get("gender", "")
                        rank = name_entry.get("rank", 0)
                        name = ""

                        romanized = name_entry.get("romanized")
                        if isinstance(romanized, list) and romanized and isinstance(romanized[0], str):
                            name = romanized[0]
                        else:
                            localized = name_entry.get("localized")
                            if isinstance(localized, list) and localized and isinstance(localized[0], str):
                                name = localized[0]

                        if not name:
                            continue  # Skip if no valid name found

                        # --- 5. Write row ---
                        writer.writerow([country_code, region, gender, rank, name])

    except IOError:
        print("❌ Could not create output.csv")
        return

    print("✅ Conversion complete: output.csv created.")

if __name__ == "__main__":
    main()
