<template>
  <BaseModal :open="open" @close="emit('close')">
    <div class="px-6 py-5">
      <h2 class="text-xl font-bold text-red-600 mb-4">Удалить заказ №{{ order?.id }}?</h2>

      <div class="space-y-2 mb-6">
        <p class="text-base text-gray-700">
          <span class="font-semibold">Клиент:</span> {{ order?.client_name }}
        </p>
        <p class="text-base text-gray-700">
          <span class="font-semibold">Техника:</span> {{ order?.device }}
        </p>
      </div>

      <p class="text-sm text-gray-500 mb-6">
        Это действие нельзя отменить!
      </p>

      <div class="flex justify-end space-x-3">
        <button
            @click="emit('close')"
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300"
        >
          Отмена
        </button>
        <button
            @click="emit('delete')"
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
        >
          Удалить
        </button>
      </div>
    </div>
  </BaseModal>
</template>

<script setup>
import BaseModal from '@/components/layout/BaseModal.vue'

defineProps({
  open: Boolean,
  order: Object
})

const emit = defineEmits(['close', 'delete'])
</script>