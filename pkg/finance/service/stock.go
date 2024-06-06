package finance_service

import (
	"context"
	"errors"

	v1 "core/api/v1"
	biz "core/pkg/finance/biz"
)

type FinanceService struct {
	v1.UnimplementedFinanceServer

	action *biz.StockQuoteAction
}

func NewFinanceService(action *biz.StockQuoteAction) *FinanceService {
	return &FinanceService{action: action}
}

func (s *FinanceService) GetStockQuote(ctx context.Context, req *v1.GetStockQuoteRequest) (*v1.GetStockQuoteReply, error) {
	if req.Symbol == "" {
		return &v1.GetStockQuoteReply{Quote: nil, Symbol: ""}, errors.New("id not supplied")
	}
	quote, err := s.action.GetStockQuote(ctx, req.Symbol)
	if err != nil {
		return nil, err
	}
	return &v1.GetStockQuoteReply{
		Symbol: quote.Symbol,
		Quote:  quote.Quote,
	}, nil
}

func (s *FinanceService) GetUSASpending(ctx context.Context, req *v1.GetUSASpendingRequest) (*v1.GetUSASpendingReply, error) {
	// TODO: Implement the logic for retrieving USA spending data
	return &v1.GetUSASpendingReply{}, nil
}

func (s *FinanceService) GetSenateLobbying(ctx context.Context, req *v1.GetSenateLobbyingRequest) (*v1.GetSenateLobbyingReply, error) {
	// TODO: Implement the logic for retrieving Senate lobbying data
	return &v1.GetSenateLobbyingReply{}, nil
}

func (s *FinanceService) WatchTrades(req *v1.SyncTradesRequest, stream v1.Finance_WatchTradesServer) error {
	if req.Symbol == "" {
		return nil
	}

	ctx := stream.Context()
	err := s.action.WatchTrades(req.Symbol, stream)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}
