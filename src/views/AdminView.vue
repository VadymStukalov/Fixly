<script setup>
import { computed, onMounted, ref } from "vue";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

const orders = ref([]);
const loading = ref(false);
const error = ref("");

const selected = ref(null); // выбранный заказ (для деталей)

async function fetchOrders() {
  loading.value = true;
  error.value = "";
  try {
    const res = await fetch(`${API_BASE_URL}/api/orders`);
    if (!res.ok) throw new Error(await res.text());
    orders.value = await res.json();
  } catch (e) {
    error.value = e?.message || "Failed to load orders";
  } finally {
    loading.value = false;
  }
}

function openDetails(order) {
  selected.value = { ...order }; // копия, чтобы менять локально
}

function closeDetails() {
  selected.value = null;
}

async function saveStatus() {
  if (!selected.value) return;

  try {
    const id = selected.value.id;

    const res = await fetch(`${API_BASE_URL}/api/orders?id=${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(selected.value),
    });

    if (!res.ok) throw new Error(await res.text());

    await fetchOrders();
    closeDetails();
  } catch (e) {
    alert(e?.message || "Failed to update");
  }
}

async function deleteOrder(id) {
  if (!confirm(`Delete order #${id}?`)) return;

  try {
    const res = await fetch(`${API_BASE_URL}/api/orders?id=${id}`, {
      method: "DELETE",
    });
    if (!res.ok && res.status !== 204) throw new Error(await res.text());
    await fetchOrders();
    if (selected.value?.id === id) closeDetails();
  } catch (e) {
    alert(e?.message || "Failed to delete");
  }
}

const sortedOrders = computed(() =>
    [...orders.value].sort((a, b) => b.id - a.id)
);

onMounted(fetchOrders);
</script>

<template>
  <div class="max-w-6xl mx-auto p-4">
    <div class="flex items-center justify-between gap-3">
      <h1 class="text-xl font-semibold">Admin · Orders</h1>

      <button
          class="rounded-md px-4 py-2 text-sm font-semibold text-white bg-black disabled:opacity-50"
          :disabled="loading"
          @click="fetchOrders"
      >
        Refresh
      </button>
    </div>

    <p v-if="error" class="mt-3 text-sm text-red-600">{{ error }}</p>
    <p v-if="loading" class="mt-3 text-sm text-gray-600">Loading...</p>

    <div v-if="!loading" class="mt-4 overflow-x-auto rounded-lg border border-gray-200">
      <table class="w-full text-left text-sm">
        <thead class="bg-gray-50">
        <tr>
          <th class="p-3">ID</th>
          <th class="p-3">Client</th>
          <th class="p-3">Phone</th>
          <th class="p-3">Device</th>
          <th class="p-3">ZIP</th>
          <th class="p-3">Status</th>
          <th class="p-3 text-right">Actions</th>
        </tr>
        </thead>

        <tbody>
        <tr v-for="o in sortedOrders" :key="o.id" class="border-t hover:bg-gray-50">
          <td class="p-3">#{{ o.id }}</td>
          <td class="p-3">{{ o.client_name }}</td>
          <td class="p-3">{{ o.phone }}</td>
          <td class="p-3">{{ o.device }}</td>
          <td class="p-3">{{ o.zip_code }}</td>
          <td class="p-3">
            <span class="rounded px-2 py-1 text-xs bg-gray-100">{{ o.status }}</span>
          </td>
          <td class="p-3 text-right">
            <button class="text-blue-600 hover:underline mr-3" @click="openDetails(o)">
              View
            </button>
            <button class="text-red-600 hover:underline" @click="deleteOrder(o.id)">
              Delete
            </button>
          </td>
        </tr>

        <tr v-if="sortedOrders.length === 0">
          <td class="p-6 text-gray-500" colspan="7">No orders yet.</td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Details modal (простая, без BaseModal, чтобы быстрее) -->
    <div
        v-if="selected"
        class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4"
        @click.self="closeDetails"
    >
      <div class="w-full max-w-lg rounded-2xl bg-white p-5">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Order #{{ selected.id }}</h2>
          <button class="text-sm text-gray-600" @click="closeDetails">Close</button>
        </div>

        <div class="mt-4 space-y-3 text-sm">
          <div><b>Client:</b> {{ selected.client_name }}</div>
          <div><b>Phone:</b> {{ selected.phone }}</div>
          <div><b>Device:</b> {{ selected.device }}</div>
          <div><b>ZIP:</b> {{ selected.zip_code }}</div>

          <div>
            <b>Problem:</b>
            <div class="mt-1 rounded bg-gray-50 p-3 whitespace-pre-wrap">
              {{ selected.problem || "—" }}
            </div>
          </div>

          <div class="flex items-center gap-3">
            <b>Status:</b>
            <select v-model="selected.status" class="h-10 rounded border px-3">
              <option value="new">new</option>
              <option value="in_progress">in_progress</option>
              <option value="done">done</option>
            </select>
          </div>

          <div class="flex items-center gap-3">
            <b>Price:</b>
            <input
                v-model.number="selected.price"
                type="number"
                class="h-10 rounded border px-3 w-32"
                min="0"
            />
          </div>
        </div>

        <div class="mt-5 flex justify-end gap-3">
          <button class="rounded px-4 py-2 text-sm border" @click="closeDetails">
            Cancel
          </button>
          <button class="rounded px-4 py-2 text-sm font-semibold text-white bg-black" @click="saveStatus">
            Save
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
