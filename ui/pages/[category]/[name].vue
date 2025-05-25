<script setup lang="ts">
const route = useRoute();
const router = useRouter();
const currentPage = ref(1);
currentPage.value = parseInt((route.query.page as string) || "1");

const validCategories = [
  "tags",
  "artists",
  "characters",
  "groups",
  "parodies",
  "languages",
  "categories",
];

const category = route.params.category as string;
const name = route.params.name as string;

if (!validCategories.includes(category)) {
  throw createError({
    statusCode: 404,
    statusMessage: "invalid category",
  });
}

const { data } = useFetch<IndexResponse>(
  `http://localhost:8033/api/search/${route.params.category}: ${route.params?.name}`,
  {
    query: {
      page: currentPage,
    },
  },
);

watch(currentPage, () => {
  router.push({ query: { page: currentPage.value } });
});
const getThumb = (id: string) => {
  return `http://localhost:8033/downs/thumbs/${id}.jpg`;
};
</script>

<template>
  <div>
    <div v-if="data == null">No data</div>
    <main v-else class="flex flex-col">
      <h2 class="text-center text-xl m-4">
        Characters: {{ route.params.name }} with
        {{ data.pagination.totalResults }} results
      </h2>
      <UPagination v-model="currentPage" :page-count="data.pagination.pageSize" :total="data.pagination.totalResults"
        class="m-auto" v-if="data.pagination.totalResults > data.pagination.pageSize" />
      <UContainer class="py-3">
        <div class="grid grid-cols-1 md:grid-cols-4 lg:grid-cols-5 gap-3">
          <NuxtLink v-for="doj in data.dojs" key="data" :to="`/g/${doj.id}`"
            class="bg-secondary border border-gray-200 rounded-md flex justify-center items-center">
            <img :src="getThumb(doj.id)" class="rounded-md" />
          </NuxtLink>
        </div>
      </UContainer>
      <UPagination v-model="currentPage" :page-count="data.pagination.pageSize" :total="data.pagination.totalResults"
        class="m-auto" v-if="data.pagination.totalResults > data.pagination.pageSize" />
    </main>
  </div>
</template>
