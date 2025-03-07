// Question: Integrate a new payment gateway into an existing e-commerce application without changing the existing code.
// Scenario: Create an Adapter class that allows the new payment gateway to be used as if it were the existing one.

package main

import "fmt"

// PaymentGateway interface
type PaymentGateway interface {
	Pay(amount float64) string
}

// PaymentGatewayAdapter struct
type PaymentGatewayAdapter struct {
	paymentGateway PaymentGateway
}

// Pay function
func (pga *PaymentGatewayAdapter) Pay(amount float64) string {

	return pga.paymentGateway.Pay(amount)
}

// NewPaymentGatewayAdapter function
func NewPaymentGatewayAdapter(paymentGateway PaymentGateway) *PaymentGatewayAdapter {
	return &PaymentGatewayAdapter{paymentGateway: paymentGateway}
}

// PayPal struct
type PayPal struct{}

// Pay function
func (p *PayPal) Pay(amount float64) string {
	return "Paid using PayPal"
}

// Stripe struct
type Stripe struct{}

// Pay function
func (s *Stripe) Pay(amount float64) string {
	return "Paid using Stripe"
}

func main() {
	paypal := &PayPal{}
	adapter := NewPaymentGatewayAdapter(paypal)

	// Print the payment result
	result := adapter.Pay(100)
	fmt.Println(result)

	stripe := &Stripe{}
	stripeAdapter := NewPaymentGatewayAdapter(stripe)
	result = stripeAdapter.Pay(200)
	fmt.Println(result)
}
