from prometheus_client import Gauge
import time

# Assuming the structure of the data based on the provided example
# We will create a Gauge for each field in the blob data

def create_metrics():
    """
    Create Prometheus Gauges for various fields in the blob data.
    """
    metrics = {
        'blob_index': Gauge('blob_index', 'Index of the blob'),
        'reference_block_number': Gauge('reference_block_number', 'Reference block number'),
        'batch_id': Gauge('batch_id', 'Batch ID'),
        'confirmation_block_number': Gauge('confirmation_block_number', 'Confirmation block number'),
        'time_difference': Gauge('time_difference_hours', 'Difference in hours from the requested_at time to current time')
        # Additional metrics can be added here based on the data fields
    }
    return metrics

def update_metrics(metrics, data):
    """
    Update the Prometheus metrics with the latest data from the API.

    Args:
    metrics (dict): The dictionary of Prometheus Gauges.
    data (dict): The data fetched from the API.
    """
    current_time = time.time()
    for blob in data.get('result', {}).get('data', {}).get('json', {}).get('data', []):
        requested_at = blob.get('requested_at', 0)
        time_difference = (current_time - requested_at) / 3600  # Convert seconds to hours

        metrics['blob_index'].set(blob.get('blob_index', 0))
        metrics['reference_block_number'].set(blob.get('reference_block_number', 0))
        metrics['batch_id'].set(blob.get('batch_id', 0))
        metrics['confirmation_block_number'].set(blob.get('confirmation_block_number', 0))
        metrics['time_difference'].set(time_difference)
        # Update additional metrics here

