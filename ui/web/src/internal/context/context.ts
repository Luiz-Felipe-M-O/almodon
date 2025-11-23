export interface Context {
	Build(): Promise<HTMLElement>
}

export class Orquestrator {
	#room: HTMLElement

	#contexts: Record<string, Context>
	#current: string | undefined

	constructor(switchable: HTMLElement) {
		this.#room = switchable

		this.#contexts = {}
		this.#current = undefined
	}

	Current(): string | undefined {
		return this.#current
	}

	Link(namespace: string, context: Context): void {
		this.#contexts[namespace] = context
	}

	Unlink(namespace: string): void {
		delete this.#contexts[namespace]
	}

	SwapTo(namespace: string): boolean {
		if (this.#current === namespace) {
			return true
		}

		const context = this.#contexts[namespace]
		if (context === undefined) {
			return false
		}

		context.Build().then((content) => {
			this.#room.replaceChildren(content)
		}).catch(() => {
			this.#room.replaceChildren("something went wrong")
			this.Unlink(namespace)
		})

		return true
	}
}

export class context {
	private static parser = new DOMParser()

	#url: string
	#content: HTMLElement | undefined

	constructor(url: string) {
		this.#content = undefined
		this.#url = url
	}

	async Build(): Promise<HTMLElement> {
		if (this.#content !== undefined) {
			return this.#content
		}

		const result = await fetch(this.#url)
		if (!result.ok) {
			throw new Error(result.statusText)
		}

		const page = await result.text()
		const document = context.parser.parseFromString(page, "text/html")

		const content = document.getElementById("almodon")
		if (content === null) {
			throw new Error("No content")
		}

		this.#content = content
		return content
	}
}