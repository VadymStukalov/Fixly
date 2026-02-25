<script setup>
import { ref, watch } from 'vue'
// import { RouterLink, RouterView } from 'vue-router'
import TheHeader from "@/components/layout/TheHeader.vue";
import TheHero from "@/components/sections/TheHero.vue";
import HowItWorks from "@/components/sections/HowItWorks.vue";
import AppliencesToFix from "@/components/sections/AppliencesToFix.vue";
import TheAboutFixly from "@/components/sections/TheAboutFixly.vue";
import TheFooter from "@/components/layout/TheFooter.vue";
import RequestQuoteModal from "@/components/RequestQuoteModal.vue";
import SuccessModal from "@/components/modals/SuccessModal.vue";
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
import ErrorModal from "@/components/modals/ErrorModal.vue";



const isFormOpen = ref(false)

const selectedAppliance = ref("")

const isSuccessOpen= ref(false)

const isErrorOpen = ref(false)
const errorMessage = ref("")

const requestModalRef = ref(null)

function openError(message) {
  console.log("OPEN ERROR MODAL:", message)
  errorMessage.value = message || "Request failed"
  isErrorOpen.value = true
}

function openForm (appliance = '') {
  closeAllModals()
selectedAppliance.value = appliance
  isFormOpen.value = true
}

function closeAllModals() {
  isFormOpen.value = false
  isSuccessOpen.value = false
}


function closeForm() {
  isFormOpen.value = false
}

async function handleSubmit(payload) {
  console.log("FORM SUBMIT:", payload)

  try{
    // Отправляем на наш Go бэкенд
    const response = await fetch(`${API_BASE_URL}/api/orders`, {
      method:"POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        client_name:payload.name,
        phone:payload.phone,
        device:payload.appliance,
        problem:payload.description,
        zip_code: payload.zipCode,
        status:"new",
        price: 0,
      }),

    })

    if (!response.ok) {
      const text = await response.text().catch(() => "")
      throw new Error(`Server error ${response.status}: ${text}`)
    }

    const created = await response.json()

    console.log("Создан заказ:", created)

    requestModalRef.value?.resetForm?.()

    // Показываем успех
    closeAllModals()
    isSuccessOpen.value = true

    setTimeout(() => {
      isSuccessOpen.value = false
    }, 3000)
  } catch (e) {
    console.error("Ошибка отправки:", e)
    openError(e?.message || "Network error")
  }
}

watch(
    () => isFormOpen.value || isSuccessOpen.value || isErrorOpen.value,
    (isAnyModalOpen) => {
      document.body.style.overflow = isAnyModalOpen ? "hidden" : "";
    }
);



</script>

<template>
  <TheHeader @open-form = "openForm" />
  <TheHero/>
  <HowItWorks/>
  <AppliencesToFix @open-form = "openForm" />
  <TheAboutFixly/>
  <TheFooter/>

  <RequestQuoteModal
      ref="requestModalRef"
      :open="isFormOpen"
      :default-appliance = "selectedAppliance"
      @close="isFormOpen = false"
      @submit="handleSubmit"/>

  <SuccessModal
  :open = "isSuccessOpen"
  @close="isSuccessOpen = false"
  />
  <ErrorModal
      :open="isErrorOpen"
      :message="errorMessage"
      @close="isErrorOpen = false"
  />


</template>
