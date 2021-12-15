FROM cypress/included:5.2.0

WORKDIR /root

RUN npm install "cypress-failed-log"

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://pro-deputy-hub:8888

COPY cypress.json .
COPY cypress cypress