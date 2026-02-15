<template>
  <BaseModal :open="open" @close="emit('close')">
    <div class="px-6 py-8 text-center">
      <h3 class="text-lg font-semibold text-gray-800">
        Request sent successfully 🎉
      </h3>

      <p class="mt-3 text-sm text-gray-600">
        Our specialist will contact you shortly.
      </p>

      <button
          class="mt-6 mx-auto block rounded-full bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] px-6 py-3 text-sm font-semibold text-white"
          @click="emit('close')"
      >
        Close
      </button>
    </div>
  </BaseModal>
</template>

<script setup>
import BaseModal from "@/components/layout/BaseModal.vue"

defineProps({
  open: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])
</script>
