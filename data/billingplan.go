package data

const (
	// Plan201612 is for plans from Dec 2016 to Jul 2017
	Plan201612 = "201612"
	// Plan201711 current pricing set
	Plan201711 = "201707"
)

// BillingFlags is used to set which integrations a plan is authorize to use
type BillingFlags int

// BillingPlan defines what one plan to have access to and set limitations
type BillingPlan struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Price       float32 `json:"price"`
	YearlyPrice float32 `json:"yearly"`
}

var plans map[string]BillingPlan

func init() {
	plans = make(map[string]BillingPlan)

	plans["free"] = BillingPlan{
		ID:      "free",
		Name:    "Free",
		Version: "201612",
	}

	plans["starter-201612"] = BillingPlan{
		ID:          "starter-201612",
		Name:        "Starter",
		Version:     "201612",
		Price:       25,
		YearlyPrice: 15,
	}

	plans["pro-201612"] = BillingPlan{
		ID:          "pro-201612",
		Name:        "Pro",
		Version:     "201612",
		Price:       55,
		YearlyPrice: 35,
	}

	plans["enterprise-201612"] = BillingPlan{
		ID:          "enterprise-201612",
		Name:        "Enterprise",
		Version:     "201612",
		Price:       95,
		YearlyPrice: 65,
	}

	plans["starter-201707"] = BillingPlan{
		ID:          "starter-201707",
		Name:        "Starter",
		Version:     "201707",
		Price:       39,
		YearlyPrice: 29,
	}

	plans["pro-201707"] = BillingPlan{
		ID:          "pro-201707",
		Name:        "Pro",
		Version:     "201707",
		Price:       99,
		YearlyPrice: 79,
	}

	plans["enterprise-201707"] = BillingPlan{
		ID:          "enterprise-201707",
		Name:        "Enterprise",
		Version:     "201707",
		Price:       129,
		YearlyPrice: 159,
	}
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
func GetPlansVersion(plan string) []BillingPlan {
	if p, ok := plans[plan]; ok {
		return GetPlans(p.Version)
	}
	// we are returning current plan since we could not find this plan
	return GetPlans(Plan201711)
}
