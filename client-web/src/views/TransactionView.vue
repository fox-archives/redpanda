<template>
	<h1>Information</h1>
	<p>NAME: {{ transaction.name }}</p>
	<p style="display: inline-block; margin-right: 5px">REPOSITORIES:</p>
	<span
		v-for="repo in transaction.repos"
		:key="repo.Name"
		style="margin-right: 10px"
		>{{ repo.Name }}</span
	>

	<h1>Transformers</h1>
	<input v-model="transformerAddText" /><button>Add Transformer</button>
	<br />
	<br />
	<hr />
	<div v-for="transformer in transaction.transformers" :key="transformer.name">
		<h3>type: {{ transformer.type }}</h3>
		<h3>
			name: <b>{{ transformer.name }}</b>
		</h3>
		<textarea
			cols="40"
			@input="saveTransformerContent(transformer.name)"
			v-model="transformer.content"
		>
		</textarea>
		<hr />
	</div>
	<h2>Steps</h2>
	<button style="margin-right: 5px" @click="stepDiff">Update</button>
	<button @click="stepCommit">Commit</button>
	<input v-model="commitMessage" />
	<br /><button @click="stepPush">Push</button>
	<br />
	<textarea rows="20" cols="79">{{ diffText }}</textarea>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { debounce } from 'lodash'

export default defineComponent({
	setup() {
		const transaction = ref({})
		const router = useRoute()
		const transformerAddText = ref('')
		const diffText = ref('')
		const commitMessage = ref('')

		async function stepDiff() {
			await fetch('/api/step/idempotent-apply', {
				method: 'POST',
				body: JSON.stringify({
					transaction: transaction.value.name,
				}),
			}).then(async (res) => {
				console.log(await res.text())
			})

			await fetch('/api/step/diff', {
				method: 'POST',
				body: JSON.stringify({
					transaction: transaction.value.name,
				}),
			}).then(async (res) => {
				diffText.value = (await res.json()).contents
				console.log(diffText.value)
			})
		}

		async function stepCommit() {
			if (commitMessage.value.length < 4) {
				throw new Error('Commit message too small')
			}

			await fetch('/api/step/commit', {
				method: 'POST',
				body: JSON.stringify({
					transaction: transaction.value.name,
					commitMessage: commitMessage.value,
				}),
			}).then(async (res) => {
				let text = await res.text()
				console.log(text)
			})
		}

		async function stepPush() {
			await fetch('/api/step/push', {
				method: 'POST',
				body: JSON.stringify({
					transaction: transaction.value.name,
				}),
			}).then(async (res) => {
				let text = await res.text()
				console.log(text)
			})
		}

		function saveTransformerContent(value: string) {
			const transformer = transaction.value.transformers.filter((item) => {
				return item.name === value
			})[0]

			const obj = {
				transaction: transaction.value.name,
				transformer: value,
				newContent: transformer.content,
			}

			fetch('/api/transformer/edit', {
				method: 'POST',
				body: JSON.stringify(obj),
			})
		}

		onMounted(async () => {
			await fetch('/api/transaction/list', {
				method: 'POST',
				body: '{}',
			})
				.then((res) => {
					return res.json()
				})
				.then((data) => {
					transaction.value = data.transactions.filter((item) => {
						return item.name == router.params['transaction']
					})[0]
				})

			await stepDiff()
		})
		return {
			transaction,
			transformerAddText,
			saveTransformerContent,
			stepDiff,
			diffText,
			stepCommit,
			commitMessage,
			stepPush,
		}
	},
})
</script>

<style></style>
