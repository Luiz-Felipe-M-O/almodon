import { AsyncTry } from "../../../module/errors/try.ts"
import jsxmm from "../../../module/jsxmm/element.ts"
import { StatusPage } from "../../component/status.ts"
import UserPresenter from "./presenter.ts"

export default class UserView {
    #gateway: user.Gateway
    #build: HTMLElement
    #done: boolean

    constructor(gateway: user.Gateway) {
        this.#gateway = gateway

        this.#build = jsxmm.Element("div", { id: "almodon" })
        this.#done = false
    }

    Final(): boolean {
        return this.#done
    }

    HTML(): HTMLElement {
        if (!this.#done) {
            this.load().then(result => {
                this.#build.replaceChildren(result)
                this.#done = true
            }).catch(error => {
                this.#build.replaceChildren(error)
            })
        }

        return this.#build
    }

    private async load(): Promise<HTMLElement> {
        const [response, error] = await AsyncTry(() => this.#gateway.List(0, Number.MAX_SAFE_INTEGER))
        if (error !== null) {
            throw StatusPage(403)
        }

        return UserPresenter.Table(response)
    }
}  