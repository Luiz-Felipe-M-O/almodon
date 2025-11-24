namespace batch {
    interface Gateway {
        async List(offset: number, limit: number): Promise<ListResponse>
        async Get(uuid: UUID): Promise<Response>
        async Create(req: Entity): Promise<UUID>
        async Patch(uuid: UUID, req: PartialEntity): Promise<void>
        async Delete(uuid: UUID): Promise<void>
    }

    type Entity = {
        idBatch: number
        amount: number
        codeBatch: string
        observationsNote: string
        specifications: string
        valueUnitary: number
        dateEntry: string
        dateExpiration: string
    }

    type PartialEntity = {
        amount?: number
        codeBatch?: string
        observationsNote?: string
        specifications?: string
        valueUnitary?: number
        dateEntry?: string
        dateExpiration?: string
    }

    type ListResponse = {
        offset: number
        length: number
        records: Response[]
        total_records: number
    }

    type Response = {
        uuid: UUID
        idBatch: number
        amount: number;
        codeBatch: string;
        observationsNote: string;
        specifications: string;
        valueUnitary: number;
        dateEntry: string;
        dateExpiration: string;
    }
}