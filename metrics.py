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
        'requested_at': Gauge('requested_at', 'Time when the blob was requested'),
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
        metrics['blob_index'].set(blob.get('blob_index', 0))
        metrics['reference_block_number'].set(blob.get('reference_block_number', 0))
        metrics['batch_id'].set(blob.get('batch_id', 0))
        metrics['confirmation_block_number'].set(blob.get('confirmation_block_number', 0))
        metrics['requested_at'].set(blob.get('requested_at', 0)
        # Update additional metrics here

