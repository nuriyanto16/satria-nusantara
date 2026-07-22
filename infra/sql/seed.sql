-- Satria Nusantara - Seed Complete Dummy Data with Performance Tuning Indices

-- ── 1. Performance Tuning Indices ─────────────────────────────
CREATE INDEX IF NOT EXISTS idx_cabang_provinsi ON cabang(provinsi_id);
CREATE INDEX IF NOT EXISTS idx_unit_latihan_cabang ON unit_latihan(cabang_id);
CREATE INDEX IF NOT EXISTS idx_sesi_latihan_unit_tanggal ON sesi_latihan(unit_id, tanggal);
CREATE INDEX IF NOT EXISTS idx_anggota_status ON anggota(status);
CREATE INDEX IF NOT EXISTS idx_anggota_nama_trgm ON anggota USING gin (nama_lengkap gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_cabang_nama_trgm ON cabang USING gin (nama gin_trgm_ops);

-- ── 2. Clear Existing Dummy Data ────────────────────────────────
DELETE FROM kehadiran;
DELETE FROM absensi_qr;
DELETE FROM sesi_latihan;
DELETE FROM pelatih;
DELETE FROM tagihan_blba;
DELETE FROM anggota;
DELETE FROM unit_latihan;
DELETE FROM cabang;

-- ── 3. Insert Cabang (dibuat_oleh references default Admin: fc2adf16-ec21-4f4a-83ea-8b337f6b4b42)
INSERT INTO cabang (id, kode, nama, provinsi_id, kota_kab, alamat, status, dibuat_oleh) VALUES
('ca000000-0000-0000-0000-000000000001', 'YGY', 'Kota Yogyakarta', 1, 'Kota Yogyakarta', 'Jl. Malioboro No. 42, Yogyakarta', 'aktif', 'fc2adf16-ec21-4f4a-83ea-8b337f6b4b42'),
('ca000000-0000-0000-0000-000000000002', 'BDG', 'Kota Bandung', 2, 'Kota Bandung', 'Jl. Dago No. 102, Bandung', 'aktif', 'fc2adf16-ec21-4f4a-83ea-8b337f6b4b42'),
('ca000000-0000-0000-0000-000000000003', 'JKT', 'Jakarta Pusat', 5, 'Jakarta Pusat', 'Jl. Menteng Raya No. 15, Jakarta', 'aktif', 'fc2adf16-ec21-4f4a-83ea-8b337f6b4b42'),
('ca000000-0000-0000-0000-000000000004', 'SBY', 'Kota Surabaya', 4, 'Kota Surabaya', 'Jl. Rungkut Industri No. 8, Surabaya', 'aktif', 'fc2adf16-ec21-4f4a-83ea-8b337f6b4b42'),
('ca000000-0000-0000-0000-000000000005', 'SMG', 'Kota Semarang', 3, 'Kota Semarang', 'Jl. Pandanaran No. 56, Semarang', 'aktif', 'fc2adf16-ec21-4f4a-83ea-8b337f6b4b42'),
('ca000000-0000-0000-0000-000000000006', 'MND', 'Kota Manado', 7, 'Kota Manado', 'Jl. Sam Ratulangi No. 12, Manado', 'aktif', 'fc2adf16-ec21-4f4a-83ea-8b337f6b4b42'),
('ca000000-0000-0000-0000-000000000007', 'DPS', 'Kota Denpasar', 9, 'Kota Denpasar', 'Jl. Teuku Umar No. 88, Denpasar', 'aktif', 'fc2adf16-ec21-4f4a-83ea-8b337f6b4b42');

-- ── 4. Insert Unit Latihan
INSERT INTO unit_latihan (id, cabang_id, nama, lokasi_nama, lokasi_alamat, status) VALUES
('de000000-0000-0000-0000-000000000001', 'ca000000-0000-0000-0000-000000000001', 'Unit Malioboro', 'Lapangan Malioboro', 'Jl. Malioboro, Yogyakarta', 'aktif'),
('de000000-0000-0000-0000-000000000002', 'ca000000-0000-0000-0000-000000000001', 'Unit Kotagede', 'Lapangan Kotagede', 'Jl. Kotagede, Yogyakarta', 'aktif'),
('de000000-0000-0000-0000-000000000003', 'ca000000-0000-0000-0000-000000000002', 'Unit Balkot', 'Lapangan Dago', 'Jl. Dago, Bandung', 'aktif'),
('de000000-0000-0000-0000-000000000004', 'ca000000-0000-0000-0000-000000000003', 'Unit Menteng', 'Taman Menteng', 'Jl. Menteng, Jakarta', 'aktif'),
('de000000-0000-0000-0000-000000000005', 'ca000000-0000-0000-0000-000000000004', 'Unit Rungkut', 'Lapangan Rungkut', 'Jl. Rungkut, Surabaya', 'aktif'),
('de000000-0000-0000-0000-000000000006', 'ca000000-0000-0000-0000-000000000005', 'Unit Simpang Lima', 'Lapangan Simpang Lima', 'Jl. Pandanaran, Semarang', 'aktif'),
('de000000-0000-0000-0000-000000000007', 'ca000000-0000-0000-0000-000000000007', 'Unit Renon', 'Lapangan Renon', 'Jl. Renon, Denpasar', 'aktif');

-- ── 5. Insert Anggota (Verified / Aktif)
INSERT INTO anggota (id, nomor_anggota, nama_lengkap, jenis_kelamin, no_hp, unit_id, tingkatan, jurus_saat_ini, counter_kehadiran, status, tanggal_daftar) VALUES
('da000000-0000-0000-0000-000000000001', 'YO-YGY-00034', 'Ahmad Santoso', 'L', '081234567890', 'de000000-0000-0000-0000-000000000001', 'PH', 6, 88, 'aktif', NOW() - INTERVAL '60 days'),
('da000000-0000-0000-0000-000000000002', 'JB-BDG-00142', 'Rina Wulandari', 'P', '081311112222', 'de000000-0000-0000-0000-000000000003', 'Gabungan', 3, 71, 'aktif', NOW() - INTERVAL '45 days'),
('da000000-0000-0000-0000-000000000003', 'JI-SBY-00055', 'Budi Hartono', 'L', '081455556666', 'de000000-0000-0000-0000-000000000005', 'Dasar', 8, 92, 'aktif', NOW() - INTERVAL '30 days'),
('da000000-0000-0000-0000-000000000004', 'JT-SMG-00089', 'Dewi Wardani', 'P', '081577778888', 'de000000-0000-0000-0000-000000000002', 'PK', 2, 65, 'aktif', NOW() - INTERVAL '90 days'),
('da000000-0000-0000-0000-000000000005', 'JK-JKT-00201', 'Hendra Nugraha', 'L', '081299990000', 'de000000-0000-0000-0000-000000000004', 'Pra Dasar', 5, 55, 'aktif', NOW() - INTERVAL '15 days'),
('da000000-0000-0000-0000-000000000010', 'YO-YGY-00035', 'Cahyo Wibowo', 'L', '081234567891', 'de000000-0000-0000-0000-000000000001', 'Dasar', 3, 40, 'aktif', NOW() - INTERVAL '20 days'),
('da000000-0000-0000-0000-000000000011', 'JB-BDG-00143', 'Siska Lestari', 'P', '081311112223', 'de000000-0000-0000-0000-000000000003', 'Pra Dasar', 1, 10, 'aktif', NOW() - INTERVAL '12 days'),
('da000000-0000-0000-0000-000000000012', 'JI-SBY-00056', 'Fajar Pratama', 'L', '081455556667', 'de000000-0000-0000-0000-000000000005', 'PH', 7, 75, 'aktif', NOW() - INTERVAL '40 days'),
('da000000-0000-0000-0000-000000000013', 'JK-JKT-00202', 'Reza Pahlevi', 'L', '081299990001', 'de000000-0000-0000-0000-000000000004', 'Gabungan', 4, 82, 'aktif', NOW() - INTERVAL '35 days'),
('da000000-0000-0000-0000-000000000014', 'BA-DPS-00012', 'Wayan Sudarsana', 'L', '081822223333', 'de000000-0000-0000-0000-000000000007', 'Dasar', 4, 30, 'aktif', NOW() - INTERVAL '50 days'),
('da000000-0000-0000-0000-000000000021', 'YO-YGY-00036', 'Guntur Wibisono', 'L', '081234567801', 'de000000-0000-0000-0000-000000000001', 'PK', 5, 45, 'aktif', NOW() - INTERVAL '40 days'),
('da000000-0000-0000-0000-000000000022', 'JB-BDG-00144', 'Novi Herawati', 'P', '081311112202', 'de000000-0000-0000-0000-000000000003', 'Dasar', 2, 60, 'aktif', NOW() - INTERVAL '30 days'),
('da000000-0000-0000-0000-000000000023', 'JI-SBY-00057', 'Eko Prasetyo', 'L', '081455556603', 'de000000-0000-0000-0000-000000000005', 'Gabungan', 3, 80, 'aktif', NOW() - INTERVAL '25 days');

-- ── 6. Insert Anggota (Menunggu Verifikasi / Pending)
INSERT INTO anggota (id, nama_lengkap, jenis_kelamin, no_hp, unit_id, tingkatan, jurus_saat_ini, counter_kehadiran, status, tanggal_daftar) VALUES
('da000000-0000-0000-0000-000000000006', 'Siti Rahayu', 'P', '082122223333', 'de000000-0000-0000-0000-000000000001', 'Pra Dasar', 1, 0, 'pending', NOW() - INTERVAL '2 hours'),
('da000000-0000-0000-0000-000000000007', 'Wisnu Saputro', 'L', '082133334444', 'de000000-0000-0000-0000-000000000002', 'Pra Dasar', 1, 0, 'pending', NOW() - INTERVAL '5 hours'),
('da000000-0000-0000-0000-000000000008', 'Ratna Ningrum', 'P', '082144445555', 'de000000-0000-0000-0000-000000000003', 'Pra Dasar', 1, 0, 'pending', NOW() - INTERVAL '1 day'),
('da000000-0000-0000-0000-000000000015', 'I Gusti Ngurah', 'L', '082155556666', 'de000000-0000-0000-0000-000000000007', 'Pra Dasar', 1, 0, 'pending', NOW() - INTERVAL '12 hours'),
('da000000-0000-0000-0000-000000000016', 'Mega Utami', 'P', '082166667777', 'de000000-0000-0000-0000-000000000004', 'Pra Dasar', 1, 0, 'pending', NOW() - INTERVAL '1 day');

-- ── 7. Insert Pelatih
INSERT INTO pelatih (id, anggota_id, cabang_id, jenis, kategori_transport, is_active) VALUES
('ee000000-0000-0000-0000-000000000001', 'da000000-0000-0000-0000-000000000001', 'ca000000-0000-0000-0000-000000000001', 'cabang', 'Jarak Dekat', TRUE),
('ee000000-0000-0000-0000-000000000002', 'da000000-0000-0000-0000-000000000002', 'ca000000-0000-0000-0000-000000000002', 'cabang', 'Jarak Sedang', TRUE),
('ee000000-0000-0000-0000-000000000003', 'da000000-0000-0000-0000-000000000003', 'ca000000-0000-0000-0000-000000000004', 'cabang', 'Jarak Jauh', TRUE);

-- ── 8. Insert Sesi Latihan (Jadwal)
INSERT INTO sesi_latihan (id, unit_id, pelatih_id, jenis, tanggal, jam_mulai, jam_selesai, lokasi, materi) VALUES
('fe000000-0000-0000-0000-000000000001', 'de000000-0000-0000-0000-000000000001', 'ee000000-0000-0000-0000-000000000001', 'Latihan Rutin', '2026-07-02', '07:00:00', '09:00:00', 'Lapangan Malioboro', 'Jurus 1 - 3, Napas Pembinaan'),
('fe000000-0000-0000-0000-000000000002', 'de000000-0000-0000-0000-000000000002', 'ee000000-0000-0000-0000-000000000001', 'Latihan Khusus', '2026-07-05', '06:30:00', '08:30:00', 'Lapangan Kotagede', 'Pemantapan Jurus 6 - 9'),
('fe000000-0000-0000-0000-000000000003', 'de000000-0000-0000-0000-000000000003', 'ee000000-0000-0000-0000-000000000002', 'Latihan Pelatih', '2026-07-08', '16:00:00', '18:00:00', 'Lapangan Dago', 'Evaluasi Jurus & Meditasi'),
('fe000000-0000-0000-0000-000000000004', 'de000000-0000-0000-0000-000000000004', 'ee000000-0000-0000-0000-000000000002', 'Latihan Rutin', '2026-07-12', '19:00:00', '21:00:00', 'Taman Menteng', 'Jurus Dasar, Pernapasan Segitiga'),
('fe000000-0000-0000-0000-000000000005', 'de000000-0000-0000-0000-000000000007', 'ee000000-0000-0000-0000-000000000003', 'Latihan Rutin', '2026-07-14', '07:00:00', '09:00:00', 'Lapangan Renon', 'Jurus PH, Pernapasan Dada'),
('fe000000-0000-0000-0000-000000000006', 'de000000-0000-0000-0000-000000000001', 'ee000000-0000-0000-0000-000000000001', 'Latihan Rutin', '2026-07-19', '07:00:00', '09:00:00', 'Lapangan Malioboro', 'Jurus 4 - 6, Napas Pembinaan II'),
('fe000000-0000-0000-0000-000000000007', 'de000000-0000-0000-0000-000000000003', 'ee000000-0000-0000-0000-000000000002', 'Latihan Khusus', '2026-07-21', '06:00:00', '08:00:00', 'Lapangan Dago', 'Pemantapan Jurus Gabungan'),
('fe000000-0000-0000-0000-000000000008', 'de000000-0000-0000-0000-000000000005', 'ee000000-0000-0000-0000-000000000003', 'Latihan Rutin', '2026-07-25', '07:00:00', '09:00:00', 'Lapangan Rungkut', 'Jurus 1 - 5 & Pernapasan Khusus'),
('fe000000-0000-0000-0000-000000000009', 'de000000-0000-0000-0000-000000000006', 'ee000000-0000-0000-0000-000000000003', 'Latihan Rutin', '2026-07-26', '08:00:00', '10:00:00', 'Simpang Lima', 'Jurus Dasar, Napas Pembinaan'),
('fe000000-0000-0000-0000-000000000010', 'de000000-0000-0000-0000-000000000002', 'ee000000-0000-0000-0000-000000000001', 'Latihan Pelatih', '2026-07-28', '19:00:00', '21:00:00', 'Lapangan Kotagede', 'Meditasi Tingkat Lanjut');

-- ── 9. Insert Tagihan BLBA (Finance) for Members
INSERT INTO tagihan_blba (id, anggota_id, bulan, tahun, nominal, status, tanggal_bayar) VALUES
(uuid_generate_v4(), 'da000000-0000-0000-0000-000000000001', 5, 2026, 50000, 'lunas', NOW() - INTERVAL '10 days'),
(uuid_generate_v4(), 'da000000-0000-0000-0000-000000000002', 5, 2026, 50000, 'lunas', NOW() - INTERVAL '8 days'),
(uuid_generate_v4(), 'da000000-0000-0000-0000-000000000003', 5, 2026, 50000, 'lunas', NOW() - INTERVAL '5 days'),
(uuid_generate_v4(), 'da000000-0000-0000-0000-000000000004', 5, 2026, 50000, 'belum_bayar', NULL),
(uuid_generate_v4(), 'da000000-0000-0000-0000-000000000005', 5, 2026, 50000, 'belum_bayar', NULL),
(uuid_generate_v4(), 'da000000-0000-0000-0000-000000000010', 5, 2026, 50000, 'lunas', NOW() - INTERVAL '3 days'),
(uuid_generate_v4(), 'da000000-0000-0000-0000-000000000011', 5, 2026, 50000, 'lunas', NOW() - INTERVAL '2 days'),
(uuid_generate_v4(), 'da000000-0000-0000-0000-000000000012', 5, 2026, 50000, 'lunas', NOW() - INTERVAL '1 day'),
(uuid_generate_v4(), 'da000000-0000-0000-0000-000000000013', 5, 2026, 50000, 'belum_bayar', NULL);
