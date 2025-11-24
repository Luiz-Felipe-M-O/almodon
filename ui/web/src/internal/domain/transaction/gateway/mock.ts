import { APIError } from "../../../../module/errors/error.ts"

type Transaction = {
    uuid: UUID
    amount: number
    typeTransaction: string
    dateHour: Date
}

export class TransactionGateway implements transaction.Gateway {
    #transactions: Transaction[]

    constructor() {
        this.#transactions = [
            {
                uuid: "00000000-0000-0000-0000-000000000001",
                amount: 120.50,
                typeTransaction: "loan",
                dateHour: new Date("2025-12-24T10:30:35") 
            },
            {
                uuid: "00000000-0000-0000-0000-000000000002",
                amount: -40.00,
                typeTransaction: "usable",
                dateHour: new Date("2025-12-27T18:40:24")
            }
        ]
    }

    async List(offset: number, limit: number): Promise<transaction.ListResponse> {
        const lo = offset
        const hi = offset + limit

        const records = this.#transactions.slice(lo, hi).map(t => ({
            uuid: t.uuid,
            amount: t.amount,
            typeTransaction: t.typeTransaction,
            dateHour: t.dateHour
        }))

        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#transactions.length,
        }
    }

    async Get(uuid: UUID): Promise<transaction.Response> {
        for (const t of this.#transactions) {
            if (t.uuid === uuid) {
                return {
                    uuid: t.uuid,
                    amount: t.amount,
                    typeTransaction: t.typeTransaction,
                    dateHour: t.dateHour
                }
            }
        }

        throw APIError.New("not found", "transaction-not-found", `transaction ${uuid} not found`)
    }

    async Create(req: transaction.Entity): Promise<UUID> {
        const uuid = crypto.randomUUID()

        this.#transactions.push({
            uuid: uuid,
            amount: req.amount,
            typeTransaction: req.typeTransaction,
            dateHour: req.dateHour
        })

        return uuid
    }

    async Patch(uuid: UUID, req: transaction.PartialEntity): Promise<void> {
        for (const t of this.#transactions) {
            if (t.uuid === uuid) {
                if (req.amount !== undefined) t.amount = req.amount
                if (req.typeTransaction !== undefined) t.typeTransaction = req.typeTransaction
                if (req.dateHour !== undefined) t.dateHour = req.dateHour
                return
            }
        }

        throw APIError.New("not found", "transaction-not-found", `transaction ${uuid} not found`)
    }

    async Delete(uuid: UUID): Promise<void> {
        for (let i = 0; i < this.#transactions.length; i++) {
            if (this.#transactions[i].uuid === uuid) {
                this.#transactions.splice(i, 1)
                return
            }
        }
    }
}