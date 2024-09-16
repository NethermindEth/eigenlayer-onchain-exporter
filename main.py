from prometheus_client import start_http_server
from api_client import fetch_data_from_api
from metrics import create_metrics, update_metrics
import time
import os

# URL of the API endpoint
API_URL = os.environ.get('API_URL', "https://blobs-goerli.eigenda.xyz/api/trpc/blobs.getBlobs")

# How often to fetch new data and update metrics (in seconds)
FETCH_INTERVAL = int(os.environ.get('FETCH_INTERVAL', 60))

def main():
    """
    Main function to start the server, fetch data periodically, and update metrics.
    """
    # Start up the server to expose the metrics.
    start_http_server(9600)
    print("Prometheus metrics server running on port 9600")

    # Create metrics
    metrics = create_metrics()

    last_timestamp = 0

    while True:
        try:
            # Fetch new data from the API
            data = fetch_data_from_api(API_URL)
            
            # Update the metrics with the new data
            last_timestamp = update_metrics(metrics, data, last_timestamp)

            print("Metrics updated.")
        except Exception as e:
            print(f"Error fetching data or updating metrics: {e}")

        # Wait for the next fetch interval
        time.sleep(FETCH_INTERVAL)

if __name__ == "__main__":
    main()
