# Use the official Node.js image as the base image
FROM node:latest

# Set the working directory inside the container
WORKDIR /app/client

# Copy package.json and yarn.lock files to the container
COPY package.json yarn.lock ./

# Install project dependencies
RUN yarn install

# Copy the entire project directory to the container
COPY . .

# Build the production-ready React application
RUN yarn build

# Expose the port on which the application will run
EXPOSE 3000

# Specify the command to start the application when the container starts
CMD ["yarn", "dev"]
