package main

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func main() {
	principal := decimal.NewFromInt(5_000_000)
	interestRate := decimal.NewFromFloat(0.10)

	loan := NewLoan(100, principal, interestRate, 50)
	fmt.Println("Loan: ", *loan)

	fmt.Println("Loan created")
	fmt.Println("Weekly installment:", loan.WeeklyInstallment)

	// make first advance
	_ = loan.AdvanceWeek()

	// ---- Week 1 ----
	fmt.Println("\n== Week 1 ==")
	fmt.Println("Outstanding before first payment:", loan.GetOutstanding())
	if err := loan.MakePayment(loan.WeeklyInstallment); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Outstanding after first payment:", loan.GetOutstanding())
	fmt.Println("Delinquent:", loan.IsDelinquent())
}