package midtrans_test

import (
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-midtrans-sdk"
	"testing"
	"time"
)

func Test_api_ChargeQRIS(t *testing.T) {

}

func Test_api_ChargeQRIS_Validator(t *testing.T) {
	validator, _ := midtrans.NewValidator()
	tests := []struct {
		name    string
		request midtrans.ChargeRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: midtrans.ChargeRequest{
				PaymentType: "qris",
				TransactionDetail: midtrans.TransactionDetail{
					OrderID:     "order123",
					GrossAmount: 100000,
				},
				CustomExpiry: &midtrans.CustomExpiry{
					OrderTime:      time.Now(),
					ExpiryDuration: 60,
					Unit:           "minutes",
				},
				ItemDetails: []midtrans.ItemDetail{
					{
						ID:       "item1",
						Name:     "Item 1",
						Price:    50000,
						Qty:      2,
						Brand:    "BrandA",
						Category: "CategoryA",
					},
				},
				CustomerDetail: &midtrans.CustomerDetail{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					Phone:     "123456789",
					BillingAddress: &midtrans.CustomerBillingAddress{
						FirstName:   "John",
						LastName:    "Doe",
						Phone:       "123456789",
						Address:     "123 Street",
						City:        "CityX",
						PostalCode:  "12345",
						CountryCode: "IDN",
					},
					ShippingAddress: &midtrans.CustomerShippingAddress{
						FirstName:   "John",
						LastName:    "Doe",
						Phone:       "123456789",
						Address:     "456 Avenue",
						City:        "CityY",
						PostalCode:  "67890",
						CountryCode: "IDN",
					},
				},
				QRIS: &midtrans.QRIS{
					Acquirer: "gopay",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid PaymentType",
			request: midtrans.ChargeRequest{
				PaymentType: "",
				TransactionDetail: midtrans.TransactionDetail{
					OrderID:     "order123",
					GrossAmount: 100000,
				},
				QRIS: &midtrans.QRIS{
					Acquirer: "gopay",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid CustomerDetail email",
			request: midtrans.ChargeRequest{
				PaymentType: "qris",
				TransactionDetail: midtrans.TransactionDetail{
					OrderID:     "order123",
					GrossAmount: 100000,
				},
				CustomerDetail: &midtrans.CustomerDetail{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "invalid-email",
					Phone:     "123456789",
				},
				QRIS: &midtrans.QRIS{
					Acquirer: "gopay",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid QRIS Acquirer",
			request: midtrans.ChargeRequest{
				PaymentType: "qris",
				TransactionDetail: midtrans.TransactionDetail{
					OrderID:     "order123",
					GrossAmount: 100000,
				},
				QRIS: &midtrans.QRIS{
					Acquirer: "invalid_acquirer",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.request)
			fmt.Println(err)
			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none for test %s", tt.name)
			} else if !tt.wantErr && err != nil {
				t.Errorf("did not expect error but got %v for test %s", err, tt.name)
			}
		})
	}
}
