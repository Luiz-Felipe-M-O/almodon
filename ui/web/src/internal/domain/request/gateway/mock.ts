import { APIError } from "../../../../module/errors/error.ts"

type RequestData = {
    uuid: UUID
    status: string
    dateRequest: Date
    reason: string
}

export class RequestGateway implements request.Gateway {
    #requests: RequestData[]

    constructor() {
        this.#requests = [
            {
                uuid: "00000000-0000-0000-0000-000000000001",
                status: "PENDING",
                dateRequest: new Date("2025-12-01T10:00:00"),
                reason: "Solicitação de acrílico."
            },
            {
                uuid: "00000000-0000-0000-0000-000000000002",
                status: "APPROVED",
                dateRequest: new Date("2025-12-05T14:30:00"),
                reason: "Solicitação de resina composta."
            }
        ]
    }

    async List(offset: number, limit: number): Promise<request.ListResponse> {
        const lo = offset
        const hi = offset + limit

        const records = this.#requests.slice(lo, hi).map(r => ({
            uuid: r.uuid,
            status: r.status,
            dateRequest: r.dateRequest,
            reason: r.reason,
        }))

        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#requests.length,
        }
    }

    async Get(uuid: UUID): Promise<request.Response> {
        for (const r of this.#requests) {
            if (r.uuid === uuid) {
                return {
                    uuid: r.uuid,
                    status: r.status,
                    dateRequest: r.dateRequest,
                    reason: r.reason,
                }
            }
        }

        throw APIError.New("not found", "request-not-found", `request with UUID ${uuid} not found`)
    }

    async Create(req: request.Entity): Promise<UUID> {
        const uuid = crypto.randomUUID()
        this.#requests.push({
            uuid: uuid,
            status: req.status,
            dateRequest: req.dateRequest,
            reason: req.reason,
        })

        return uuid
    }

    async Patch(uuid: UUID, req: request.PartialEntity): Promise<void> {
        for (const r of this.#requests) {
            if (r.uuid === uuid) {
                if (req.status !== undefined) { r.status = req.status }
                if (req.dateRequest !== undefined) { r.dateRequest = req.dateRequest }
                if (req.reason !== undefined) { r.reason = req.reason }
                return
            }
        }

        throw APIError.New("not found", "request-not-found", `request with UUID ${uuid} not found`)
    }

    async Delete(uuid: UUID): Promise<void> {
        for (let i = 0; i < this.#requests.length; i++) {
            if (this.#requests[i].uuid === uuid) {
                this.#requests.splice(i, 1)
                return
            }
        }
    }
}