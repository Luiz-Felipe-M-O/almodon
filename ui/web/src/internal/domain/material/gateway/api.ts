import { APIFetch } from "../../../support/gateway.ts"

export class MaterialGateway implements material.Gateway {
    #base: string
    #siads: string
    #catmat: string
    #ecampus: string

    constructor(resource: string) {
        this.#base = resource + "/"
        this.#siads = resource + "/siads/"
        this.#catmat = resource + "/catmat/"
        this.#ecampus = resource + "/ecampus/"
    }

    async List(offset: number, limit: number): Promise<material.ListResponse> {
        const response: material.ListResponse = await APIFetch(this.#base, "GET", { query: { offset: `${offset}`, limit: `${limit}` } })

        response.records.forEach(record => {
            record.created = new Date(record.created as any)
            record.updated = new Date(record.updated as any)
        })

        return response
    }

    async ListBySIADS(siads: string, offset: number, limit: number): Promise<material.ListResponse> {
        const response: material.ListResponse = await APIFetch(this.#siads + siads, "GET", { query: { offset: `${offset}`, limit: `${limit}` } })

        response.records.forEach(record => {
            record.created = new Date(record.created as any)
            record.updated = new Date(record.updated as any)
        })

        return response
    }

    async ListByCATMAT(catmat: string, offset: number, limit: number): Promise<material.ListResponse> {
        const response: material.ListResponse = await APIFetch(this.#catmat + catmat, "GET", { query: { offset: `${offset}`, limit: `${limit}` } })

        response.records.forEach(record => {
            record.created = new Date(record.created as any)
            record.updated = new Date(record.updated as any)
        })

        return response
    }

    async ListByECampus(ecampus: string, offset: number, limit: number): Promise<material.ListResponse> {
        const response: material.ListResponse = await APIFetch(this.#ecampus + ecampus, "GET", { query: { offset: `${offset}`, limit: `${limit}` } })

        response.records.forEach(record => {
            record.created = new Date(record.created as any)
            record.updated = new Date(record.updated as any)
        })

        return response
    }

    async Get(uuid: UUID): Promise<material.Response> {
        const response: material.Response = await APIFetch(this.#base + uuid, "GET")

        response.created = new Date(response.created as any)
        response.updated = new Date(response.updated as any)

        return response
    }

    async Create(req: material.Entity): Promise<UUID> {
        const result: { uuid: UUID } = await APIFetch(this.#base, "POST", { body: JSON.stringify(req) })
        return result.uuid
    }

    async Patch(uuid: UUID, req: material.PartialEntity): Promise<void> {
        return await APIFetch(this.#base + uuid, "PATCH", { body: JSON.stringify(req) })
    }

    async Delete(uuid: UUID): Promise<void> {
        return await APIFetch(this.#base + uuid, "DELETE")
    }
}
