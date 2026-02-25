<template>
  <BaseModal :open="open" @close="emit('close')">
    <div class="px-6 py-5">
      <h2 class="text-xl font-bold text-gray-900 mb-4">Заказ №{{ order?.id }}</h2>

      <div class="space-y-3">
        <div>
          <span class="text-sm font-semibold text-gray-600">Клиент:</span>
          <p class="text-base text-gray-900">{{ order?.client_name }}</p>
        </div>

        <div>
          <span class="text-sm font-semibold text-gray-600">Телефон:</span>
          <p class="text-base text-gray-900">{{ order?.phone }}</p>
        </div>

        <div>
          <span class="text-sm font-semibold text-gray-600">ZIP код:</span>
          <p class="text-base text-gray-900">{{ order?.zip_code }}</p>
        </div>

        <div>
          <span class="text-sm font-semibold text-gray-600">Техника:</span>
          <p class="text-base text-gray-900">{{ order?.device }}</p>
        </div>

        <div>
          <span class="text-sm font-semibold text-gray-600">Проблема:</span>
          <p class="text-base text-gray-900">{{ order?.problem || 'Не указана' }}</p>
        </div>

        <div>
          <span class="text-sm font-semibold text-gray-600">Статус:</span>
          <p class="text-base text-gray-900">{{ order?.status }}</p>
        </div>

        <div>
          <span class="text-sm font-semibold text-gray-600">Цена:</span>
          <p class="text-base text-gray-900">{{ order?.price || 0 }} грн</p>
        </div>
      </div>

      <div class="mt-6 flex justify-end">
        <button
            @click="emit('close')"
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300"
        >
          Закрыть
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

const emit = defineEmits(['close'])
</script>