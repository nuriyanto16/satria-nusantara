# Stage 1: Build the Flutter Web App
FROM debian:stable-slim AS build-env

# Install system dependencies
RUN apt-get update && apt-get install -y \
  curl \
  git \
  unzip \
  xz-utils \
  zip \
  libglu1-mesa \
  && rm -rf /var/lib/apt/lists/*

# Clone Flutter stable SDK
RUN git clone https://github.com/flutter/flutter.git -b stable /usr/local/flutter
ENV PATH="/usr/local/flutter/bin:/usr/local/flutter/bin/cache/dart-sdk/bin:${PATH}"

# Accept Android licenses or perform flutter prep
RUN flutter doctor

# Enable web platform
RUN flutter config --enable-web

# Set working directory and copy source code
WORKDIR /app
COPY . .

# Initialize/re-create the platform template files for Web
RUN flutter create --platforms=web .

# Build Flutter Web application
RUN flutter build web --release

# Stage 2: Serve using Nginx
FROM nginx:alpine
COPY --from=build-env /app/build/web /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
