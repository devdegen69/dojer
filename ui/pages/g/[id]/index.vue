<script setup lang="ts">
import type { Reader } from '~/shared/types/DojsList';

const route = useRoute()

const { data } = await useFetch<Reader>(`http://localhost:8033/api/${route.params.id}`)

</script>


<template>
	<main>
		<div class="flex bg-neutral-800 px-4 py-4 justify-evenly" v-if="data !== null">
			<img :src="'http://localhost:8033' + data?.images[0]" :width="data?.sizes[0].split('x')[0]" class="rounded-xl" />
			<div class="flex flex-col justify-around items-center w-1/3">
				<h1 class="text-2xl font-bold">{{ data?.title }}</h1>
				<div class="max-w-full">
					<BadgeList :content="data.doj.parodies" title="Parodies" :counters="data.counters['parodies']"/>
					<BadgeList :content="data.doj.characters" title="Characters" :counters="data.counters['characters']"/>
					<BadgeList :content="data.doj.tags" title="Tags" :counters="data.counters['tags']"/>
					<BadgeList :content="data.doj.artists" title="Artists" :counters="data.counters['artists']"/>
					<BadgeList :content="data.doj.groups" title="Groups" :counters="data.counters['groups']"/>
					<BadgeList :content="data.doj.languages" title="Languages" :counters="data.counters['languages']"/>
					<BadgeList :content="data.doj.categories" title="Categories" :counters="data.counters['categories']"/>
					Pages: <UBadge class="mx-1">{{ data?.doj.pages }}</UBadge>
				</div>
			</div>
		</div>

		<div class="grid grid-cols-5 gap-4 p-4 w-2/3 m-auto">
			<NuxtLink v-for="(d, index) in data?.images" :key="d"
				:to="{ name: 'g-id-page', params: { id: data?.doj.id, page: index + 1 } }"
				class="bg-neutral-600 p-2 rounded hover:bg-neutral-700 transition-colors ">
				<img :src="'http://localhost:8033' + d" class="rounded-xl" />
			</NuxtLink>
		</div>
	</main>
</template>
