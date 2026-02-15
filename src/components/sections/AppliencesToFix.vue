<template>

  <section class="py-6">
    <div class="max-w-5xl mx-auto flex flex-col items-center">
      <h2 class="text-xl font-semibold mb-5">What Do You Need Fixed?</h2>

      <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 sm:gap-8 md:gap-10 px-4">

          <button
              v-for="app in appliences"
              :key="app.id"
              class="flex flex-col shadow-md gap-4 p-3 items-center bg-[#C9E5E4] rounded-2xl sm:p-6 md:p-8 hover:scale-105"
              @click= "emit('open-form', app.name)"
              type="button"
          >
            <div class="flex-1 flex items-center">
              <img :src="app.image" :alt="app.name"  class="w-full max-w-[60px] sm:max-w-[120px] md:max-w-[160px]">
            </div>
            <p class="text-xs font-medium sm:text-base">{{app.name}}</p>
          </button>

      </div>

    </div>

  </section>
</template>


<script setup>

const appliences = {
  washingMachine: {
    id:1,
    name: 'Washing Machine',
    image: 'img/washing-mashine-2.png'
  },
  refrigerator: {
    id:2,
    name: 'Refrigerator',
    image: 'img/refrigerator-2.png'
  },
  cooktop: {
    id:3,
    name: 'Cooktop',
    image: 'img/cooktop-2.png'
  },
  oven: {
    id:4,
    name: 'Oven',
    image: 'img/oven-2.png'
  },
  dishWasher: {
    id:5,
    name: 'Dishwasher',
    image: 'img/dish-washer-2.png'
  },
  dryer: {
    id:6,
    name: 'Dryer',
    image: 'img/dryer-2.png'
  },
  airConditioner: {
    id:7,
    name: 'Air Conditioner',
    image: 'img/air-conditioner-2.png'
  },
  microwaveOven: {
    id:8,
    name: 'Microwave Oven',
    image: 'img/microwave-oven-2.png'
  },
  boiler: {
    id:9,
    name: 'Boiler',
    image: 'img/boiler-2.png'
  }
}


const emit = defineEmits(['open-form'])

</script>