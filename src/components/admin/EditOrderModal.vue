<template>
  <BaseModal :open="open" @close="emit('close')">
    <div class="px-6 py-5">
      <h2 class="text-xl font-bold text-gray-900 mb-4">Изменить статус заказа №{{ order?.id }}</h2>

      <div class="space-y-3">
        <label class="flex items-center space-x-3 cursor-pointer">
          <input
              type="radio"
              v-model="selectedStatus"
              value="new"
              class="w-4 h-4 text-blue-600"
          />
          <span class="text-base">🔵 New</span>
        </label>

        <label class="flex items-center space-x-3 cursor-pointer">
          <input
              type="radio"
              v-model="selectedStatus"
              value="in_progress"
              class="w-4 h-4 text-yellow-600"
          />
          <span class="text-base">🟡 In Progress</span>
        </label>

        <label class="flex items-center space-x-3 cursor-pointer">
          <input
              type="radio"
              v-model="selectedStatus"
              value="completed"
              class="w-4 h-4 text-green-600"
          />
          <span class="text-base">🟢 Ready</span>
        </label>

        <label class="flex items-center space-x-3 cursor-pointer">
          <input
              type="radio"
              v-model="selectedStatus"
              value="cancelled"
              class="w-4 h-4 text-red-600"
          />
          <span class="text-base">🔴 Cancel</span>
        </label>
      </div>

      <div class="mt-6 flex justify-end space-x-3">
        <button
            @click="emit('close')"
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300"
        >
          Отмена
        </button>
        <button
            @click="saveStatus"
            class="px-4 py-2 bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white rounded-md hover:opacity-90"
        >
          Сохранить
        </button>
      </div>
    </div>
  </BaseModal>
</template>

<script setup>
import { ref, watch } from 'vue'
import BaseModal from '@/components/layout/BaseModal.vue'

const props = defineProps({
  open: Boolean,
  order: Object
})

const emit = defineEmits(['close', 'save'])

const selectedStatus = ref('')

watch(() => props.order, (newOrder) => {
  if (newOrder) {
    selectedStatus.value = newOrder.status
  }
})

function saveStatus() {
  emit('save', selectedStatus.value)
}
</script>