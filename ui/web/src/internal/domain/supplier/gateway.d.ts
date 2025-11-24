namespace supplier {
    interface Gateway {
        async List(offset: number, limit: number): Promise<ListResponse>
        async Get(uuid: UUID): Promise<Response>
        async Create(req: Entity): Promise<UUID>
        async Patch(uuid: UUID, req: PartialEntity): Promise<void>
        async Delete(uuid: UUID): Promise<void>
    }

    type Entity = {
        name: string
        cnpj: string
        contact: number
    }

    type PartialEntity = {
        name?: string
        cnpj?: string
        contact?: number
    }

    type ListResponse = {
        offset: number
        length: number
        records: Response[]
        total_records: number
    }

    type Response = {
        uuid: UUID
        name: string
        cnpj: string
        contact: number
    }
}