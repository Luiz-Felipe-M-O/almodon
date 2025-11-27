import { APIError } from "../../../../module/errors/error.js";
export class UserGateway {
    #users;
    #session;
    constructor() {
        this.#users = [
            {
                uuid: "019a921e-02ff-7cfb-6c98-a9d888ebe4a2",
                siape: 1,
                name: "Alan Lima",
                email: "alan-lima.al@ufvjm.edu.br",
                password: "12345678",
                role: "chief"
            },
            {
                uuid: "019a921e-3fb1-7d33-5b88-9b569416db4c",
                siape: 2,
                name: "Breno",
                email: "breno@ufvjm.edu.br",
                password: "12345678",
                role: "admin"
            },
            {
                uuid: "019a9220-d660-7e60-70d5-7dc1dc580a02",
                siape: 3,
                name: "Luiz",
                email: "lf@ufvjm.edu.br",
                password: "12345678",
                role: "user"
            },
            {
                uuid: "019a9798-2c46-74e0-5dea-3c7c98a45599",
                siape: 4,
                name: "Rafael",
                email: "r@ufvjm.edu.br",
                password: "12345678",
                role: "user"
            },
            {
                uuid: "019a9f7c-ad1b-70a7-49a3-58828f3b66d6",
                siape: 5,
                name: "Otávio Calazans",
                email: "tavinhogomesoficial@hotmail.com",
                password: "12345678",
                role: "user"
            },
            {
                uuid: "019aa844-0fbf-78d9-7422-caf8961a8ccb",
                siape: 6,
                name: "Lucas",
                email: "rocha@ufvjm.edu.br",
                password: "12345678",
                role: "admin"
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
        this.#session = session;
        return {
            uuid: session,
            user: user.uuid,
            expires: new Date(Date.now() + expires),
        };
    }
    async Me() {
        if (this.#session === null) {
            throw APIError.New("unauthorized", "no-active-session", "no active session");
        }
        return await this.Get(this.#session);
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibW9jay5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uL3NyYy9pbnRlcm5hbC9kb21haW4vdXNlci9nYXRld2F5L21vY2sudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUEsT0FBTyxFQUFFLFFBQVEsRUFBRSxNQUFNLG9DQUFvQyxDQUFBO0FBVzdELE1BQU0sT0FBTyxXQUFXO0lBQ3BCLE1BQU0sQ0FBUTtJQUNkLFFBQVEsQ0FBYTtJQUVyQjtRQUNJLElBQUksQ0FBQyxNQUFNLEdBQUc7WUFDVjtnQkFDSSxJQUFJLEVBQUUsc0NBQXNDO2dCQUM1QyxLQUFLLEVBQUUsQ0FBQztnQkFDUixJQUFJLEVBQUUsV0FBVztnQkFDakIsS0FBSyxFQUFFLDJCQUEyQjtnQkFDbEMsUUFBUSxFQUFFLFVBQVU7Z0JBQ3BCLElBQUksRUFBRSxPQUFPO2FBQ2hCO1lBQ0Q7Z0JBQ0ksSUFBSSxFQUFFLHNDQUFzQztnQkFDNUMsS0FBSyxFQUFFLENBQUM7Z0JBQ1IsSUFBSSxFQUFFLE9BQU87Z0JBQ2IsS0FBSyxFQUFFLG9CQUFvQjtnQkFDM0IsUUFBUSxFQUFFLFVBQVU7Z0JBQ3BCLElBQUksRUFBRSxPQUFPO2FBQ2hCO1lBQ0Q7Z0JBQ0ksSUFBSSxFQUFFLHNDQUFzQztnQkFDNUMsS0FBSyxFQUFFLENBQUM7Z0JBQ1IsSUFBSSxFQUFFLE1BQU07Z0JBQ1osS0FBSyxFQUFFLGlCQUFpQjtnQkFDeEIsUUFBUSxFQUFFLFVBQVU7Z0JBQ3BCLElBQUksRUFBRSxNQUFNO2FBQ2Y7WUFDRDtnQkFDSSxJQUFJLEVBQUUsc0NBQXNDO2dCQUM1QyxLQUFLLEVBQUUsQ0FBQztnQkFDUixJQUFJLEVBQUUsUUFBUTtnQkFDZCxLQUFLLEVBQUUsZ0JBQWdCO2dCQUN2QixRQUFRLEVBQUUsVUFBVTtnQkFDcEIsSUFBSSxFQUFFLE1BQU07YUFDZjtZQUNEO2dCQUNJLElBQUksRUFBRSxzQ0FBc0M7Z0JBQzVDLEtBQUssRUFBRSxDQUFDO2dCQUNSLElB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