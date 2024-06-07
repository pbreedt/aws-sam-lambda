FROM public.ecr.aws/docker/library/golang:1.19 as build-image
WORKDIR /src
COPY go.mod go.sum main.go ./
# required for alpine
# ENV CGO_ENABLED=0
RUN go build -o lambda-handler
FROM public.ecr.aws/lambda/provided:al2023
# FROM alpine:3.9
COPY --from=build-image /src/lambda-handler .
ENTRYPOINT ./lambda-handler
