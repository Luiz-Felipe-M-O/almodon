namespace user {
    interface Gateway {
        async List(offset: number, limit: number): Promise<ListResponse>
        async Get(uuid: UUID): Promise<Response>
        async Create(req: Entity): Promise<UUID>
        async Patch(uuid: UUID, req: PartialEntity): Promise<void>
        async Delete(uuid: UUID): Promise<void>
        async Autheticate(siape: number, password: string): Promise<AuthResponse>
    }

    type Entity = {
        siape: number
        name: string
        email: string
        password: string
        role: string
    }

    type PartialEntity = {
        siape?: number
        name?: string
        email?: string
    }

    type ListResponse = {
        offset: number
        length: number
        records: Response[]
        total_records: number
    }

    type Response = {
        uuid: UUID
        siape: number
        name: string
        email: string
        role: string
    }

    type AuthResponse = {
        uuid: UUID
        user: UUID
        expires: Date
    }
}
