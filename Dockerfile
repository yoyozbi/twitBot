FROM node:lts-alpine3.15

COPY . /app
WORKDIR /app

RUN npm install

CMD ["node", "src/index.js"]