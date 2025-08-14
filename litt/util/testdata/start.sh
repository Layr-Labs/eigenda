#!/bin/bash

# Create container-controlled workspace directories
# Instead of trying to change ownership of mounted directories (which often fails),
# we create subdirectories that the container user fully owns
if [ -d "/mnt/data" ]; then
    # Create a container workspace that testuser fully owns
    mkdir -p /mnt/data/container_workspace/work
    chown -R testuser:testuser /mnt/data/container_workspace
    chmod -R 755 /mnt/data/container_workspace
fi

if [ -d "/mnt/test" ]; then
    # Create a container workspace that testuser fully owns  
    mkdir -p /mnt/test/container_workspace/work
    chown -R testuser:testuser /mnt/test/container_workspace
    chmod -R 755 /mnt/test/container_workspace
fi

# Start SSH daemon in background
/usr/sbin/sshd -D &
SSHD_PID=$!

# Self-destruct after 5 minutes (300 seconds)
(
  sleep 300
  echo "SSH test container self-destructing after 5 minutes..."
  kill $SSHD_PID
  exit 0
) &

# Wait for SSH daemon to finish
wait $SSHD_PID