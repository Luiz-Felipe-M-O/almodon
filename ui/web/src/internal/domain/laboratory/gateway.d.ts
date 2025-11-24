namespace laboratory {
    interface Gateway {
        async List(offset: number, limit: number): Promise<ListResponse>
        async Get(uuid: UUID): Promise<Response>
        async Create(req: Entity): Promise<UUID>
        async Patch(uuid: UUID, req: PartialEntity): Promise<void>
        async Delete(uuid: UUID): Promise<void>
    }

    type Entity = {
        idLaboratory: number
        name: string
    }

    type PartialEntity = {
        idLaboratory?: number
        name?: string
    }

    type ListResponse = {
        offset: number
        length: number
        records: Response[]
        total_records: number
    }

    type Response = {
        uuid: UUID
        idLaboratory: number
        name: string
    }
}