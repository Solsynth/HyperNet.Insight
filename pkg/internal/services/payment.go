package services

import (
	"context"
	"fmt"
	"time"

	"git.solsynth.dev/hypernet/insight/pkg/internal/gap"
	wproto "git.solsynth.dev/hypernet/wallet/pkg/proto"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// PlaceOrder create a transaction if needed for user
// Pricing is 128 words input cost 10 source points, round up.
func PlaceOrder(user uint, inputLength int) error {
	amount := float64(inputLength+128-1) / 128

	conn, err := gap.Nx.GetClientGrpcConn("wa")
	if err != nil {
		return fmt.Errorf("unable to connect wallet: %v", err)
	}

	wc := wproto.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := wc.MakeTransactionWithAccount(ctx, &wproto.MakeTransactionWithAccountRequest{
		PayerAccountId: lo.ToPtr(uint64(user)),
		Amount:         amount,
		Remark:         "Insight Thinking Fee",
	})
	if err != nil {
		return err
	}

	log.Info().
		Uint64("transaction", resp.Id).Float64("amount", amount).
		Msg("Order placed for charge insight thinking fee...")

	return nil
}

// MakeRefund to user who got error in generating insight
func MakeRefund(user uint, inputLength int) error {
	amount := float64(inputLength+128-1) / 128

	conn, err := gap.Nx.GetClientGrpcConn("wa")
	if err != nil {
		return fmt.Errorf("unable to connect wallet: %v", err)
	}

	wc := wproto.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := wc.MakeTransactionWithAccount(ctx, &wproto.MakeTransactionWithAccountRequest{
		PayeeAccountId: lo.ToPtr(uint64(user)),
		Amount:         amount,
		Remark:         "Insight Thinking Failed - Refund",
	})
	if err != nil {
		return err
	}

	log.Info().
		Uint64("transaction", resp.Id).Float64("amount", amount).
		Msg("Refund placed for insight thinking fee...")

	return nil
}
