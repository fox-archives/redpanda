<template>
	<div>
		<h1>Transaction List</h1>
		<div v-if="loadState === 'loaded'">
			<div v-for="transaction in transactions" :key="transaction.name">
				<div>
					<router-link :to="`/transaction/${transaction.name}`">
						<h2>{{ transaction.name }}</h2>
					</router-link>
					<RepoList
						:repositories="transaction.repos.map((item) => item.name)"
					></RepoList>
				</div>
			</div>
		</div>
		<div v-else-if="loadState === 'loading'">
			<div>Loading</div>
		</div>
		<div v-else-if="loadState === 'error'">
			<h1>Error</h1>
		</div>
	</div>
</template>

<script lang="ts">
import type { Transaction } from '../types'
import { defineComponent, onDeactivated, onMounted, ref } from 'vue'
import RepoList from '../components/RepoList.vue'

export default defineComponent({
	setup() {
		const loadState = ref<'loading' | 'error' | 'loaded'>('loading')

		const abortController = new AbortController()
		const abortTimer = setTimeout(() => {
			abortController.abort()
		}, 3000)

		const transactions = ref<Transaction[]>([])
		onMounted(async () => {
			try {
				const res = await fetch('/api/transaction/list', {
					method: 'POST',
					body: '{}',
					signal: abortController.signal,
				})
				const data = await res.json()
				transactions.value = data.transactions
				loadState.value = 'loaded'
			} catch (error: unknown) {
				console.log(error)
				loadState.value = 'error'
			}
		})
		onDeactivated(() => {
			clearTimeout(abortTimer)
		})

		return { transactions, loadState }
	},
	components: { RepoList },
})
</script>
