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
import BaseModal from "@/components/layout/BaseModal.vue";
import SuccessModal from "@/components/modals/SuccessModal.vue";


const isFormOpen = ref(false)

const selectedAppliance = ref("")

const isSuccessOpen= ref(false)

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
  closeAllModals()
  isSuccessOpen.value = true

  try {
    await fetch("https://example.com/api/quote", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    })
  } catch (e) {
    console.error("Send error", e);
  }


  setTimeout(() => {
    isSuccessOpen.value = false
  }, 3000)
}

watch(
    () => isFormOpen.value || isSuccessOpen.value,
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
      :open="isFormOpen"
      :default-appliance = "selectedAppliance"
      @close="isFormOpen = false"
      @submit="handleSubmit"/>

  <SuccessModal
  :open = "isSuccessOpen"
  @close="isSuccessOpen = false"
  />

</template>
