FROM node:20-alpine as build
WORKDIR /app
ARG REACT_APP_API_URL=http://backend:8080
ENV REACT_APP_API_URL=${REACT_APP_API_URL}
COPY package.json ./
RUN npm install
COPY public ./public
COPY src ./src
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
