import csv
from collections import defaultdict

def get_countries_from_csv(filename):
    countries = set()
    with open(filename, 'r', encoding='utf-8') as f:
        reader = csv.reader(f)
        header = next(reader)  # Skip header
        for row in reader:
            if row:
                countries.add(row[0])
    return countries

def filter_csv(input_filename, output_filename, common_countries):
    with open(input_filename, 'r', encoding='utf-8') as infile, \
         open(output_filename, 'w', encoding='utf-8', newline='') as outfile:
        reader = csv.reader(infile)
        writer = csv.writer(outfile, quoting=csv.QUOTE_MINIMAL)
        
        header = next(reader)
        writer.writerow(header)
        
        for row in reader:
            if row and row[0] in common_countries:
                writer.writerow(row)

def main():
    forenames_file = '/developer/familytree/main/data/forenames.csv'
    surnames_file = '/developer/familytree/main/data/surnames.csv'
    
    forenames_countries = get_countries_from_csv(forenames_file)
    surnames_countries = get_countries_from_csv(surnames_file)
    
    common_countries = forenames_countries & surnames_countries
    
    if not common_countries:
        print("No common countries found.")
        return
    
    # Overwrite originals - use temp files to avoid issues
    filter_csv(forenames_file, 'temp_forenames.csv', common_countries)
    filter_csv(surnames_file, 'temp_surnames.csv', common_countries)
    
    # Rename temp to original
    import os
    os.replace('temp_forenames.csv', forenames_file)
    os.replace('temp_surnames.csv', surnames_file)
    
    print("Files updated with only common countries.")

if __name__ == "__main__":
    main()
