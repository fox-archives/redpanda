export async function fetchWrapper<T = Record<string, unknown>>(
	url: string,
	body: Record<string, unknown> = {},
	options?: { timeout?: number } = {},
): Promise<[true, T] | [false, Error]> {
	const abortController = new AbortController()
	setTimeout(() => {
		abortController.abort()
	}, (options.timeout || 3) * 1000)

	const req = await fetch(url, {
		method: 'POST',
		body: JSON.stringify(body),
		signal: abortController.signal,
	})
	if (req.ok) {
		const text = await req.text()
		if (text.length > 0) {
			try {
				return [true, JSON.parse(text)]
			} catch (err: unknown) {
				if (!(err instanceof Error)) {
					return [false, new Error('Failed to parse JSON')]
				} else {
					return [false, err]
				}
			}
		} else {
			return [true, JSON.parse('{}')]
		}
	} else {
		try {
			const text = await req.text()
			return [false, new Error(text)]
		} catch (err: unknown) {
			return [false, new Error(`${req.status}: ${req.statusText}`)]
		}
	}
}
