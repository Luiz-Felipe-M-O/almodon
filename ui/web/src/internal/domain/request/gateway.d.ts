namespace request {
    interface Gateway {
        async List(offset: number, limit: number): Promise<ListResponse>
        async Get(uuid: UUID): Promise<Response>
        async Create(req: Entity): Promise<UUID>
        async Patch(uuid: UUID, req: PartialEntity): Promise<void>
        async Delete(uuid: UUID): Promise<void>
    }

    type Entity = {
        status: string
        dateRequest: Date
        reason: string
    }

    type PartialEntity = {
        status?: string
        dateRequest?: Date
        reason?: string
    }

    type ListResponse = {
        offset: number
        length: number
        records: Response[]
        total_records: number
    }

    type Response = {
        uuid: UUID
        status: string
        dateRequest: Date
        reason: string
    }
}