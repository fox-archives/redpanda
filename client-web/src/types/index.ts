export type Repo = {
	name: string
	url: string
	dir: string
	status: string
}

export type Transformer = {
	type: 'command'
	name: string
	content: string
}

export type Transaction = {
	name: string
	repos: Repo[]
	transformers: Transformer[]
}
