# Satria Nusantara — Sistem Informasi Manajemen Anggota

> Sistem informasi terintegrasi berbasis Microservices untuk LSP Satria Nusantara — organisasi seni pernapasan nasional Indonesia.

---

## 🗂️ Struktur Proyek

```
.
├── backend/                    # Golang Microservices
│   ├── api-gateway/            # API Gateway (port 8080)
│   ├── auth-service/           # Auth & JWT Service (port 8081)
│   ├── organization-service/   # Cabang, Unit, Anggota (port 8082)
│   ├── training-service/       # Jadwal & Absensi QR (port 8083)
│   ├── event-service/          # Latgab, EKT, Pelatnas (port 8084)
│   ├── finance-service/        # BLBA & Midtrans (port 8085)
│   ├── content-service/        # Konten & FCM Notif (port 8086)
│   └── kebugaran-service/      # Tes Kebugaran (port 8087)
│
├── frontend/                   # Nuxt.js 3 Web Admin
│   ├── assets/css/main.css     # Design System CSS
│   ├── layouts/default.vue     # Sidebar + Topbar Layout
│   └── pages/
│       ├── index.vue           # Dashboard
│       ├── login.vue           # Login Page
│       ├── cabang.vue          # Kelola Cabang
│       ├── anggota.vue         # Kelola Anggota
│       ├── jadwal.vue          # Jadwal Latihan
│       ├── absensi.vue         # Absensi & Kehadiran
│       ├── latgab.vue          # Latgab / EKT / Pelatnas
│       ├── kebugaran.vue       # Tes Kebugaran
│       ├── iuran.vue           # BLBA (Keuangan)
│       ├── laporan.vue         # Laporan & Export
│       ├── konten.vue          # Konten & Artikel
│       ├── user.vue            # Manajemen User
│       └── pengaturan.vue      # Pengaturan Sistem
│
├── mobile/                     # Flutter App (Android + iOS)
│   ├── lib/main.dart
│   └── pubspec.yaml
│
├── infra/
│   └── sql/init.sql            # Skema PostgreSQL lengkap
│
└── docker-compose.yml          # Orkestrasi semua service
```

---

## 🚀 Cara Menjalankan (Docker)

### Prerequisites
- Docker Desktop terinstall dan berjalan
- Port 8080, 3000, 5432, 5050, 6379 tersedia

### 1. Clone & Setup
```bash
git clone <repo>
cd "SATRIA NUSANTARA"
```

### 2. Konfigurasi Environment
Salin `.env.example` menjadi `.env` dan isi kredensial:
```bash
# Kredensial Midtrans (sandbox untuk dev)
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxx

# Firebase FCM
FCM_SERVER_KEY=your-firebase-server-key

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
```

### 3. Jalankan Infrastructure (Database dulu)
```bash
docker-compose up -d postgres redis pgadmin
```
> Tunggu ~10 detik agar PostgreSQL siap, lalu skema `init.sql` akan otomatis dijalankan.

### 4. Jalankan Semua Service
```bash
docker-compose up -d
```

### 5. Akses Aplikasi
| Layanan        | URL                           |
|----------------|-------------------------------|
| Web Admin      | http://localhost:3000         |
| API Gateway    | http://localhost:8080         |
| pgAdmin        | http://localhost:5050         |
| Auth Service   | http://localhost:8081/health  |

### 6. Login pgAdmin
- https://troops-gui-ink-listing.trycloudflare.com
- Email: `admin@satria-nusantara.org`
- Password: `admin2026`

---

## 🧑‍💻 Pengembangan Lokal (Tanpa Docker)

### Backend — Golang
```bash
# Install Go 1.21+ dahulu: https://go.dev/dl/
cd backend/auth-service
go mod tidy
go run main.go
```

### Frontend — Nuxt.js
```bash
# Install Node.js 20+ dahulu
cd frontend
npm install
npm run dev
# Buka: http://localhost:3000
```

### Mobile — Flutter
```bash
# Install Flutter SDK dahulu: https://flutter.dev/docs/get-started/install
cd mobile
flutter pub get
flutter run
```

---

## 🏗️ Tech Stack

| Layer       | Teknologi                                           |
|-------------|-----------------------------------------------------|
| Backend     | Go 1.21, go-chi, JWT, bcrypt                        |
| Frontend    | Nuxt 3 / Vue 3, Pinia, Tabler Icons                 |
| Mobile      | Flutter / Dart, BLoC/Riverpod, Dio                  |
| Database    | PostgreSQL 16 + Redis 7                             |
| DevOps      | Docker, Docker Compose, multi-stage builds          |
| Payment     | Midtrans (DANA, OVO, GoPay, ShopeePay)              |
| Notifikasi  | Firebase Cloud Messaging (FCM)                      |
| Auth Mobile | Google OAuth 2.0 + JWT                              |

---

## 📋 Microservices Port Map

| Service              | Port | Fungsi                              |
|----------------------|------|-------------------------------------|
| API Gateway          | 8080 | Routing & rate limiting             |
| Auth Service         | 8081 | Login, JWT, Role management         |
| Organization Service | 8082 | Cabang, Unit, Anggota               |
| Training Service     | 8083 | Jadwal & Absensi QR                 |
| Event Service        | 8084 | Latgab, EKT Jurus, Pelatnas         |
| Finance Service      | 8085 | BLBA & Midtrans                     |
| Content Service      | 8086 | Artikel & FCM Notifications         |
| Kebugaran Service    | 8087 | Tes Kebugaran 5 jenis               |
| Web Admin (Nuxt)     | 3000 | Frontend Web                        |
| PostgreSQL           | 5432 | Database utama                      |
| Redis                | 6379 | Cache & token management            |
| pgAdmin              | 5050 | Database management GUI             |

---

## 📁 Referensi SRS
Implementasi mengacu pada **SRS v5 (Juni 2026)** yang mencakup:
- 32 layar aplikasi mobile
- 12 halaman web admin
- 9 tingkatan keanggotaan (Pra Dasar → Meditasi)
- 13 jenis notifikasi (N-01 s/d N-13)
- 48+ parameter konfigurasi sistem
- Integrasi Midtrans, FCM, Google OAuth, WhatsApp Deep Link

---

*Dibuat untuk LSP Satria Nusantara — 2026*




