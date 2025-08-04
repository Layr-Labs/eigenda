#!/bin/bash

# Fix permissions on mounted directories
# The host mounts will override directory ownership, so we need to fix it at runtime
if [ -d "/mnt/data" ]; then
    chown testuser:testuser /mnt/data
    chmod 755 /mnt/data
fi

if [ -d "/mnt/test" ]; then
    chown testuser:testuser /mnt/test
    chmod 755 /mnt/test
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