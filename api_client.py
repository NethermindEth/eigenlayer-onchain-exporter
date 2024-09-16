import requests

def fetch_data_from_api(url):
    """
    Fetch data from the given API URL.
    
    Args:
    url (str): The URL of the API endpoint.

    Returns:
    dict: The JSON response data.
    """
    response = requests.get(url)
    if response.status_code == 200:
        return response.json()
    else:
        raise Exception(f"Failed to fetch data from API. Status code: {response.status_code}")

# # Test the function with the provided URL
# test_url = "https://blobs-goerli.eigenda.xyz/api/trpc/blobs.getBlobs"
# try:
#     data = fetch_data_from_api(test_url)
#     print("Data fetched successfully.")
# except Exception as e:
#     print(str(e))
