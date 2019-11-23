package gosaas

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/model"
	"github.com/dstpierre/gosaas/queue"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/invoice"
	"github.com/stripe/stripe-go/sub"
)

// SetStripeKey sets the Stripe Key.
func SetStripeKey(key string) {
	stripe.Key = key
}

// Billing handles everything related to the billing requests
type Billing struct {
	DB *data.DB
}

// BillingOverview represents if an account is a paid customer or not
type BillingOverview struct {
	Account        *model.Account    `json:"account"`
	StripeID       string            `json:"stripeId"`
	Plan           string            `json:"plan"`
	IsYearly       bool              `json:"isYearly"`
	IsNew          bool              `json:"isNew"`
	Cards          []BillingCardData `json:"cards"`
	CostForNewUser int               `json:"costForNewUser"`
	CurrentPlan    *data.BillingPlan `json:"currentPlan"`
	Seats          int               `json:"seats"`
	Logins         []model.User      `json:"logins"`
	NextInvoice    *stripe.Invoice   `json:"nextInvoice"`
}

// BillingCardData represents a Stripe credit card
type BillingCardData struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Number     string `json:"number"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	CVC        string `json:"cvc"`
	Brand      string `json:"brand"`
	Expiration string `json:"expiration"`
}

func newBilling() *Route {
	var b interface{} = Billing{}
	return &Route{
		AllowCrossOrigin: true,
		Logger:           true,
		MinimumRole:      model.RoleFree,
		Handler:          b.(http.Handler),
	}
}

func (b Billing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	if r.Method == http.MethodGet {
		if head == "invoices" {
			head, r.URL.Path = ShiftPath(r.URL.Path)
			if head == "" {
				b.invoices(w, r)
			} else if head == "next" {
				b.getNextInvoice(w, r)
				return
			}
		}
	} else if r.Method == http.MethodPost {
		if head == "changeplan" {
			b.changePlan(w, r)
		} else if head == "webhooks" {
			b.stripe(w, r)
		}
	} else if r.Method == http.MethodDelete {
		if head == "card" {
			b.deleteCard(w, r)
		}
	}
}

func (b Billing) Overview(accountID int64) (*BillingOverview, error) {
	// this struct will be returned should we be a paid customer or not
	ov := &BillingOverview{}

	// Get the current account
	account, err := b.DB.Users.GetDetail(accountID)
	if err != nil {
		return nil, fmt.Errorf("unable to find this account: %v", err)
	}

	ov.Account = account
	ov.Logins = account.Users

	// get all logins for user roles
	for _, l := range account.Users {
		if l.Role < model.RoleFree {
			ov.Seats++
		}
	}

	if len(account.StripeID) == 0 {
		ov.IsNew = true

		// if they are on trial, we set the current plan to
		// that so the UI can based permissions on that plan.
		if account.TrialInfo.IsTrial {
			if p, ok := data.GetPlan(account.TrialInfo.Plan); ok {
				ov.CurrentPlan = &p
			}
		}

		ov.Cards = make([]BillingCardData, 0)

		return ov, nil
	}

	// getting stripe customer
	cus, err := customer.Get(account.StripeID, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get stripe customer: %v", err)
	}

	ov.StripeID = cus.ID
	ov.Plan = account.Plan
	ov.IsYearly = account.IsYearly

	if p, ok := data.GetPlan(account.Plan); ok {
		ov.CurrentPlan = &p
	}

	cards := card.List(&stripe.CardListParams{Customer: &account.StripeID})
	for cards.Next() {
		c := cards.Card()
		if !c.Deleted {
			ov.Cards = append(ov.Cards, BillingCardData{
				ID:         c.ID,
				Name:       c.Name,
				Number:     c.Last4,
				Month:      fmt.Sprintf("%d", c.ExpMonth),
				Year:       fmt.Sprintf("%d", c.ExpYear),
				Expiration: fmt.Sprintf("%d / %d", c.ExpMonth, c.ExpYear),
				Brand:      string(c.Brand),
			})
		}
	}

	i, err := invoice.GetNext(&stripe.InvoiceParams{Customer: &account.StripeID})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch next invoice: %v", err)
	}

	ov.NextInvoice = i

	return ov, nil
}

func (b Billing) changeQuantity(stripeID, subID string, qty int) error {
	q := int64(qty)
	p := &stripe.SubscriptionParams{Customer: &stripeID, Quantity: &q}
	_, err := sub.Update(subID, p)
	return err
}

func (b Billing) userRoleChanged(db data.DB, accountID int64, oldRole, newRole model.Roles) (paid bool, err error) {
	acct, err := db.Users.GetDetail(accountID)
	if err != nil {
		return false, err
	}

	// if this is a paid account
	if acct.IsPaid() {
		// if they were a free user
		if oldRole == model.RoleFree {
			// and are now a paid user, we need to +1 qty and prepare the invoice
			if newRole == model.RoleAdmin || newRole == model.RoleUser {
				paid = true

				// we increase the seats number for this account
				acct.Seats++

				// try to change their subscription (+1 qty)
				if err = b.changeQuantity(acct.StripeID, acct.SubscriptionID, acct.Seats); err != nil {
					return
				}

				// ensure that the charges will be immediate and not on next billing date
				if err := queue.Enqueue(queue.TaskCreateInvoice, acct.StripeID); err != nil {
					return paid, err
				}

				if err = db.Users.SetSeats(acct.ID, acct.Seats); err != nil {
					return false, err
				}
			}
		} else {
			// they were a paid user, now they are set as free
			if newRole == model.RoleFree {
				acct.Seats--

				if err = b.changeQuantity(acct.StripeID, acct.SubscriptionID, acct.Seats); err != nil {
					return
				}

				if err = db.Users.SetSeats(acct.ID, acct.Seats); err != nil {
					return false, err
				}
			}
		}
	}
	return false, nil
}

// BillingNewCustomer represents data sent to api for creating a new customer
type BillingNewCustomer struct {
	AccountID   int64
	Email       string
	Plan        string
	StripeToken string
	Coupon      string
	IsPerSeat   bool
	IsYearly    bool
	TrialDays   int
	Quantity    int
}

func (b Billing) Start(bc BillingNewCustomer) error {
	p := &stripe.CustomerParams{Email: stripe.String(bc.Email)}
	p.SetSource(bc.StripeToken)

	c, err := customer.New(p)
	if err != nil {
		return fmt.Errorf("unable to create the customer: %v", err)
	}

	acct, err := b.DB.Users.GetDetail(bc.AccountID)
	if err != nil {
		return fmt.Errorf("unable to get the account for this account ID: %d -> %v", bc.AccountID, err)
	}

	seats := 1
	if bc.IsPerSeat {
		seats = 0
		for _, u := range acct.Users {
			if u.Role < model.RoleFree {
				seats++
			}
		}
	}

	plan := bc.Plan
	if bc.IsYearly {
		plan += "_yearly"
	}

	// we get the current set of pricing plans
	currentPlans := data.GetPlans("current")
	var bp data.BillingPlan
	for _, p := range currentPlans {
		if p.Name == plan {
			bp = p
			break
		}
	}

	if len(bp.ID) == 0 {
		return fmt.Errorf("unable to find this plan %s in the 'current' pricing se", plan)
	}

	seatsptr := int64(seats)
	subp := &stripe.SubscriptionParams{
		Customer: &c.ID,
		Plan:     stripe.String(bp.StripeID),
		Quantity: &seatsptr,
	}

	if bc.TrialDays > 0 {
		subp.TrialPeriodDays = stripe.Int64(int64(bc.TrialDays))
	}

	if len(bc.Coupon) > 0 {
		subp.Coupon = &bc.Coupon
	}

	s, err := sub.New(subp)
	if err != nil {
		return fmt.Errorf("unable to create the subscription: %v", err)
	}

	acct.TrialInfo.IsTrial = false
	if err := b.DB.Users.ConvertToPaid(acct.ID, c.ID, s.ID, bc.Plan, bc.IsYearly, seats); err != nil {
		return fmt.Errorf("unable to convert the account to a paid account: %v", err)
	}

	//TODO: Trigger a new customer event

	return nil
}

// Convert upgrades an existing trial account to a paid account.
func (b Billing) Convert(bc BillingNewCustomer) error {
	acct, err := b.DB.Users.GetDetail(bc.AccountID)
	if err != nil {
		return fmt.Errorf("unable to get the account for this account ID: %d -> %v", bc.AccountID, err)
	}

	p := &stripe.CustomerParams{Email: stripe.String(bc.Email)}
	p.SetSource(bc.StripeToken)

	c, err := customer.New(p)
	if err != nil {
		return fmt.Errorf("unable to create the customer: %v", err)
	}

	if bc.Quantity <= 0 {
		bc.Quantity = 1
	}

	seats := bc.Quantity
	// IsPerSeat means based on number of users.
	// Otherwise, we're using the bc.Quantity field
	if bc.IsPerSeat {
		seats = 0
		for _, u := range acct.Users {
			if u.Role < model.RoleFree {
				seats++
			}
		}
	}

	plan := bc.Plan
	if bc.IsYearly {
		plan += "_yearly"
	}

	// we get the current set of pricing plans
	currentPlans := data.GetPlans("current")
	var bp data.BillingPlan
	for _, p := range currentPlans {
		if p.Name == plan {
			bp = p
			break
		}
	}

	if len(bp.ID) == 0 {
		return fmt.Errorf("unable to find this plan %s in the 'current' pricing se", plan)
	}

	seatsptr := int64(seats)
	subp := &stripe.SubscriptionParams{
		Customer: &c.ID,
		Plan:     stripe.String(bp.StripeID),
		Quantity: &seatsptr,
	}

	if bc.TrialDays > 0 {
		subp.TrialPeriodDays = stripe.Int64(int64(bc.TrialDays))
	}

	if len(bc.Coupon) > 0 {
		subp.Coupon = &bc.Coupon
	}

	s, err := sub.New(subp)
	if err != nil {
		return fmt.Errorf("unable to create the subscription: %v", err)
	}

	acct.TrialInfo.IsTrial = false
	if err := b.DB.Users.ConvertToPaid(acct.ID, c.ID, s.ID, bc.Plan, bc.IsYearly, seats); err != nil {
		return fmt.Errorf("unable to convert the account to a paid account: %v", err)
	}

	//TODO: Trigger a new customer event

	return nil
}

func (b Billing) changePlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	var data = new(struct {
		Plan     string `json:"plan"`
		IsYearly bool   `json:"isYearly"`
	})
	if err := ParseBody(r.Body, &data); err != nil {
		Respond(w, r, http.StatusBadRequest, err)
		return
	}

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	plan := data.Plan

	newLevel, currentLevel := 0, 0
	if len(plan) == 0 || plan == "free" {
		newLevel = 0
	} else if strings.HasPrefix(plan, "starter") {
		newLevel = 1
	} else if strings.HasPrefix(plan, "pro") {
		newLevel = 2
	} else {
		newLevel = 3
	}

	if strings.HasPrefix(account.Plan, "starter") {
		currentLevel = 1
	} else if strings.HasPrefix(account.Plan, "pro") {
		currentLevel = 2
	} else {
		currentLevel = 3
	}

	// did they cancelled
	if newLevel == 0 {
		// we need to cancel their subscriptions
		if _, err := sub.Cancel(account.SubscriptionID, nil); err != nil {
			Respond(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := db.Users.Cancel(account.ID); err != nil {
			Respond(w, r, http.StatusInternalServerError, err)
			return
		}
	} else {
		if data.IsYearly {
			plan += "_yearly"
		}

		seats := 0
		for _, u := range account.Users {
			if u.Role < model.RoleFree {
				seats++
			}
		}

		seatsptr := int64(seats)
		subParams := &stripe.SubscriptionParams{
			Customer: &account.StripeID,
			Plan:     &plan,
			Quantity: &seatsptr,
		}

		// if we upgrade we need to change billing cycle date
		upgraded := false
		if newLevel > currentLevel {
			upgraded = true
		} else if account.IsYearly == false && data.IsYearly {
			upgraded = true
		}

		if upgraded {
			// queue an invoice create for this upgrade
			queue.Enqueue(queue.TaskCreateInvoice, account.StripeID)
		}

		if _, err := sub.Update(account.SubscriptionID, subParams); err != nil {
			Respond(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := db.Users.ChangePlan(account.ID, plan, data.IsYearly); err != nil {
			Respond(w, r, http.StatusInternalServerError, err)
			return
		}
		Respond(w, r, http.StatusOK, true)
	}
}

func (b Billing) updateCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		Respond(w, r, http.StatusBadRequest, err)
		return
	}

	var data BillingCardData
	if err := ParseBody(r.Body, &data); err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	if c, err := card.Update(data.ID, &stripe.CardParams{
		Customer: &account.StripeID,
		ExpMonth: &data.Month,
		ExpYear:  &data.Month,
		CVC:      &data.CVC,
	}); err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
	} else {
		card := BillingCardData{
			ID:         c.ID,
			Name:       c.Name,
			Number:     c.Last4,
			Month:      fmt.Sprintf("%d", c.ExpMonth),
			Year:       fmt.Sprintf("%d", c.ExpYear),
			Expiration: fmt.Sprintf("%d / %d", c.ExpMonth, c.ExpYear),
			Brand:      string(c.Brand),
		}
		Respond(w, r, http.StatusOK, card)
	}
}

func (b Billing) addCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		Respond(w, r, http.StatusBadRequest, err)
		return
	}

	var data BillingCardData
	if err := ParseBody(r.Body, &data); err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	if c, err := card.New(&stripe.CardParams{
		Customer: &account.StripeID,
		Name:     &data.Name,
		Number:   &data.Number,
		ExpMonth: &data.Month,
		ExpYear:  &data.Year,
		CVC:      &data.CVC}); err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
	} else {
		card := BillingCardData{
			ID:         c.ID,
			Number:     c.Last4,
			Expiration: fmt.Sprintf("%d / %d", c.ExpMonth, c.ExpYear),
			Brand:      string(c.Brand),
		}
		Respond(w, r, http.StatusOK, card)
	}
}

func (b Billing) deleteCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		Respond(w, r, http.StatusBadRequest, err)
		return
	}

	cardID, _ := ShiftPath(r.URL.Path)

	if _, err := card.Del(cardID, &stripe.CardParams{Customer: &account.StripeID}); err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
	} else {
		Respond(w, r, http.StatusOK, true)
	}
}

func (b Billing) invoices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	var invoices []*stripe.Invoice

	iter := invoice.List(&stripe.InvoiceListParams{Customer: &account.StripeID})
	for iter.Next() {
		invoices = append(invoices, iter.Invoice())
	}

	Respond(w, r, http.StatusOK, invoices)
}

func (b Billing) getNextInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	i, err := invoice.GetNext(&stripe.InvoiceParams{Customer: &account.StripeID})
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, i)
}

// StripeWebhook is used to grab data sent by Stripe for a webhook
type StripeWebhook struct {
	Event stripe.Event `json:"event"`
}

// WebhookData used when stripe webhook is call
type WebhookData struct {
	ID   string            `json:"id"`
	Type string            `json:"type"`
	Data WebhookDataObject `json:"data"`
}

// WebhookDataObject is the container for the object received
type WebhookDataObject struct {
	Object WebhookDataObjectData `json:"object"`
}

// WebhookDataObjectData is the object being sent by stripe
type WebhookDataObjectData struct {
	ID           string `json:"id"`
	Customer     string `json:"customer"`
	Subscription string `json:"subscription"`
	Closed       bool   `json:"closed"`
}

func (b Billing) stripe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(ContextDatabase).(*data.DB)

	// no matter what happen, Stripe wants us to send a 200
	defer w.Write([]byte("ok"))

	var data WebhookData
	if err := ParseBody(r.Body, &data); err != nil {
		log.Println(err)
		return
	}

	if data.Type == "customer.subscription.deleted" {
		subID := data.Data.Object.ID
		if len(subID) == 0 {
			log.Println(fmt.Errorf("no subscription found to this customer.subscription.deleted %s", data.ID))
			return
		}

		stripeID := data.Data.Object.Customer
		if len(stripeID) == 0 {
			log.Println(fmt.Errorf("no customer found to this invoice.payment_succeeded %s", data.ID))
			return
		}

		// check if it's a failed payment_succeeded
		account, err := db.Users.GetByStripe(stripeID)
		if err != nil {
			log.Println(fmt.Errorf("no customer matches stripe id %s", stripeID))
			return
		}

		if len(account.SubscriptionID) > 0 {
			//TODO: Send emails

			if err := db.Users.Cancel(account.ID); err != nil {
				log.Println(fmt.Errorf("unable to cancel this account %v", account.ID))
				return
			}
		}
	}
}

func (b Billing) cancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	var data = new(struct {
		Reason string `json:"reason"`
	})
	if err := ParseBody(r.Body, &data); err != nil {
		Respond(w, r, http.StatusBadRequest, err)
		return
	}

	// SendMail would be call here passing the reason

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	if _, err := sub.Cancel(account.SubscriptionID, nil); err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := db.Users.Cancel(account.ID); err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	Respond(w, r, http.StatusOK, true)
}
