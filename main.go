package main

import (
	"context"

	"github.com/Nine9full/workshop-deployment/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	r := gin.Default()

	r.GET("/todos", handler.GetTodos(nil))
	r.GET("/todos/:id", handler.GetTodoByID(nil))
	r.POST("/todos", handler.CreateTodo(nil))
	r.DELETE("/todos/:id", handler.DeleteTodo(nil))

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
