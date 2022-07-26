<template>
	<div class="group">
		<span v-for="repository in repositories" :key="repository">
			{{ repository }}
		</span>
	</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
	props: {
		repositories: Array,
	},
})
</script>

<style scoped>
.group {
	display: flex;
	gap: 10px;
}
</style>
