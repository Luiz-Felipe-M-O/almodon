namespace product {
    interface Gateway {
        async List(offset: number, limit: number): Promise<ListResponse>
        async Get(uuid: UUID): Promise<Response>
        async Create(req: Entity): Promise<UUID>
        async Patch(uuid: UUID, req: PartialEntity): Promise<void>
        async Delete(uuid: UUID): Promise<void>
    }

    type Entity = {
        catmat: string
        siads: string
        codeEcampus: string
        description: string
        name: string
        stockMinimum: number
    }

    type PartialEntity = {
        catmat?: string
        siads?: string
        codeEcampus?: string
        description?: string
        name?: string
        stockMinimum?: number
    }

    type ListResponse = {
        offset: number
        length: number
        records: Response[]
        total_records: number
    }

    type Response = {
        uuid: UUID
        catmat: string
        siads: string
        codeEcampus: string
        description: string
        name: string
        stockMinimum: number
    }
}