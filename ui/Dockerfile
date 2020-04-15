FROM node:12-alpine3.10

ARG version=v0.0.0-dev
ENV NODE_ENV production

WORKDIR /app
COPY package.json .
COPY package-lock.json .
RUN npm install
COPY . .
RUN echo ${version} > ./VERSION
RUN npm run build

CMD [ "npm", "start" ]
