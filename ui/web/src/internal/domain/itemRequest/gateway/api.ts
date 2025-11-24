import { APIFetch } from "../../../support/gateway.ts"

export class ItemRequestGateway implements itemRequest.Gateway {
    #base: string

    constructor(resource: string) {
        this.#base = resource + "/"
    }

    async List(offset: number, limit: number): Promise<itemRequest.ListResponse> {
        return await APIFetch(this.#base, "GET", {
            query: { offset: `${offset}`, limit: `${limit}` }
        })
    }

    async Get(uuid: UUID): Promise<itemRequest.Response> {
        return await APIFetch(this.#base + uuid, "GET")
    }

    async Create(req: itemRequest.Entity): Promise<UUID> {
        return await APIFetch(this.#base, "POST", {
            body: JSON.stringify(req)
        })
    }

    async Patch(uuid: UUID, req: itemRequest.PartialEntity): Promise<void> {
        return await APIFetch(this.#base + uuid, "PATCH", {
            body: JSON.stringify(req)
        })
    }

    async Delete(uuid: UUID): Promise<void> {
        return await APIFetch(this.#base + uuid, "DELETE")
    }
}