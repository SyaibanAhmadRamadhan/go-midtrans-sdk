package midtrans

type TransactionStatus string

var (
	TransactionStatusSettlement TransactionStatus = "settlement"
	TransactionStatusExpired    TransactionStatus = "expire"
	TransactionStatusPending    TransactionStatus = "pending"
)

type ChargeRequest struct {
	PaymentType       string            `json:"payment_type" validate:"required"`
	TransactionDetail TransactionDetail `json:"transaction_details" validate:"required"`
	CustomExpiry      *CustomExpiry     `json:"custom_expiry,omitempty" validate:"omitempty"`
	ItemDetails       []ItemDetail      `json:"item_details,omitempty" validate:"omitempty,dive"`
	CustomerDetail    *CustomerDetail   `json:"customer_details,omitempty" validate:"omitempty"`
	SellerDetail      *SellerDetail     `json:"seller_details,omitempty" validate:"omitempty"`
	GoPay             *GoPay            `json:"gopay,omitempty" validate:"omitempty"`
	CreditCard        *CreditCard       `json:"credit_card,omitempty" validate:"omitempty"`
	MetaData          map[string]string `json:"meta_data,omitempty"`
	QRIS              *QRIS             `json:"qris,omitempty" validate:"omitempty"`
	ShopeePay         *ShopeePay        `json:"shopeepay,omitempty" validate:"omitempty"`
}

// CUSTOMER DETAIL
// REF https://docs.midtrans.com/reference/customer-details-object
type CustomerDetail struct {
	FirstName       string                   `json:"first_name,omitempty" validate:"omitempty,max=30"`
	LastName        string                   `json:"last_name,omitempty" validate:"omitempty,max=30"`
	Email           string                   `json:"email,omitempty" validate:"omitempty,email"`
	Phone           string                   `json:"phone,omitempty" validate:"omitempty,max=255"`
	BillingAddress  *CustomerBillingAddress  `json:"billing_address,omitempty" validate:"omitempty"`
	ShippingAddress *CustomerShippingAddress `json:"shipping_address,omitempty" validate:"omitempty"`
}
type CustomerBillingAddress struct {
	FirstName   string `json:"first_name,omitempty" validate:"omitempty,max=255"`
	LastName    string `json:"last_name,omitempty" validate:"omitempty,max=255"`
	Phone       string `json:"phone,omitempty" validate:"omitempty,max=255"`
	Address     string `json:"address,omitempty" validate:"omitempty,max=255"`
	City        string `json:"city,omitempty" validate:"omitempty,max=255"`
	PostalCode  string `json:"postal_code,omitempty" validate:"omitempty,max=255,alphanumunicode"`
	CountryCode string `json:"country_code,omitempty" validate:"omitempty,eq=IDN,len=3"`
}
type CustomerShippingAddress struct {
	FirstName   string `json:"first_name,omitempty" validate:"omitempty,max=255"`
	LastName    string `json:"last_name,omitempty" validate:"omitempty,max=255"`
	Phone       string `json:"phone,omitempty" validate:"omitempty,max=255"`
	Address     string `json:"address,omitempty" validate:"omitempty,max=255"`
	City        string `json:"city,omitempty" validate:"omitempty,max=255"`
	PostalCode  string `json:"postal_code,omitempty" validate:"omitempty,max=255,alphanumunicode"`
	CountryCode string `json:"country_code,omitempty" validate:"omitempty,eq=IDN,len=3"`
}

// ITEM DETAIL
// REF https://docs.midtrans.com/reference/item-details-object
type ItemDetail struct {
	ID           string `json:"id,omitempty" validate:"omitempty"`
	Name         string `json:"name" validate:"required"`
	Price        int64  `json:"price" validate:"required,gt=0"`
	Qty          int32  `json:"quantity" validate:"required,gt=0"`
	Brand        string `json:"brand,omitempty" validate:"omitempty"`
	Category     string `json:"category,omitempty" validate:"omitempty"`
	MerchantName string `json:"merchant_name,omitempty" validate:"omitempty"`
	Tenor        int    `json:"tenor,omitempty" validate:"omitempty,numeric,len=2"`
	CodePlan     int    `json:"code_plan,omitempty" validate:"omitempty,numeric,len=3"`
	MID          int    `json:"mid,omitempty" validate:"omitempty,numeric,len=9"`
	URL          string `json:"url,omitempty" validate:"omitempty,url"`
}

// TRANSACTION DETAIL
// REF https://docs.midtrans.com/reference/transaction-details-object
type TransactionDetail struct {
	OrderID     string `json:"order_id" validate:"required,max=50"` // Order ID, required
	GrossAmount int64  `json:"gross_amount" validate:"required"`    // Total transaction amount, required
}

// CUSTOM EXPIRY
// REF https://docs.midtrans.com/reference/custom-expiry-object
type CustomExpiry struct {
	//OrderTime      time.Time `json:"order_time,omitempty"`
	ExpiryDuration int32  `json:"expiry_duration,omitempty"`
	Unit           string `json:"unit,omitempty"`
}

// SELLER DETAIL
// REF https://docs.midtrans.com/reference/seller-details-object
type SellerDetail struct {
	ID      string         `json:"id,omitempty"`
	Name    string         `json:"name,omitempty"`
	Email   string         `json:"email,omitempty"`
	URL     string         `json:"url,omitempty"`
	Address *SellerAddress `json:"address,omitempty"`
}
type SellerAddress struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}

// GOPAY
// REF https://docs.midtrans.com/reference/gopay-object
type GoPay struct {
	EnableCallback     bool     `json:"enable_callback,omitempty" validate:"omitempty"`
	CallbackURL        string   `json:"callback_url,omitempty" validate:"omitempty,url"`
	AccountID          string   `json:"account_id,omitempty" validate:"omitempty,uuid4"`
	PaymentOptionToken string   `json:"payment_option_token,omitempty" validate:"omitempty"`
	PreAuth            bool     `json:"pre_auth,omitempty" validate:"omitempty"`
	Recurring          bool     `json:"recurring,omitempty" validate:"omitempty"`
	PromotionIDs       []string `json:"promotion_ids,omitempty" validate:"omitempty,dive,required"`
}
type GoPayResponse struct {
	StatusCode             string           `json:"status_code"`
	StatusMessage          string           `json:"status_message"`
	TransactionID          string           `json:"transaction_id"`
	OrderID                string           `json:"order_id"`
	GrossAmount            string           `json:"gross_amount"`
	PaymentType            string           `json:"payment_type"`
	TransactionTime        string           `json:"transaction_time"`
	TransactionStatus      string           `json:"transaction_status"`
	Actions                []ActionResponse `json:"actions"`
	ChannelResponseCode    string           `json:"channel_response_code"`
	ChannelResponseMessage string           `json:"channel_response_message"`
	Currency               string           `json:"currency"`
	SignatureKey           string           `json:"signature_key"`

	// fill by code
	ActionGenerateQRCode   ActionResponse `json:"-"`
	ActionDeepLinkRedirect ActionResponse `json:"-"`
	ActionGetStatus        ActionResponse `json:"-"`
	ActionCancel           ActionResponse `json:"-"`
}

// CREDIT CARD
// REF https://docs.midtrans.com/reference/credit-card-object
type CreditCard struct {
	TokenID         string   `json:"token_id" validate:"required"`
	Bank            string   `json:"bank,omitempty" validate:"omitempty,oneof=mandiri bni cimb bca maybank bri"`
	InstallmentTerm int      `json:"installment_term,omitempty" validate:"omitempty,min=1"`
	Bins            []string `json:"bins,omitempty" validate:"omitempty,dive,len=4"`
	Type            string   `json:"type,omitempty" validate:"omitempty,oneof=authorize"`
	SaveTokenID     bool     `json:"save_token_id,omitempty" validate:"omitempty"`
	Channel         string   `json:"channel,omitempty" validate:"omitempty,oneof=dragon mti migs cybersource braintree mpgs"`
}

// SHOPEEPAY
// REF https://docs.midtrans.com/reference/shopeepay-object
type ShopeePay struct {
	CallbackURL string `json:"callback_url" validate:"required,url"`
}
type ShopeePayResponse struct {
	StatusCode             string           `json:"status_code"`
	StatusMessage          string           `json:"status_message"`
	ChannelResponseCode    string           `json:"channel_response_code"`
	ChannelResponseMessage string           `json:"channel_response_message"`
	TransactionID          string           `json:"transaction_id"`
	OrderID                string           `json:"order_id"`
	MerchantID             string           `json:"merchant_id"`
	GrossAmount            string           `json:"gross_amount"`
	Currency               string           `json:"currency"`
	PaymentType            string           `json:"payment_type"`
	TransactionTime        string           `json:"transaction_time"`
	TransactionStatus      string           `json:"transaction_status"`
	FraudStatus            string           `json:"fraud_status"`
	Actions                []ActionResponse `json:"actions"`

	ActionDeepLinkRedirect ActionResponse `json:"-"`
}

// QRIS
// REF https://docs.midtrans.com/reference/qris-object
type QRIS struct {
	Acquirer string `json:"acquirer,omitempty" validate:"omitempty,oneof=airpay shopee gopay"`
}
type ChargeQRISResponse struct {
	StatusCode        string           `json:"status_code"`
	StatusMessage     string           `json:"status_message"`
	TransactionID     string           `json:"transaction_id"`
	OrderID           string           `json:"order_id"`
	MerchantID        string           `json:"merchant_id"`
	GrossAmount       string           `json:"gross_amount"`
	Currency          string           `json:"currency"`
	PaymentType       string           `json:"payment_type"`
	TransactionTime   string           `json:"transaction_time"`
	TransactionStatus string           `json:"transaction_status"`
	FraudStatus       string           `json:"fraud_status"`
	Acquirer          string           `json:"acquirer"`
	Actions           []ActionResponse `json:"actions"`

	// fill by code
	ActionGenerateQRCode ActionResponse `json:"-"`
}

// VIRTUAL ACCOUNT
// REF https://docs.midtrans.com/reference/bank-transfer-object
type BankTransfer struct {
	Bank     string          `json:"bank" validate:"required,oneof=permata bni bri bca cimb"`
	VANumber string          `json:"va_number,omitempty" validate:"omitempty,min=6,max=18"`
	FreeText *FreeTextObject `json:"free_text,omitempty"`
	BCA      *BCAOptions     `json:"bca,omitempty" validate:"required_if=bank bca"`
	Permata  *PermataOptions `json:"permata,omitempty" validate:"required_if=bank permata"`
}
type FreeTextObject struct {
	Inquiry []InquiryPayment `json:"inquiry,omitempty" validate:"omitempty,dive,max=10"`
	Payment []InquiryPayment `json:"payment,omitempty" validate:"omitempty,dive,max=10"`
}
type InquiryPayment struct {
	ID string `json:"id" validate:"required,max=50"` // Free text message in Bahasa Indonesia
	EN string `json:"en" validate:"required,max=50"` // Free text message in English
}
type BCAOptions struct {
	SubCompanyCode string `json:"sub_company_code,omitempty" validate:"max=5"` // Default is 00000
}
type PermataOptions struct {
	RecipientName string `json:"recipient_name,omitempty" validate:"max=20"` // Uppercase string
}

// ECHANNEL MANDIRI BILL
// REF https://docs.midtrans.com/reference/e-channel-object
type EChannel struct {
	BillInfo1 string `json:"bill_info1" validate:"required,max=10"`  // Label 1, max 10 characters
	BillInfo2 string `json:"bill_info2" validate:"required,max=30"`  // Value for Label 1, max 30 characters
	BillInfo3 string `json:"bill_info3,omitempty" validate:"max=10"` // Label 2, max 10 characters, optional
	BillInfo4 string `json:"bill_info4,omitempty" validate:"max=30"` // Value for Label 2, max 30 characters, optional
	BillInfo5 string `json:"bill_info5,omitempty" validate:"max=10"` // Label 3, max 10 characters, optional
	BillInfo6 string `json:"bill_info6,omitempty" validate:"max=30"` // Value for Label 3, max 30 characters, optional
	BillInfo7 string `json:"bill_info7,omitempty" validate:"max=10"` // Label 4, max 10 characters, optional
	BillInfo8 string `json:"bill_info8,omitempty" validate:"max=30"` // Value for Label 4, max 30 characters, optional
	BillKey   string `json:"bill_key,omitempty" validate:"max=12"`   // Custom bill key, max 12 characters, optional
}

type ActionResponse struct {
	Name   string   `json:"name"`
	Method string   `json:"method"`
	URL    string   `json:"url"`
	Fields []string `json:"fields,omitempty"`
}

type ErrorBadReqResponse struct {
	HttpStatus         int      `json:"-"`
	StatusCode         string   `json:"status_code"`
	StatusMessage      string   `json:"status_message"`
	ValidationMessages []string `json:"validation_messages,omitempty"`
	ID                 string   `json:"id,omitempty"`
}
