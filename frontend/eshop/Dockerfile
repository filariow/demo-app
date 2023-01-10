FROM registry.access.redhat.com/ubi9/nginx-120:latest as server
LABEL "Maintainer" "Red Hat"
ADD nginx.conf "${NGINX_CONFIGURATION_PATH}"
CMD nginx -g "daemon off;"

FROM node:18.12.1 as builder
WORKDIR /app
COPY package.json .
COPY package-lock.json .
RUN npm install
COPY . .
RUN npm run build

FROM server
COPY --from=builder /app/build .
