package main

import (
	"testing"

	"github.com/shopspring/decimal"
)

func newTestLoan() *Loan {
	return NewLoan(
		1,
		decimal.NewFromInt(5_000_000),
		decimal.NewFromFloat(0.10),
		50,
	)
}

func TestNewLoanInitialization(t *testing.T) {
	loan := newTestLoan()

	expectedTotal := decimal.NewFromInt(5_500_000)
	if !loan.TotalAmount.Equal(expectedTotal) {
		t.Fatalf("expected total %s, got %s", expectedTotal, loan.TotalAmount)
	}

	if loan.CurrentWeek != 0 {
		t.Fatalf("expected currentWeek 0, got %d", loan.CurrentWeek)
	}
}

func TestOutstandingAtWeekZero(t *testing.T) {
	loan := newTestLoan()

	if !loan.GetOutstanding().Equal(loan.TotalAmount) {
		t.Fatalf("outstanding should equal total amount at week 0")
	}
}

func TestMakePaymentReducesOutstanding(t *testing.T) {
	loan := newTestLoan()
	loan.AdvanceWeek() // week 1

	err := loan.MakePayment(loan.WeeklyInstallment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := loan.TotalAmount.Sub(loan.WeeklyInstallment)
	if !loan.GetOutstanding().Equal(expected) {
		t.Fatalf("expected outstanding %s, got %s", expected, loan.GetOutstanding())
	}
}

func TestMakePaymentWrongAmountFails(t *testing.T) {
	loan := newTestLoan()
	loan.AdvanceWeek()

	err := loan.MakePayment(decimal.NewFromInt(1))
	if err == nil {
		t.Fatalf("expected error for wrong payment amount")
	}
}

func TestMakePaymentWhenNoUnpaidWeeks(t *testing.T) {
	loan := newTestLoan()

	// Move to week 1
	loan.AdvanceWeek()

	// Pay week 1
	err := loan.MakePayment(loan.WeeklyInstallment)
	if err != nil {
		t.Fatalf("unexpected error paying week 1: %v", err)
	}

	// Try to pay again with no unpaid weeks
	err = loan.MakePayment(loan.WeeklyInstallment)
	if err == nil {
		t.Fatalf("expected error when no unpaid weeks exist")
	}
}
func TestCatchUpPaymentsAppliedCorrectly(t *testing.T) {
	loan := newTestLoan()

	// Week 1 paid
	loan.AdvanceWeek()
	_ = loan.MakePayment(loan.WeeklyInstallment)

	// Week 2 missed
	loan.AdvanceWeek()

	// Week 3 missed
	loan.AdvanceWeek()

	// Week 4
	loan.AdvanceWeek()

	// Catch up: pay 3 times
	_ = loan.MakePayment(loan.WeeklyInstallment)
	_ = loan.MakePayment(loan.WeeklyInstallment)
	_ = loan.MakePayment(loan.WeeklyInstallment)

	expectedPaid := loan.WeeklyInstallment.Mul(decimal.NewFromInt(4))
	expectedOutstanding := loan.TotalAmount.Sub(expectedPaid)

	if !loan.GetOutstanding().Equal(expectedOutstanding) {
		t.Fatalf("expected outstanding %s, got %s",
			expectedOutstanding, loan.GetOutstanding())
	}
}

func TestDelinquencyLogic(t *testing.T) {
	loan := newTestLoan()

	loan.AdvanceWeek() // week 1
	loan.AdvanceWeek() // week 2 (missed)
	loan.AdvanceWeek() // week 3 (missed)

	if !loan.IsDelinquent() {
		t.Fatalf("expected borrower to be delinquent")
	}
}

func TestDelinquencyClearedAfterCatchUp(t *testing.T) {
	loan := newTestLoan()

	loan.AdvanceWeek() // 1
	loan.AdvanceWeek() // 2 missed
	loan.AdvanceWeek() // 3 missed

	if !loan.IsDelinquent() {
		t.Fatalf("expected delinquent before payment")
	}

	// catch up
	_ = loan.MakePayment(loan.WeeklyInstallment)
	_ = loan.MakePayment(loan.WeeklyInstallment)

	if loan.IsDelinquent() {
		t.Fatalf("expected delinquency to be cleared")
	}
}

func TestAdvanceWeekLimit(t *testing.T) {
	loan := newTestLoan()

	for i := 0; i < loan.TotalWeeks; i++ {
		_ = loan.AdvanceWeek()
	}

	err := loan.AdvanceWeek()
	if err == nil {
		t.Fatalf("expected error advancing beyond total weeks")
	}
}