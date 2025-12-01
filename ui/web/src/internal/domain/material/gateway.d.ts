namespace material {
    interface Gateway {
        List(offset: number, limit: number): Promise<ListResponse>
        ListBySIADS(siads: string, offset: number, limit: number): Promise<ListResponse>
        ListByCATMAT(catmat: string, offset: number, limit: number): Promise<ListResponse>
        ListByECampus(ecampus: string, offset: number, limit: number): Promise<ListResponse>
        Get(uuid: UUID): Promise<Response>
        Create(req: Entity): Promise<UUID>
        Patch(uuid: UUID, req: PartialEntity): Promise<void>
        Delete(uuid: UUID): Promise<void>
    }

    type Entity = {
        name: string
        siads: string
        catmat: string
        ecampus: string
        description: string
        unit: string
        min_quantity: number
    }

    type PartialEntity = {
        name?: string
        siads?: string
        catmat?: string
        ecampus?: string
        description?: string
        unit?: string
        min_quantity?: number
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
        siads: string
        catmat: string
        ecampus: string
        description: string
        unit: string
        min_quantity: number
        created: Date
        updated: Date
    }
}