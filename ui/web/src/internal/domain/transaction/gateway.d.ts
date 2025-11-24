namespace transaction {
    interface Gateway {
        async List(offset: number, limit: number): Promise<ListResponse>
        async Get(uuid: UUID): Promise<Response>
        async Create(req: Entity): Promise<UUID>
        async Patch(uuid: UUID, req: PartialEntity): Promise<void>
        async Delete(uuid: UUID): Promise<void>
    }

    type Entity = {
        amount: number
        typeTransaction: string
        dateHour: Date
    }

    type PartialEntity = {
        amount?: number
        typeTransaction?: string
        dateHour?: Date
    }

    type ListResponse = {
        offset: number
        length: number
        records: Response[]
        total_records: number
    }

    type Response = {
        uuid: UUID
        amount: number
        typeTransaction: string
        dateHour: Date
    }
}