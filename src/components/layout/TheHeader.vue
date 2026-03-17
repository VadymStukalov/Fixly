<template>
  <header class="flex max-w-5xl mx-auto justify-between px-5 py-5 items-center">
    <img src="../../assets/logo.png" alt="logo" class="w-14 md:w-20">

    <nav class="hidden md:flex gap-4 font-brand">
      <a href="#" class="font-brand font-medium hover:opacity-70">About Us</a>
      <a href="#" class="font-brand font-medium hover:opacity-70">Contact</a>
    </nav>

    <button
        class="px-[23px] py-[6px] text-sm shadow-sm font-brand font-bold rounded-full text-white bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] hover:opacity-90 transition-opacity"
        @click="showPopup = true"
    >Get a Quote</button>

    <!-- Попап с полем телефона -->
    <div
        v-if="showPopup"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 px-4"
        @click.self="showPopup = false"
    >
      <div class="bg-white rounded-2xl shadow-xl p-8 w-full max-w-sm">
        <h2 class="text-xl font-bold text-gray-900 mb-2">Get a Free Quote</h2>
        <p class="text-gray-500 text-sm mb-6">Enter your phone number and our AI assistant will call you right now to take your request.</p>

        <input
            v-model="phone"
            type="tel"
            placeholder="+1 (555) 000-0000"
            class="w-full border border-gray-300 rounded-xl px-4 py-3 text-sm mb-4 focus:outline-none focus:border-[#3BA3A9]"
            @keyup.enter="submitQuote"
        />

        <div v-if="error" class="text-red-500 text-sm mb-3">{{ error }}</div>

        <div v-if="success" class="bg-green-50 rounded-xl p-4 text-center">
          <p class="text-green-700 font-semibold">📞 Calling you now!</p>
          <p class="text-green-600 text-sm mt-1">Our AI assistant will help you with your repair request.</p>
        </div>

        <template v-else>
          <button
              @click="submitQuote"
              :disabled="loading"
              class="w-full py-3 bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white font-bold rounded-xl hover:opacity-90 disabled:opacity-50"
          >
            {{ loading ? 'Calling...' : '📞 Call Me Now' }}
          </button>
          <button
              @click="showPopup = false"
              class="w-full mt-2 py-3 text-gray-500 text-sm hover:text-gray-700"
          >
            Cancel
          </button>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup>
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
import { ref } from 'vue'

const emit = defineEmits(['open-form'])

const showPopup = ref(false)
const phone = ref('')
const loading = ref(false)
const success = ref(false)
const error = ref('')

async function submitQuote() {
  error.value = ''
  if (!phone.value.trim()) {
    error.value = 'Please enter your phone number'
    return
  }

  loading.value = true
  try {
    const response = await fetch(`${API_BASE_URL}/api/call-quote`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ phone: phone.value.trim() })
    })

    if (!response.ok) {
      error.value = 'Something went wrong. Please try again.'
      return
    }

    success.value = true
    // Закрываем попап через 4 секунды
    setTimeout(() => {
      showPopup.value = false
      success.value = false
      phone.value = ''
    }, 4000)
  } catch (e) {
    error.value = 'Connection error. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>