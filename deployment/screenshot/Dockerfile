FROM zenika/alpine-chrome:with-puppeteer

COPY --chown=chrome ./ /app
RUN ls -lia /app

WORKDIR /app/deployment/screenshot
RUN yarn install

ENTRYPOINT [ "/app/deployment/screenshot/node_modules/.bin/ts-node", "screenshot.ts" ]
