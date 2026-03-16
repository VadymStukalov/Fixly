<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center p-4">

    <!-- Loading -->
    <div v-if="loading" class="text-center">
      <p class="text-gray-500">Loading job details...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="bg-white rounded-2xl shadow p-8 max-w-md w-full text-center">
      <h2 class="text-xl font-bold text-red-600 mb-2">Link not valid</h2>
      <p class="text-gray-500">{{ error }}</p>
    </div>

    <!-- Job details -->
    <div v-else-if="order" class="bg-white rounded-2xl shadow p-8 max-w-md w-full">
      <div class="text-center mb-6">
        <h1 class="text-2xl font-bold text-gray-900">New Job Available</h1>
        <p class="text-gray-500 text-sm mt-1">First to accept gets the job</p>
      </div>

      <div class="space-y-3 mb-8">
        <div class="flex justify-between py-2 border-b">
          <span class="text-gray-500">Device</span>
          <span class="font-semibold text-gray-900">{{ order.device }}</span>
        </div>
        <div class="flex flex-col py-2 border-b">
          <span class="text-gray-500 mb-1">Problem</span>
          <span class="font-semibold text-gray-900">{{ order.problem || 'Not specified' }}</span>
        </div>
        <div class="flex justify-between py-2 border-b">
          <span class="text-gray-500">ZIP Code</span>
          <span class="font-semibold text-gray-900">{{ order.zip_code }}</span>
        </div>
      </div>

      <!-- Success state -->
      <div v-if="accepted" class="text-center">
        <div class="text-4xl mb-3">✅</div>
        <h2 class="text-xl font-bold text-green-600 mb-2">Job Accepted!</h2>
        <p class="text-gray-600 text-sm mb-4">Contact the client:</p>
        <div class="bg-green-50 rounded-xl p-4 space-y-2 mb-4">
          <p class="font-bold text-gray-900 text-lg">{{ order.client_name }}</p>
          <p class="text-gray-700">{{ order.zip_code }}</p>
        </div>

        <!-- Таймер 15 минут -->
        <div class="bg-yellow-50 rounded-xl p-3 mb-4">
          <p class="text-yellow-700 text-sm font-semibold">⏱ Call the client within 15 minutes</p>
        </div>

        <!-- Клиент недоступен — администратор разбирается -->
        <div v-if="clientUnreachable" class="bg-orange-50 border border-orange-200 rounded-xl p-4 mb-4">
          <p class="text-orange-700 font-semibold text-base">📵 Client is not answering</p>
          <p class="text-orange-600 text-sm mt-2">
            Our administrator has been notified and will try to reach the client manually.
            We will update you on this order shortly.
          </p>
        </div>

        <!-- Кнопка звонка — скрываем если клиент недоступен -->
        <template v-else>
          <button
              v-if="!callInitiated"
              @click="callClient"
              :disabled="calling"
              class="w-full py-4 bg-green-600 hover:bg-green-700 text-white font-bold text-lg rounded-xl disabled:opacity-50"
          >
            {{ calling ? 'Connecting...' : '📞 Call Client' }}
          </button>

          <div v-if="callInitiated" class="bg-blue-50 rounded-xl p-4 mt-4">
            <p class="text-blue-700 font-semibold">📞 Calling you now...</p>
            <p class="text-blue-500 text-sm mt-1">Answer your phone to connect with the client</p>
            <button
                @click="callAgain"
                class="mt-3 w-full py-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl text-sm"
            >
              Call Again
            </button>
          </div>
        </template>
      </div>

      <!-- Accept button -->
      <button
          v-else
          @click="acceptJob"
          :disabled="accepting"
          class="w-full py-4 bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white font-bold text-lg rounded-xl hover:opacity-90 disabled:opacity-50"
      >
        {{ accepting ? 'Accepting...' : '✓ Accept Job' }}
      </button>
    </div>

  </div>
</template>

<script setup>
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const token = route.params.token

const order = ref(null)
const loading = ref(true)
const error = ref('')
const accepting = ref(false)
const accepted = ref(false)
const calling = ref(false)
const callInitiated = ref(false)

const contractorPhone = ref('')
const contractorId = ref(0)
const clientUnreachable = ref(false)

// Polling статуса заказа каждые 15 сек — чтобы подрядчик видел обновление без перезагрузки
let pollingInterval = null

function startPolling() {
  pollingInterval = setInterval(async () => {
    if (!order.value) return
    try {
      const response = await fetch(`${API_BASE_URL}/accept/${token}`)
      if (response.ok) {
        const data = await response.json()
        if (data.status === 'client_unreachable') {
          clientUnreachable.value = true
          clearInterval(pollingInterval)
        }
      }
    } catch (e) {
      // тихо игнорируем ошибки polling
    }
  }, 15000)
}

onMounted(async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/accept/${token}`)
    if (!response.ok) {
      error.value = 'This link is no longer valid or has already been used.'
      return
    }
    order.value = await response.json()
  } catch (e) {
    error.value = 'Connection error. Please try again.'
  } finally {
    loading.value = false
  }
})

async function acceptJob() {
  accepting.value = true
  try {
    const response = await fetch(`${API_BASE_URL}/accept/${token}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    })

    if (!response.ok) {
      const text = await response.text()
      error.value = text || 'Failed to accept job. Someone else may have taken it.'
      return
    }

    const data = await response.json()
    contractorPhone.value = data.contractor_phone
    contractorId.value = data.contractor_id
    accepted.value = true
    startPolling()
  } catch (e) {
    error.value = 'Connection error.'
  } finally {
    accepting.value = false
  }
}

async function callClient() {
  calling.value = true
  try {
    const response = await fetch(`${API_BASE_URL}/api/call`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        contractor_phone: contractorPhone.value,
        client_phone: order.value.phone,
        order_id: order.value.id,
        contractor_id: contractorId.value  // ← передаём для StatusCallback
      })
    })

    if (!response.ok) {
      alert('Failed to initiate call')
      return
    }

    callInitiated.value = true
  } catch (e) {
    alert('Connection error')
  } finally {
    calling.value = false
  }
}

// Повторный звонок — сбрасываем callInitiated чтобы показать кнопку снова
function callAgain() {
  callInitiated.value = false
}

onUnmounted(() => {
  if (pollingInterval) clearInterval(pollingInterval)
})
</script>