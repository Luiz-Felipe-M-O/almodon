import { APIError } from "../../../../module/errors/error.ts"

type Product = {
    uuid: UUID
    catmat: string
    siads: string
    codeEcampus: string
    description: string
    name: string
    stockMinimum: number
}

export class ProductGateway implements product.Gateway {
    #products: Product[]

    constructor() {
        this.#products = [
            {
                uuid: "00000001",
                catmat: "123456",
                siads: "987654",
                codeEcampus: "MAT-001",
                description: "Álcool Etílico 92% 1L",
                name: "Álcool 92%",
                stockMinimum: 10
            },
            {
                uuid: "00000002",
                catmat: "654321",
                siads: "456789",
                codeEcampus: "MAT-002",
                description: "Luva de Látex Procedimento M",
                name: "Luva Látex M",
                stockMinimum: 100
            }
        ]
    }

    async List(offset: number, limit: number): Promise<product.ListResponse> {
        const lo = offset
        const hi = offset + limit

        const records = this.#products.slice(lo, hi).map(p => ({
            uuid: p.uuid,
            catmat: p.catmat,
            siads: p.siads,
            codeEcampus: p.codeEcampus,
            description: p.description,
            name: p.name,
            stockMinimum: p.stockMinimum
        }))

        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#products.length,
        }
    }

    async Get(uuid: UUID): Promise<product.Response> {
        for (const p of this.#products) {
            if (p.uuid === uuid) {
                return { ...p }
            }
        }

        throw APIError.New("not found", "product-not-found", `product with UUID ${uuid} not found`)
    }

    private findByCodeEcampus(code: string): Product | undefined {
        return this.#products.find(p => p.codeEcampus === code)
    }

    async Create(req: product.Entity): Promise<UUID> {
        if (this.findByCodeEcampus(req.codeEcampus)) {
             throw APIError.New("conflict", "code-exists", `product with codeEcampus ${req.codeEcampus} already exists`)
        }

        const uuid = crypto.randomUUID()
        this.#products.push({
            uuid: uuid,
            catmat: req.catmat,
            siads: req.siads,
            codeEcampus: req.codeEcampus,
            description: req.description,
            name: req.name,
            stockMinimum: req.stockMinimum
        })

        return uuid
    }

    async Patch(uuid: UUID, req: product.PartialEntity): Promise<void> {
        for (const p of this.#products) {
            if (p.uuid === uuid) {
                if (req.catmat !== undefined) { p.catmat = req.catmat }
                if (req.siads !== undefined) { p.siads = req.siads }
                if (req.codeEcampus !== undefined) { 
                    const conflict = this.findByCodeEcampus(req.codeEcampus)
                    if (conflict && conflict.uuid !== uuid) {
                        throw APIError.New("conflict", "code-exists", `product with codeEcampus ${req.codeEcampus} already exists`)
                    }
                    p.codeEcampus = req.codeEcampus 
                }
                if (req.description !== undefined) { p.description = req.description }
                if (req.name !== undefined) { p.name = req.name }
                if (req.stockMinimum !== undefined) { p.stockMinimum = req.stockMinimum }
                return
            }
        }

        throw APIError.New("not found", "product-not-found", `product with UUID ${uuid} not found`)
    }

    async Delete(uuid: UUID): Promise<void> {
        for (let i = 0; i < this.#products.length; i++) {
            if (this.#products[i].uuid === uuid) {
                this.#products.splice(i, 1)
                return
            }
        }
    }
}