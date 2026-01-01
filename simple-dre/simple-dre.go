// Simple Deterministic Rules Engine for recommendation based on user's constraint
package main

import (
	"fmt"
	"sort"
	"time"
)

type Location struct {
	Name string
	Lat  float32
	Long float32
}

type Accommodation struct {
	Name          string
	PricePerNight float64
	Location      Location
	AvailableFrom string
	AvailableTo   string
	Rating        float64
	TotalScore    float64
	Distance      float64
}

type Constraint struct {
	Budget          float64
	ArrivalCity     string
	FromDate        time.Time
	ToDate          time.Time
	RemainingBudget float64
}

var accommodations = []Accommodation{
	{
		Name: "Movenpick Mecca",
		Location: Location{
			Name: "Mecca",
			Lat:  0.000,
			Long: 0.000,
		},
		PricePerNight: 4000000,
		AvailableFrom: "01/01/2026",
		AvailableTo:   "31/12/2026",
		Rating:        5.0,
		Distance:      0.5,
	},
	{
		Name: "Hilton Mecca",
		Location: Location{
			Name: "Mecca",
			Lat:  0.000,
			Long: 0.000,
		},
		PricePerNight: 6000000,
		AvailableFrom: "01/01/2026",
		AvailableTo:   "31/12/2026",
		Rating:        4.5,
		Distance:      0.5,
	},
	{
		Name: "Near ring 1",
		Location: Location{
			Name: "Mecca",
			Lat:  0.000,
			Long: 0.000,
		},
		PricePerNight: 2000000,
		AvailableFrom: "01/01/2026",
		AvailableTo:   "31/12/2026",
		Rating:        4.0,
		Distance:      1.5,
	},
	{
		Name: "closer to mosque",
		Location: Location{
			Name: "Mecca",
			Lat:  0.000,
			Long: 0.000,
		},
		PricePerNight: 2800000,
		AvailableFrom: "01/01/2026",
		AvailableTo:   "31/12/2026",
		Rating:        3.8,
		Distance:      0.4,
	},
	{
		Name: "motel somewhere",
		Location: Location{
			Name: "Mecca",
			Lat:  0.000,
			Long: 0.000,
		},
		PricePerNight: 800000,
		AvailableFrom: "01/01/2026",
		AvailableTo:   "31/12/2026",
		Rating:        2.4,
		Distance:      2.0,
	},
}

func main() {
	from, _ := time.Parse("02/01/2006", "25/05/2026")
	to, _ := time.Parse("02/01/2006", "04/06/2026")
	travelDuration := to.Sub(from).Hours() / 24

	user := Constraint{
		Budget:      35000000,
		ArrivalCity: "Mecca",
		// FromDate:    from,
		// ToDate:      to,
	}

	// var AlHaram = Location{Name: "Al-Haram", Lat: 21.4225, Long: 39.8262}

	// Filtering
	var filteredAccommodations []Accommodation
	var minPrice, maxPrice float64

	for _, accommodation := range accommodations {
		totalCost := accommodation.PricePerNight * travelDuration

		if accommodation.Location.Name == user.ArrivalCity {
			if accommodation.PricePerNight > maxPrice {
				maxPrice = accommodation.PricePerNight
			}
			if minPrice == 0 || accommodation.PricePerNight < minPrice {
				minPrice = accommodation.PricePerNight
			}
		}

		if accommodation.Location.Name == user.ArrivalCity && totalCost <= user.Budget {
			filteredAccommodations = append(filteredAccommodations, accommodation)
		}
	}

	fmt.Printf("Max price?? %.0f", maxPrice)

	minimumBudget := minPrice * travelDuration
	fmt.Printf("Minimum budget for this trip is %.0f\n", minimumBudget)

	if user.Budget < minimumBudget {
		fmt.Println("You don't have enough budget for this trip!")
		return
	}
	// Scoring
	weight := GetWeight(user.Budget, minPrice, maxPrice, travelDuration)
	fmt.Printf("User weight: %.0f \n", weight)

	for i := range filteredAccommodations {
		filteredAccommodations[i].TotalScore = CalculateScore(filteredAccommodations[i], minPrice, maxPrice, weight)
	}

	// Sorting
	sort.Slice(filteredAccommodations, func(i, j int) bool {
		return filteredAccommodations[i].TotalScore > filteredAccommodations[j].TotalScore
	})

	if len(filteredAccommodations) > 0 {
		bestMatch := filteredAccommodations[0]
		fmt.Printf("Recommended stay: %s, (Score: %.2f)\n", bestMatch.Name, bestMatch.TotalScore)
		fmt.Printf("Total Cost for %.0f nights: %.2f\n", travelDuration, bestMatch.PricePerNight*travelDuration)
		fmt.Println("Other hotels that might suit you: ")
		for i, v := range filteredAccommodations {
			fmt.Printf("%d. %s \n", i+1, v.Name)
		}
	} else {
		fmt.Println("No accommodations found within your budget")
	}
}

func CalculateScore(accomodation Accommodation, minPrice, maxPrice float64, userWeight float64) float64 {

	priceScore := (maxPrice - accomodation.PricePerNight) / (maxPrice - minPrice)
	distanceScore := CalculateDistanceToHaram(accomodation.Distance)
	ratingScore := accomodation.Rating / 5.0
	// Add penalty for rating below 3.0 so it's not recommended on top
	if accomodation.Rating < 3.0 {
		ratingScore *= 0.5
	}

	if maxPrice == minPrice {
		priceScore = 1.0
	}

	weightPrice := userWeight
	weightRating := (1 - userWeight) * 0.5
	weightDistance := (1.0 - userWeight) * 0.5

	return (weightPrice * priceScore) + (weightRating * ratingScore) + (weightDistance * distanceScore)
}

func GetWeight(userBudget, minPricePerNight, maxPricePerNight, travelDays float64) float64 {
	const MAX_WEIGHT = 0.9
	const MIN_WEIGHT = 0.1

	totalMinPrice := minPricePerNight * travelDays
	totalMaxPrice := maxPricePerNight * travelDays

	if totalMaxPrice == totalMinPrice {
		return 0.5
	}

	if userBudget >= totalMaxPrice {
		return MIN_WEIGHT
	}

	postition := (userBudget - totalMinPrice) / (totalMaxPrice - totalMinPrice)

	return MAX_WEIGHT - (postition * (MAX_WEIGHT - MIN_WEIGHT))
}

func CalculateDistanceToHaram(distance float64) float64 {
	const MaxRadius = 1.0 // in km
	if distance >= MaxRadius {
		return 0.0
	}

	// inverse logic: Closer to 0 (in distance) means closer to 1.0 score (good)
	return (MaxRadius - distance) / MaxRadius
}
