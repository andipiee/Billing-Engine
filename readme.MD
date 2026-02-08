# Billing Engine – Take Home Test

This repository contains a simple **Billing Engine** implementation for a loan system, created as part of a Senior Software Engineer take-home assignment.

The solution focuses on **correctness, clarity, simplicity**, and **test coverage**, while implementing only what is required by the problem statement.

---

## Problem Overview

We are required to build a billing system that supports:

- A **50-week loan** of **Rp 5,000,000**
- A **flat interest rate of 10% per annum**
- Fixed weekly repayments
- Ability to determine:
  - Weekly repayment schedule
  - Outstanding loan amount
  - Whether a borrower is delinquent

### Business Rules

- Weekly installment amount is fixed
- Borrower can either:
  - Pay **exactly one weekly installment**, or
  - Not pay at all
- If a borrower misses **2 consecutive repayments**, they are considered **delinquent**
- Missed payments must be **caught up later**
  - Each payment clears exactly **one unpaid week**

---

## Assumptions

The following reasonable assumptions were made where the problem statement was not explicit:

1. Flat interest is applied once at loan creation  
2. Weekly installment = `(principal + interest) / totalWeeks`
3. Borrower may make **multiple payments in the same week** to catch up
4. Payments are always applied to the **oldest unpaid scheduled week**
5. Delinquency depends only on the **last two scheduled weeks**
6. Billing progression is **explicitly week-based**, not inferred from wall-clock time
7. Payments are modeled as an **append-only ledger** (in-memory for this exercise)

---

## Project Structure

```
billing-engine/
├── go.mod
├── go.sum
├── model.go          # Domain models
├── business.go       # Core billing logic
├── business_test.go  # Unit tests
└── main.go           # Example usage / demo
```

---

## Design Overview

### Schedule-Driven Logic

- Loan progression is based on **billing weeks**, not system time
- This keeps the logic deterministic and easy to test

### Derived State

- **Outstanding balance** is derived from:
  - Total loan amount
  - Sum of recorded payments
- No mutable accounting state is stored

### Catch-Up Payments

- Payments always clear the **oldest unpaid week**
- This supports recovery from missed repayments

### Separation of Concerns

- `model.go` → domain entities
- `business.go` → business rules
- `main.go` → orchestration / demo
- `business_test.go` → verification of behavior

---

## Public API

### Create a Loan

```go
loan := NewLoan(id, principal, interestRate, totalWeeks)
```

### Advance Billing Week

```go
loan.AdvanceWeek()
```

### Make a Payment

```go
loan.MakePayment(loan.WeeklyInstallment)
```

### Get Outstanding Balance

```go
loan.GetOutstanding()
```

### Check Delinquency Status

```go
loan.IsDelinquent()
```

---

## Example Scenario

1. Borrower pays Week 1  
2. Misses Week 2 and Week 3  
3. Borrower becomes delinquent  
4. At Week 4, borrower makes 3 payments  
   - Clears Week 2, Week 3, and Week 4  
5. Outstanding balance is reduced correctly  
6. Delinquency status is cleared

This behavior matches the problem requirements.

---

## Tests

Unit tests are provided for:

- Loan initialization
- Outstanding balance calculation
- Payment validation
- Catch-up payment behavior
- Delinquency detection
- Edge cases (invalid payment, advancing beyond loan duration)

Run tests with:

```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
```

- Business logic is fully covered by tests
- `main.go` is intentionally excluded from coverage

---

## How to Run

```bash
go run .
```

---

## Notes for Reviewer

- The implementation avoids over-engineering and focuses on the specified requirements
- No database, cron job, or external services are included, as they are outside the scope of this task
- The design is intentionally extensible for real-world use

---

## Author

**Andi**  
Senior Software Engineer Candidate