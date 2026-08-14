package finance

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ─── XENDIT CONFIG ────────────────────────────────────────────────────────────

func xenditSecretKey() string {
	k := os.Getenv("XENDIT_SECRET_KEY")
	if k == "" {
		// fallback: use public key from env or hardcoded dev key
		k = os.Getenv("XENDIT_API_KEY")
	}
	if k == "" {
		k = "xnd_public_development_7wB51nwFxZasCXmyAe4x_bxA3vBeqqQfKvGCFdFHeVaPha3LgQkRDRmGoA8GqrOp"
	}
	return k
}

func xenditWebhookToken() string {
	t := os.Getenv("XENDIT_WEBHOOK_TOKEN")
	if t == "" {
		t = "tA4z5aeBchmLxkkTV5d0Hynwcbn2ovuyD1nHSDylhbbGgrT5"
	}
	return t
}

const xenditBaseURL = "https://api.xendit.co"

// ─── XENDIT REQUEST / RESPONSE TYPES ─────────────────────────────────────────

// CreateInvoiceRequest mirrors Xendit Invoice v2 payload
type CreateInvoiceRequest struct {
	ExternalID         string   `json:"external_id"`
	Amount             int      `json:"amount"`
	Description        string   `json:"description"`
	PayerEmail         string   `json:"payer_email,omitempty"`
	CustomerName       string   `json:"customer[given_names],omitempty"`
	Currency           string   `json:"currency"`
	InvoiceDuration    int      `json:"invoice_duration"`          // seconds, default 86400 (24h)
	SuccessRedirectURL string   `json:"success_redirect_url,omitempty"`
	FailureRedirectURL string   `json:"failure_redirect_url,omitempty"`
	PaymentMethods     []string `json:"payment_methods,omitempty"`
}

// XenditInvoice is the response from Xendit
type XenditInvoice struct {
	ID                 string    `json:"id"`
	ExternalID         string    `json:"external_id"`
	Status             string    `json:"status"` // PENDING, SETTLED, EXPIRED
	Amount             float64   `json:"amount"`
	Description        string    `json:"description"`
	InvoiceURL         string    `json:"invoice_url"`
	ExpiryDate         time.Time `json:"expiry_date"`
	Currency           string    `json:"currency"`
	MerchantName       string    `json:"merchant_name"`
	PaymentMethod      string    `json:"payment_method,omitempty"`
	PaymentChannel     string    `json:"payment_channel,omitempty"`
	PaidAt             *string   `json:"paid_at,omitempty"`
	MidLabel           string    `json:"mid_label,omitempty"`
	PaymentDestination string    `json:"payment_destination,omitempty"`
}

// XenditWebhookPayload represents the Xendit webhook body
type XenditWebhookPayload struct {
	// Invoice webhook fields
	ID             string  `json:"id"`
	ExternalID     string  `json:"external_id"`
	UserID         string  `json:"user_id"`
	IsHigh         bool    `json:"is_high"`
	Status         string  `json:"status"`
	MerchantName   string  `json:"merchant_name"`
	Amount         float64 `json:"amount"`
	PayerEmail     string  `json:"payer_email"`
	Description    string  `json:"description"`
	PaidAmount     float64 `json:"paid_amount"`
	Updated        string  `json:"updated"`
	Created        string  `json:"created"`
	Currency       string  `json:"currency"`
	PaidAt         string  `json:"paid_at"`
	PaymentMethod  string  `json:"payment_method"`
	PaymentChannel string  `json:"payment_channel"`
	PaymentDestination string `json:"payment_destination"`

	// eWallet/payment request fields
	ReferenceID     string `json:"reference_id"`
	BusinessID      string `json:"business_id"`
	ChannelCode     string `json:"channel_code"`
	Type            string `json:"type"` // for payment_request webhooks
}

// ─── XENDIT SERVICE ───────────────────────────────────────────────────────────

// XenditService handles all Xendit API calls
type XenditService struct {
	secretKey    string
	webhookToken string
	httpClient   *http.Client
}

func NewXenditService() *XenditService {
	return &XenditService{
		secretKey:    xenditSecretKey(),
		webhookToken: xenditWebhookToken(),
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

// basicAuthHeader encodes the key for Basic Auth (Xendit uses key:empty-password)
func (x *XenditService) basicAuthHeader() string {
	raw := x.secretKey + ":"
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}

// CreateInvoice calls Xendit to create a payment invoice
func (x *XenditService) CreateInvoice(req CreateInvoiceRequest) (*XenditInvoice, error) {
	if req.Currency == "" {
		req.Currency = "IDR"
	}
	if req.InvoiceDuration == 0 {
		req.InvoiceDuration = 86400 // 24 hours
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	httpReq, err := http.NewRequest("POST", xenditBaseURL+"/v2/invoices", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", x.basicAuthHeader())

	resp, err := x.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("xendit error [%d]: %s", resp.StatusCode, string(respBody))
	}

	var invoice XenditInvoice
	if err := json.Unmarshal(respBody, &invoice); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &invoice, nil
}

// VerifyWebhookToken validates the incoming Xendit webhook token
func (x *XenditService) VerifyWebhookToken(r *http.Request) bool {
	incoming := r.Header.Get("x-callback-token")
	return incoming == x.webhookToken
}

// VerifyWebhookSignature validates HMAC-SHA256 signature (for newer webhook versions)
func (x *XenditService) VerifyWebhookSignature(r *http.Request, body []byte) bool {
	sig := r.Header.Get("X-Webhook-Id")
	if sig == "" {
		// Fall back to token verification
		return x.VerifyWebhookToken(r)
	}
	mac := hmac.New(sha256.New, []byte(x.webhookToken))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	incoming := r.Header.Get("X-Signature")
	return hmac.Equal([]byte(incoming), []byte(expected))
}
