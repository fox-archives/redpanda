<script lang="ts">
import TheWelcome from '@/components/TheWelcome.vue'
import { defineComponent, onMounted, reactive, ref } from 'vue'

export default defineComponent({
	setup() {
		let transactions = ref([])

		onMounted(() => {
			fetch('/api/transaction/list', {
				method: 'POST',
				body: '{}',
			})
				.then((res) => {
					return res.json()
				})
				.then((data) => {
					console.log(data.transactions)
					transactions.value = data.transactions
				})
		})

		return { transactions, TheWelcome }
	},
})
</script>

<template>
	<div v-for="transaction in transactions" :key="transaction.name">
		<div>
			<router-link :to="'/transaction/' + transaction.name"
				><h2>{{ transaction.name }}</h2></router-link
			>

			<span v-for="repo in transaction.repos" style="margin-right: 10px">{{
				repo.Name
			}}</span>
			<div v-for="transformer in transaction.transformers"></div>
		</div>
	</div>
	<main>
		<TheWelcome />
	</main>
</template>
