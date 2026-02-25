<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
    <div class="max-w-md w-full space-y-8">
      <!-- Header -->
      <div>
        <h2 class="text-center text-3xl font-bold text-gray-900">
          Contractor Registration
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Already have an account?
          <router-link to="/contractors/login" class="font-medium text-[#3BA3A9] hover:underline">
            Login
          </router-link>
        </p>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleRegister" class="mt-8 space-y-6">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">Full Name</label>
            <input
                v-model="form.name"
                type="text"
                required
                class="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#3BA3A9]"
                placeholder="John Smith"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700">Email</label>
            <input
                v-model="form.email"
                type="email"
                required
                class="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#3BA3A9]"
                placeholder="john@example.com"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700">Phone</label>
            <input
                v-model="form.phone"
                type="tel"
                required
                class="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#3BA3A9]"
                placeholder="+1 (555) 123-4567"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700">Password</label>
            <input
                v-model="form.password"
                type="password"
                required
                minlength="6"
                class="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#3BA3A9]"
                placeholder="Minimum 6 characters"
            />
          </div>
        </div>

        <div v-if="error" class="text-red-600 text-sm">
          {{ error }}
        </div>

        <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 px-4 bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] text-white font-semibold rounded-md hover:opacity-90 disabled:opacity-50"
        >
          {{ loading ? 'Signing up...' : 'Sign Up' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = ref({
  name: '',
  email: '',
  phone: '',
  password: ''
})

const loading = ref(false)
const error = ref('')

async function handleRegister() {
  loading.value = true
  error.value = ''

  try {
    const response = await fetch('http://localhost:8080/api/contractors/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(form.value)
    })

    if (!response.ok) {
      const text = await response.text()
      error.value = text || 'Registration failed'
      return
    }

    const contractor = await response.json()

    // Save contractor ID
    localStorage.setItem('contractorId', contractor.id)
    localStorage.setItem('contractorName', contractor.name)

    // Redirect to dashboard
    router.push('/contractors/dashboard')

  } catch (e) {
    error.value = 'Connection error'
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>