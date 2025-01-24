# First stage: Build stage
FROM node:18 AS builder

WORKDIR /app

# Copy source code to the working directory
COPY . .

# Install dependencies
RUN npm install

# Build the application, need .env.prod file
RUN npm run build:prod

# Second stage: Serve stage
FROM nginx:alpine

# Copy the build output to the Nginx HTML directory
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]
