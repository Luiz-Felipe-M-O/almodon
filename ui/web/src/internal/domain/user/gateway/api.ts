import { APIFetch } from "../../../support/gateway.ts"

export class UserGateway implements user.Gateway {
    #base: string
    #auth: string
    #me: string

    constructor(resource: string) {
        this.#base = resource + "/"
        this.#auth = resource + "/auth/"
        this.#me = resource + "/me/"
    }

    async List(offset: number, limit: number): Promise<user.ListResponse> {
        return await APIFetch(this.#base, "GET", { query: { offset: `${offset}`, limit: `${limit}` } })
    }

    async Get(uuid: UUID): Promise<user.Response> {
        return await APIFetch(this.#base + uuid, "GET")
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

    async Autheticate(siape: number, password: string): Promise<user.AuthResponse> {
        return await APIFetch(this.#auth, "POST", { body: JSON.stringify({ siape: siape, password: password }) })
    }

    async Me(): Promise<user.Response> {
        return await APIFetch(this.#me, "GET")
    }
}