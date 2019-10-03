package data

// BillingFlags is used to set which integrations a plan is authorize to use
type BillingFlags int

// BillingPlan defines what one plan to have access to and set limitations
type BillingPlan struct {
	ID          string                 `json:"id"`
	StripeID    string                 `json:"stripeId"`
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Price       float32                `json:"price"`
	YearlyPrice float32                `json:"yearly"`
	Params      map[string]interface{} `json:"params"`
}

var plans map[string]BillingPlan

func init() {
	plans = make(map[string]BillingPlan)
}

func AddPlan(plan BillingPlan) {
	plans[plan.ID] = plan
}

// GetPlan returns a specific plan by ID
func GetPlan(id string) (BillingPlan, bool) {
	v, ok := plans[id]
	return v, ok
}

// GetPlans returns a slice of the desired version plans
func GetPlans(v string) []BillingPlan {
	var list []BillingPlan
	for k, p := range plans {
		if k == "free" {
			// the free plan is available on all versions
			list = append(list, p)
		} else if p.Version == v {
			// this is a plan for the requested version
			list = append(list, p)
		}
	}
	return list
}

// GetPlansVersion returns a slice of the plans matching a current plan
func GetPlansVersion(plan string, defaultVersion string) []BillingPlan {
	if p, ok := plans[plan]; ok {
		return GetPlans(p.Version)
	}
	// we are returning current plan since we could not find this plan
	return GetPlans(defaultVersion)
}
