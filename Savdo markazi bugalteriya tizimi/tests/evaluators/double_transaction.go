package evaluators

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/darcys22/godbledger/proto/transaction"

	"github.com/darcys22/godbledger/tests/types"

	"google.golang.org/grpc"
)

// DoubleTransaction submits a two transactions to the server and queries the trial balance, it expects both transactions to be in the trial balance and no errors as a response
var DoubleTransaction = types.Evaluator{
	Name:       "Double Transaction",
	Evaluation: doubleTransaction,
}

func doubleTransaction(conns ...*grpc.ClientConn) error {
	client := transaction.NewTransactorClient(conns[0])

	date, _ := time.Parse("2006-01-02", "2011-03-15")
	desc := "Whole Food Market"

	transactionLines := make([]*transaction.LineItem, 2)

	transactionLines[0] = &transaction.LineItem{
		Accountname: "Expenses:Groceries",
		Description: "Groceries",
		Amount:      7500,
		Currency:    "USD",
	}

	transactionLines[1] = &transaction.LineItem{
		Accountname: "Assets:Checking",
		Description: "Groceries",
		Amount:      -7500,
		Currency:    "USD",
	}

	req := &transaction.TransactionRequest{
		Date:        date.Format("2006-01-02"),
		Description: desc,
		Lines:       transactionLines,
	}
	_, err := client.AddTransaction(context.Background(), req)
	if err != nil {
		return err
	}

	_, err = client.AddTransaction(context.Background(), req)
	if err != nil {
		return err
	}

	res, err := client.GetTB(context.Background(), &transaction.TBRequest{Date: time.Now().Format("2006-01-02")})
	if err != nil {
		return err
	}

	// Initialise a variable to check that the trial balance balances
	balance := int64(0)
	// Check to ensure the Trial Balance Matches.
	for i := 0; i < len(res.Lines); i++ {
		balance += res.Lines[i].Amount
		switch res.Lines[i].Accountname {
		case "Assets:Checking":
			if res.Lines[i].Amount != int64(-15000) {
				return errors.New("Trial Balance Checking Account Incorrect")
			}
		case "Expenses:Groceries":
			if res.Lines[i].Amount != int64(15000) {
				return errors.New("Trial Balance Groceries Account Incorrect")
			}
		default:
			return fmt.Errorf("Unknown Account %s", res.Lines[i].Accountname)
		}
	}

	if balance != int64(0) {
		return errors.New("Trial Balance does not balance")
	}

	return nil
}
