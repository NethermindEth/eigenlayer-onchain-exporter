# Use an official Python runtime as a parent image
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /usr/src/app

# Copy the current directory contents into the container at /usr/src/app
COPY . .

# Install any needed packages specified in requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Make port 9600 available to the world outside this container
EXPOSE 9600 

# Define arguments for build-time configuration
ARG BUILD_FETCH_INTERVAL=60
ARG BUILD_API_URL=https://blobs-goerli.eigenda.xyz/api/trpc/blobs.getBlobs

# Define environment variable
ENV FETCH_INTERVAL=${BUILD_FETCH_INTERVAL}
ENV API_URL=${BUILD_API_URL}

# Run main.py when the container launches
CMD ["python", "./main.py"]
