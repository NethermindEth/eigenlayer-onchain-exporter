# eigenda-blob-scraper

A Python application that scrapes the Eigenda website for the latest blob data and exposes it via Prometheus metrics.

## Project Description

The application gets the blob data in a json form from an EigenDA public endpoint and transforms it into Prometheus metrics every `FETCH_INTERVAL` seconds. The metrics are then exposed via a Prometheus server at port 9600.

The application is deployed via Docker image to DockerHub. The image can be found here: https://hub.docker.com/repository/docker/nethermind/eigenda-blob-scraper.

## Environment Variables

The application uses the following environment variables:

| Variable | Description | Default Value |
| --- | --- | --- |
| API_URL | Endpoint to get the blob json data from | https://blobs-goerli.eigenda.xyz/api/trpc/blobs.getBlobs |
| FETCH_INTERVAL | Fetch interval in seconds | 60 |

## Docker image

`nethermind/eigenda-blob-scraper:latest`

## How to Run the Application

### Docker

The application can be run via Docker. The Docker image is available on DockerHub. To run the application, execute the following command:

```bash
docker run -p 9600:9600 nethermind/eigenda-blob-scraper:latest
```

or alternatively, if you want to modify the env variables:

```bash
docker run -p 9600:9600 -e API_URL=https://blobs-goerli.eigenda.xyz/api/trpc/blobs.getBlobs -e FETCH_INTERVAL=60 nethermind/eigenda-blob-scraper:latest
```

### CLI

The application can also be run via the CLI. To run the application, execute the following command:

```bash
python3 main.py
```

or alternatively, if you want to modify the env variables:

```bash
API_URL=https://blobs-goerli.eigenda.xyz/api/trpc/blobs.getBlobs FETCH_INTERVAL=60 python3 main.py
```

## How to Update the Code and Trigger CI/CD

The repository uses GitHub Actions for Continuous Integration and Continuous Deployment. The workflow is defined in the `.github/workflows/publish-docker.yml` file. To update the code and trigger the CI/CD pipeline, follow these steps:

1. Make your code changes locally.
2. Commit and push the changes to the main branch.
3. GitHub Actions will automatically build and push the docker image to DockerHub.

The CI/CD pipeline is triggered on every push to the main branch, but can also be triggered manually by clicking on the "Actions" tab in the repository and selecting the "Publish Docker image" workflow. Then click on the "Run workflow" button.