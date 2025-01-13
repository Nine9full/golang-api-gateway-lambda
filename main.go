package main

import (
	"context"

	config "github.com/Nine9full/workshop-deployment/db"
	"github.com/Nine9full/workshop-deployment/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	r := gin.Default()

	configs := config.ReadConfig()

	conn, err := config.Connect(configs)
	if err != nil {
		panic(err)
	}

	// , err := conn.DB()
	// if err != nil {
	// 	panic(err)
	// }

	// defer db.Close()

	r.GET("/todos", handler.GetTodos(conn))
	r.GET("/todos/:id", handler.GetTodoByID(conn))
	r.POST("/todos", handler.CreateTodo(conn))
	r.DELETE("/todos/:id", handler.DeleteTodo(conn))

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
