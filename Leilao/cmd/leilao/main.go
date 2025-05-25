package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rafael0502/curso-go/Leilao/configuration/database/mongodb"
	"github.com/rafael0502/curso-go/Leilao/internal/infra/api/web/controller/bid_controller"
	"github.com/rafael0502/curso-go/Leilao/internal/infra/api/web/controller/leilao_controller"
	"github.com/rafael0502/curso-go/Leilao/internal/infra/api/web/controller/user_controller"
	"github.com/rafael0502/curso-go/Leilao/internal/infra/database/bid"
	"github.com/rafael0502/curso-go/Leilao/internal/infra/database/leilao"
	"github.com/rafael0502/curso-go/Leilao/internal/infra/database/user"
	"github.com/rafael0502/curso-go/Leilao/internal/usecase/bid_usecase"
	"github.com/rafael0502/curso-go/Leilao/internal/usecase/leilao_usecase"
	"github.com/rafael0502/curso-go/Leilao/internal/usecase/user_usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/leilao/.env"); err != nil {
		panic("Error loading .env file")
	}

	database, err := mongodb.NewMongoDBConnection(ctx)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := gin.Default()

	userController, bidController, leilaoController := initDependencies(database)

	router.POST("/leiloes", leilaoController.CreateLeilao)
	router.GET("/leiloes", leilaoController.FindLeiloes)
	router.GET("/leiloes/:leilaoId", leilaoController.FindLeilaoById)
	router.GET("/leilao/winner/:leilaoId", leilaoController.FindWinningBidByLeilaoid)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:leilaoId", bidController.FindBidByLeilaoId)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	leilaoController *leilao_controller.LeilaoController) {

	leilaoRepository := leilao.NewLeilaoRepository(database)
	bidRepository := bid.NewBidRepository(database, leilaoRepository)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(user_usecase.NewUserUseCase(userRepository))
	leilaoController = leilao_controller.NewLeilaoController(leilao_usecase.NewLeilaoUseCase(leilaoRepository, bidRepository))
	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}
