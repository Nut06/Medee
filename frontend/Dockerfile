# FROM dhi.io/node:25.8.1-alpine3.22-dev as base
# FROM dhi.io/node:25.8.1-alpine3.22 as base
FROM dhi.io/node:25.8-debian13-dev as base
# FROM node:25.8-bookworm-slim as base
WORKDIR /app

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY . /app

FROM base AS build
ARG VITE_API_URL=/api
ENV VITE_API_URL=$VITE_API_URL
WORKDIR /app
COPY pnpm-lock.yaml package.json ./
RUN pnpm install --frozen-lockfile
RUN pnpm exec playwright install --with-deps chromium
COPY . .
RUN pnpm run build

FROM build as run-test-stage
RUN pnpm test

FROM dhi.io/nginx:1.29.4-debian13 as runner
COPY nginx.conf /etc/nginx/nginx.conf
COPY --chown=nginx:nginx --from=build /app/dist /usr/share/nginx/html
USER nginx
EXPOSE 8080

ENTRYPOINT [ "nginx", "-c", "/etc/nginx/nginx.conf" ]
CMD [ "-g", "daemon off;" ]