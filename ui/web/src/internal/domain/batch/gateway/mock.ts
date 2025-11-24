import { APIError } from "../../../../module/errors/error.ts"

type Batch = {
    uuid: UUID
    idBatch: number
    amount: number
    codeBatch: string
    observationsNote: string
    specifications: string
    valueUnitary: number
    dateEntry: string
    dateExpiration: string
}

export class BatchGateway implements batch.Gateway {
    #batches: Batch[]

    constructor() {
        this.#batches = [
            {
                uuid: crypto.randomUUID(),
                idBatch: 1,
                amount: 120,
                codeBatch: "LE-2512",
                observationsNote: "Lote de exemplo",
                specifications: "Material A",
                valueUnitary: 10.50,
                dateEntry: "2025-12-20",
                dateExpiration: "2027-08-20",
            }
        ]
    }

    async List(offset: number, limit: number): Promise<batch.ListResponse> {
        const lo = offset
        const hi = offset + limit

        const records = this.#batches.slice(lo, hi).map(batch => ({
            uuid: batch.uuid,
            idBatch: batch.idBatch,
            amount: batch.amount,
            codeBatch: batch.codeBatch,
            observationsNote: batch.observationsNote,
            specifications: batch.specifications,
            valueUnitary: batch.valueUnitary,
            dateEntry: batch.dateEntry,
            dateExpiration: batch.dateExpiration,
        }))

        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#batches.length
        }
    }

    async Get(uuid: UUID): Promise<batch.Response> {
        const batch = this.#batches.find(b => b.uuid === uuid)

        if (!batch) {
            throw APIError.New("not found", "batch-not-found", `batch with UUID ${uuid} not found`)
        }

        return {
            uuid: batch.uuid,
            idBatch: batch.idBatch,
            amount: batch.amount,
            codeBatch: batch.codeBatch,
            observationsNote: batch.observationsNote,
            specifications: batch.specifications,
            valueUnitary: batch.valueUnitary,
            dateEntry: batch.dateEntry,
            dateExpiration: batch.dateExpiration,
        }
    }

    async Create(req: batch.Entity): Promise<UUID> {
        const uuid = crypto.randomUUID()

        this.#batches.push({
            uuid: uuid,
            idBatch: req.idBatch,
            amount: req.amount,
            codeBatch: req.codeBatch,
            observationsNote: req.observationsNote,
            specifications: req.specifications,
            valueUnitary: req.valueUnitary,
            dateEntry: req.dateEntry,
            dateExpiration: req.dateExpiration,
        })

        return uuid
    }

    async Patch(uuid: UUID, req: batch.PartialEntity): Promise<void> {
        const batch = this.#batches.find(b => b.uuid === uuid)

        if (!batch) {
            throw APIError.New("not found", "batch-not-found", `batch with UUID ${uuid} not found`)
        }

        if (req.amount !== undefined) batch.amount = req.amount
        if (req.codeBatch !== undefined) batch.codeBatch = req.codeBatch
        if (req.observationsNote !== undefined) batch.observationsNote = req.observationsNote
        if (req.specifications !== undefined) batch.specifications = req.specifications
        if (req.valueUnitary !== undefined) batch.valueUnitary = req.valueUnitary
        if (req.dateEntry !== undefined) batch.dateEntry = req.dateEntry
        if (req.dateExpiration !== undefined) batch.dateExpiration = req.dateExpiration
    }

    async Delete(uuid: UUID): Promise<void> {
        const index = this.#batches.findIndex(b => b.uuid === uuid)

        if (index === -1) {
            throw APIError.New("not found", "batch-not-found", `batch with UUID ${uuid} not found`)
        }

        this.#batches.splice(index, 1)
    }
}
