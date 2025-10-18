package bid_usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/ivandersr/go-auction/config/logger"
	"github.com/ivandersr/go-auction/internal/entity/bid"
	"github.com/ivandersr/go-auction/internal/internal_errors"
)

type BidInputDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05"`
}

type BidUseCase struct {
	BidRepository       bid.BidRepositoryInterface
	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bid.Bid
}

type BidUseCaseInterface interface {
	CreateBid(ctx context.Context, input *BidInputDTO) *internal_errors.InternalError
	FindWinningBidByAuctionId(ctx context.Context, id string) (*BidOutputDTO, *internal_errors.InternalError)
	FindBidsByAuctionId(ctx context.Context, id string) ([]BidOutputDTO, *internal_errors.InternalError)
}

func NewBidUseCase(bidRepository bid.BidRepositoryInterface) *BidUseCase {
	maxBatchSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()

	bu := &BidUseCase{
		BidRepository:       bidRepository,
		timer:               time.NewTimer(maxBatchSizeInterval),
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxBatchSizeInterval,
		bidChannel:          make(chan bid.Bid, maxBatchSize),
	}
	bu.triggerCreateRoutine(context.Background())
	return bu
}

var bidBatch []bid.Bid

func (bu *BidUseCase) CreateBid(ctx context.Context, input *BidInputDTO) *internal_errors.InternalError {
	bidEntity, err := bid.CreateBid(input.UserId, input.AuctionId, input.Amount)
	if err != nil {
		return err
	}

	bu.bidChannel <- *bidEntity
	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}
	return duration
}

func getMaxBatchSize() int {
	maxBatchSize, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}
	return maxBatchSize
}

func (bu *BidUseCase) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)

		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						if err := bu.BidRepository.CreateBids(ctx, bidBatch); err != nil {
							logger.Error("error when trying to create bids", err)
						}
					}
					return
				}

				bidBatch = append(bidBatch, bidEntity)
				if len(bidBatch) >= bu.maxBatchSize {
					if err := bu.BidRepository.CreateBids(ctx, bidBatch); err != nil {
						logger.Error("error when trying to create bids", err)
					}
					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			case <-bu.timer.C:
				if err := bu.BidRepository.CreateBids(ctx, bidBatch); err != nil {
					logger.Error("error when trying to create bids", err)
				}
				bidBatch = nil
				bu.timer.Reset(bu.batchInsertInterval)
			}
		}
	}()
}
