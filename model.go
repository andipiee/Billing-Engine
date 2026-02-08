package main

import (
	"time"

	"github.com/shopspring/decimal"
)

type Loan struct {
	ID                int
	Principal         decimal.Decimal
	InterestRate      decimal.Decimal
	WeeklyInstallment decimal.Decimal
	Outstanding       decimal.Decimal
	TotalWeeks        int
	StartDate         time.Time
	CurrentWeek       int
	Payments          map[int]Payment // week -> payment
	TotalAmount       decimal.Decimal
}

type Payment struct {
	Week   int
	Amount decimal.Decimal
	PaidAt time.Time
}

// Loan creation
func NewLoan(id int, principal decimal.Decimal, interestRate decimal.Decimal, weeks int) *Loan {
	// totalAmount = principal + (principal * interestRate)
	totalAmount := principal.Mul(decimal.NewFromInt(1).Add(interestRate))

	// weeklyInstallment = totalAmount / weeks
	weekly := totalAmount.Div(decimal.NewFromInt(int64(weeks)))

	return &Loan{
    	ID:                id,
		Principal:         principal,
		InterestRate:      interestRate,
		TotalAmount:       totalAmount,
		WeeklyInstallment: weekly,
		Outstanding:       totalAmount,
		TotalWeeks:        weeks,
		StartDate:         time.Now(),
		CurrentWeek:       0,
		Payments:          make(map[int]Payment),
	}
}
