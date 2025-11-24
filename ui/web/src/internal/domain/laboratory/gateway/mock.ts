import { APIError } from "../../../../module/errors/error.ts"

type Laboratory = {
    uuid: UUID
    idLaboratory: number
    name: string
}

export class LaboratoryGateway implements laboratory.Gateway {
    #labs: Laboratory[]

    constructor() {
        this.#labs = [
            {
                uuid: "00000000-0000-0000-0000-000000000001",
                idLaboratory: 20,
                name: "Laboratório de Materiais Dentários",
            },
            {
                uuid: "00000000-0000-0000-0000-000000000002",
                idLaboratory: 21,
                name: "Laboratório de Dentística",
            },
            {
                uuid: "00000000-0000-0000-0000-000000000003",
                idLaboratory: 22,
                name: "Laboratório de Prótese",
            },
            {
                uuid: "00000000-0000-0000-0000-000000000004",
                idLaboratory: 09,
                name: "Laboratório de Orto/Oclusão",
            }
        ]
    }

    async List(offset: number, limit: number): Promise<laboratory.ListResponse> {
        const lo = offset
        const hi = offset + limit

        const records = this.#labs.slice(lo, hi).map(lab => ({
            uuid: lab.uuid,
            idLaboratory: lab.idLaboratory,
            name: lab.name,
        }))

        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#labs.length,
        }
    }

    async Get(uuid: UUID): Promise<laboratory.Response> {
        for (const lab of this.#labs) {
            if (lab.uuid === uuid) {
                return {
                    uuid: lab.uuid,
                    idLaboratory: lab.idLaboratory,
                    name: lab.name,
                }
            }
        }

        throw APIError.New("not found", "laboratory-not-found", `laboratory with UUID ${uuid} not found`)
    }

    async GetByIdLaboratory(id: number): Promise<laboratory.Response> {
        for (const lab of this.#labs) {
            if (lab.idLaboratory === id) {
                return {
                    uuid: lab.uuid,
                    idLaboratory: lab.idLaboratory,
                    name: lab.name,
                }
            }
        }
        throw APIError.New("not found", "laboratory-not-found", `laboratory with ID ${id} not found`)
    }

    async Create(req: laboratory.Entity): Promise<UUID> {
        try {
            await this.GetByIdLaboratory(req.idLaboratory)
            throw APIError.New("conflict", "id-exists", `laboratory with ID ${req.idLaboratory} already exists`)
        } catch (e) {
            if (e instanceof APIError && e.code !== "laboratory-not-found") {
               throw e;
            }
             if (e instanceof APIError && e.code === "id-exists") {
                throw e;
             }
        }

        const uuid = crypto.randomUUID()
        this.#labs.push({
            uuid: uuid,
            idLaboratory: req.idLaboratory,
            name: req.name,
        })

        return uuid
    }

    async Patch(uuid: UUID, req: laboratory.PartialEntity): Promise<void> {
        for (const lab of this.#labs) {
            if (lab.uuid === uuid) {
                if (req.idLaboratory !== undefined) { lab.idLaboratory = req.idLaboratory }
                if (req.name !== undefined) { lab.name = req.name }
                return
            }
        }

        throw APIError.New("not found", "laboratory-not-found", `laboratory with UUID ${uuid} not found`)
    }

    async Delete(uuid: UUID): Promise<void> {
        for (let i = 0; i < this.#labs.length; i++) {
            if (this.#labs[i].uuid === uuid) {
                this.#labs.splice(i, 1)
                return
            }
        }
    }
}