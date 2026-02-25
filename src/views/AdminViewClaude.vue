<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Шапка -->
    <header class="bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white py-6">
      <div class="container mx-auto px-4">
        <h1 class="text-2xl font-bold">FIXLY Admin Panel</h1>
        <div class="mt-2 flex gap-6 text-sm">
          <span>📊 Всего заказов: {{ orders.length }}</span>
          <span>🆕 Новых: {{ newOrdersCount }}</span>
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
                  {{ order.status }}
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
import { ref, computed, onMounted } from 'vue'
import ViewOrderModal from '@/components/admin/ViewOrderModal.vue'
import EditOrderModal from '@/components/admin/EditOrderModal.vue'
import DeleteOrderModal from '@/components/admin/DeleteOrderModal.vue'

const orders = ref([])
const selectedOrder = ref(null)
const isViewModalOpen = ref(false)
const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)

async function loadOrders() {
  try {
    const response = await fetch('http://localhost:8080/api/orders')
    orders.value = await response.json()
  } catch (e) {
    console.error('Ошибка загрузки:', e)
  }
}

const newOrdersCount = computed(() => {
  return orders.value.filter(o => o.status === 'новый').length
})

function getStatusClass(status) {
  const classes = {
    'new': 'bg-blue-100 text-blue-800',
    'confirmed': 'bg-blue-100 text-blue-800',
    'in_progress': 'bg-yellow-100 text-yellow-800',
    'completed': 'bg-green-100 text-green-800',
    'cancelled': 'bg-red-100 text-red-800'
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

  try{
    //Формируем данные для запроса

    const updateData={
      ...selectedOrder.value,
      status:newStatus
    }

    const response = await fetch(`http://localhost:8080/api/orders/${selectedOrder.value.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(updateData)
    })
    if (!response.ok) {
      console.error('Ошибка обновления:', response.status)
      return
    }

    // Получаем обновлённый заказ
    const updated = await response.json()
    console.log('✅ Заказ обновлён:', updated)

    // Закрываем модалку
    isEditModalOpen.value = false

    // Перезагружаем список
    await loadOrders()
  } catch (e){
    console.error('Ошибка:', e)
  }

}

async function handleDelete() {

  try {
    const response = await fetch(`http://localhost:8080/api/orders/${selectedOrder.value.id}`, {
      method: 'DELETE'
    })
    console.log('Удалить:', selectedOrder.value.id)
    if (!response.ok) {
      console.error('Ошибка удаления:', response.status)
      return
    }

    // Успех! Закрываем модалку
    isDeleteModalOpen.value = false

    // Перезагружаем список заказов
    await loadOrders()

    console.log('✅ Заказ удалён!')
  } catch (e) {
    console.error('Ошибка:', e)
  }

  //
  // // TODO: здесь будет API запрос на бэк
  // isDeleteModalOpen.value = false
}

onMounted(() => {
  loadOrders()
})
</script>