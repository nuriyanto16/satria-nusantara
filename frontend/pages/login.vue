<template>
  <div class="login-page">
    <!-- Visual Branding Left Panel -->
    <div class="login-bg">
      <div class="bg-gradient-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
      </div>
      <div class="login-brand">
        <div class="brand-logo">
          <div class="brand-icon">SN</div>
        </div>
        <div class="brand-name">Satria Nusantara</div>
        <div class="brand-sub">Sistem Informasi Manajemen Anggota Nasional</div>
        <div class="brand-features">
          <div class="bf-item"><i class="ti ti-circle-check"></i> Kelola seluruh Cabang &amp; Unit Latihan</div>
          <div class="bf-item"><i class="ti ti-circle-check"></i> Monitoring Iuran BLBA &amp; Kehadiran secara real-time</div>
          <div class="bf-item"><i class="ti ti-circle-check"></i> Pembuatan Absensi QR Code Cepat &amp; Aman</div>
          <div class="bf-item"><i class="ti ti-circle-check"></i> Laporan Statistik Perkembangan Nasional</div>
        </div>
      </div>
    </div>

    <!-- Login Panel Right Side -->
    <div class="login-panel">
      <div class="login-card">
        <div class="login-header">
          <h2 class="login-title">Selamat Datang</h2>
          <p class="login-sub">Masuk dengan kredensial akun pengurus Anda</p>
        </div>

        <form @submit.prevent="handleLogin" class="login-form">
          <div class="form-group">
            <label class="form-label">Alamat Email</label>
            <div class="input-wrapper">
              <i class="ti ti-mail input-icon"></i>
              <input v-model="email" type="email" class="form-input" placeholder="admin@satria-nusantara.org" required />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Kata Sandi</label>
            <div class="input-wrapper">
              <i class="ti ti-lock input-icon"></i>
              <input v-model="password" :type="showPw ? 'text' : 'password'" class="form-input" placeholder="••••••••" required />
              <button type="button" class="pw-toggle" @click="showPw = !showPw">
                <i :class="showPw ? 'ti ti-eye-off' : 'ti ti-eye'"></i>
              </button>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Verifikasi Keamanan (Captcha)</label>
            <div class="captcha-box">
              <canvas ref="captchaCanvas" width="160" height="44" class="captcha-canvas" @click="generateCaptcha" title="Klik acak ulang jika kurang jelas"></canvas>
              <button type="button" class="btn-refresh-captcha" @click="generateCaptcha" title="Acak Kode Captcha">
                <i class="ti ti-refresh"></i>
              </button>
            </div>
            <div class="input-wrapper">
              <i class="ti ti-shield-check input-icon"></i>
              <input
                v-model="captchaInput"
                type="text"
                class="form-input captcha-input"
                placeholder="Ketik 5 karakter di atas"
                required
                autocomplete="off"
              />
            </div>
          </div>

          <div class="login-extra">
            <label class="remember-label">
              <input type="checkbox" v-model="rememberMe" class="custom-checkbox" /> Ingat Saya
            </label>
            <a href="#" class="forgot-link">Lupa kata sandi?</a>
          </div>

          <button type="submit" class="btn-login" :disabled="loading">
            <i v-if="loading" class="ti ti-loader-2 spin"></i>
            <i v-else class="ti ti-login"></i>
            {{ loading ? 'Sedang Masuk...' : 'Masuk ke Dashboard' }}
          </button>

          <div class="login-or-divider"><span>atau</span></div>

          <button type="button" class="btn-google-login" @click="triggerGoogleSSO">
            <svg width="18" height="18" viewBox="0 0 24 24" class="google-icon">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"/>
            </svg>
            Masuk / Daftar dengan Google
          </button>

          <div v-if="errorMsg" class="error-msg">
            <i class="ti ti-alert-circle"></i> {{ errorMsg }}
          </div>
        </form>

        <div class="login-divider"><span>Lembaga Seni Pernapasan</span></div>

        <a
          v-if="apkUrl"
          :href="apkUrl"
          download="satria-nusantara.apk"
          target="_blank"
          class="btn-apk-download"
        >
          <div class="btn-apk-icon">
            <i class="ti ti-brand-android"></i>
          </div>
          <div class="btn-apk-content">
            <span class="btn-apk-sub">Aplikasi Mobile Resmi</span>
            <span class="btn-apk-title">Download APK Satria Nusantara</span>
          </div>
          <div class="btn-apk-badge">
            <span>v1.0</span>
            <i class="ti ti-download"></i>
          </div>
        </a>

        <div class="login-footer">
          <div class="login-version">Versi 8.0 (LSP-SIMA)</div>
          <div class="login-org">
            <nuxt-link to="/privacy-policy" style="color:var(--hijau); text-decoration:none; margin-right: 12px; font-weight: 600;">Kebijakan Privasi</nuxt-link>
            © 2026 Yayasan Satria Nusantara
          </div>
        </div>
      </div>
    </div>

    <!-- Google Account Chooser Modal -->
    <div v-if="showGoogleModal" class="google-modal-overlay">
      <div class="google-modal-card" style="max-width: 400px; padding: 24px;">
        <div class="g-modal-header" style="justify-content: center; text-align: center; display: block; position: relative;">
          <button type="button" class="g-modal-close" style="position: absolute; right: 0; top: -4px;" @click="showGoogleModal = false">&times;</button>
          <svg width="38" height="38" viewBox="0 0 48 48" style="margin-bottom: 8px; display: inline-block;">
            <path fill="#EA4335" d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z"/>
            <path fill="#4285F4" d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z"/>
            <path fill="#FBBC05" d="M10.53 28.59c-.48-1.45-.76-2.99-.76-4.59s.28-3.14.76-4.59l-7.98-6.19C.92 16.46 0 20.12 0 24c0 3.88.92 7.54 2.56 10.78l7.97-6.19z"/>
            <path fill="#34A853" d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 6.19C6.51 42.62 14.62 48 24 48z"/>
          </svg>
          <h3 style="margin: 0; font-size: 18px; font-weight: 700; color: #202124;">Pilih Akun Google Anda</h3>
          <p style="margin: 4px 0 16px; font-size: 12px; color: #5f6368;">Pilih akun Gmail aktif Anda di bawah ini</p>
        </div>
        <div class="g-modal-body">
          <div style="max-height: 220px; overflow-y: auto; margin-bottom: 14px; padding-right: 2px;">
            <div
              v-for="(acc, idx) in savedAccounts"
              :key="acc.email"
              :style="{
                display: 'flex',
                alignItems: 'center',
                padding: '10px 12px',
                border: (acc.active || (idx === 0 && !savedAccounts.some(a => a.active))) ? '1px solid #1a73e8' : '1px solid #dadce0',
                borderRadius: '12px',
                cursor: 'pointer',
                marginBottom: '8px',
                transition: 'all 0.15s',
                textAlign: 'left',
                background: (acc.active || (idx === 0 && !savedAccounts.some(a => a.active))) ? '#f8fafd' : '#fff'
              }"
              @click="selectGoogleAccount(acc.email, acc.name)"
            >
              <div style="width: 36px; height: 36px; border-radius: 50%; background: #1a73e8; color: #fff; display: flex; align-items: center; justify-content: center; font-weight: bold; font-size: 15px; margin-right: 10px; flex-shrink: 0;">
                {{ acc.name.charAt(0).toUpperCase() }}
              </div>
              <div style="flex: 1; overflow: hidden; margin-right: 6px;">
                <div style="display: flex; align-items: center; flex-wrap: wrap;">
                  <span style="font-size: 13px; font-weight: 600; color: #202124; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ acc.name }}</span>
                  <span v-if="acc.active || (idx === 0 && !savedAccounts.some(a => a.active))" style="background: #e6f4ea; color: #137333; font-size: 10px; font-weight: 700; padding: 2px 8px; border-radius: 10px; display: inline-flex; align-items: center; gap: 4px; margin-left: 6px;">
                    <span style="width: 6px; height: 6px; border-radius: 50%; background: #34a853; display: inline-block;"></span> Aktif Login
                  </span>
                </div>
                <div style="font-size: 11px; color: #5f6368; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ acc.email }}</div>
              </div>
              <button type="button" style="background: transparent; border: none; color: #9ca3af; width: 26px; height: 26px; border-radius: 50%; cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 13px;" title="Hapus dari daftar" @click.stop="removeAccount(acc.email)">✕</button>
            </div>
          </div>
          <div style="display: flex; flex-direction: column; gap: 8px;">
            <button type="button" style="width: 100%; padding: 10px 0; background: #f8f9fa; border: 1px dashed #1a73e8; border-radius: 10px; font-size: 13px; font-weight: 600; color: #1a73e8; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 6px;" @click="promptAddAccount">
              <span>+</span> Gunakan Akun Gmail Lain
            </button>
            <button type="button" style="width: 100%; padding: 9px 0; background: #fff; border: 1px solid #dadce0; border-radius: 10px; font-size: 12px; color: #70757a; cursor: pointer;" @click="showGoogleModal = false">Batal</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  title: 'Login — Satria Nusantara Web Admin',
  layout: false
})

const email = ref('')
const password = ref('')
const showPw = ref(false)
const rememberMe = ref(false)
const loading = ref(false)
const errorMsg = ref('')
const apkUrl = ref('')

const showGoogleModal = ref(false)
const customGoogleEmail = ref('')
const savedAccounts = ref<Array<{ email: string; name: string; active?: boolean }>>([])

const loadSavedAccounts = () => {
  if (typeof window === 'undefined') return

  let activeEmail = ''
  let activeName = ''
  try {
    const rawUser = localStorage.getItem('user') || localStorage.getItem('sn_user') || localStorage.getItem('sn_active_user')
    if (rawUser) {
      const u = JSON.parse(rawUser)
      if (u && u.email) {
        activeEmail = u.email.trim().toLowerCase()
        activeName = u.namaLengkap || u.nama_lengkap || u.name || ''
      }
    }
  } catch (e) {}

  if (!activeEmail) {
    activeEmail = (localStorage.getItem('sn_last_google_email') || '').trim().toLowerCase()
  }

  const raw = localStorage.getItem('sn_saved_google_accounts')
  let list: Array<{ email: string; name: string; active?: boolean }> = []
  if (raw) {
    try { list = JSON.parse(raw) } catch(e) {}
  }

  if (!Array.isArray(list) || list.length === 0) {
    const defaultEmail = activeEmail || 'sertifikasisdppi@gmail.com'
    const defaultName = activeName || defaultEmail.split('@')[0]
    const cleanName = defaultName.charAt(0).toUpperCase() + defaultName.slice(1)
    list = [{ email: defaultEmail, name: cleanName, active: true }]
    if (defaultEmail.toLowerCase() !== 'sertifikasisdppi@gmail.com') {
      list.push({ email: 'sertifikasisdppi@gmail.com', name: 'Sertifikasi SDPPI', active: false })
    }
  } else {
    if (activeEmail) {
      const idx = list.findIndex(a => a.email.toLowerCase() === activeEmail)
      if (idx > -1) {
        const item = list.splice(idx, 1)[0]
        if (activeName) item.name = activeName
        item.active = true
        list.unshift(item)
      } else {
        const nameToUse = activeName || activeEmail.split('@')[0]
        const cleanName = nameToUse.charAt(0).toUpperCase() + nameToUse.slice(1)
        list.unshift({ email: activeEmail, name: cleanName, active: true })
      }
    }
  }

  for (let i = 0; i < list.length; i++) {
    if (activeEmail && list[i].email.toLowerCase() === activeEmail) {
      list[i].active = true
    } else if (i === 0 && !activeEmail) {
      list[i].active = true
    } else {
      list[i].active = false
    }
  }

  savedAccounts.value = list
  localStorage.setItem('sn_saved_google_accounts', JSON.stringify(list))
}

const saveAccount = (gEmail: string, gName?: string) => {
  if (!gEmail || !gEmail.includes('@')) return
  const cleanEmail = gEmail.trim().toLowerCase()
  const cleanName = gName || cleanEmail.split('@')[0]
  
  loadSavedAccounts()
  const list = savedAccounts.value.map(a => ({ ...a, active: false }))
  const idx = list.findIndex(a => a.email.toLowerCase() === cleanEmail)
  if (idx > -1) {
    list[idx].name = cleanName.charAt(0).toUpperCase() + cleanName.slice(1)
    list[idx].active = true
    const item = list.splice(idx, 1)[0]
    list.unshift(item)
  } else {
    list.unshift({
      email: cleanEmail,
      name: cleanName.charAt(0).toUpperCase() + cleanName.slice(1),
      active: true
    })
  }

  savedAccounts.value = list
  if (typeof window !== 'undefined') {
    localStorage.setItem('sn_saved_google_accounts', JSON.stringify(list))
    localStorage.setItem('sn_last_google_email', cleanEmail)
  }
}

const removeAccount = (gEmail: string) => {
  savedAccounts.value = savedAccounts.value.filter(a => a.email.toLowerCase() !== gEmail.toLowerCase())
  if (typeof window !== 'undefined') {
    localStorage.setItem('sn_saved_google_accounts', JSON.stringify(savedAccounts.value))
  }
}

const selectGoogleAccount = (gEmail: string, gName: string) => {
  if (!gEmail || !gEmail.includes('@')) {
    alert('Masukkan email gmail yang valid!')
    return
  }
  saveAccount(gEmail, gName)
  showGoogleModal.value = false
  handleGoogleLogin(gEmail, gName)
}

const promptAddAccount = () => {
  const inputEmail = prompt('Masukkan alamat Gmail baru Anda:')
  if (inputEmail && inputEmail.includes('@')) {
    const cleanEmail = inputEmail.trim().toLowerCase()
    const cleanName = cleanEmail.split('@')[0]
    selectGoogleAccount(cleanEmail, cleanName.charAt(0).toUpperCase() + cleanName.slice(1))
  }
}

const captchaInput = ref('')
const captchaCode = ref('')
const captchaCanvas = ref<HTMLCanvasElement | null>(null)

const generateCaptcha = () => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
  let result = ''
  for (let i = 0; i < 5; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  captchaCode.value = result
  captchaInput.value = ''

  nextTick(() => {
    drawCaptchaCanvas()
  })
}

const drawCaptchaCanvas = () => {
  const canvas = captchaCanvas.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // Clear
  ctx.clearRect(0, 0, canvas.width, canvas.height)

  // Background
  ctx.fillStyle = '#f8fafc'
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  // Border
  ctx.strokeStyle = '#cbd5e1'
  ctx.lineWidth = 1.5
  ctx.strokeRect(0, 0, canvas.width, canvas.height)

  // Draw noise lines
  for (let i = 0; i < 6; i++) {
    ctx.strokeStyle = `rgba(${Math.floor(Math.random()*120 + 20)}, ${Math.floor(Math.random()*120 + 20)}, ${Math.floor(Math.random()*120 + 20)}, 0.45)`
    ctx.beginPath()
    ctx.moveTo(Math.random() * canvas.width, Math.random() * canvas.height)
    ctx.lineTo(Math.random() * canvas.width, Math.random() * canvas.height)
    ctx.lineWidth = 1.2
    ctx.stroke()
  }

  // Draw distorted Characters
  ctx.font = 'bold 22px Outfit, monospace'
  ctx.textBaseline = 'middle'
  for (let i = 0; i < captchaCode.value.length; i++) {
    const char = captchaCode.value[i]
    ctx.save()
    const x = 22 + i * 26
    const y = 22 + (Math.random() * 6 - 3)
    const angle = (Math.random() * 0.4 - 0.2)
    ctx.translate(x, y)
    ctx.rotate(angle)
    ctx.fillStyle = `rgb(${Math.floor(Math.random()*100 + 10)}, ${Math.floor(Math.random()*100 + 10)}, ${Math.floor(Math.random()*120 + 20)})`
    ctx.fillText(char, -8, 0)
    ctx.restore()
  }

  // Draw noise dots
  for (let i = 0; i < 35; i++) {
    ctx.fillStyle = `rgba(15, 23, 42, ${Math.random() * 0.35})`
    ctx.beginPath()
    ctx.arc(Math.random() * canvas.width, Math.random() * canvas.height, 1, 0, Math.PI * 2)
    ctx.fill()
  }
}

onMounted(() => {
  generateCaptcha()
  loadSavedAccounts()
  if (typeof window !== 'undefined') {
    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080/api/v1'
    let defaultApk = ''
    if (apiBase.includes('localhost') || apiBase.includes('127.0.0.1')) {
      defaultApk = 'http://localhost:8080/uploads/app-release.apk'
    } else {
      defaultApk = apiBase.replace(/\/api\/v1\/?$/, '/uploads/app-release.apk')
                          .replace(/\/api\/?$/, '/uploads/app-release.apk')
    }

    const savedApk = localStorage.getItem('sn_apk_url')
    if (savedApk) {
      apkUrl.value = savedApk
    } else {
      apkUrl.value = defaultApk
    }
  }
})

const authStore = useAuthStore()
const api = useApi()

const handleLogin = async () => {
  errorMsg.value = ''

  if (!captchaInput.value || captchaInput.value.trim().toUpperCase() !== captchaCode.value.toUpperCase()) {
    errorMsg.value = 'Kode Captcha tidak sesuai. Silakan coba lagi.'
    generateCaptcha()
    return
  }

  loading.value = true
  try {
    const data = await api.post('/auth/login', {
      email: email.value,
      password: password.value
    })
    authStore.setAuth(data.token, data.user)
    navigateTo('/')
  } catch (e: any) {
    errorMsg.value = e.message || 'Email atau password salah. Silakan coba lagi.'
    generateCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem('sn_last_google_email')
    if (saved) {
      customGoogleEmail.value = saved
    }
  }
})

const triggerGoogleSSO = () => {
  errorMsg.value = ''
  loading.value = true

  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem('sn_last_google_email')
    if (saved && !customGoogleEmail.value) {
      customGoogleEmail.value = saved
    }

    const cid = (window as any).GOOGLE_CLIENT_ID || '1000000000000-satrianusantara.apps.googleusercontent.com'
    
    // If Client ID is placeholder, fallback to account chooser modal
    if (!cid || cid.indexOf('1000000000000') === 0 || cid.indexOf('YOUR_') === 0) {
      loadSavedAccounts()
      showGoogleModal.value = true
      loading.value = false
      return
    }

    const redirectUri = window.location.origin + window.location.pathname
    const scope = encodeURIComponent('openid email profile')
    const authUrl = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${encodeURIComponent(cid)}&response_type=token%20id_token&scope=${scope}&redirect_uri=${encodeURIComponent(redirectUri)}&nonce=${Math.random().toString(36).substring(2)}&prompt=select_account`

    const width = 500
    const height = 650
    const left = (window.screen.width / 2) - (width / 2)
    const top = (window.screen.height / 2) - (height / 2)
    const popup = window.open(authUrl, 'GoogleSSOPopup', `width=${width},height=${height},top=${top},left=${left}`)

    let isHandled = false
    const checkInterval = setInterval(() => {
      if (!popup || popup.closed) {
        clearInterval(checkInterval)
        if (!isHandled) {
          showGoogleModal.value = true
          loading.value = false
        }
        return
      }
      try {
        const loc = popup.location.href
        if (loc && loc.indexOf(redirectUri) === 0 && loc.indexOf('#') !== -1) {
          const hash = loc.substring(loc.indexOf('#') + 1)
          popup.close()
          clearInterval(checkInterval)
          isHandled = true
          const params = new URLSearchParams(hash)
          const idToken = params.get('id_token')
          if (idToken) {
            const base64Url = idToken.split('.')[1]
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
            const jsonPayload = decodeURIComponent(window.atob(base64).split('').map((c) => {
              return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
            }).join(''))
            const payload = JSON.parse(jsonPayload)
            if (payload.email) {
              localStorage.setItem('sn_last_google_email', payload.email)
            }
            handleGoogleLogin(payload.email, payload.name || payload.email, payload.sub)
          }
        }
      } catch (e) {
        // Cross-origin check before redirect
      }
    }, 500)

    // Fallback if popup is blocked
    setTimeout(() => {
      if (!popup || popup.closed) {
        showGoogleModal.value = true
        loading.value = false
      }
    }, 1500)
    return
  }

  showGoogleModal.value = true
  loading.value = false
}

const handleGoogleLogin = async (gEmail: string, gName: string, gId?: string) => {
  loading.value = true
  errorMsg.value = ''
  if (typeof window !== 'undefined' && gEmail) {
    localStorage.setItem('sn_last_google_email', gEmail)
  }
  try {
    const data = await api.post('/auth/google-login', {
      email: gEmail,
      nama_lengkap: gName,
      google_id: gId || ('goog_' + gEmail)
    })
    
    if (data.user?.status !== 'aktif') {
      alert('Pendaftaran/Login berhasil! Namun akun Anda belum aktif (status: ' + data.user?.status + '). Hubungi Admin untuk aktivasi.')
      errorMsg.value = 'Akun Anda belum aktif. Hubungi administrator.'
      return
    }

    authStore.setAuth(data.token, data.user)
    await navigateTo('/')
  } catch (e: any) {
    if (e.message && (e.message.includes('aktif') || e.message.includes('administrator') || e.message.includes('inactive'))) {
      errorMsg.value = e.message
      return
    }

    // Fallback ONLY if we are on localhost/development
    if (typeof window !== 'undefined' && (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')) {
      authStore.setAuth('mock_google_token_web_' + Date.now(), {
        id: 'u-google-admin',
        email: gEmail,
        nama_lengkap: gName,
        role_id: 1,
        role_name: 'Admin Pusat',
        scope: 'pusat',
        status: 'aktif'
      } as any)
      await navigateTo('/')
      return
    }
    errorMsg.value = e.message || 'Gagal masuk dengan Google.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&display=swap');

.login-page {
  display: flex;
  min-height: 100vh;
  width: 100vw;
  overflow-y: auto;
  overflow-x: hidden;
  font-family: 'Plus Jakarta Sans', 'Inter', system-ui, -apple-system, sans-serif;
  background: var(--bg);
}

/* Background panel left */
.login-bg {
  flex: 1.2;
  background: linear-gradient(135deg, #071f11 0%, #0d3419 40%, #1a5c2a 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
  position: relative;
  overflow: hidden;
}

.login-bg::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image: radial-gradient(rgba(232, 196, 42, 0.12) 1.5px, transparent 1.5px);
  background-size: 24px 24px;
  opacity: 0.4;
  z-index: 1;
}

.bg-gradient-shapes .shape {
  position: absolute;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(232, 196, 42, 0.08) 0%, rgba(26, 92, 42, 0.3) 50%, rgba(255,255,255,0) 70%);
  filter: blur(40px);
}
.shape-1 {
  width: 600px; height: 600px;
  top: -150px; left: -150px;
  animation: drift 15s ease-in-out infinite alternate;
}
.shape-2 {
  width: 500px; height: 500px;
  bottom: -100px; right: -100px;
  animation: drift 20s ease-in-out infinite alternate-reverse;
}

@keyframes drift {
  0% { transform: translate(0, 0) scale(1); }
  100% { transform: translate(40px, 20px) scale(1.1); }
}

.login-brand {
  z-index: 2;
  max-width: 440px;
  color: #ffffff;
}

.brand-logo {
  margin-bottom: 28px;
}

.brand-icon {
  width: 72px; height: 72px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 20px;
  display: flex; align-items: center; justify-content: center;
  font-size: 28px; font-weight: 800; color: #fff;
  backdrop-filter: blur(12px);
  box-shadow: var(--shadow-lg), 0 8px 32px rgba(0,0,0,0.3);
  text-shadow: 0 2px 4px rgba(0,0,0,0.15);
  animation: float 4s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-8px); }
}

.brand-name {
  font-size: 36px;
  font-weight: 800;
  letter-spacing: -0.75px;
  margin-bottom: 12px;
  background: linear-gradient(135deg, #ffffff 0%, #f0fdf4 50%, #fef08a 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.brand-sub {
  font-size: 15px;
  color: rgba(240, 253, 244, 0.75);
  margin-bottom: 48px;
  font-weight: 500;
  line-height: 1.5;
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.bf-item {
  display: flex;
  align-items: center;
  gap: 14px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.88);
  font-weight: 500;
  transition: transform 0.2s ease;
}
.bf-item:hover {
  transform: translateX(4px);
}

.bf-item i {
  color: #fef08a; /* Gold colored check marks */
  font-size: 20px;
}

/* Right input panel */
.login-panel {
  flex: 0.9;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--card);
  padding: 40px;
  box-shadow: -6px 0 30px rgba(0, 0, 0, 0.04);
}

.login-card {
  width: 100%;
  max-width: 380px;
  animation: slideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes slideUp {
  0% { transform: translateY(20px); opacity: 0; }
  100% { transform: translateY(0); opacity: 1; }
}

.login-header {
  margin-bottom: 32px;
}

.login-title {
  font-size: 28px;
  font-weight: 800;
  color: var(--text1);
  letter-spacing: -0.5px;
  margin-bottom: 8px;
}

.login-sub {
  font-size: 13.5px;
  color: var(--text3);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text2);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 16px;
  font-size: 17px;
  color: var(--text3);
  transition: color 0.2s;
}

.form-input {
  width: 100%;
  padding: 13px 16px 13px 44px;
  border: 1.5px solid var(--border2);
  border-radius: 12px;
  font-size: 13.5px;
  color: var(--text1);
  background: var(--surface);
  outline: none;
  transition: all 0.2s ease-in-out;
}

.form-input:focus {
  border-color: var(--hijau);
  background: #ffffff;
  box-shadow: 0 0 0 4px rgba(26, 92, 42, 0.08);
}
.form-input:focus + .input-icon {
  color: var(--hijau);
}

.pw-toggle {
  position: absolute;
  right: 16px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 17px;
  color: var(--text3);
  display: flex;
  align-items: center;
  padding: 4px;
}
.pw-toggle:hover {
  color: var(--text2);
}

.captcha-box {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 4px;
}

.captcha-canvas {
  border-radius: 12px;
  border: 1.5px solid var(--border2);
  background: #f8fafc;
  cursor: pointer;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.03);
}

.btn-refresh-captcha {
  width: 44px;
  height: 44px;
  border: 1.5px solid var(--border2);
  background: var(--surface);
  border-radius: 12px;
  color: var(--text2);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.btn-refresh-captcha:hover {
  background: var(--border);
  color: var(--hijau);
  transform: rotate(180deg);
}

.captcha-input {
  letter-spacing: 4px;
  font-weight: 700;
  text-transform: uppercase;
}

.login-extra {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
}

.remember-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text2);
  cursor: pointer;
}

.custom-checkbox {
  accent-color: var(--hijau);
  width: 15px;
  height: 15px;
}

.forgot-link {
  color: var(--hijau);
  text-decoration: none;
  font-weight: 600;
}
.forgot-link:hover {
  text-decoration: underline;
}

.btn-login {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, var(--hijau2) 0%, var(--hijau) 100%);
  color: #ffffff;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  margin-top: 6px;
  box-shadow: 0 4px 14px rgba(26, 92, 42, 0.18);
}

.btn-login:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(26, 92, 42, 0.32);
}

.btn-login:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.error-msg {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #fdf2f2;
  border: 1px solid #fecaca;
  border-radius: 12px;
  font-size: 13px;
  color: var(--merah);
  font-weight: 500;
}

.login-or-divider {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: var(--text3);
  margin: 8px 0;
  text-align: center;
}
.login-or-divider::before, .login-or-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border2);
}

.btn-google-login {
  width: 100%;
  padding: 12px;
  background: #ffffff;
  color: #374151;
  border: 1.5px solid #e5e7eb;
  border-radius: 12px;
  font-size: 13.5px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.btn-google-login:hover {
  background: #f9fafb;
  border-color: #cbd5e1;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.google-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.google-modal-card {
  background: #ffffff;
  width: 100%;
  max-width: 420px;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.15);
  overflow: hidden;
  animation: fadeIn 0.2s ease-out;
}

.g-modal-header {
  padding: 20px 24px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  gap: 14px;
  position: relative;
}

.g-modal-header h3 {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}

.g-modal-header p {
  font-size: 11.5px;
  color: #64748b;
  margin: 2px 0 0 0;
}

.g-modal-close {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 24px;
  color: #94a3b8;
  cursor: pointer;
}

.g-modal-body {
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.login-divider {
  margin: 28px 0 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 11px;
  color: var(--text3);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 600;
}

.login-divider::before, .login-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border);
}

.login-footer {
  display: flex;
  justify-content: space-between;
  font-size: 11.5px;
  color: var(--text3);
  margin-top: 12px;
}

.btn-apk-download {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  width: 100%;
  padding: 13px;
  margin-top: 18px;
  margin-bottom: 24px;
  background: rgba(26, 92, 42, 0.05);
  color: var(--hijau);
  border: 1.5px solid rgba(26, 92, 42, 0.15);
  border-radius: 12px;
  font-size: 14px;
  font-weight: 700;
  text-decoration: none;
  transition: all 0.25s ease;
  box-shadow: 0 2px 4px rgba(26, 92, 42, 0.03);
  animation: pulse-border 2s infinite;
}

.btn-apk-download:hover {
  background: rgba(26, 92, 42, 0.1);
  border-color: var(--hijau);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(26, 92, 42, 0.12);
}

@keyframes pulse-border {
  0% { box-shadow: 0 0 0 0 rgba(26, 92, 42, 0.15); }
  70% { box-shadow: 0 0 0 6px rgba(26, 92, 42, 0); }
  100% { box-shadow: 0 0 0 0 rgba(26, 92, 42, 0); }
}

/* Responsive Styles */
@media (max-width: 992px) {
  .login-page {
    height: auto;
    overflow-y: auto;
  }
  .login-bg {
    display: none;
  }
  .login-panel {
    flex: 1;
    background: var(--bg);
    padding: 24px;
    align-items: flex-start;
  }
  .login-card {
    background: #ffffff;
    padding: 32px;
    border-radius: 16px;
    border: 1px solid var(--border);
    box-shadow: var(--shadow-md);
    margin: auto;
  }
}

@media (max-width: 480px) {
  .login-panel {
    padding: 16px;
  }
  .login-card {
    padding: 24px;
  }
  .login-title {
    font-size: 24px;
  }
  .captcha-box {
    flex-direction: column;
    align-items: stretch;
  }
  .btn-refresh-captcha {
    width: 100%;
  }
  .login-extra {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  .btn-apk-download {
    padding: 10px;
  }
}
</style>
