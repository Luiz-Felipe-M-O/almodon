namespace user {
    interface Gateway {
        List(offset: number, limit: number): Promise<ListResponse>
        Get(uuid: UUID): Promise<Response>
        GetBySIAPE(siape: string): Promise<Response>
        Create(req: Entity): Promise<UUID>
        Patch(uuid: UUID, req: PartialEntity): Promise<void>
        Delete(uuid: UUID): Promise<void>
        Autheticate(siape: string, password: string): Promise<AuthResponse>
        Logout(): Promise<void>
        Me(): Promise<Response>
    }

    type Entity = {
        siape: string
        name: string
        email: string
        password: string
        role: Role
    }

    type PartialEntity = {
        siape?: string
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
        siape: string
        name: string
        email: string
        role: Role
        created: Date
        updated: Date
    }

    type AuthResponse = {
        uuid: UUID
        user: UUID
        expires: Date
    }
}
