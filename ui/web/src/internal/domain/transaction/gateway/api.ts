import { APIFetch } from "../../../support/gateway.ts"

export class TransactionGateway implements transaction.Gateway {
    #base: string

    constructor(resource: string) {
        this.#base = resource + "/"
    }

    async List(offset: number, limit: number): Promise<transaction.ListResponse> {
        return await APIFetch(this.#base, "GET", { 
            query: { offset: `${offset}`, limit: `${limit}` } 
        })
    }

    async Get(uuid: UUID): Promise<transaction.Response> {
        return await APIFetch(this.#base + uuid, "GET")
    }

    async Create(req: transaction.Entity): Promise<UUID> {
        return await APIFetch(this.#base, "POST", { 
            body: JSON.stringify(req) 
        })
    }

    async Patch(uuid: UUID, req: transaction.PartialEntity): Promise<void> {
        return await APIFetch(this.#base + uuid, "PATCH", { 
            body: JSON.stringify(req) 
        })
    }

    async Delete(uuid: UUID): Promise<void> {
        return await APIFetch(this.#base + uuid, "DELETE")
    }
}