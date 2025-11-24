import { APIError } from "../../../../module/errors/error.ts"

type Supplier = {
    uuid: UUID
    name: string
    cnpj: string
    contact: number
}

export class SupplierGateway implements supplier.Gateway {
    #suppliers: Supplier[]

    constructor() {
        this.#suppliers = [
            {
                uuid: "00000000-0000-0000-0000-000000000001",
                name: "Papelaria Central Ltda",
                cnpj: "12.345.678/0001-90",
                contact: 11999999999
            },
            {
                uuid: "00000000-0000-0000-0000-000000000002",
                name: "Tech Hardware Solutions",
                cnpj: "98.765.432/0001-10",
                contact: 4133333333
            }
        ]
    }

    async List(offset: number, limit: number): Promise<supplier.ListResponse> {
        const lo = offset
        const hi = offset + limit

        const records = this.#suppliers.slice(lo, hi).map(s => ({
            uuid: s.uuid,
            name: s.name,
            cnpj: s.cnpj,
            contact: s.contact,
        }))

        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#suppliers.length,
        }
    }

    async Get(uuid: UUID): Promise<supplier.Response> {
        for (const s of this.#suppliers) {
            if (s.uuid === uuid) {
                return {
                    uuid: s.uuid,
                    name: s.name,
                    cnpj: s.cnpj,
                    contact: s.contact,
                }
            }
        }

        throw APIError.New("not found", "supplier-not-found", `supplier with UUID ${uuid} not found`)
    }

    private findByCNPJ(cnpj: string): Supplier | undefined {
        return this.#suppliers.find(s => s.cnpj === cnpj)
    }

    async Create(req: supplier.Entity): Promise<UUID> {
        if (this.findByCNPJ(req.cnpj)) {
            throw APIError.New("conflict", "cnpj-exists", `supplier with CNPJ ${req.cnpj} already exists`)
        }

        const uuid = crypto.randomUUID()
        this.#suppliers.push({
            uuid: uuid,
            name: req.name,
            cnpj: req.cnpj,
            contact: req.contact,
        })

        return uuid
    }

    async Patch(uuid: UUID, req: supplier.PartialEntity): Promise<void> {
        for (const s of this.#suppliers) {
            if (s.uuid === uuid) {
                if (req.name !== undefined) { s.name = req.name }
                
                if (req.cnpj !== undefined) { 
                    const conflict = this.findByCNPJ(req.cnpj)
                    if (conflict && conflict.uuid !== uuid) {
                         throw APIError.New("conflict", "cnpj-exists", `supplier with CNPJ ${req.cnpj} already exists`)
                    }
                    s.cnpj = req.cnpj 
                }
                
                if (req.contact !== undefined) { s.contact = req.contact }
                return
            }
        }

        throw APIError.New("not found", "supplier-not-found", `supplier with UUID ${uuid} not found`)
    }

    async Delete(uuid: UUID): Promise<void> {
        for (let i = 0; i < this.#suppliers.length; i++) {
            if (this.#suppliers[i].uuid === uuid) {
                this.#suppliers.splice(i, 1)
                return
            }
        }
    }
}