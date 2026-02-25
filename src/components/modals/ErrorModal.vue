<template>
  <BaseModal :open="open" @close="emit('close')">
    <div class="px-6 py-8 text-center">
      <h3 class="text-lg font-semibold text-gray-800">
        Something went wrong ❌
      </h3>

      <p class="mt-3 text-sm text-gray-600">
        {{ message || "Failed to send the request. Please try again." }}
      </p>

      <button
          class="mt-6 mx-auto block rounded-full bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] px-6 py-3 text-sm font-semibold text-white"
          @click="emit('close')"
      >
        OK
      </button>
    </div>
  </BaseModal>
</template>

<script setup>
import BaseModal from "@/components/layout/BaseModal.vue";

defineProps({
  open: { type: Boolean, default: false },
  message: { type: String, default: "" },
});

const emit = defineEmits(["close"]);
</script>
