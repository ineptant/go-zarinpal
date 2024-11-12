# go-zarinpal

It's an implementation of the [Zarinpal](https://www.zarinpal.com/) payment gateway with Golang.

## ⚒️ Installation

```bash
go get github.com/ineptant/go-zarinpal
```

## 💻 Usage

```go
import github.com/ineptant/go-zarinpal
```

### New

```go
merchantID := "XXXX-XXXX-XXXX-XXXX" // Your merchant id
z, err := zarinpal.New(merchantID, false)
if err != nil {
	log.Fatal(err)
}
```

### Create new payment request

The payment amount, the description displayed on the payment gateway, and a callback url to the application from the payment gateway. We use this callback url to get authority to verify a transaction.

```go
paymentURL, authority, statusCode, err := z.NewPayment(10000, "http://callback/url", "Description")
if err != nil {
	log.Fatal(err)
}
```

Save the authority and the payment amount, and then redirect to the gateway using paymentUrl.

### Handle callback

Get the authority from callback url. Get the payment amount from your database using the authority, and then verify your transaction. Don't forget to verify all successful transactions:

```go
verified, refID, statusCode, err := z.VerifyTransaction(10000, "authority")
if err != nil {
	log.Fatal(err)
}
```

If successfully verified, save refID.

### Get unverified transactions

You can periodically check if there are unverified transactions, and then verify them.

```go
authorities, statusCode, err := z.UnverifiedTransactions()
if err != nil {
	log.Fatal(err)
}
```

### Get transaction inquiry

To get status of a transaction, use:

```go
inquiry, statusCode, err := z.TransactionInquiry("authority")
if err != nil {
	log.Fatal(err)
}
```

### Reverse transaction

To reverse an unsuccessful transaction, use:

```go
statusCode, err := z.ReverseTransaction("authority")
if err != nil {
	log.Fatal(err)
}
```

note that if you want to reverse transactions, you should first set your ip address on your Zarinpal panel.