<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from '@/components/HelloWorld.vue'
</script>

<template>
	<header>
		<h1>Transaction List</h1>
	</header>
	<RouterView />
</template>

<style>
@import '@/assets/base.css';

#app {
	max-width: 1280px;
	margin: 0 auto;
	padding: 2rem;

	font-weight: normal;
}
</style>
