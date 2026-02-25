<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white py-6">
      <div class="container mx-auto px-4 flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold">Contractor Dashboard</h1>
          <p class="text-sm mt-1">Welcome, {{ contractorName }}!</p>
        </div>
        <button
            @click="logout"
            class="px-4 py-2 bg-white/20 hover:bg-white/30 rounded-md text-sm font-medium"
        >
          Logout
        </button>
      </div>
    </header>

    <!-- Content -->
    <main class="container mx-auto px-4 py-8 space-y-12">

      <!-- Available Orders Section -->
      <section>
        <h2 class="text-xl font-bold text-gray-900 mb-6">
          Available Orders ({{ availableOrders.length }})
        </h2>

        <!-- Loading -->
        <div v-if="loadingOrders" class="text-center py-12">
          <p class="text-gray-500">Loading orders...</p>
        </div>

        <!-- No orders -->
        <div v-else-if="availableOrders.length === 0" class="text-center py-12 bg-white rounded-lg">
          <p class="text-gray-500">No orders available at the moment</p>
        </div>

        <!-- Orders grid -->
        <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <div
              v-for="order in availableOrders"
              :key="order.id"
              class="bg-white rounded-lg shadow p-6 space-y-4"
          >
            <div class="flex justify-between items-start">
              <h3 class="text-lg font-semibold text-gray-900">Order #{{ order.id }}</h3>
              <span class="px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                {{ order.status }}
              </span>
            </div>

            <div class="space-y-2 text-sm">
              <div>
                <span class="font-medium text-gray-700">Client:</span>
                <span class="text-gray-900 ml-2">{{ order.client_name }}</span>
              </div>
              <div>
                <span class="font-medium text-gray-700">Device:</span>
                <span class="text-gray-900 ml-2">{{ order.device }}</span>
              </div>
              <div>
                <span class="font-medium text-gray-700">Problem:</span>
                <span class="text-gray-900 ml-2">{{ order.problem || 'Not specified' }}</span>
              </div>
              <div>
                <span class="font-medium text-gray-700">ZIP:</span>
                <span class="text-gray-900 ml-2">{{ order.zip_code }}</span>
              </div>
            </div>

            <!-- Bid form -->
            <div class="pt-4 border-t">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                When can you complete this job?
              </label>
              <select
                  v-model="bidTimes[order.id]"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md mb-3 text-sm"
              >
                <option value="">Select time...</option>
                <option value="Today">Today</option>
                <option value="Tomorrow">Tomorrow</option>
                <option value="In 2 days">In 2 days</option>
                <option value="In 3 days">In 3 days</option>
              </select>

              <button
                  @click="placeBid(order.id)"
                  :disabled="!bidTimes[order.id] || submitting[order.id]"
                  class="w-full py-2 px-4 bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white font-semibold rounded-md hover:opacity-90 disabled:opacity-50"
              >
                {{ submitting[order.id] ? 'Submitting...' : 'Place Bid' }}
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- My Active Bids Section -->
      <section>
        <h2 class="text-xl font-bold text-gray-900 mb-6">
          My Active Bids ({{ myBids.length }})
        </h2>

        <!-- Loading -->
        <div v-if="loadingBids" class="text-center py-12">
          <p class="text-gray-500">Loading your bids...</p>
        </div>

        <!-- No bids -->
        <div v-else-if="myBids.length === 0" class="text-center py-12 bg-white rounded-lg">
          <p class="text-gray-500">You haven't placed any bids yet</p>
        </div>

        <!-- Bids grid -->
        <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <div
              v-for="item in myBids"
              :key="item.bid.id"
              class="bg-white rounded-lg shadow p-6 space-y-4 border-l-4 border-yellow-400"
          >
            <div class="flex justify-between items-start">
              <h3 class="text-lg font-semibold text-gray-900">Order #{{ item.order.id }}</h3>
              <span class="px-2 py-1 text-xs font-semibold rounded-full bg-yellow-100 text-yellow-800">
                {{ item.order.status }}
              </span>
            </div>

            <div class="space-y-2 text-sm">
              <div>
                <span class="font-medium text-gray-700">Client:</span>
                <span class="text-gray-900 ml-2">{{ item.order.client_name }}</span>
              </div>
              <div>
                <span class="font-medium text-gray-700">Device:</span>
                <span class="text-gray-900 ml-2">{{ item.order.device }}</span>
              </div>
              <div>
                <span class="font-medium text-gray-700">Problem:</span>
                <span class="text-gray-900 ml-2">{{ item.order.problem || 'Not specified' }}</span>
              </div>
              <div>
                <span class="font-medium text-gray-700">ZIP:</span>
                <span class="text-gray-900 ml-2">{{ item.order.zip_code }}</span>
              </div>
            </div>

            <!-- Bid info -->
            <div class="pt-4 border-t bg-yellow-50 -mx-6 -mb-6 px-6 py-4 rounded-b-lg">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-xs text-gray-600">Your bid:</p>
                  <p class="font-semibold text-gray-900">{{ item.bid.proposed_time }}</p>
                </div>
                <div class="text-right">
                  <p class="text-xs text-gray-600">Status:</p>
                  <p class="font-semibold text-yellow-600">Waiting</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const contractorId = localStorage.getItem('contractorId')
const contractorName = ref(localStorage.getItem('contractorName') || 'Contractor')

const allOrders = ref([])
const myBids = ref([])
const loadingOrders = ref(true)
const loadingBids = ref(true)
const bidTimes = ref({})
const submitting = ref({})

// Check auth
if (!contractorId) {
  router.push('/contractors/login')
}

// Available orders = все заказы "confirmed" минус те на которые уже сделали ставку
const availableOrders = computed(() => {
  const bidOrderIds = myBids.value.map(item => item.order.id)
  return allOrders.value.filter(order => !bidOrderIds.includes(order.id))
})

async function loadOrders() {
  loadingOrders.value = true
  try {
    const response = await fetch('http://localhost:8080/api/orders/available')
    allOrders.value = await response.json()
  } catch (e) {
    console.error('Error loading orders:', e)
  } finally {
    loadingOrders.value = false
  }
}

async function loadMyBids() {
  loadingBids.value = true
  try {
    const response = await fetch(`http://localhost:8080/api/contractors/${contractorId}/bids`)
    myBids.value = await response.json() || []
  } catch (e) {
    console.error('Error loading bids:', e)
  } finally {
    loadingBids.value = false
  }
}

async function placeBid(orderId) {
  const proposedTime = bidTimes.value[orderId]
  if (!proposedTime) return

  submitting.value[orderId] = true

  try {
    const response = await fetch('http://localhost:8080/api/bids', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        order_id: orderId,
        contractor_id: parseInt(contractorId),
        proposed_time: proposedTime
      })
    })

    if (!response.ok) {
      const text = await response.text()
      alert(text || 'Failed to place bid')
      return
    }

    alert('Bid placed successfully!')

    // Reload bids
    await loadMyBids()

    // Clear selection
    delete bidTimes.value[orderId]

  } catch (e) {
    alert('Connection error')
    console.error(e)
  } finally {
    submitting.value[orderId] = false
  }
}

function logout() {
  localStorage.removeItem('contractorId')
  localStorage.removeItem('contractorName')
  router.push('/contractors/login')
}

onMounted(() => {
  loadOrders()
  loadMyBids()
})
</script>








<!--<template>-->
<!--  <div class="min-h-screen bg-gray-50">-->
<!--    &lt;!&ndash; Header &ndash;&gt;-->
<!--    <header class="bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white py-6">-->
<!--      <div class="container mx-auto px-4 flex justify-between items-center">-->
<!--        <div>-->
<!--          <h1 class="text-2xl font-bold">Contractor Dashboard</h1>-->
<!--          <p class="text-sm mt-1">Welcome, {{ contractorName }}!</p>-->
<!--        </div>-->
<!--        <button-->
<!--            @click="logout"-->
<!--            class="px-4 py-2 bg-white/20 hover:bg-white/30 rounded-md text-sm font-medium"-->
<!--        >-->
<!--          Logout-->
<!--        </button>-->
<!--      </div>-->
<!--    </header>-->

<!--    &lt;!&ndash; Content &ndash;&gt;-->
<!--    <main class="container mx-auto px-4 py-8">-->
<!--      <h2 class="text-xl font-bold text-gray-900 mb-6">Available Orders</h2>-->

<!--      &lt;!&ndash; Loading &ndash;&gt;-->
<!--      <div v-if="loading" class="text-center py-12">-->
<!--        <p class="text-gray-500">Loading orders...</p>-->
<!--      </div>-->

<!--      &lt;!&ndash; No orders &ndash;&gt;-->
<!--      <div v-else-if="orders.length === 0" class="text-center py-12">-->
<!--        <p class="text-gray-500">No orders available at the moment</p>-->
<!--      </div>-->

<!--      &lt;!&ndash; Orders list &ndash;&gt;-->
<!--      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">-->
<!--        <div-->
<!--            v-for="order in orders"-->
<!--            :key="order.id"-->
<!--            class="bg-white rounded-lg shadow p-6 space-y-4"-->
<!--        >-->
<!--          <div class="flex justify-between items-start">-->
<!--            <h3 class="text-lg font-semibold text-gray-900">Order #{{ order.id }}</h3>-->
<!--            <span class="px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">-->
<!--              {{ order.status }}-->
<!--            </span>-->
<!--          </div>-->

<!--          <div class="space-y-2 text-sm">-->
<!--            <div>-->
<!--              <span class="font-medium text-gray-700">Client:</span>-->
<!--              <span class="text-gray-900">{{ order.client_name }}</span>-->
<!--            </div>-->
<!--            <div>-->
<!--              <span class="font-medium text-gray-700">Device:</span>-->
<!--              <span class="text-gray-900">{{ order.device }}</span>-->
<!--            </div>-->
<!--            <div>-->
<!--              <span class="font-medium text-gray-700">Problem:</span>-->
<!--              <span class="text-gray-900">{{ order.problem || 'Not specified' }}</span>-->
<!--            </div>-->
<!--            <div>-->
<!--              <span class="font-medium text-gray-700">ZIP:</span>-->
<!--              <span class="text-gray-900">{{ order.zip_code }}</span>-->
<!--            </div>-->
<!--          </div>-->

<!--          &lt;!&ndash; Bid form &ndash;&gt;-->
<!--          <div class="pt-4 border-t">-->
<!--            <label class="block text-sm font-medium text-gray-700 mb-2">-->
<!--              When can you complete this job?-->
<!--            </label>-->
<!--            <select-->
<!--                v-model="bidTimes[order.id]"-->
<!--                class="w-full px-3 py-2 border border-gray-300 rounded-md mb-3"-->
<!--            >-->
<!--              <option value="">Select time...</option>-->
<!--              <option value="Today">Today</option>-->
<!--              <option value="Tomorrow">Tomorrow</option>-->
<!--              <option value="In 2 days">In 2 days</option>-->
<!--              <option value="In 3 days">In 3 days</option>-->
<!--            </select>-->

<!--            <button-->
<!--                @click="placeBid(order.id)"-->
<!--                :disabled="!bidTimes[order.id] || submitting[order.id]"-->
<!--                class="w-full py-2 px-4 bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white font-semibold rounded-md hover:opacity-90 disabled:opacity-50"-->
<!--            >-->
<!--              {{ submitting[order.id] ? 'Submitting...' : 'Place Bid' }}-->
<!--            </button>-->
<!--          </div>-->
<!--        </div>-->
<!--      </div>-->
<!--    </main>-->
<!--  </div>-->
<!--</template>-->

<!--<script setup>-->
<!--import { ref, onMounted } from 'vue'-->
<!--import { useRouter } from 'vue-router'-->

<!--const router = useRouter()-->

<!--const contractorId = localStorage.getItem('contractorId')-->
<!--const contractorName = ref(localStorage.getItem('contractorName') || 'Contractor')-->

<!--const orders = ref([])-->
<!--const loading = ref(true)-->
<!--const bidTimes = ref({})-->
<!--const submitting = ref({})-->

<!--// Check auth-->
<!--if (!contractorId) {-->
<!--  router.push('/contractors/login')-->
<!--}-->

<!--async function loadOrders() {-->
<!--  loading.value = true-->
<!--  try {-->
<!--    const response = await fetch('http://localhost:8080/api/orders/available')-->
<!--    orders.value = await response.json()-->
<!--  } catch (e) {-->
<!--    console.error('Error loading orders:', e)-->
<!--  } finally {-->
<!--    loading.value = false-->
<!--  }-->
<!--}-->

<!--async function placeBid(orderId) {-->
<!--  const proposedTime = bidTimes.value[orderId]-->
<!--  if (!proposedTime) return-->

<!--  submitting.value[orderId] = true-->

<!--  try {-->
<!--    const response = await fetch('http://localhost:8080/api/bids', {-->
<!--      method: 'POST',-->
<!--      headers: {-->
<!--        'Content-Type': 'application/json'-->
<!--      },-->
<!--      body: JSON.stringify({-->
<!--        order_id: orderId,-->
<!--        contractor_id: parseInt(contractorId),-->
<!--        proposed_time: proposedTime-->
<!--      })-->
<!--    })-->

<!--    if (!response.ok) {-->
<!--      const text = await response.text()-->
<!--      alert(text || 'Failed to place bid')-->
<!--      return-->
<!--    }-->

<!--    alert('Bid placed successfully!')-->

<!--    // Remove order from list (already bid)-->
<!--    orders.value = orders.value.filter(o => o.id !== orderId)-->
<!--    delete bidTimes.value[orderId]-->

<!--  } catch (e) {-->
<!--    alert('Connection error')-->
<!--    console.error(e)-->
<!--  } finally {-->
<!--    submitting.value[orderId] = false-->
<!--  }-->
<!--}-->

<!--function logout() {-->
<!--  localStorage.removeItem('contractorId')-->
<!--  localStorage.removeItem('contractorName')-->
<!--  router.push('/contractors/login')-->
<!--}-->

<!--onMounted(() => {-->
<!--  loadOrders()-->
<!--})-->
<!--</script>-->