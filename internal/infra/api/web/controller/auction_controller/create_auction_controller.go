package auction_controller

import (
	"context"
	"net/http"

	"github.com/dpcamargo/distributed-auction-system-go/configuration/rest_err"
	"github.com/dpcamargo/distributed-auction-system-go/internal/infra/api/web/validation"
	"github.com/dpcamargo/distributed-auction-system-go/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
)

type AuctionController struct {
	AuctionUseCase auction_usecase.AuctionUseCaseInterface
}

func NewAuctionController(AuctionUseCase auction_usecase.AuctionUseCaseInterface) *AuctionController {
	return &AuctionController{
		AuctionUseCase: AuctionUseCase,
	}
}

func (u *AuctionController) CreateAuction(c *gin.Context) {
	var auctionInputDTO auction_usecase.AuctionInputDTO

	if err := c.ShouldBindJSON(&auctionInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.AuctionUseCase.CreateAuction(context.Background(), auctionInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
