import { APIError } from "../../../../module/errors/error.ts"

type ItemRequest = {
    uuid: UUID
    amountRequest: number
}

export class ItemRequestGateway implements itemRequest.Gateway {
    #items: ItemRequest[]

    constructor() {
        this.#items = [
            {
                uuid: crypto.randomUUID(),
                amountRequest: 5
            }
        ]
    }

    async List(offset: number, limit: number): Promise<itemRequest.ListResponse> {
        const lo = offset
        const hi = offset + limit

        const records = this.#items.slice(lo, hi).map(item => ({
            uuid: item.uuid,
            amountRequest: item.amountRequest
        }))

        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#items.length
        }
    }

    async Get(uuid: UUID): Promise<itemRequest.Response> {
        const item = this.#items.find(i => i.uuid === uuid)

        if (!item) {
            throw APIError.New("not found", "itemRequest-not-found", `itemRequest with UUID ${uuid} not found`)
        }

        return {
            uuid: item.uuid,
            amountRequest: item.amountRequest
        }
    }

    async Create(req: itemRequest.Entity): Promise<UUID> {
        const uuid = crypto.randomUUID()

        this.#items.push({
            uuid: uuid,
            amountRequest: req.amountRequest
        })

        return uuid
    }

    async Patch(uuid: UUID, req: itemRequest.PartialEntity): Promise<void> {
        const item = this.#items.find(i => i.uuid === uuid)

        if (!item) {
            throw APIError.New("not found", "itemRequest-not-found", `itemRequest with UUID ${uuid} not found`)
        }

        if (req.amountRequest !== undefined) {
            item.amountRequest = req.amountRequest
        }
    }

    async Delete(uuid: UUID): Promise<void> {
        const index = this.#items.findIndex(i => i.uuid === uuid)

        if (index === -1) {
            throw APIError.New("not found", "itemRequest-not-found", `itemRequest with UUID ${uuid} not found`)
        }

        this.#items.splice(index, 1)
    }
}