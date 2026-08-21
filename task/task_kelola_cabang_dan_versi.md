# Task List: Admin Panel CRUD & App Versioning

- [x] **Database Setup**
  - [x] Tambahkan tabel `app_versions` ke schema database (`infra/sql/init.sql`).
  - [x] Eksekusi DDL `app_versions` di Postgres.
- [x] **Backend (API)**
  - [x] Lengkapi CRUD API untuk `Unit Latihan` (`updateUnit`, `deleteUnit`).
  - [x] Buat CRUD API lengkap untuk `Pelatih` (`createPelatih`, `updatePelatih`, `deletePelatih`).
  - [x] Perbaiki struct/binding untuk payload `CreatePengurusRequest` (tambahkan nomor jika perlu).
  - [x] Buat modul baru `version` dengan endpoints: `GET /versions`, `POST /versions`, `GET /public/check-version`.
- [x] **Frontend (Admin Panel)**
  - [x] Perbaiki konfigurasi Leaflet Map agar marker/pin bisa berjalan dengan baik.
  - [x] Hubungkan `saveUnit` dan `deleteUnit` di `cabang.vue` ke API backend.
  - [x] Hubungkan `saveTrainer` dan `deleteTrainer` di `cabang.vue` ke API backend.
  - [x] Perbaiki pemetaan form `pengurusForm` ke API `savePengurus`.
  - [x] Buat halaman `versi.vue` (atau tambahkan di pengaturan) untuk manajemen *App Version*.
- [x] **Verifikasi**
  - [x] Uji coba Create/Update/Delete Unit dan Pelatih di halaman Kelola Cabang.
  - [x] Uji fungsi marker peta saat tambah Cabang/Unit.
  - [x] Uji API `/public/check-version` dan fungsionalitas halaman Kelola Versi.
