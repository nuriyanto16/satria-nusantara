# Perbaikan Panel Admin & Fitur Versioning

Rencana implementasi untuk memperbaiki permasalahan pada Kelola Cabang (CRUD dan Maps) serta penambahan fitur Kelola Versi Aplikasi.

## User Review Required

> [!IMPORTANT]
> Pada penambahan versi, apakah Anda ingin sistem auto-increment untuk build number secara otomatis ketika admin menekan tombol "Tambah Versi", atau admin bisa menginput nomor versi secara manual (misal 1.0.1)? (Default pada plan ini: Admin bisa input teks versi dan build number).

## Open Questions

- Apakah *endpoint* untuk pengecekan versi ini nantinya akan dipanggil secara publik oleh aplikasi mobile tanpa perlu *token authentication*? (Saya akan set endpoint check-version menjadi publik).

## Proposed Changes

---

### Database Backend (PostgreSQL)

#### [NEW] `infra/sql/init.sql` (Update)
- Menambahkan tabel `app_versions` untuk menyimpan daftar versi aplikasi.
- Field yang disarankan: `id`, `version_name` (misal 1.0.1), `build_number` (misal 10), `release_notes`, `is_mandatory` (boolean), `created_at`.

### Backend Go (API)

#### [MODIFY] `backend/internal/organization/handler.go`, `service.go`, `repository.go`
- Menambahkan dan melengkapi *handler* CRUD untuk `unit_latihan` dan `pelatih` yang saat ini belum ada (seperti `updateUnit`, `deleteUnit`, `createPelatih`, `updatePelatih`, `deletePelatih`).
- Memperbaiki Payload `CreatePengurusRequest` jika terdapat *mismatch* dengan data frontend.

#### [NEW] `backend/internal/version/`
- Membuat *module/package* baru untuk manajemen versi aplikasi.
- Endpoint: 
  - `GET /versions` (Admin list)
  - `POST /versions` (Admin tambah versi)
  - `GET /public/check-version` (Public API untuk mobile ngecek update terbaru).

### Frontend Vue (Admin Panel)

#### [MODIFY] `frontend/pages/cabang.vue`
- Memperbaiki metode `saveUnit` dan `saveTrainer` yang saat ini hanya berupa *mock* (menyimpan di *state* lokal) menjadi request HTTP ke API backend yang sebenarnya.
- Mengintegrasikan *library* `leaflet` dengan menambahkan tag script / css Nuxt Config atau import NPM agar peta bisa tampil dan berinteraksi saat *add/edit form* lokasi.
- Menyesuaikan _payload_ `pengurusForm` dengan *request body* yang dibutuhkan backend saat disubmit.

#### [NEW] `frontend/pages/versi.vue` (atau masuk di pengaturan)
- Membuat halaman UI / komponen baru untuk **Kelola Versi Aplikasi**.
- Fitur CRUD sederhana untuk mencatat rilis versi, build number, dan tombol untuk menandai _update_ bersifat wajib (Mandatory Update).
- Menambahkannya pada sidebar navigasi Admin jika belum ada.

## Verification Plan

### Automated Tests
- Menjalankan `go build` untuk backend memastikan tidak ada error _syntax_ / _type mismatch_.
- Menjalankan kompilasi frontend Nuxt untuk memverifikasi _type-checking_.

### Manual Verification
- Melakukan percobaan _add_, _edit_, _delete_ Pengurus, Unit Latihan, dan Pelatih di menu Kelola Cabang.
- Membuka form yang mengandung peta dan memastikan _pin/marker_ Peta bisa digeser (_drag_) dan koordinatnya memperbarui nilai form.
- Membuka halaman Kelola Versi, menambahkan versi baru (misal 1.0.1), dan mengecek Endpoint GET `/public/check-version` menggunakan _curl_ atau _browser_.
