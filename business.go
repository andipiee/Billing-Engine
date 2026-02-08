package main

import (
	"fmt"
	"time"
	"github.com/shopspring/decimal"
)

func (l *Loan) GetOutstanding() decimal.Decimal {
	return l.Outstanding
}

// IsDelinquent returns true if last 2 scheduled payments were missed
func (l *Loan) IsDelinquent() bool {
	missed := 0
	for w := l.CurrentWeek; w >= 1 && missed < 2; w-- {
		if _, paid := l.Payments[w]; !paid {
			missed++
		} else {
			break
		}
	}
	return missed >= 2
}

// MakePayment processes a payment for the current week
func (l *Loan) MakePayment(amount decimal.Decimal) error {
	if !amount.Equal(l.WeeklyInstallment) {
		return fmt.Errorf("payment must be exactly %s", l.WeeklyInstallment)
	}

	// Find the oldest unpaid week
	for w := 1; w <= l.CurrentWeek; w++ {
		if _, paid := l.Payments[w]; !paid {
			l.Payments[w] = Payment{
				Week:   w,
				Amount: amount,
				PaidAt: time.Now(),
			}
			l.Outstanding = l.Outstanding.Sub(amount)
			return nil
		}
	}

	return fmt.Errorf("no unpaid weeks to pay")
}

// AdvanceWeek moves the loan to the next billing week
func (l *Loan) AdvanceWeek() error {
	if l.CurrentWeek >= l.TotalWeeks {
		return fmt.Errorf("cannot advance beyond total weeks")
	}
	l.CurrentWeek++
	return nil
}