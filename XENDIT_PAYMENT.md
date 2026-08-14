# 💳 Xendit Payment Service — Satria Nusantara

> Dokumentasi implementasi Payment Gateway Xendit untuk aplikasi **Satria Nusantara**  
> Mode: **Development / Test Mode** · Dibuat: 2026-08-12

---

## 1. Konfigurasi & Credentials

| Parameter | Value |
|---|---|
| **Public API Key** | `xnd_public_development_7wB51nwFxZasCXmyAe4x_bxA3vBeqqQfKvGCFdFHeVaPha3LgQkRDRmGoA8GqrOp` |
| **Webhook Token** | `tA4z5aeBchmLxkkTV5d0Hynwcbn2ovuyD1nHSDylhbbGgrT5` |
| **Allowed IP** | `43.133.150.102` |
| **Mode** | `TEST / Development` |

> ⚠️ Key ini adalah **Public Development Key**. Untuk produksi, ganti dengan **Secret Key** dari Xendit Dashboard → Settings → API Keys.

### Environment Variables Backend (docker-compose.yml / VPS)

```env
XENDIT_API_KEY=xnd_public_development_7wB51nwFxZasCXmyAe4x_bxA3vBeqqQfKvGCFdFHeVaPha3LgQkRDRmGoA8GqrOp
XENDIT_SECRET_KEY=<ganti_dengan_secret_key_xendit>
XENDIT_WEBHOOK_TOKEN=tA4z5aeBchmLxkkTV5d0Hynwcbn2ovuyD1nHSDylhbbGgrT5
```

---

## 2. Webhook URLs — Xendit Dashboard

Masuk ke **Xendit Dashboard → Settings → Developers → Webhook** dan isi URL berikut:

### INVOICES
| Event | URL |
|---|---|
| Invoices paid | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/invoice` |

> ✅ Centang "Also notify my application when a payment has been received after expiry"

### E-WALLETS
| Event | URL |
|---|---|
| eWallet Payment Status | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/ewallet` |
| eWallet Reconciliation Update | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/ewallet` |

### PAYMENT REQUESTS V3 (v3/payment_requests)
| Event | URL |
|---|---|
| Payment Status | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/payment-request` |
| Payment Request Status | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/payment-request` |
| Payment Verified | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/payment-request` |

### PAYMENT REQUESTS V2 (v2/payment_requests)
| Event | URL |
|---|---|
| Payment Succeeded | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/payment-request` |
| Payment Awaiting Capture | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/payment-request` |
| Payment Failed | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/payment-request` |
| Captured Succeeded | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/payment-request` |

### QR CODES
| Event | URL |
|---|---|
| QR Code paid & refunded | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/payment-request` |

### PAYMENT TOKENS V3 (v3/payment_tokens)
| Event | URL |
|---|---|
| Payment Token Status | `https://nfmtech.my.id/product/satrianusantara/api/v1/finance/webhook/xendit/payment-request` |

> 📝 Kategori **Disbursement, Cards, Direct Debit, XenPlatform, Report** tidak digunakan dan dapat dikosongkan.

---

## 3. API Endpoints

**Base URL:** `https://nfmtech.my.id/product/satrianusantara/api/v1`

### POST `/finance/iuran/xendit/create-invoice`

Membuat invoice Xendit untuk pembayaran BLBA.

**Request:**
```json
{
  "userId": "user-id-anggota",
  "nama": "Budi Santoso",
  "email": "budi@example.com",
  "bulan": "Agustus 2026",
  "amount": 40000,
  "transactionId": "iuran-id-internal"
}
```

**Response (200):**
```json
{
  "status": "success",
  "message": "Invoice berhasil dibuat",
  "data": {
    "invoiceUrl": "https://checkout.xendit.co/web/xyz123",
    "xenditId": "645e6e16bfbf38...",
    "externalId": "SN-BLBA-user123-1723478400000",
    "amount": 40000,
    "expiry": "2026-08-13T16:00:00Z",
    "status": "PENDING",
    "trxId": "pay_20260812163000"
  }
}
```

### POST `/finance/webhook/xendit/invoice`
Callback dari Xendit saat invoice dibayar. Header wajib: `x-callback-token`.

**Payload Xendit (SETTLED):**
```json
{
  "id": "645e6e16bfbf38ade9a2b3c4",
  "external_id": "SN-BLBA-user123-1723478400000",
  "status": "SETTLED",
  "amount": 40000,
  "payment_channel": "BCA",
  "paid_at": "2026-08-12T09:15:00Z"
}
```

### POST `/finance/webhook/xendit/ewallet`
Callback untuk eWallet (OVO, DANA, GoPay). Status `SUCCEEDED` → transaksi lunas.

### POST `/finance/webhook/xendit/payment-request`
Callback untuk Payment Request V2/V3 (QRIS, VA). Status `SUCCEEDED` → transaksi lunas.

---

## 4. Flow Pembayaran

```
User Mobile / Admin
    │
    ▼
POST create-invoice → Backend → Xendit API
                                    │
                            Return invoiceUrl
                                    │
                     Open invoiceUrl di browser
                                    │
                         User pilih metode bayar
                                    │
                         Xendit proses pembayaran
                                    │
                    Xendit POST webhook → Backend
                    Header: x-callback-token: <token>
                                    │
                    Backend update status → "lunas"
                                    │
                         App refresh → ✅ Lunas
```

---

## 5. File yang Diimplementasi

| File | Keterangan |
|---|---|
| `backend/internal/finance/xendit.go` | **[BARU]** Xendit service: CreateInvoice, VerifyWebhookToken |
| `backend/internal/finance/handler.go` | Tambah 4 route webhook + field XenditID/InvoiceURL |
| `frontend/pages/iuran.vue` | Tombol "Bayar via Xendit" + Modal create invoice |
| `mobile/lib/presentation/screens.dart` | Kartu XENDIT di payment screen + launch browser |
| `mobile/pubspec.yaml` | Tambah dependency `url_launcher: ^6.2.5` |

---

## 6. Metode Pembayaran

| Kategori | Provider |
|---|---|
| eWallet | OVO, DANA, GoPay, ShopeePay, LinkAja |
| Virtual Account | BCA, BNI, BRI, Mandiri, Permata |
| Gerai | Alfamart, Indomaret |
| QRIS | Semua bank |
| Kartu | Visa/Mastercard (perlu aktivasi) |

---

## 7. Checklist Produksi

### Xendit Dashboard
- [ ] Upgrade ke Live Mode (Verify Business)
- [ ] Ganti API Key ke Secret Key Live (`xnd_production_*`)
- [ ] Set semua webhook URL (lihat bagian 2)
- [ ] Klik "Test and Save" untuk verifikasi setiap URL
- [ ] Verifikasi IP `43.133.150.102` di Allowed IPs

### Backend
- [ ] Set env `XENDIT_SECRET_KEY` dengan live secret key
- [ ] Pindah transaksi dari in-memory ke database PostgreSQL
- [ ] Tambah idempotency check di webhook (cegah double-update)
- [ ] Tambah logging webhook

### Mobile
- [ ] `flutter pub get` (tambah url_launcher)
- [ ] Tambah permission di AndroidManifest.xml:
  ```xml
  <queries>
    <intent>
      <action android:name="android.intent.action.VIEW" />
      <data android:scheme="https" />
    </intent>
  </queries>
  ```
- [ ] Build ulang APK release

---

## 8. Referensi

- Xendit API Docs: https://docs.xendit.co
- Invoice API: https://docs.xendit.co/apidocs/invoice
- Webhook Behavior: https://docs.xendit.co/apidocs/webhook-behavior
- eWallet API: https://docs.xendit.co/apidocs/ewallet
- Xendit Dashboard: https://dashboard.xendit.co

---
*Dokumen dibuat otomatis: 2026-08-12 oleh Antigravity AI*
