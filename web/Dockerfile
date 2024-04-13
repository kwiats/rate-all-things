FROM node:18-alpine

WORKDIR /app

COPY package.json .

RUN npm install

COPY . .


EXPOSE 8082
CMD [ "npm", "run", "dev" ]