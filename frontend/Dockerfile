FROM dhi.io/node:25-alpine3.23-sfw-ent-dev as base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN corepack enable

FROM base AS prod
COPY pnpm-lock.yaml /app
WORKDIR /app
RUN pnpm fetch --prod
COPY . /app
RUN pnpm run build

FROM dhi.io/nginx:1.29.4-debian13 as runner
COPY nginx.conf /etc/nginx/nginx.conf
COPY --chown=nginx:nginx --from=prod /app/dist /usr/share/nginx/html
USER nginx
EXPOSE 8080

ENTRYPOINT [ "nginx", "-c", "/etc/nginx/nginx.conf" ]
CMD [ "-g", "daemon off;" ]