# First build the Golang build.
FROM golang:1.15.2-alpine as gobuilder
RUN mkdir -p /app/backend
WORKDIR /app/backend
COPY backend ./
COPY go.mod /app
COPY go.sum /app
RUN go build -o main main.go

# Now we build the js
FROM node:10.22.1-alpine3.11 as jsbuilder
RUN mkdir -p /app/frontend
WORKDIR /app/frontend
COPY frontend ./
RUN npm install
RUN npm run-script build

# Now start the actual Docker dontainer.
FROM alpine:3.13.0
RUN mkdir -p /app/conf
WORKDIR /app
COPY --from=gobuilder /app/backend/main /app/main
COPY --from=jsbuilder /app/frontend/build /app/public
RUN ls -lR /app
EXPOSE 8080
CMD ["./main", "--hidetreelog","--scrape-interval=6000"]
