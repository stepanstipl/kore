FROM node:12-alpine3.10

ARG version=v0.0.0-dev
ENV NODE_ENV production
ENV UI_VERSION ${version}

WORKDIR /app
COPY package.json .
COPY package-lock.json .
RUN npm install
COPY . .
RUN npm run build

CMD [ "npm", "start" ]
