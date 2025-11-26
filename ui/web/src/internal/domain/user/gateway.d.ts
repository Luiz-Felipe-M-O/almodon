namespace user {
    interface Gateway {
        List(offset: number, limit: number): Promise<ListResponse>
        Get(uuid: UUID): Promise<Response>
        Create(req: Entity): Promise<UUID>
        Patch(uuid: UUID, req: PartialEntity): Promise<void>
        Delete(uuid: UUID): Promise<void>
        Autheticate(siape: number, password: string): Promise<AuthResponse>
        Me(): Promise<Response>
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
