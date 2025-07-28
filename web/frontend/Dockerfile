FROM node:20-alpine as build
WORKDIR /app
ARG REACT_APP_API_URL=http://backend:8080
ENV REACT_APP_API_URL=${REACT_APP_API_URL}
COPY package.json ./
RUN npm install
COPY public ./public
COPY src ./src
RUN npm run build

FROM node:20-alpine
RUN npm install -g serve
WORKDIR /app
COPY --from=build /app/build ./build
EXPOSE 80
CMD ["serve", "-s", "build", "-l", "80"]
