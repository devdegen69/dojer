<script setup lang="ts">
const { content, title, counters } = defineProps({
  content: String,
  title: String,
  counters: Object,
});

const ct = (k: string) => k.split("__")[0];
const ctv = (k: string) => k.split("__")[1];
const st = (s: string) => {
  if (s.split(" ").length > 1) {
    return `\"${s}\"`;
  } else {
    return s;
  }
};
</script>

<template>
  <div class="max-w-full flex my-1 flex-wrap" v-if="content">
    <span>{{ title }}: </span>
    <NuxtLink
      :to="`/${title?.toLowerCase()}/${st(ctv(k))}`"
      v-for="(k, _v) in counters"
      class="p-0"
    >
      <UBadge :label="ct(k) + ' | ' + ctv(k)" class="mx-1 text-nowrap">
      </UBadge>
    </NuxtLink>
  </div>
</template>
