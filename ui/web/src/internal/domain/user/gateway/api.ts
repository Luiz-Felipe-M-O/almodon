import { APIFetch } from "../../../support/gateway.ts"

export class UserGateway implements user.Gateway {
    #base: string
    #siape: string
    #auth: string
    #me: string

    constructor(resource: string) {
        this.#base = resource + "/"
        this.#siape = resource + "/siape/"
        this.#auth = resource + "/auth/"
        this.#me = resource + "/me/"
    }

    async List(offset: number, limit: number): Promise<user.ListResponse> {
        const response: user.ListResponse = await APIFetch(this.#base, "GET", { query: { offset: `${offset}`, limit: `${limit}` } })

        response.records.forEach(record => {
            record.created = new Date(record.created as any)
            record.updated = new Date(record.updated as any)
        })

        return response
    }

    async Get(uuid: UUID): Promise<user.Response> {
        const response: user.Response = await APIFetch(this.#base + uuid, "GET")

        response.created = new Date(response.created as any)
        response.updated = new Date(response.updated as any)

        return response
    }

    async GetBySIAPE(siape: string): Promise<user.Response> {
        const response: user.Response = await APIFetch(this.#siape + siape, "GET")

        response.created = new Date(response.created as any)
        response.updated = new Date(response.updated as any)

        return response
    }

    async Create(req: user.Entity): Promise<UUID> {
        return await APIFetch(this.#base, "POST", { body: JSON.stringify(req) })
    }

    async Patch(uuid: UUID, req: user.PartialEntity): Promise<void> {
        return await APIFetch(this.#base + uuid, "PATCH", { body: JSON.stringify(req) })
    }

    async Delete(uuid: UUID): Promise<void> {
        return await APIFetch(this.#base + uuid, "DELETE")
    }

    async Autheticate(siape: string, password: string): Promise<user.AuthResponse> {
        const response: user.AuthResponse = await APIFetch(this.#auth, "POST", { body: JSON.stringify({ siape: siape, password: password }) })
        response.expires = new Date(response.expires as any)

        return response
    }

    async Logout(): Promise<void> {
        return await APIFetch(this.#auth, "DELETE")
    }

    async Me(): Promise<user.Response> {
        const response: user.Response = await APIFetch(this.#me, "GET")

        response.created = new Date(response.created as any)
        response.updated = new Date(response.updated as any)

        return response
    }
}