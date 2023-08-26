# ==================
#   Builder image
# ==================

FROM node:19-alpine3.16 as builder

# Install git and bash
RUN apk update \
    && apk upgrade \
    && apk add --no-cache bash git openssh

# Create app directory
RUN mkdir -p /app
WORKDIR /app

# Install dependencies
COPY package.json /app
RUN npm install

# Copy source files and Build
COPY . /app
RUN npm run build


# ==================
#  Production image
# ==================

FROM nginx:stable-alpine as production

# Copy nginx config
COPY ./nginx/nginx.conf /etc/nginx/nginx.conf

# Clean nginx then Copy build files
RUN rm -rf /usr/share/nginx/html/*
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port 8080
EXPOSE 8080

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
