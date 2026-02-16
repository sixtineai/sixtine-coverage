FROM node:22-alpine AS build

WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY . .
RUN pnpm build


FROM nginx:alpine AS deploy

RUN rm -rf /usr/share/nginx/html/*
COPY --from=build /app/dist /usr/share/nginx/html
COPY --from=build /app/docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY --from=build /app/docker/nginx.conf /etc/nginx/conf.d/default.conf
RUN chmod -R 755 /usr/share/nginx/html

# Create mount points for coverage data and source code
RUN mkdir -p /data/coverage /data/code

EXPOSE 80

# Mount points:
#   /data/coverage/cover.out  - the Go coverage profile
#   /data/code/               - the Go source code root
CMD ["/entrypoint.sh"]
