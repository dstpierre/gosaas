package queue

import (
	"fmt"
	"os"
	"time"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/invoice"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_KEY")
}

type Billing struct{}

func (b *Billing) Run(qt QueueTask) error {
	id, ok := qt.Data.(string)
	if !ok {
		return fmt.Errorf("the data should be a stripe customer ID")
	}

	// we delay execution for 2 hours to let add/remove
	// operations in between creating the invoice
	// since we're on a go routine we can use a time.Sleep
	time.Sleep(2 * time.Hour)

	p := &stripe.InvoiceParams{Customer: id}
	_, err := invoice.New(p)
	return err
}
