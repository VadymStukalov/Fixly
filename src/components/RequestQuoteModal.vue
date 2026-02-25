<template>
  <BaseModal :open="open" @close="emit('close')">
    <!-- Header -->
    <div class="flex items-center justify-between px-5 pt-5">
      <div class="flex items-center">
        <div class="max-w-[40px]"><img src="../assets/logo-modal.png" alt="Fixly logo"></div>
        <span class="text-sm font-semibold">FIXLY Service</span>
      </div>

      <button
          type="button"
          class="px-[23px] py-[6px] text-xs shadow-sm font-brand font-bold rounded-full text-white bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] hover:opacity-90 transition-opacity"
          @click="emit('close')"
      >
        Back to the page
      </button>
    </div>

    <!-- Body -->

    <div class="px-5 pb-5">
      <h3 class="mt-6 text-center text-sm font-semibold text-gray-800">
        Please fill in the fields
      </h3>

      <form @submit.prevent="onSubmit" class="mt-5 space-y-4">
        <div>
          <label class="block text-xs font-semibold text-gray-700" >Selected appliance</label>
          <select
              v-model="form.appliance"
              ref="applianceRef"
              :class="['mt-1 w-full h-12 rounded-md border border-gray-300 bg-white px-3 text-base font-medium outline-none', errors.appliance ? 'border-red-500' : 'border-gray-300']">
            <option value="" disabled>Select appliance</option>
            <option>Refrigerator</option>
            <option>Oven</option>
            <option>Cooktop</option>
            <option>Air Conditioner</option>
            <option>Washing Machine</option>
            <option>Dryer</option>
            <option>Microwave Oven</option>
            <option>Boiler</option>
            <option>Dishwasher</option>
          </select>
        </div>
        <div>
          <label class="block text-xs font-semibold text-gray-700">Your Name</label>
          <input
              ref="nameInputRef"
              v-model.trim="form.name"
              type="text"
              placeholder="Your Name"
              :class="['mt-1 w-full h-10 rounded-md border border-gray-300 bg-white px-3 text-sm outline-none placeholder:text-gray-400',errors.name ? 'border-red-500' : 'border-gray-300']"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-gray-700">Phone Number</label>
          <input
              ref="phoneInputRef"
              v-model.trim="form.phone"
              type="tel"
              inputmode="tel"
              placeholder="Your Phone Number"
              :class="['mt-1 w-full h-10 rounded-md border border-gray-300 bg-white px-3 text-sm outline-none placeholder:text-gray-400',errors.phone ? 'border-red-500' : 'border-gray-300']"
          />
        </div>
        <div>
          <label class="block text-xs font-semibold text-gray-700">ZIP Code</label>
          <input
              v-model.trim="form.zipCode"
              type="text"
              inputmode="numeric"
              placeholder="ZIP Code of Your Area"
              class="mt-1 w-full h-10 rounded-md border border-gray-300 bg-white px-3 text-sm outline-none placeholder:text-gray-400"
          />
        </div>
        <div>
          <label class="block text-xs font-semibold text-gray-700">Brief Description</label>
          <textarea
              v-model.trim="form.description"
              rows="4"
              placeholder="Please tell us briefly about the issue."
              class="mt-1 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm outline-none placeholder:text-gray-400 resize-none"
          />
        </div>
        <div class="pt-2 pb-6">
          <button
              :disabled="!isFormValid || isSubmitting"
              type="submit"
              class="mx-auto block w-[220px] rounded-full bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] px-6 py-3 text-sm font-semibold text-white"
          >
            Submit Request
          </button>
        </div>

      </form>

      <div class="mt-8 flex items-end justify-between gap-4">
        <p class="max-w-[160px] text-xs text-gray-600">
          An AI Agent will contact you shortly to confirm the details
        </p>

        <div class="h-20 w-20 rounded bg-gray-200">
          <img src="../assets/ai-agent.png" alt="Fixly AI Agent">
        </div>
      </div>
    </div>
  </BaseModal>
</template>






<script setup>
import {computed, reactive, watch, ref, nextTick} from "vue";
import BaseModal from "@/components/layout/BaseModal.vue";
import { useQuoteForm } from "@/composables/useQuoteForm";

const { form, isFormValid, resetForm } = useQuoteForm();


const props =defineProps({
  open: {
    type: Boolean,
    default: false
  },
  defaultAppliance:{
    type: String,
    default:''

}
})


const emit = defineEmits(['close', 'submit'])

// const form = reactive({
//   appliance: '',
//   name: '',
//   phone: '',
//   zipCode: '',
//   description: ''
// })

const nameInputRef = ref(null);
const applianceRef = ref(null);
const phoneInputRef = ref(null);
const isSubmitting = ref(false);





function scrollToFirstError() {
  if (!form.appliance) {
    applianceRef.value?.scrollIntoView({ behavior: "smooth", block: "center" });
    applianceRef.value?.focus();
    return;
  }

  if (!form.name) {
    nameInputRef.value?.scrollIntoView({ behavior: "smooth", block: "center" });
    nameInputRef.value?.focus();
    return;
  }

  if (!form.phone) {
    phoneInputRef.value?.scrollIntoView({ behavior: "smooth", block: "center" });
    phoneInputRef.value?.focus();
    return;
  }
}
const errors = computed(() => ({
  appliance: !form.appliance,
  name: !form.name,
  phone: !form.phone,
}));



function onSubmit() {
  if (!isFormValid.value || isSubmitting.value) {
    scrollToFirstError();
    return;
  }

  isSubmitting.value = true;

  emit('submit', {...form})
  // resetForm()
  // emit('close')

  setTimeout(() => {
    isSubmitting.value = false;
  }, 1000);


}

defineExpose({ resetForm })
// function resetForm() {
//   form.appliance = ''
//   form.name = ''
//   form.phone = ''
//   form.zipCode = ''
//   form.description = ''
// }

// const isFormValid = computed( () => {
//   return (
//       form.appliance &&
//       form.name &&
//       form.phone &&
//       form.zipCode
//   );
// })

watch(
    () => props.open,
    async (isOpen) => {
      if (!isOpen) return

      isSubmitting.value = false;
        form.appliance = props.defaultAppliance || "";


      await nextTick();
      nameInputRef.value?.focus();
    }
);
</script>