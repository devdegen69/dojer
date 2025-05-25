<script setup lang="ts">
import type { Reader } from '~/shared/types/DojsList';

const route = useRoute()

const { data } = await useFetch<Reader>(`http://localhost:8033/api/${route.params.id}`)

const baseURL = 'http://localhost:8033'
const currentPage = parseInt(route.params.page as string)


// Gest current image
const cImage = () => {
	return baseURL + data.value?.images[currentPage - 1]
}

const cSize = () => {
	return data.value?.sizes[currentPage - 1].split("x")[0]
}
</script>

<template>
	<div>
		<div v-if="data == null">
			No data
		</div>
		<main v-else class="flex flex-col items-center justify-center">
			<UHorizontalNavigation :links="[
				[
					{ label: '', icon: 'material-symbols:chevron-left-rounded', to: '/' },
					{ label: '', icon: 'material-symbols:chevron-right', to: '/' },
					{ label: 'Random', to: '/random' },
					{ label: 'Shuffle', to: '/shuffle' }
				]
			]" />
			<img :src="cImage()" width="600" alt="" />
		</main>
	</div>
</template>
