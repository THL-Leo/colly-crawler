import pandas as pd

# Read the CSV file
df = pd.read_csv('fetch_latimes.csv')

# Count the total number of duplicates
duplicate_count = df.duplicated(subset=['URL']).sum()
print(f"Total number of duplicate URLs: {duplicate_count}")

# Count occurrences of each URL
url_counts = df['URL'].value_counts()

# Filter for URLs that appear more than once
duplicate_urls = url_counts[url_counts > 1]

print("\nUnique Duplicate URLs and their counts:")
for url, count in duplicate_urls.items():
    print(f"{url}: {count} times")

# Optional: Save the results to a CSV file
duplicate_urls_df = pd.DataFrame({'URL': duplicate_urls.index, 'Count': duplicate_urls.values})