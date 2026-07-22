<template>
  <div class="pagination-wrapper">
    <div class="pagination-info">
      Menampilkan <strong>{{ startItem }}</strong> - <strong>{{ endItem }}</strong> dari <strong>{{ totalItems }}</strong> data
    </div>
    
    <div class="pagination-controls">
      <!-- Items per page selector -->
      <div class="per-page-selector">
        <span>Baris per halaman:</span>
        <select :value="itemsPerPage" @change="onLimitChange($event)">
          <option v-for="opt in limitOptions" :key="opt" :value="opt">{{ opt }}</option>
        </select>
      </div>

      <!-- Page navigation buttons -->
      <div class="page-buttons">
        <button 
          class="page-nav-btn" 
          :disabled="currentPage === 1" 
          @click="changePage(1)"
          title="Halaman Pertama"
        >
          <i class="ti ti-chevrons-left"></i>
        </button>
        <button 
          class="page-nav-btn" 
          :disabled="currentPage === 1" 
          @click="changePage(currentPage - 1)"
          title="Halaman Sebelumnya"
        >
          <i class="ti ti-chevron-left"></i>
        </button>
        
        <span class="page-indicator">
          Halaman <strong>{{ currentPage }}</strong> dari <strong>{{ totalPages }}</strong>
        </span>

        <button 
          class="page-nav-btn" 
          :disabled="currentPage === totalPages || totalPages === 0" 
          @click="changePage(currentPage + 1)"
          title="Halaman Berikutnya"
        >
          <i class="ti ti-chevron-right"></i>
        </button>
        <button 
          class="page-nav-btn" 
          :disabled="currentPage === totalPages || totalPages === 0" 
          @click="changePage(totalPages)"
          title="Halaman Terakhir"
        >
          <i class="ti ti-chevrons-right"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
  currentPage: { type: Number, required: true },
  totalItems: { type: Number, required: true },
  itemsPerPage: { type: Number, required: true },
  limitOptions: { type: Array as () => number[], default: () => [5, 10, 20, 50, 100] }
})

const emit = defineEmits(['update:currentPage', 'update:itemsPerPage'])

const totalPages = computed(() => {
  if (props.totalItems === 0) return 1
  return Math.ceil(props.totalItems / props.itemsPerPage)
})

const startItem = computed(() => {
  if (props.totalItems === 0) return 0
  return (props.currentPage - 1) * props.itemsPerPage + 1
})

const endItem = computed(() => {
  const calculatedEnd = props.currentPage * props.itemsPerPage
  return calculatedEnd > props.totalItems ? props.totalItems : calculatedEnd
})

const changePage = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  emit('update:currentPage', page)
}

const onLimitChange = (event: Event) => {
  const select = event.target as HTMLSelectElement
  emit('update:itemsPerPage', parseInt(select.value))
  emit('update:currentPage', 1)
}
</script>

<style scoped>
.pagination-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: var(--card);
  border-top: 1px solid var(--border);
  flex-wrap: wrap;
  gap: 12px;
  user-select: none;
}

.pagination-info {
  font-size: 12px;
  color: var(--text2);
}

.pagination-info strong {
  color: var(--text1);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.per-page-selector {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text2);
}

.per-page-selector select {
  padding: 4px 8px;
  border: 1px solid var(--border2);
  border-radius: var(--r6);
  background: var(--surface);
  color: var(--text1);
  font-size: 11px;
  outline: none;
  cursor: pointer;
}

.page-buttons {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-nav-btn {
  width: 28px;
  height: 28px;
  border: 1px solid var(--border2);
  background: var(--surface);
  color: var(--text1);
  border-radius: var(--r6);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 12px;
}

.page-nav-btn:hover:not(:disabled) {
  border-color: var(--hijau);
  color: var(--hijau);
  background: var(--hijau3);
}

.page-nav-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--border);
}

.page-indicator {
  font-size: 12px;
  color: var(--text2);
  margin: 0 4px;
}

.page-indicator strong {
  color: var(--text1);
}
</style>
