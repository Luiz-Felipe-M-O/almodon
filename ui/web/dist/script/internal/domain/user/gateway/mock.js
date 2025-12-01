import { APIError } from "../../../../module/errors/error.js";
export class UserGateway {
    #users;
    #session;
    constructor() {
        this.#users = [
            {
                uuid: "019a921e-02ff-7cfb-6c98-a9d888ebe4a2",
                siape: "0000001",
                name: "Alan Lima",
                email: "alan-lima.al@ufvjm.edu.br",
                password: "12345678",
                role: "chief",
                created: new Date(),
                updated: new Date(),
            },
            {
                uuid: "019a921e-3fb1-7d33-5b88-9b569416db4c",
                siape: "0000002",
                name: "Breno",
                email: "breno@ufvjm.edu.br",
                password: "12345678",
                role: "admin",
                created: new Date(),
                updated: new Date(),
            },
            {
                uuid: "019a9220-d660-7e60-70d5-7dc1dc580a02",
                siape: "0000003",
                name: "Luiz",
                email: "lf@ufvjm.edu.br",
                password: "12345678",
                role: "user",
                created: new Date(),
                updated: new Date(),
            },
            {
                uuid: "019a9798-2c46-74e0-5dea-3c7c98a45599",
                siape: "0000004",
                name: "Rafael",
                email: "r@ufvjm.edu.br",
                password: "12345678",
                role: "user",
                created: new Date(),
                updated: new Date(),
            },
            {
                uuid: "019a9f7c-ad1b-70a7-49a3-58828f3b66d6",
                siape: "0000005",
                name: "Otávio Calazans",
                email: "tavinhogomesoficial@hotmail.com",
                password: "12345678",
                role: "user",
                created: new Date(),
                updated: new Date(),
            },
            {
                uuid: "019aa844-0fbf-78d9-7422-caf8961a8ccb",
                siape: "0000006",
                name: "Lucas",
                email: "rocha@ufvjm.edu.br",
                password: "12345678",
                role: "admin",
                created: new Date(),
                updated: new Date(),
            }
        ];
        this.#session = null;
    }
    async List(offset, limit) {
        const lo = offset;
        const hi = offset + limit;
        const records = this.#users.slice(lo, hi).map(user => ({
            uuid: user.uuid,
            siape: user.siape,
            name: user.name,
            email: user.email,
            role: user.role,
            created: user.created,
            updated: user.updated,
        }));
        return {
            offset: lo,
            length: records.length,
            records: records,
            total_records: this.#users.length,
        };
    }
    async Get(uuid) {
        for (const user of this.#users) {
            if (user.uuid === uuid) {
                return {
                    uuid: user.uuid,
                    siape: user.siape,
                    name: user.name,
                    email: user.email,
                    role: user.role,
                    created: user.created,
                    updated: user.updated,
                };
            }
        }
        throw APIError.New("not found", "user-not-found", `user with UUID ${uuid} not found`);
    }
    async GetBySIAPE(siape) {
        for (const user of this.#users) {
            if (user.siape === siape) {
                return {
                    uuid: user.uuid,
                    siape: user.siape,
                    name: user.name,
                    email: user.email,
                    role: user.role,
                    created: user.created,
                    updated: user.updated,
                };
            }
        }
        throw APIError.New("not found", "user-not-found", `user with SIAPE ${siape} not found`);
    }
    async Create(req) {
        try {
            await this.GetBySIAPE(req.siape);
            throw APIError.New("not found", "siape-in-user", `user with SIAPE ${req.siape} was found`);
        }
        catch { }
        const uuid = crypto.randomUUID();
        this.#users.push({
            uuid: uuid,
            siape: req.siape,
            name: req.name,
            email: req.email,
            password: req.password,
            role: req.role,
            created: new Date(),
            updated: new Date(),
        });
        return uuid;
    }
    async Patch(uuid, req) {
        for (const user of this.#users) {
            if (user.uuid === uuid) {
                if (req.siape !== undefined) {
                    user.siape = req.siape;
                }
                if (req.name !== undefined) {
                    user.name = req.name;
                }
                if (req.email !== undefined) {
                    user.email = req.email;
                }
            }
            user.updated = new Date();
            return;
        }
        throw APIError.New("not found", "user-not-found", `user with UUID ${uuid} not found`);
    }
    async Delete(uuid) {
        for (let i = 0; i < this.#users.length; i++) {
            if (this.#users[i].uuid === uuid) {
                this.#users.splice(i, 1);
            }
        }
    }
    async Autheticate(siape, password) {
        let user = undefined;
        for (const u of this.#users) {
            if (u.siape === siape) {
                if (u.password !== password) {
                    throw APIError.New("unauthorized", "incorrect-password", "invalid SIAPE or password");
                }
                user = u;
            }
        }
        if (user === undefined) {
            throw APIError.New("not found", "user-not-found", `user with SIAPE ${siape} not found`);
        }
        const session = crypto.randomUUID();
        const expires = 3600 * 1000;
        setTimeout(() => { this.#session = null; }, expires);
        this.#session = user.uuid;
        return {
            uuid: session,
            user: user.uuid,
            expires: new Date(Date.now() + expires),
        };
    }
    async Logout() {
        this.#session = null;
    }
    async Me() {
        if (this.#session === null) {
            throw APIError.New("unauthorized", "no-active-session", "no active session");
        }
        return await this.Get(this.#session);
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibW9jay5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uL3NyYy9pbnRlcm5hbC9kb21haW4vdXNlci9nYXRld2F5L21vY2sudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUEsT0FBTyxFQUFFLFFBQVEsRUFBRSxNQUFNLG9DQUFvQyxDQUFBO0FBYzdELE1BQU0sT0FBTyxXQUFXO0lBQ3BCLE1BQU0sQ0FBUTtJQUNkLFFBQVEsQ0FBYTtJQUVyQjtRQUNJLElBQUksQ0FBQyxNQUFNLEdBQUc7WUFDVjtnQkFDSSxJQUFJLEVBQUUsc0NBQXNDO2dCQUM1QyxLQUFLLEVBQUUsU0FBUztnQkFDaEIsSUFBSSxFQUFFLFdBQVc7Z0JBQ2pCLEtBQUssRUFBRSwyQkFBMkI7Z0JBQ2xDLFFBQVEsRUFBRSxVQUFVO2dCQUNwQixJQUFJLEVBQUUsT0FBTztnQkFDYixPQUFPLEVBQUUsSUFBSSxJQUFJLEVBQUU7Z0JBQ25CLE9BQU8sRUFBRSxJQUFJLElBQUksRUFBRTthQUN0QjtZQUNEO2dCQUNJLElBQUksRUFBRSxzQ0FBc0M7Z0JBQzVDLEtBQUssRUFBRSxTQUFTO2dCQUNoQixJQUFJLEVBQUUsT0FBTztnQkFDYixLQUFLLEVBQUUsb0JBQW9CO2dCQUMzQixRQUFRLEVBQUUsVUFBVTtnQkFDcEIsSUFBSSxFQUFFLE9BQU87Z0JBQ2IsT0FBTyxFQUFFLElBQUksSUFBSSxFQUFFO2dCQUNuQixPQUFPLEVBQUUsSUFBSSxJQUFJLEVBQUU7YUFDdEI7WUFDRDtnQkFDSSxJQUFJLEVBQUUsc0NBQXNDO2dCQUM1QyxLQUFLLEVBQUUsU0FBUztnQkFDaEIsSUFBSSxFQUFFLE1BQU07Z0JBQ1osS0FBSyxFQUFFLGlCQUFpQjtnQkFDeEIsUUFBUSxFQUFFLFVBQVU7Z0JBQ3BCLElBQUksRUFBRSxNQUFNO2dCQUNaLE9BQU8sRUFBRSxJQUFJLElBQUksRUFBRTtnQkFDbkIsT0FBTyxFQUFFLElBQUksSUFBSSxFQUFFO2FBQ3RCO1lBQ0Q7Z0JB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