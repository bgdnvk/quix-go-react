# Use a Node.js image for building the frontend
FROM node:18 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and yarn.lock for dependency installation
COPY package.json yarn.lock ./

# Install dependencies
RUN yarn install

# Copy the rest of the frontend code
COPY . .

# Build the frontend for production
RUN yarn build

# Serve the static files using an nginx server
FROM nginx:alpine

# Copy the built frontend files to the nginx html folder
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose the Nginx port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
