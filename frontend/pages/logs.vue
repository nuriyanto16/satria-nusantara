<template>
  <div class="content-wrapper p-6">
    <div class="header-section mb-8">
      <h1 class="text-3xl font-bold mb-2">Sistem Log & Audit</h1>
      <p class="text-gray-600">Pantau seluruh aktivitas aplikasi dan riwayat koneksi ke Gateway Eksternal</p>
    </div>

    <!-- Tab Navigation -->
    <div class="flex border-b border-gray-200 mb-6">
      <button 
        class="py-3 px-6 text-sm font-medium focus:outline-none"
        :class="activeTab === 'audit' ? 'text-[var(--hijau)] border-b-2 border-[var(--hijau)]' : 'text-gray-500 hover:text-gray-700'"
        @click="switchTab('audit')"
      >
        <i class="ti ti-user-scan mr-2"></i> Audit Trail (Akses)
      </button>
      <button 
        class="py-3 px-6 text-sm font-medium focus:outline-none"
        :class="activeTab === 'payment' ? 'text-[var(--hijau)] border-b-2 border-[var(--hijau)]' : 'text-gray-500 hover:text-gray-700'"
        @click="switchTab('payment')"
      >
        <i class="ti ti-receipt mr-2"></i> Log Payment Gateway
      </button>
      <button 
        class="py-3 px-6 text-sm font-medium focus:outline-none"
        :class="activeTab === 'error' ? 'text-[var(--merah)] border-b-2 border-[var(--merah)]' : 'text-gray-500 hover:text-gray-700'"
        @click="switchTab('error')"
      >
        <i class="ti ti-alert-triangle mr-2"></i> Error Aplikasi
      </button>
    </div>

    <!-- Card Container -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
      
      <!-- Loading State -->
      <div v-if="loading" class="py-12 text-center text-gray-400">
        <i class="ti ti-loader animate-spin text-4xl mb-3 inline-block"></i>
        <p>Memuat log...</p>
      </div>

      <!-- AUDIT TAB -->
      <div v-else-if="activeTab === 'audit'" class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-200">
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Waktu</th>
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Pengguna</th>
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Aksi</th>
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Entitas / IP</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in auditLogs" :key="log.id" class="border-b border-gray-100 hover:bg-gray-50/50">
              <td class="p-4 text-sm text-gray-600 whitespace-nowrap">{{ formatDate(log.createdAt) }}</td>
              <td class="p-4">
                <div class="font-medium text-gray-900">{{ log.userName || 'Sistem' }}</div>
                <div class="text-xs text-gray-400 font-mono">{{ log.userId || '-' }}</div>
              </td>
              <td class="p-4">
                <span class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-blue-50 text-blue-700">
                  {{ log.action }}
                </span>
              </td>
              <td class="p-4">
                <div class="text-sm text-gray-900">{{ log.entity || '-' }} ({{ log.entityId || '-' }})</div>
                <div class="text-xs text-gray-400"><i class="ti ti-network"></i> {{ log.ipAddress || 'Unknown' }}</div>
              </td>
            </tr>
            <tr v-if="auditLogs.length === 0">
              <td colspan="4" class="p-8 text-center text-gray-400">Belum ada catatan aktivitas</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- PAYMENT TAB -->
      <div v-else-if="activeTab === 'payment'" class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-200">
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Waktu</th>
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Gateway / TxID</th>
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Status HTTP</th>
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase text-center">Payload JSON</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in paymentLogs" :key="log.id" class="border-b border-gray-100 hover:bg-gray-50/50">
              <td class="p-4 text-sm text-gray-600 whitespace-nowrap">{{ formatDate(log.createdAt) }}</td>
              <td class="p-4">
                <div class="font-bold text-gray-800">{{ log.provider }}</div>
                <div class="text-xs text-gray-400 font-mono" :title="log.endpoint">{{ log.endpoint.substring(0, 30) }}...</div>
                <div class="text-xs text-gray-400 font-mono mt-1">Tx: {{ log.transactionId || '-' }}</div>
              </td>
              <td class="p-4">
                <span 
                  class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium"
                  :class="log.statusCode >= 200 && log.statusCode < 300 ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'"
                >
                  {{ log.statusCode || 'N/A' }}
                </span>
                <div v-if="log.errorMessage" class="text-xs text-red-500 mt-1 max-w-xs truncate" :title="log.errorMessage">
                  {{ log.errorMessage }}
                </div>
              </td>
              <td class="p-4 text-center">
                <button 
                  @click="viewJson(log)" 
                  class="text-xs border border-gray-200 rounded px-3 py-1.5 text-gray-600 hover:bg-gray-50 hover:text-[var(--hijau)] transition-colors"
                >
                  <i class="ti ti-brackets-contain mr-1"></i> Lihat Data
                </button>
              </td>
            </tr>
            <tr v-if="paymentLogs.length === 0">
              <td colspan="4" class="p-8 text-center text-gray-400">Belum ada log dari Gateway Pembayaran</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- ERROR TAB -->
      <div v-else-if="activeTab === 'error'" class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-200">
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Waktu</th>
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Service / Type</th>
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase">Error Message</th>
              <th class="p-4 text-xs font-semibold text-gray-500 uppercase text-center">Stack Trace</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in appErrors" :key="log.id" class="border-b border-gray-100 hover:bg-gray-50/50">
              <td class="p-4 text-sm text-gray-600 whitespace-nowrap">{{ formatDate(log.createdAt) }}</td>
              <td class="p-4">
                <div class="font-medium text-gray-900">{{ log.service || '-' }}</div>
                <div class="text-xs text-gray-500 font-mono mt-1">{{ log.errorType || '-' }}</div>
              </td>
              <td class="p-4">
                <div class="text-sm text-red-600 font-medium line-clamp-2" :title="log.message">{{ log.message }}</div>
              </td>
              <td class="p-4 text-center">
                <button 
                  @click="viewStackTrace(log)" 
                  class="text-xs border border-gray-200 rounded px-3 py-1.5 text-gray-600 hover:bg-gray-50 hover:text-[var(--merah)] transition-colors"
                >
                  <i class="ti ti-bug mr-1"></i> Detail Trace
                </button>
              </td>
            </tr>
            <tr v-if="appErrors.length === 0">
              <td colspan="4" class="p-8 text-center text-gray-400">Sistem berjalan mulus tanpa error terdeteksi</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Controls (Sederhana) -->
      <div class="flex justify-between items-center mt-6 pt-4 border-t border-gray-100" v-if="!loading">
        <div class="text-sm text-gray-500">Halaman {{ page }}</div>
        <div class="flex space-x-2">
          <button @click="changePage(-1)" :disabled="page <= 1" class="px-3 py-1 border rounded hover:bg-gray-50 disabled:opacity-50">Prev</button>
          <button @click="changePage(1)" class="px-3 py-1 border rounded hover:bg-gray-50">Next</button>
        </div>
      </div>
    </div>

    <!-- JSON Viewer Modal -->
    <div v-if="showJsonModal" class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-xl shadow-2xl w-full max-w-4xl max-h-[90vh] flex flex-col">
        <div class="p-5 border-b flex justify-between items-center">
          <h3 class="font-bold text-lg">Detail Payload & Respons</h3>
          <button @click="showJsonModal = false" class="text-gray-400 hover:text-gray-600"><i class="ti ti-x text-xl"></i></button>
        </div>
        <div class="p-5 overflow-y-auto flex-1 grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <h4 class="font-semibold text-sm mb-2 flex items-center"><i class="ti ti-arrow-right text-[var(--hijau)] mr-1"></i> Request Payload (Dikirim)</h4>
            <pre class="bg-gray-900 text-gray-300 p-4 rounded-lg text-xs overflow-x-auto h-[400px] shadow-inner font-mono whitespace-pre-wrap">{{ JSON.stringify(selectedLog?.requestPayload || {}, null, 2) }}</pre>
          </div>
          <div>
            <h4 class="font-semibold text-sm mb-2 flex items-center"><i class="ti ti-arrow-left text-blue-500 mr-1"></i> Response Payload (Diterima)</h4>
            <pre class="bg-gray-900 text-gray-300 p-4 rounded-lg text-xs overflow-x-auto h-[400px] shadow-inner font-mono whitespace-pre-wrap">{{ JSON.stringify(selectedLog?.responsePayload || {}, null, 2) }}</pre>
          </div>
        </div>
        <div class="p-4 border-t bg-gray-50 rounded-b-xl flex justify-end">
          <button @click="showJsonModal = false" class="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded font-medium transition-colors">Tutup</button>
        </div>
      </div>
    </div>

    <!-- Stack Trace Modal -->
    <div v-if="showTraceModal" class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-xl shadow-2xl w-full max-w-3xl max-h-[90vh] flex flex-col">
        <div class="p-5 border-b border-red-100 bg-red-50 rounded-t-xl flex justify-between items-center">
          <h3 class="font-bold text-red-700 text-lg flex items-center"><i class="ti ti-alert-octagon mr-2 text-xl"></i> Detail Error App</h3>
          <button @click="showTraceModal = false" class="text-red-400 hover:text-red-700"><i class="ti ti-x text-xl"></i></button>
        </div>
        <div class="p-5 overflow-y-auto flex-1 bg-gray-900">
          <div class="mb-4 bg-red-900/30 p-3 rounded border border-red-500/30 text-red-200 text-sm font-mono whitespace-pre-wrap">
            {{ selectedLog?.message }}
          </div>
          <h4 class="text-gray-400 text-xs uppercase tracking-wider mb-2">Stack Trace</h4>
          <pre class="text-gray-300 text-xs font-mono whitespace-pre-wrap">{{ selectedLog?.stackTrace || 'No stack trace available' }}</pre>
          
          <h4 v-if="selectedLog?.context" class="text-gray-400 text-xs uppercase tracking-wider mt-6 mb-2">Contextual Data</h4>
          <pre v-if="selectedLog?.context" class="text-blue-300 text-xs font-mono bg-gray-800 p-3 rounded whitespace-pre-wrap">{{ JSON.stringify(selectedLog.context, null, 2) }}</pre>
        </div>
        <div class="p-4 border-t bg-gray-50 rounded-b-xl flex justify-end">
          <button @click="showTraceModal = false" class="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded font-medium transition-colors">Tutup</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

definePageMeta({
  layout: 'default'
})

const activeTab = ref('audit') // 'audit' | 'payment' | 'error'
const loading = ref(true)
const page = ref(1)

const auditLogs = ref<any[]>([])
const paymentLogs = ref<any[]>([])
const appErrors = ref<any[]>([])

const showJsonModal = ref(false)
const showTraceModal = ref(false)
const selectedLog = ref<any>(null)

// API Config
const api = useApi()

const fetchLogs = async () => {
  loading.value = true
  try {
    const endpoint = `/admin/logs/${activeTab.value === 'error' ? 'errors' : activeTab.value === 'audit' ? 'audit' : 'payments'}?page=${page.value}`
    const { data } = await api.get(endpoint)
    
    if (activeTab.value === 'audit') {
      auditLogs.value = data.data || []
    } else if (activeTab.value === 'payment') {
      paymentLogs.value = data.data || []
    } else {
      appErrors.value = data.data || []
    }
  } catch (e) {
    console.error('Failed fetching logs', e)
    // If backend is not ready, mock empty array
    if (activeTab.value === 'audit') auditLogs.value = []
    if (activeTab.value === 'payment') paymentLogs.value = []
    if (activeTab.value === 'error') appErrors.value = []
  } finally {
    loading.value = false
  }
}

const switchTab = (tab: string) => {
  if (activeTab.value !== tab) {
    activeTab.value = tab
    page.value = 1
    fetchLogs()
  }
}

const changePage = (delta: number) => {
  const newPage = page.value + delta
  if (newPage >= 1) {
    page.value = newPage
    fetchLogs()
  }
}

const formatDate = (isoString: string) => {
  if (!isoString) return '-'
  const d = new Date(isoString)
  return d.toLocaleString('id-ID', { 
    year: 'numeric', month: 'short', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  })
}

const viewJson = (log: any) => {
  selectedLog.value = log
  showJsonModal.value = true
}

const viewStackTrace = (log: any) => {
  selectedLog.value = log
  showTraceModal.value = true
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
/* Scoped overrides to ensure clean transitions */
</style>
