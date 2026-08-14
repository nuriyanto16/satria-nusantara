<template>
  <div class="policy-page">
    <div class="policy-container">
      <div class="policy-header">
        <h1 class="policy-title">Kebijakan Privasi (Privacy Policy)</h1>
        <p class="policy-date">Terakhir diperbarui: 14 Agustus 2026</p>
      </div>

      <div class="policy-content">
        <h2>1. Pendahuluan</h2>
        <p>
          Selamat datang di platform Sistem Informasi Manajemen Anggota Nasional (SIMA) Satria Nusantara. 
          Kami sangat menghargai privasi Anda dan berkomitmen untuk melindungi data pribadi Anda. 
          Kebijakan privasi ini menjelaskan bagaimana kami mengumpulkan, menggunakan, dan melindungi informasi Anda.
        </p>

        <h2>2. Informasi yang Kami Kumpulkan</h2>
        <p>
          Saat Anda menggunakan layanan kami, kami dapat mengumpulkan informasi berikut:
        </p>
        <ul>
          <li><strong>Data Profil:</strong> Nama, alamat email, nomor telepon, tanggal lahir, dan informasi identitas lainnya yang relevan dengan keanggotaan.</li>
          <li><strong>Data Akses:</strong> Informasi login, riwayat akses (termasuk akses via akun Google SSO), serta data perangkat yang Anda gunakan.</li>
          <li><strong>Data Aktivitas:</strong> Rekam jejak kehadiran latihan, status iuran (BLBA), dan perkembangan sabuk/tingkat.</li>
        </ul>

        <h2>3. Penggunaan Informasi</h2>
        <p>
          Informasi yang dikumpulkan digunakan untuk:
        </p>
        <ul>
          <li>Mengelola dan memverifikasi keanggotaan Anda di Yayasan Satria Nusantara.</li>
          <li>Memberikan akses ke fitur-fitur pada sistem, seperti laporan presensi dan pembayaran.</li>
          <li>Berkomunikasi dengan Anda mengenai pembaruan sistem, jadwal latihan, atau pengumuman penting lainnya.</li>
          <li>Menganalisis data statistik perkembangan cabang untuk keperluan internal yayasan.</li>
        </ul>

        <h2>4. Berbagi Data dengan Pihak Ketiga</h2>
        <p>
          Kami tidak akan menjual, menyewakan, atau menukar data pribadi Anda kepada pihak ketiga. 
          Data Anda hanya akan dibagikan jika diwajibkan oleh hukum atau untuk keperluan operasional teknis dengan mitra yang terikat perjanjian kerahasiaan (contoh: layanan hosting dan pembayaran gateway seperti Xendit).
        </p>

        <h2>5. Keamanan Data</h2>
        <p>
          Kami menerapkan standar keamanan teknis dan administratif yang ketat untuk melindungi data Anda dari akses, perubahan, atau pengungkapan yang tidak sah.
        </p>

        <h2>6. Hak Pengguna</h2>
        <p>
          Anda memiliki hak untuk meminta akses, perbaikan, atau penghapusan data pribadi Anda dari sistem kami. 
          Silakan hubungi administrator cabang atau pusat untuk melakukan pengajuan tersebut.
        </p>

        <h2>7. Hubungi Kami</h2>
        <p>
          Jika Anda memiliki pertanyaan mengenai kebijakan privasi ini, silakan hubungi kami melalui:<br>
          Email: <strong>admin@satria-nusantara.org</strong><br>
          Alamat: Sekretariat Pusat Yayasan Satria Nusantara
        </p>
      </div>

      <div class="policy-footer">
        <button class="btn-back" @click="goBack">
          <i class="ti ti-arrow-left"></i> Kembali ke Halaman Login
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

definePageMeta({
  title: 'Kebijakan Privasi',
  layout: false
})

const router = useRouter()
const goBack = () => {
  router.back()
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&display=swap');

.policy-page {
  min-height: 100vh;
  background-color: #f8fafc;
  font-family: 'Plus Jakarta Sans', 'Inter', sans-serif;
  padding: 40px 20px;
  display: flex;
  justify-content: center;
}

.policy-container {
  background: #ffffff;
  max-width: 800px;
  width: 100%;
  border-radius: 16px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.05);
  overflow: hidden;
}

.policy-header {
  background: linear-gradient(135deg, #0d3419 0%, #1a5c2a 100%);
  color: #fff;
  padding: 40px 40px;
  text-align: center;
}

.policy-title {
  margin: 0 0 10px 0;
  font-size: 28px;
  font-weight: 800;
}

.policy-date {
  margin: 0;
  font-size: 14px;
  opacity: 0.8;
}

.policy-content {
  padding: 40px;
  color: #334155;
  line-height: 1.6;
}

.policy-content h2 {
  font-size: 18px;
  color: #1e293b;
  margin-top: 32px;
  margin-bottom: 12px;
  font-weight: 700;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 8px;
}

.policy-content h2:first-child {
  margin-top: 0;
}

.policy-content p {
  margin-bottom: 16px;
  font-size: 15px;
}

.policy-content ul {
  padding-left: 24px;
  margin-bottom: 24px;
}

.policy-content li {
  margin-bottom: 8px;
  font-size: 15px;
}

.policy-footer {
  padding: 24px 40px;
  background: #f1f5f9;
  border-top: 1px solid #e2e8f0;
  text-align: center;
}

.btn-back {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  padding: 10px 20px;
  border-radius: 8px;
  font-family: inherit;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn-back:hover {
  background: #f8fafc;
  color: #1a5c2a;
  border-color: #1a5c2a;
}
</style>
