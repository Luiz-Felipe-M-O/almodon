import { APIError } from "../../../../module/errors/error.ts"

type User = {
    uuid: UUID
    siape: number
    name: string
    email: string
    password: string
    role: string
}

export class UserGateway implements user.Gateway {
    #users: User[]
    #session: UUID | null

    constructor() {
        this.#users = [
            {
                uuid: "00000000-0000-0000-0000-000000000001",
                siape: 123456,
                name: "Alan Lima",
                email: "example@example.com",
                password: "12345678",
                role: "admin",
            }
        ]

        this.#session = null
    }

    async List(offset: number, limit: number): Promise<user.ListResponse> {
        const lo = offset
        const hi = offset + limit

        const records = this.#users.slice(lo, hi).map(user => ({
            uuid: user.uuid,
            siape: user.siape,
            name: user.name,
            email: user.email,
            role: user.role,
        }))

        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#users.length,
        }
    }

    async Get(uuid: UUID): Promise<user.Response> {
        for (const user of this.#users) {
            if (user.uuid === uuid) {
                return {
                    uuid: user.uuid,
                    siape: user.siape,
                    name: user.name,
                    email: user.email,
                    role: user.role,
                }
            }
        }

        throw APIError.New("not found", "user-not-found", `user with UUID ${uuid} not found`)
    }

    async GetBySIAPE(siape: number): Promise<user.Response> {
        for (const user of this.#users) {
            if (user.siape === siape) {
                return {
                    uuid: user.uuid,
                    siape: user.siape,
                    name: user.name,
                    email: user.email,
                    role: user.role,
                }
            }
        }

        throw APIError.New("not found", "user-not-found", `user with SIAPE ${siape} not found`)
    }

    async Create(req: user.Entity): Promise<UUID> {
        try {
            await this.GetBySIAPE(req.siape)
            throw APIError.New("not found", "siape-in-user", `user with SIAPE ${req.siape} was found`)
        } catch { }

        const uuid = crypto.randomUUID()
        this.#users.push({
            uuid: uuid,
            siape: req.siape,
            name: req.name,
            email: req.email,
            password: req.password,
            role: req.role,
        })

        return uuid
    }

    async Patch(uuid: UUID, req: user.PartialEntity): Promise<void> {
        for (const user of this.#users) {
            if (user.uuid === uuid) {
                if (req.siape !== undefined) { user.siape = req.siape }
                if (req.name !== undefined) { user.name = req.name }
                if (req.email !== undefined) { user.email = req.email }
            }
        }

        throw APIError.New("not found", "user-not-found", `user with UUID ${uuid} not found`)
    }

    async Delete(uuid: UUID): Promise<void> {
        for (let i = 0; i < this.#users.length; i++) {
            if (this.#users[i].uuid === uuid) {
                this.#users.splice(i, 1)
            }
        }
    }

    async Autheticate(siape: number, password: string): Promise<user.AuthResponse> {
        let user: User | undefined = undefined
        for (const u of this.#users) {
            if (u.siape === siape) {
                if (u.password !== password) {
                    throw APIError.New("unauthorized", "incorrect-password", "invalid SIAPE or password")
                }

                user = u
            }
        }

        if (user === undefined) {
            throw APIError.New("not found", "user-not-found", `user with SIAPE ${siape} not found`)
        }

        const session = crypto.randomUUID()
        const expires = 3600 * 1000

        setTimeout(() => { this.#session = null }, expires)
        this.#session = session

        return {
            uuid: session,
            user: user.uuid,
            expires: new Date(Date.now() + expires),
        }
    }

    async Me(): Promise<user.Response> {
        if (this.#session === null) {
            throw APIError.New("unauthorized", "no-active-session", "no active session")
        }

        return await this.Get(this.#session)
    }
}