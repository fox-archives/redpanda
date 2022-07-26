<template>
	<div>
		<h1 class="title">
			Transaction <em>{{ transaction.name }}</em>
		</h1>
		<div v-if="transaction.repos">
			<RepoList
				:repositories="transaction.repos.map((item) => item.name)"
			></RepoList>
		</div>

		<h2>Transformers</h2>
		<div
			v-for="transformer in transaction.transformers"
			:key="transformer.name"
		>
			<div class="pure-form transformer-card">
				<h3>
					<b>{{ transformer.name }}</b> ({{ transformer.type }})
				</h3>
				<textarea
					cols="40"
					rows="6"
					@input="saveTransformerContent(transformer.name)"
					v-model="transformer.content"
				>
				</textarea>
			</div>
		</div>
		<div class="pure-form add-transformer">
			<input v-model="transformerAddText" />
			<button class="pure-button">Add Transformer</button>
		</div>

		<h2>Actions</h2>
		<div class="pure-form steps">
			<div class="groupp">
				<button class="pure-button" @click="actionApply">Apply</button>
				<button class="pure-button pull" @click="actionPullAndReapply">
					Pull and Reapply
				</button>
			</div>

			<textarea rows="20" cols="79" v-model="diffText"></textarea>
			<input v-model="commitMessage" />
			<button class="pure-button" @click="actionCommit">Commit</button>
			<br />
			<button class="pure-button" @click="actionPush">Push</button>
		</div>

		<h2>Status</h2>
		<p>{{ actionStatus }}</p>
	</div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { debounce } from 'lodash'
import type { Transaction } from '../types'
import RepoList from '../components/RepoList.vue'
import { fetchWrapper } from '../util/util'

export default defineComponent({
	setup() {
		const actionStatus = ref('')
		const transaction = ref<Transaction | Record<string, never>>({})
		const router = useRoute()
		const transformerAddText = ref('')
		const diffText = ref('')
		const commitMessage = ref('')

		onMounted(async () => {
			const [ok, data] = await fetchWrapper('/api/transaction/list')
			if (!ok) throw data

			transaction.value = data.transactions.filter((item) => {
				return item.name == router.params['transaction']
			})[0]

			await actionApply()
		})

		async function actionApply() {
			const [ok, data] = await fetchWrapper('/api/action/apply', {
				transaction: transaction.value.name,
			})
			if (!ok) throw data

			diffText.value = data.contents
		}

		async function actionPullAndReapply() {
			const [ok, data] = await fetchWrapper(
				'/api/action/refresh',
				{
					transaction: transaction.value.name,
				},
				{
					timeout: 30,
				},
			)
			if (!ok) throw data
		}

		async function actionCommit() {
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
				const text = await res.text()
				console.log(text)
			})
		}

		async function actionPush() {
			await fetch('/api/step/push', {
				method: 'POST',
				body: JSON.stringify({
					transaction: transaction.value.name,
				}),
			}).then(async (res) => {
				const text = await res.text()
				console.log(text)
			})
		}

		// util
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

		return {
			transaction,
			transformerAddText,
			saveTransformerContent,
			actionApply,
			actionPullAndReapply,
			actionCommit,
			actionPush,
			diffText,

			commitMessage,
			actionStatus,
		}
	},
	components: { RepoList },
})
</script>

<style scoped>
.title {
	margin-block: 0;
}

.transformer-card {
	background-color: beige;
	border-radius: 5px;
	margin-block-end: 5px;
	padding: 3px;
}

.pull {
	margin-block-start: auto;
}

.groupp {
	display: flex;
	justify-content: space-between;
}
</style>
