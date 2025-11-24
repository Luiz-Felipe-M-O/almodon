import { APIError } from "../../../../module/errors/error.ts"

type Clinic = {
    uuid: UUID
    idclinic: number
    name: string
}

export class ClinicGateway implements clinic.Gateway {
    #clinics: Clinic[]

    constructor() {
        this.#clinics = [
            {
                uuid: crypto.randomUUID(),
                idclinic: 1,
                name: "Clínica Central"
            }
        ]
    }

    async List(offset: number, limit: number): Promise<clinic.ListResponse> {
        const lo = offset
        const hi = offset + limit

        const records = this.#clinics.slice(lo, hi).map(clinic => ({
            uuid: clinic.uuid,
            idclinic: clinic.idclinic,
            name: clinic.name,
        }))

        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#clinics.length
        }
    }

    async Get(uuid: UUID): Promise<clinic.Response> {
        const clinic = this.#clinics.find(c => c.uuid === uuid)

        if (!clinic) {
            throw APIError.New("not found", "clinic-not-found", `clinic with UUID ${uuid} not found`)
        }

        return {
            uuid: clinic.uuid,
            idclinic: clinic.idclinic,
            name: clinic.name,
        }
    }

    async Create(req: clinic.Entity): Promise<UUID> {
        const uuid = crypto.randomUUID()

        this.#clinics.push({
            uuid: uuid,
            idclinic: req.idclinic,
            name: req.name,
        })

        return uuid
    }

    async Patch(uuid: UUID, req: clinic.PartialEntity): Promise<void> {
        const clinic = this.#clinics.find(c => c.uuid === uuid)

        if (!clinic) {
            throw APIError.New("not found", "clinic-not-found", `clinic with UUID ${uuid} not found`)
        }

        if (req.idclinic !== undefined) clinic.idclinic = req.idclinic
        if (req.name !== undefined) clinic.name = req.name
    }

    async Delete(uuid: UUID): Promise<void> {
        const index = this.#clinics.findIndex(c => c.uuid === uuid)

        if (index === -1) {
            throw APIError.New("not found", "clinic-not-found", `clinic with UUID ${uuid} not found`)
        }

        this.#clinics.splice(index, 1)
    }
}