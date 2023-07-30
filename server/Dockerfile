# Use the official Nginx image as the base image
FROM nginx:latest

# Copy your wait-for-it.sh script into the image
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh

# Make the script executable
RUN chmod +x /usr/local/bin/wait-for-it.sh

# Copy your Nginx configuration file into the image
COPY nginx.conf /etc/nginx/nginx.conf
