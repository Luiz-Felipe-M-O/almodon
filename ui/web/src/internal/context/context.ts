import { AsyncTry } from "../../module/errors/try.ts"
import jsxmm from "../../module/jsxmm/element.ts"
import { StatusPage } from "../component/status.ts"
import Source from "../support/source.ts"

export interface Swapper {
	Namespace(): string | undefined
	SwapNamespace(namespace: string): void
}

interface Context extends Component {
	Final(): boolean
}

export class Orquestrator {
	#room: HTMLElement

	#contexts: Record<string, Context>
	#swapper: Swapper

	constructor(placeholder: HTMLElement, swapper: Swapper = new hash()) {
		this.#room = placeholder

		this.#contexts = {}
		this.#swapper = swapper
	}

	Current(): string | undefined {
		return this.#swapper.Namespace()
	}

	Link(namespace: string, context: Context): void {
		this.#contexts[namespace] = context
	}

	Unlink(namespace: string): void {
		delete this.#contexts[namespace]
	}

	SwapTo(namespace: string): boolean {
		if (!Object.hasOwn(this.#contexts, namespace)) {
			return false
		}

		const context = this.#contexts[namespace]
		if (context.Final() && this.Current() === namespace) {
			return true
		}

		const content = context.HTML()
		this.#room.replaceChildren(content)
		this.#swapper.SwapNamespace(namespace)

		return true
	}
}

export class context {
	#url: string
	#content: HTMLElement
	#retry: boolean

	onpreload?: () => void
	onload?: () => void

	constructor(url: string) {
		this.#url = url
		this.#retry = true

		this.#content = jsxmm.Element("div", { id: "almodon" })
	}

	Final(): boolean {
		return !this.#retry
	}

	HTML(): HTMLElement {
		if (this.#retry) {
			try_callback(this.onpreload)

			load(this.#url).then(([result, ok]) => {
				this.#content.replaceWith(result)
				this.#content = result
				this.#retry = !ok

				try_callback(this.onload)
			})

			return this.#content
		}

		return this.#content
	}
}

const Parser = new DOMParser()

async function load(url: string): Promise<[HTMLElement, boolean]> {
	const [result, error] = await AsyncTry(fetch, url)
	if (error !== null) {
		throw error
	}
	if (!result.ok) {
		return [StatusPage(result.status), false]
	}

	const page = await result.text()
	const new_document = Parser.parseFromString(page, "text/html")

	const content = new_document.getElementById("almodon")
	if (content === null) {
		return [StatusPage(204), false]
	}

	const element = new_document.getElementById("meta-almodon")
	if (element !== null) {
		for (const property of element.children as any as HTMLElement[]) {
			switch (property.tagName) {
			case "ALMODON-SCRIPT":
				const src = property.dataset["src"]
				if (src !== undefined) {
					await import(new URL(src, url).href)
				}
				break

			case "ALMODON-STYLE":
				const href = property.dataset["href"]
				if (href !== undefined) {
					const style = jsxmm.Element("link", {
						rel: "stylesheet",
						href: Source.From(href, url),
					})

					document.head.append(style)
					await new Promise(resolve => {
						style.onload = resolve
					})
				}
				break

			default:
				throw new Error("Unrecognized property " + property.tagName.toLocaleLowerCase())
			}
		}
	}

	return [content, true]
}

export class hash implements Swapper {
	Namespace(): string | undefined {
		return location.hash.slice(1)
	}

	SwapNamespace(namespace: string): void {
		location.hash = namespace
	}
}

function try_callback(callback?: () => void): void {
	if (callback !== undefined) {
		callback()
	}
}