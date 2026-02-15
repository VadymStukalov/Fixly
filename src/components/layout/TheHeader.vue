<template>
  <header class="flex max-w-5xl mx-auto justify-between px-5 py-5 items-center">
    <img src="../../assets/logo.png" alt="logo" class="w-14 md:w-20">

    <nav class="hidden md:flex gap-4 font-brand">
      <a href="#" class=" font-brand font-medium">About Us</a>
      <a href="#" class=" font-brand font-medium">Contact</a>

    </nav>

  <button
      class="px-[23px] py-[6px] text-sm shadow-sm font-brand font-bold rounded-full text-white bg-gradient-to-r from-[#3BA3A9] to-[#6A339E] hover:opacity-90 transition-opacity"
      @click = "emit('open-form')"
  >Get a Quote</button>
  </header>
</template>

<script setup>
const emit = defineEmits(['open-form'])
</script>