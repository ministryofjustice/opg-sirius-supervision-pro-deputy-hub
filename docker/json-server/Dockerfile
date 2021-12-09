## build json-server for mock API
FROM node:16.13.0-alpine3.14 as json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install

ENTRYPOINT ["yarn", "serve"]