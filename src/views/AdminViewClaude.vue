<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Шапка -->
    <header class="bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white py-6">
      <div class="container mx-auto px-4">
        <div class="flex justify-between items-start">
          <div>
            <h1 class="text-2xl font-bold">FIXLY Admin Panel</h1>
            <div class="mt-2 flex flex-wrap gap-4 text-sm">
              <span>📊 Всего: {{ orders.length }}</span>
              <span>🆕 Новых: {{ countByStatus('new') }}</span>
              <span>✅ Confirmed: {{ countByStatus('confirmed') }}</span>
              <span>🔧 In Progress: {{ countByStatus('in_progress') }}</span>
              <span>💰 Lead Sold: {{ countByStatus('lead_sold') }}</span>
              <span>📵 Unreachable: {{ countByStatus('client_unreachable') }}</span>
            </div>
          </div>
          <router-link
              to="/admin/contractors"
              class="bg-white/20 hover:bg-white/30 text-white text-sm font-semibold px-4 py-2 rounded-lg"
          >
            👷 Contractors ({{ contractorsCount }})
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
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Клиент</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Телефон</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Техника</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Статус</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Действия</th>
          </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="order in orders" :key="order.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 text-sm font-medium text-gray-900">{{ order.id }}</td>
            <td class="px-6 py-4 text-sm text-gray-900">{{ order.client_name }}</td>
            <td class="px-6 py-4 text-sm text-gray-500">{{ order.phone }}</td>
            <td class="px-6 py-4 text-sm text-gray-500">{{ order.device }}</td>
            <td class="px-6 py-4">
              <span :class="getStatusClass(order.status)" class="px-2 py-1 text-xs font-semibold rounded-full">
                {{ getStatusLabel(order.status) }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm font-medium space-x-3">
              <button @click="viewOrder(order)" class="text-blue-600 hover:text-blue-900">👁</button>
              <button @click="editOrder(order)" class="text-yellow-600 hover:text-yellow-900">✏️</button>
              <button @click="confirmDelete(order)" class="text-red-600 hover:text-red-900">🗑</button>
            </td>
          </tr>
          </tbody>
        </table>

        <div v-if="orders.length === 0" class="text-center py-12 text-gray-500">
          Заказов пока нет
        </div>
      </div>
    </main>

    <!-- Модалки -->
    <ViewOrderModal
        :open="isViewModalOpen"
        :order="selectedOrder"
        @close="isViewModalOpen = false"
    />

    <EditOrderModal
        :open="isEditModalOpen"
        :order="selectedOrder"
        @close="isEditModalOpen = false"
        @save="handleStatusUpdate"
    />

    <DeleteOrderModal
        :open="isDeleteModalOpen"
        :order="selectedOrder"
        @close="isDeleteModalOpen = false"
        @delete="handleDelete"
    />
  </div>
</template>

<script setup>
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
import { ref, computed, onMounted, onUnmounted } from 'vue'
import ViewOrderModal from '@/components/admin/ViewOrderModal.vue'
import EditOrderModal from '@/components/admin/EditOrderModal.vue'
import DeleteOrderModal from '@/components/admin/DeleteOrderModal.vue'

const orders = ref([])
const contractorsCount = ref(0)
const selectedOrder = ref(null)
const isViewModalOpen = ref(false)
const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)

async function loadOrders() {
  try {
    const response = await fetch(`${API_BASE_URL}/api/orders`)
    orders.value = await response.json()
  } catch (e) {
    console.error('Ошибка загрузки:', e)
  }
}

async function loadContractorsCount() {
  try {
    const response = await fetch(`${API_BASE_URL}/api/contractors`)
    const data = await response.json()
    contractorsCount.value = data ? data.length : 0
  } catch (e) {
    console.error('Ошибка загрузки подрядчиков:', e)
  }
}

function countByStatus(status) {
  return orders.value.filter(o => o.status === status).length
}

function getStatusLabel(status) {
  const labels = {
    'new': 'New',
    'confirmed': 'Confirmed',
    'in_progress': 'In Progress',
    'lead_sold': 'Lead Sold',
    'completed': 'Completed',
    'cancelled': 'Cancelled',
    'client_unreachable': 'Unreachable',
    'reassign': 'Reassign',
  }
  return labels[status] || status
}

function getStatusClass(status) {
  const classes = {
    'new': 'bg-blue-100 text-blue-800',
    'confirmed': 'bg-indigo-100 text-indigo-800',
    'in_progress': 'bg-yellow-100 text-yellow-800',
    'lead_sold': 'bg-green-100 text-green-800',
    'completed': 'bg-green-100 text-green-800',
    'cancelled': 'bg-red-100 text-red-800',
    'client_unreachable': 'bg-orange-100 text-orange-800',
    'reassign': 'bg-gray-100 text-gray-800',
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

function viewOrder(order) {
  selectedOrder.value = order
  isViewModalOpen.value = true
}

function editOrder(order) {
  selectedOrder.value = order
  isEditModalOpen.value = true
}

function confirmDelete(order) {
  selectedOrder.value = order
  isDeleteModalOpen.value = true
}

async function handleStatusUpdate(newStatus) {
  try {
    const updateData = {
      ...selectedOrder.value,
      status: newStatus
    }

    const response = await fetch(`${API_BASE_URL}/api/orders/${selectedOrder.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(updateData)
    })

    if (!response.ok) {
      console.error('Ошибка обновления:', response.status)
      return
    }

    isEditModalOpen.value = false
    await loadOrders()
  } catch (e) {
    console.error('Ошибка:', e)
  }
}

async function handleDelete() {
  try {
    const response = await fetch(`${API_BASE_URL}/api/orders/${selectedOrder.value.id}`, {
      method: 'DELETE'
    })

    if (!response.ok) {
      console.error('Ошибка удаления:', response.status)
      return
    }

    isDeleteModalOpen.value = false
    await loadOrders()
  } catch (e) {
    console.error('Ошибка:', e)
  }
}

let refreshInterval = null

onMounted(() => {
  loadOrders()
  loadContractorsCount()
  refreshInterval = setInterval(() => {
    loadOrders()
  }, 15000)
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>