<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Шапка -->
    <header class="bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white py-6">
      <div class="container mx-auto px-4">
        <div class="flex justify-between items-center">
          <div>
            <h1 class="text-2xl font-bold">Contractors</h1>
            <div class="mt-2 flex gap-4 text-sm">
              <span>👷 Всего: {{ contractors.length }}</span>
              <span>🔧 Активных сейчас: {{ activeCount }}</span>
            </div>
          </div>
          <router-link
              to="/admin"
              class="bg-white/20 hover:bg-white/30 text-white text-sm font-semibold px-4 py-2 rounded-lg"
          >
            ← Orders
          </router-link>
        </div>
      </div>
    </header>

    <!-- Таблица -->
    <main class="container mx-auto px-4 py-8">
      <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Имя</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Телефон</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Зарегистрирован</th>
            <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Взял</th>
            <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Продал</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Статус</th>
          </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="c in contractors" :key="c.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 text-sm font-medium text-gray-900">{{ c.id }}</td>
            <td class="px-6 py-4 text-sm text-gray-900 font-semibold">{{ c.name }}</td>
            <td class="px-6 py-4 text-sm text-gray-500">{{ c.phone }}</td>
            <td class="px-6 py-4 text-sm text-gray-500">{{ c.email }}</td>
            <td class="px-6 py-4 text-sm text-gray-500">{{ formatDate(c.created_at) }}</td>
            <td class="px-6 py-4 text-sm text-center font-semibold text-gray-900">{{ c.orders_taken }}</td>
            <td class="px-6 py-4 text-sm text-center font-semibold text-green-600">{{ c.orders_sold }}</td>
            <td class="px-6 py-4">
              <span v-if="c.active_order_id" class="px-2 py-1 text-xs font-semibold rounded-full bg-yellow-100 text-yellow-800">
                🔧 Order #{{ c.active_order_id }}
              </span>
              <span v-else class="px-2 py-1 text-xs font-semibold rounded-full bg-gray-100 text-gray-500">
                Free
              </span>
            </td>
          </tr>
          </tbody>
        </table>

        <div v-if="contractors.length === 0" class="text-center py-12 text-gray-500">
          Подрядчиков пока нет
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
import { ref, computed, onMounted, onUnmounted } from 'vue'

const contractors = ref([])

const activeCount = computed(() =>
    contractors.value.filter(c => c.active_order_id !== null).length
)

async function loadContractors() {
  try {
    const response = await fetch(`${API_BASE_URL}/api/contractors`)
    contractors.value = await response.json() || []
  } catch (e) {
    console.error('Ошибка загрузки подрядчиков:', e)
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '—'
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

let refreshInterval = null

onMounted(() => {
  loadContractors()
  refreshInterval = setInterval(loadContractors, 15000)
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>