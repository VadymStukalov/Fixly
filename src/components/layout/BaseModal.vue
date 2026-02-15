  <template>
    <Transition name="fade">
    <div
        v-if="open"
        class="fixed inset-0 z-50 bg-black/40 p-4 flex items-center justify-center"
        @click.self = "emit('close')"
    >
      <div class="w-full max-w-[420px] md:max-w-[560px] rounded-2xl bg-white shadow-xl">
        <slot />
      </div>
    </div>
    </Transition>
  </template>


  <script setup>
  const props = defineProps({
    open: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['close'])
  </script>


  <style scoped>
  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.2s ease;
  }

  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
  }

  </style>