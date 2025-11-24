import { APIFetch } from "../../../support/gateway.ts"

export class ProductGateway implements product.Gateway {
    #base: string

    constructor(resource: string) {
        this.#base = resource + "/"
    }

    async List(offset: number, limit: number): Promise<product.ListResponse> {
        return await APIFetch(this.#base, "GET", { query: { offset: `${offset}`, limit: `${limit}` } })
    }

    async Get(uuid: UUID): Promise<product.Response> {
        return await APIFetch(this.#base + uuid, "GET")
    }

    async Create(req: product.Entity): Promise<UUID> {
        return await APIFetch(this.#base, "POST", { body: JSON.stringify(req) })
    }

    async Patch(uuid: UUID, req: product.PartialEntity): Promise<void> {
        return await APIFetch(this.#base + uuid, "PATCH", { body: JSON.stringify(req) })
    }

    async Delete(uuid: UUID): Promise<void> {
        return await APIFetch(this.#base + uuid, "DELETE")
    }
}