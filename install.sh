#!/bin/bash

# Set the destination directory
DEST_DIR="/usr/local/echo-server"

# Create the destination directory if it doesn't exist
if [ ! -d "$DEST_DIR" ]; then
  sudo mkdir -p "$DEST_DIR"
fi

# Copy all files from the current directory to the destination directory
sudo cp -r ./* "$DEST_DIR"

# Install the echo-server.service to systemd
SERVICE_FILE="echo-server.service"
SYSTEMD_DIR="/etc/systemd/system"

if [ -f "$SERVICE_FILE" ]; then
  sudo cp "$SERVICE_FILE" "$SYSTEMD_DIR"
  sudo systemctl daemon-reload
  sudo systemctl enable echo-server.service
  sudo systemctl start echo-server.service

  # Check the status of the service
  SERVICE_STATUS=$(sudo systemctl is-active echo-server.service)
  SERVICE_ENABLED=$(sudo systemctl is-enabled echo-server.service)

  echo "Installation complete."
  echo "Files have been installed to: $DEST_DIR"
  echo "To start the service: sudo systemctl start echo-server.service"
  echo "To stop the service: sudo systemctl stop echo-server.service"
  echo "Current service status: $SERVICE_STATUS"
  echo "Service enabled at startup: $SERVICE_ENABLED"

else
  echo "Service file $SERVICE_FILE not found in the current directory."
  exit 1
fi

