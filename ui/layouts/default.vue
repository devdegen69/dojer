<script setup lang="ts">
useHead({
	title: "My self-hosted doujinshi reader",
});
const links = [
	[
		{ label: "Home", to: "/", icon: "i-heroicons-home" },
		{ label: "Random", to: "/random", icon: "mdi:dice-3" },
		{ label: "Shuffle", to: "/shuffle", icon: "material-symbols:shuffle" },
	],
];
</script>
<template>
	<div>
		<UHorizontalNavigation :links="links" />
		<slot />
	</div>
</template>
