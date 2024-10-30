import pandas as pd
from http import HTTPStatus

# Load CSV files
fetch_df = pd.read_csv('fetch_latimes.csv')
urls_df = pd.read_csv('urls_latimes.csv')
visit_df = pd.read_csv('visit_latimes.csv')

# Calculate fetch statistics
fetch_attempted = len(fetch_df)
fetch_succeeded = fetch_df['Status'].value_counts().get(200, 0)
fetch_failed = fetch_attempted - fetch_succeeded

# Calculate outgoing URLs
total_urls_extracted = urls_df['URL'].count()
unique_urls_extracted = urls_df['URL'].nunique()
unique_within_news_site = urls_df[urls_df['URL'].str.contains('latimes.com')]['URL'].nunique()
unique_outside_news_site = unique_urls_extracted - unique_within_news_site

# Calculate status codes
status_counts = fetch_df['Status'].value_counts()

# Calculate file sizes
size_bins = [0, 1024, 10240, 102400, 1048576, float('inf')]
size_labels = ['< 1KB', '1KB ~ <10KB', '10KB ~ <100KB', '100KB ~ <1MB', '>= 1MB']
visit_df['Size_Category'] = pd.cut(visit_df['Size(Bytes)'], bins=size_bins, labels=size_labels)
size_counts = visit_df['Size_Category'].value_counts()

# Calculate content types
content_type_counts = visit_df['Content-Type'].value_counts()

# Write to output file
with open('CrawlReport_foxnews.txt', 'w') as f:
    f.write("Name: \n")
    f.write("USC ID: \n")
    f.write("News site crawled: latimes\n")
    f.write("Number of threads: 16\n\n")

    f.write("Fetch Statistics\n")
    f.write("================\n")
    f.write(f"# fetches attempted: {fetch_attempted}\n")
    f.write(f"# fetches succeeded: {fetch_succeeded}\n")
    f.write(f"# fetches failed or aborted: {fetch_failed}\n\n")

    f.write("Outgoing URLs:\n")
    f.write("================\n")
    f.write(f"Total URLs extracted: {total_urls_extracted}\n")
    f.write(f"# unique URLs extracted: {unique_urls_extracted}\n")
    f.write(f"# unique URLs within News Site: {unique_within_news_site}\n")
    f.write(f"# unique URLs outside News Site: {unique_outside_news_site}\n\n")

    f.write("Status Codes:\n")
    f.write("================\n")
    statusDict = {}
    for status in HTTPStatus:
        statusDict[status.value] = status.name
    for code, count in status_counts.items():
        description = statusDict.get(code, "Unknown Error")
        if code == 0:
            f.write(f"Error: {count}\n")
        else:
            f.write(f"{code} {description}: {count}\n")
    f.write("\n")

    f.write("File Sizes:\n")
    f.write("================\n")
    for label in size_labels:
        count = size_counts.get(label, 0)
        f.write(f"{label}: {count}\n")
    f.write("\n")

    f.write("Content Types:\n")
    f.write("================\n")
    for content_type, count in content_type_counts.items():
        f.write(f"{content_type}: {count}\n")