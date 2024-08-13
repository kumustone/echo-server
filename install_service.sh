#!/bin/bash

# Define the service file path
SERVICE_FILE="echo-server.service"

# Check if the service file exists
if [ ! -f "$SERVICE_FILE" ]; then
    echo "Service file $SERVICE_FILE not found!"
    exit 1
fi

# Copy the service file to the systemd directory
sudo cp "$SERVICE_FILE" /etc/systemd/system/

# Reload systemd manager configuration
sudo systemctl daemon-reload

# Enable the service to start on boot
sudo systemctl enable echo-server.service

# Start the service
sudo systemctl start echo-server.service

# Check the status of the service
sudo systemctl status echo-server.service

