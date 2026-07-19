-- ============================================================
-- DATABASE SCHEMA: Satria Nusantara
-- Generated from SRS v5 & Mockup SN_WebAdmin_v8
-- ============================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================
-- SCHEMA: AUTENTIKASI & PENGGUNA
-- ============================================================

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,  -- Admin Pusat, Pengurus Cabang, PIC Unit, Pelatih
    scope VARCHAR(30) NOT NULL          -- pusat, cabang, unit, pelatih
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),         -- NULL jika pakai Google OAuth
    google_id VARCHAR(255) UNIQUE,
    role_id INT REFERENCES roles(id),
    nama_lengkap VARCHAR(255) NOT NULL,
    no_hp VARCHAR(20),
    foto_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'aktif', -- aktif, nonaktif
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    device_info TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- SCHEMA: ORGANISASI
-- ============================================================

CREATE TABLE provinsi (
    id SERIAL PRIMARY KEY,
    kode VARCHAR(5) NOT NULL UNIQUE,    -- JO, JB, JT, dst
    nama VARCHAR(100) NOT NULL
);

CREATE TABLE cabang (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kode VARCHAR(5) NOT NULL UNIQUE,    -- YGY, BDG, JKT, SBY
    nama VARCHAR(255) NOT NULL,
    provinsi_id INT REFERENCES provinsi(id),
    kota_kab VARCHAR(100) NOT NULL,
    alamat TEXT,
    status VARCHAR(20) DEFAULT 'aktif',
    dibuat_oleh UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE pengurus_cabang (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cabang_id UUID REFERENCES cabang(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    jabatan VARCHAR(100) NOT NULL,      -- Ketua, Sekretaris, Bendahara, Kabid Latihan, dst
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE unit_latihan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cabang_id UUID REFERENCES cabang(id) ON DELETE CASCADE,
    nama VARCHAR(255) NOT NULL,
    lokasi_nama VARCHAR(255),
    lokasi_alamat TEXT,
    jadwal_tetap JSONB,                 -- {"hari": "Selasa,Jumat", "jam": "07:00"}
    pic_user_id UUID REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'aktif',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE rekening_unit (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unit_id UUID REFERENCES unit_latihan(id) ON DELETE CASCADE,
    bank VARCHAR(100) NOT NULL,
    nomor_rekening VARCHAR(50) NOT NULL,
    nama_pemilik VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- SCHEMA: KEANGGOTAAN
-- ============================================================

CREATE TYPE tingkatan_enum AS ENUM (
    'Pra Dasar', 'Dasar', 'PH', 'Gabungan', 'PK', 'GPK', 'Penjuru', 'Selangkah', 'Meditasi'
);

CREATE TYPE status_anggota AS ENUM ('pending', 'aktif', 'nonaktif');

CREATE TABLE anggota (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    nomor_anggota VARCHAR(30) UNIQUE,   -- Format: JB-BDG-00142
    nama_lengkap VARCHAR(255) NOT NULL,
    tanggal_lahir DATE,
    jenis_kelamin VARCHAR(10),
    no_hp VARCHAR(20),
    foto_url VARCHAR(500),
    unit_id UUID REFERENCES unit_latihan(id),
    tingkatan tingkatan_enum DEFAULT 'Pra Dasar',
    jurus_saat_ini INT DEFAULT 1,
    counter_kehadiran INT DEFAULT 0,    -- Reset tiap naik jurus
    status status_anggota DEFAULT 'pending',
    tanggal_daftar TIMESTAMPTZ DEFAULT NOW(),
    tanggal_aktif TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE riwayat_kenaikan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    anggota_id UUID REFERENCES anggota(id) ON DELETE CASCADE,
    tingkatan_lama tingkatan_enum,
    jurus_lama INT,
    tingkatan_baru tingkatan_enum,
    jurus_baru INT,
    jenis VARCHAR(30),                  -- EKT_Jurus, EKT_Non_Jurus
    event_id UUID,                      -- FK ke tabel events
    tanggal TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- SCHEMA: PELATIH
-- ============================================================

CREATE TYPE kategori_transport AS ENUM ('Jarak Dekat', 'Jarak Sedang', 'Jarak Jauh', 'Sukarela');

CREATE TABLE pelatih (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    anggota_id UUID REFERENCES anggota(id) ON DELETE CASCADE,
    cabang_id UUID REFERENCES cabang(id),
    jenis VARCHAR(20) DEFAULT 'cabang',  -- cabang, pusat
    kategori_transport kategori_transport DEFAULT 'Jarak Sedang',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE pelatih_unit (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pelatih_id UUID REFERENCES pelatih(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES unit_latihan(id) ON DELETE CASCADE,
    is_utama BOOLEAN DEFAULT FALSE,
    UNIQUE(pelatih_id, unit_id)
);

-- ============================================================
-- SCHEMA: JADWAL & ABSENSI
-- ============================================================

CREATE TYPE jenis_sesi AS ENUM (
    'Latihan Rutin', 'Latihan Khusus', 'Latihan Pelatih',
    'EKT Jurus', 'EKT Non Jurus', 'Latgab', 'Pelatnas'
);

CREATE TABLE sesi_latihan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unit_id UUID REFERENCES unit_latihan(id),
    pelatih_id UUID REFERENCES pelatih(id),
    jenis jenis_sesi NOT NULL,
    tanggal DATE NOT NULL,
    jam_mulai TIME NOT NULL,
    jam_selesai TIME,
    lokasi VARCHAR(255),
    materi TEXT,
    event_id UUID,                      -- FK ke events (jika EKT/Latgab/Pelatnas)
    dibuat_via VARCHAR(20) DEFAULT 'web', -- web, mobile
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE absensi_qr (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sesi_id UUID REFERENCES sesi_latihan(id) ON DELETE CASCADE,
    qr_token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE kehadiran (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sesi_id UUID REFERENCES sesi_latihan(id) ON DELETE CASCADE,
    anggota_id UUID REFERENCES anggota(id) ON DELETE CASCADE,
    waktu_scan TIMESTAMPTZ DEFAULT NOW(),
    metode VARCHAR(20) DEFAULT 'qr',   -- qr, manual
    UNIQUE(sesi_id, anggota_id)
);

-- ============================================================
-- SCHEMA: EVENT (Latgab, EKT, Pelatnas)
-- ============================================================

CREATE TYPE jenis_event AS ENUM ('Latgab', 'EKT Jurus', 'EKT Non Jurus', 'Pelatnas');
CREATE TYPE scope_event AS ENUM ('nasional', 'provinsi', 'multi_cabang', 'cabang');
CREATE TYPE jenis_ekt_non_jurus AS ENUM ('Pengayaan I', 'Pengayaan II', 'Pendalaman Wawasan I');

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    jenis jenis_event NOT NULL,
    nama VARCHAR(255) NOT NULL,
    scope scope_event DEFAULT 'cabang',
    penyelenggara_cabang_id UUID REFERENCES cabang(id),
    tanggal DATE NOT NULL,
    jam_mulai TIME,
    jam_selesai TIME,
    lokasi VARCHAR(255),
    alamat TEXT,
    batas_pendaftaran TIMESTAMPTZ,
    -- Biaya
    biaya_latgab INT DEFAULT 0,
    biaya_ekt_dasar INT DEFAULT 0,
    biaya_ekt_ph INT DEFAULT 0,
    biaya_ekt_gabungan INT DEFAULT 0,
    biaya_ekt_pk INT DEFAULT 0,
    biaya_ekt_gpk INT DEFAULT 0,
    biaya_pelatnas_pokok INT DEFAULT 0,
    -- Config EKT Non Jurus
    jenis_ekt_non_jurus jenis_ekt_non_jurus,
    -- Rekening tujuan
    rekening_bank VARCHAR(100),
    rekening_nomor VARCHAR(50),
    rekening_nama VARCHAR(255),
    -- Meta
    penerima_notifikasi TEXT[],
    status VARCHAR(20) DEFAULT 'aktif',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE event_cabang_peserta (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    cabang_id UUID REFERENCES cabang(id),
    biaya_transport INT DEFAULT 0       -- Khusus Pelatnas: biaya per cabang
);

CREATE TABLE pendaftaran_event (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    anggota_id UUID REFERENCES anggota(id) ON DELETE CASCADE,
    pilihan VARCHAR(30) DEFAULT 'latgab', -- latgab, latgab_ekt, ekt_non_jurus
    -- Status pembayaran (hanya terdaftar jika LUNAS)
    status VARCHAR(20) DEFAULT 'terdaftar',
    midtrans_order_id VARCHAR(100),
    nominal_bayar INT DEFAULT 0,
    tanggal_daftar TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(event_id, anggota_id)
);

-- ============================================================
-- SCHEMA: TES KEBUGARAN
-- ============================================================

CREATE TABLE periode_tes_kebugaran (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tahun INT NOT NULL,
    periode INT NOT NULL,               -- 1, 2, 3 (default 3 per tahun)
    tanggal_mulai DATE,
    tanggal_selesai DATE,
    UNIQUE(tahun, periode)
);

CREATE TABLE sesi_tes_kebugaran (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unit_id UUID REFERENCES unit_latihan(id),
    periode_id UUID REFERENCES periode_tes_kebugaran(id),
    sesi_id UUID REFERENCES sesi_latihan(id), -- opsional, dikaitkan ke sesi latihan
    tanggal DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE hasil_tes_kebugaran (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sesi_tes_id UUID REFERENCES sesi_tes_kebugaran(id) ON DELETE CASCADE,
    anggota_id UUID REFERENCES anggota(id) ON DELETE CASCADE,
    nafas_dalam_air INT,                -- detik
    push_up INT,                        -- repetisi
    sit_up INT,                         -- repetisi
    sit_and_reach INT,                  -- cm
    shuttle_run DECIMAL(5,2),           -- detik
    -- Kategori otomatis (Baik/Cukup/Kurang)
    kategori_nafas VARCHAR(10),
    kategori_push_up VARCHAR(10),
    kategori_sit_up VARCHAR(10),
    kategori_sit_and_reach VARCHAR(10),
    kategori_shuttle_run VARCHAR(10),
    kategori_keseluruhan VARCHAR(10),
    -- Input oleh Pelatih
    pelatih_id UUID REFERENCES pelatih(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(sesi_tes_id, anggota_id)
);

-- ============================================================
-- SCHEMA: KEUANGAN (BLBA)
-- ============================================================

CREATE type status_blba AS ENUM ('belum_bayar', 'lunas', 'dispensasi');

CREATE TABLE tagihan_blba (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    anggota_id UUID REFERENCES anggota(id) ON DELETE CASCADE,
    bulan INT NOT NULL,                 -- 1-12
    tahun INT NOT NULL,
    nominal INT NOT NULL,
    status status_blba DEFAULT 'belum_bayar',
    midtrans_order_id VARCHAR(100),
    tanggal_bayar TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(anggota_id, bulan, tahun)
);

-- ============================================================
-- SCHEMA: KONTEN & NOTIFIKASI
-- ============================================================

CREATE TYPE jenis_konten AS ENUM ('Pengumuman', 'Artikel', 'Info Kegiatan');
CREATE TYPE scope_konten AS ENUM ('nasional', 'provinsi', 'cabang', 'unit');

CREATE TABLE konten (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    judul VARCHAR(500) NOT NULL,
    slug VARCHAR(500) UNIQUE,
    jenis jenis_konten DEFAULT 'Pengumuman',
    scope scope_konten DEFAULT 'cabang',
    cabang_id UUID REFERENCES cabang(id),
    isi TEXT NOT NULL,
    thumbnail_url VARCHAR(500),
    is_published BOOLEAN DEFAULT FALSE,
    event_id UUID,                      -- CTA terkait event (opsional)
    dibuat_oleh UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE notifikasi (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    jenis VARCHAR(30),                  -- N-01 hingga N-13 sesuai SRS
    judul VARCHAR(255) NOT NULL,
    pesan TEXT NOT NULL,
    target_user_id UUID REFERENCES users(id),
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE fcm_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) NOT NULL,
    platform VARCHAR(10),               -- android, ios
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, token)
);

-- ============================================================
-- SCHEMA: PENGATURAN SISTEM
-- ============================================================

CREATE TABLE pengaturan_sistem (
    id SERIAL PRIMARY KEY,
    kunci VARCHAR(100) NOT NULL UNIQUE,
    nilai TEXT NOT NULL,
    deskripsi TEXT,
    kategori VARCHAR(50),               -- blba, transport, ekt, kebugaran, profil, notifikasi, mobile
    updated_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- DATA AWAL (SEED)
-- ============================================================

INSERT INTO roles (name, scope) VALUES
    ('Admin Pusat', 'pusat'),
    ('Pengurus Cabang', 'cabang'),
    ('PIC Unit', 'unit'),
    ('Pelatih', 'pelatih');

INSERT INTO provinsi (kode, nama) VALUES
    ('JO', 'DI Yogyakarta'),
    ('JB', 'Jawa Barat'),
    ('JT', 'Jawa Tengah'),
    ('JI', 'Jawa Timur'),
    ('JK', 'DKI Jakarta'),
    ('BS', 'Banten'),
    ('SN', 'Sulawesi Selatan'),
    ('SS', 'Sumatera Selatan'),
    ('BA', 'Bali');

-- Pengaturan Sistem Default
INSERT INTO pengaturan_sistem (kunci, nilai, deskripsi, kategori) VALUES
    ('blba_nominal', '50000', 'Nominal BLBA per bulan (Rp)', 'blba'),
    ('blba_grace_period', '10', 'Toleransi hari pembayaran setelah jatuh tempo', 'blba'),
    ('transport_jarak_dekat', '30000', 'Nominal transport jarak dekat per sesi (Rp)', 'transport'),
    ('transport_jarak_sedang', '50000', 'Nominal transport jarak sedang per sesi (Rp)', 'transport'),
    ('transport_jarak_jauh', '75000', 'Nominal transport jarak jauh per sesi (Rp)', 'transport'),
    ('ekt_syarat_hadir_pra_dasar', '10', 'Min. kehadiran untuk EKT dari Pra Dasar', 'ekt'),
    ('ekt_syarat_hadir_dasar', '8', 'Min. kehadiran untuk EKT dari Dasar', 'ekt'),
    ('ekt_syarat_hadir_ph', '8', 'Min. kehadiran untuk EKT dari PH', 'ekt'),
    ('ekt_syarat_hadir_gabungan', '12', 'Min. kehadiran untuk EKT dari Gabungan', 'ekt'),
    ('ekt_syarat_hadir_pk', '18', 'Min. kehadiran untuk EKT dari PK', 'ekt'),
    ('ekt_syarat_hadir_gpk', '18', 'Min. kehadiran untuk EKT dari GPK', 'ekt'),
    ('ekt_syarat_hadir_penjuru', '80', 'Min. kehadiran untuk EKT dari Penjuru', 'ekt'),
    ('qr_berlaku_jam', '4', 'Durasi berlaku QR Absensi (jam)', 'mobile'),
    ('tes_kebugaran_periode_per_tahun', '3', 'Jumlah periode tes kebugaran per tahun', 'kebugaran'),
    ('nama_organisasi', 'Satria Nusantara', 'Nama resmi organisasi', 'profil'),
    ('nama_singkat', 'SN', 'Nama singkat organisasi', 'profil'),
    ('kota_pusat', 'Kota Yogyakarta', 'Kota pusat organisasi', 'profil'),
    ('provinsi_pusat', 'DI Yogyakarta', 'Provinsi pusat organisasi', 'profil');

-- Indeks untuk performa
CREATE INDEX idx_anggota_unit ON anggota(unit_id);
CREATE INDEX idx_anggota_nomor ON anggota(nomor_anggota);
CREATE INDEX idx_kehadiran_sesi ON kehadiran(sesi_id);
CREATE INDEX idx_kehadiran_anggota ON kehadiran(anggota_id);
CREATE INDEX idx_tagihan_blba_anggota ON tagihan_blba(anggota_id);
CREATE INDEX idx_pendaftaran_event ON pendaftaran_event(event_id);
CREATE INDEX idx_events_tanggal ON events(tanggal);
CREATE INDEX idx_konten_slug ON konten(slug);
