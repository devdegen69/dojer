<script setup lang="ts">
import type { DojsList } from "~/shared/types/DojsList";

const config = useRuntimeConfig();
const apiUrl = config.public.apiUrl;
const currentPage = ref(1);
const router = useRouter();
const route = useRoute();

watch(currentPage, () => {
  router.push({ query: { page: currentPage.value } });
});

currentPage.value = parseInt((route.query.page as string) || "1");
const getThumb = (id: string) => {
  return `${apiUrl}/downs/thumbs/${id}.jpg`;
};

const { data, status } = useApi<DojsList>("/", {
  query: {
    page: currentPage,
  },
});
</script>

<template>
  <div>
    <div v-if="data == null">No data</div>
    <main v-else class="flex flex-col">
      <UPagination
        v-model="currentPage"
        :page-count="data.pagination.pageSize"
        :total="data.pagination.totalResults"
        class="m-auto"
        v-if="data.pagination.totalResults > data.pagination.pageSize"
      />
      <UContainer class="my-4">
        <div class="grid grid-cols-1 md:grid-cols-4 lg:grid-cols-5 gap-3">
          <NuxtLink
            v-if="status === 'success'"
            v-for="doj in data.dojs"
            key="data"
            :to="`/g/${doj.id}`"
            class="bg-secondary border border-gray-200 rounded-md flex justify-center items-center"
          >
            <img loading="lazy" :src="getThumb(doj.id)" class="rounded-md" />
          </NuxtLink>
          <div
            v-else
            v-for="n in 40"
            :key="n"
            class="flex justify-center items-center"
          >
            <USkeleton class="h-[340px] w-[320px]" />
          </div>
        </div>
      </UContainer>
      <UPagination
        v-model="currentPage"
        :page-count="data.pagination.pageSize"
        :total="data.pagination.totalResults"
        class="m-auto"
        v-if="data.pagination.totalResults > data.pagination.pageSize"
      />
    </main>
  </div>
</template>
