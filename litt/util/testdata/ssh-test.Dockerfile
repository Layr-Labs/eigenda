FROM ubuntu:22.04

# Install required packages
RUN apt-get update && apt-get install -y \
    openssh-server \
    rsync \
    && rm -rf /var/lib/apt/lists/*

# Create test group with GID 1337 and user with UID 1337
RUN groupadd -g 1337 testgroup
RUN useradd -m -s /bin/bash -u 1337 -g 1337 testuser

# Setup SSH
RUN mkdir /var/run/sshd
RUN mkdir -p /home/testuser/.ssh

# Configure SSH daemon
RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
RUN sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config

# Set proper permissions
RUN chown -R testuser:testgroup /home/testuser/.ssh
RUN chmod 700 /home/testuser/.ssh

# Create mount directories and set ownership
RUN mkdir -p /mnt/data
RUN chown testuser:testgroup /mnt/data

# Copy startup script with self-destruct mechanism
COPY start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 22
CMD ["/start.sh"]