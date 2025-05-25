<script setup lang="ts">
defineProps<{
  error: {
    statusCode: number;
    statusMessage: string;
    message: string;
  };
}>();
</script>

<template>
  <div class="p-8 text-center">
    <h1 class="text-4xl font-bold mb-2">Erro {{ error.statusCode }}</h1>
    <p class="text-xl text-gray-600">{{ error.message }}</p>
    <NuxtLink to="/" class="mt-4 inline-block text-blue-500 underline">Voltar para o início</NuxtLink>
  </div>
</template>
