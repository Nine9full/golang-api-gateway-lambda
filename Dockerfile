FROM golang:1.23.4 as build

WORKDIR /app

COPY . /app

RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o lambda-golang main.go


FROM public.ecr.aws/lambda/go:latest

COPY --from=build /app/lambda-golang ${LAMBDA_TASK_ROOT}/

CMD ["lambda-golang"]