import { APIError } from "../../module/errors/error.ts"

type APIMethods = "GET" | "HEAD" | "POST" | "PUT" | "PATCH" | "DELETE"

export async function APIFetch<M extends APIMethods>(url: string, method: M, options?: { query?: Record<string, string>, body?: BodyInit }): Promise<any> {
	const config: { query: Record<string, string> | undefined, body: BodyInit | null } = {
		query: undefined,
		body: null,
	}

	if (options !== undefined) {
		if (options.query !== undefined) { config.query = options.query }
		if (options.body !== undefined) { config.body = options.body }
	}

	const query = new URLSearchParams(config.query).toString()
	const resp = await fetch(url + query, {
		method: method,
		headers: Headers[method],
		body: config.body,
	})

	if (!resp.ok) {
		throw APIError.FromObject(await resp.json())
	}

	if (resp.status === 204) {
		return
	}

	return await resp.json()
}

const GetHeaders: Record<string, string> = {
	"Accept": "application/json",
} as const

const PostHeaders: Record<string, string> = {
	"Content-Type": "application/json",
	"Accept": "application/json",
} as const

const PutHeaders: Record<string, string> = {
	"Content-Type": "application/json",
} as const

const DeleteHeaders: Record<string, string> = {} as const

const Headers: Record<APIMethods, HeadersInit> = {
	"GET": GetHeaders,
	"HEAD": GetHeaders,
	"POST": PostHeaders,
	"PUT": PutHeaders,
	"PATCH": PutHeaders,
	"DELETE": DeleteHeaders,
} as const
