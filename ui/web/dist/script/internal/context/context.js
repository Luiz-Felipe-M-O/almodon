var __rewriteRelativeImportExtension = (this && this.__rewriteRelativeImportExtension) || function (path, preserveJsx) {
    if (typeof path === "string" && /^\.\.?\//.test(path)) {
        return path.replace(/\.(tsx)$|((?:\.d)?)((?:\.[^./]+?)?)\.([cm]?)ts$/i, function (m, tsx, d, ext, cm) {
            return tsx ? preserveJsx ? ".jsx" : ".js" : d && (!ext || !cm) ? m : (d + ext + "." + cm.toLowerCase() + "js");
        });
    }
    return path;
};
import { AsyncTry } from "../../module/errors/try.js";
import jsxmm from "../../module/jsxmm/element.js";
import { StatusPage } from "../component/status.js";
import Source from "../support/source.js";
export class Orquestrator {
    #room;
    #contexts;
    #swapper;
    constructor(placeholder, swapper = new hash()) {
        this.#room = placeholder;
        this.#contexts = {};
        this.#swapper = swapper;
    }
    Current() {
        return this.#swapper.Namespace();
    }
    Link(namespace, context) {
        this.#contexts[namespace] = context;
    }
    Unlink(namespace) {
        delete this.#contexts[namespace];
    }
    SwapTo(namespace) {
        if (!Object.hasOwn(this.#contexts, namespace)) {
            return false;
        }
        const context = this.#contexts[namespace];
        if (context.Final() && this.Current() === namespace) {
            return true;
        }
        const content = context.HTML();
        this.#room.replaceChildren(content);
        this.#swapper.SwapNamespace(namespace);
        return true;
    }
}
export class context {
    #url;
    #content;
    #retry;
    onpreload;
    onload;
    constructor(url) {
        this.#url = url;
        this.#retry = true;
        this.#content = jsxmm.Element("div", { id: "almodon" });
    }
    Final() {
        return !this.#retry;
    }
    HTML() {
        if (this.#retry) {
            try_callback(this.onpreload);
            load(this.#url).then(([result, ok]) => {
                this.#content.replaceWith(result);
                this.#content = result;
                this.#retry = !ok;
                try_callback(this.onload);
            });
            return this.#content;
        }
        return this.#content;
    }
}
const Parser = new DOMParser();
async function load(url) {
    const [result, error] = await AsyncTry(fetch, url);
    if (error !== null) {
        throw error;
    }
    if (!result.ok) {
        return [StatusPage(result.status), false];
    }
    const page = await result.text();
    const new_document = Parser.parseFromString(page, "text/html");
    const content = new_document.getElementById("almodon");
    if (content === null) {
        return [StatusPage(204), false];
    }
    const element = new_document.getElementById("meta-almodon");
    if (element !== null) {
        for (const property of element.children) {
            switch (property.tagName) {
                case "ALMODON-SCRIPT":
                    const src = property.dataset["src"];
                    if (src !== undefined) {
                        await import(__rewriteRelativeImportExtension(new URL(src, url).href));
                    }
                    break;
                case "ALMODON-STYLE":
                    const href = property.dataset["href"];
                    if (href !== undefined) {
                        const style = jsxmm.Element("link", {
                            rel: "stylesheet",
                            href: Source.From(href, url),
                        });
                        document.head.append(style);
                        await new Promise(resolve => {
                            style.onload = resolve;
                        });
                    }
                    break;
                default:
                    throw new Error("Unrecognized property " + property.tagName.toLocaleLowerCase());
            }
        }
    }
    return [content, true];
}
export class hash {
    Namespace() {
        return location.hash.slice(1);
    }
    SwapNamespace(namespace) {
        location.hash = namespace;
    }
}
function try_callback(callback) {
    if (callback !== undefined) {
        callback();
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29udGV4dC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9pbnRlcm5hbC9jb250ZXh0L2NvbnRleHQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7QUFBQSxPQUFPLEVBQUUsUUFBUSxFQUFFLE1BQU0sNEJBQTRCLENBQUE7QUFDckQsT0FBTyxLQUFLLE1BQU0sK0JBQStCLENBQUE7QUFDakQsT0FBTyxFQUFFLFVBQVUsRUFBRSxNQUFNLHdCQUF3QixDQUFBO0FBQ25ELE9BQU8sTUFBTSxNQUFNLHNCQUFzQixDQUFBO0FBV3pDLE1BQU0sT0FBTyxZQUFZO0lBQ3hCLEtBQUssQ0FBYTtJQUVsQixTQUFTLENBQXlCO0lBQ2xDLFFBQVEsQ0FBUztJQUVqQixZQUFZLFdBQXdCLEVBQUUsVUFBbUIsSUFBSSxJQUFJLEVBQUU7UUFDbEUsSUFBSSxDQUFDLEtBQUssR0FBRyxXQUFXLENBQUE7UUFFeEIsSUFBSSxDQUFDLFNBQVMsR0FBRyxFQUFFLENBQUE7UUFDbkIsSUFBSSxDQUFDLFFBQVEsR0FBRyxPQUFPLENBQUE7SUFDeEIsQ0FBQztJQUVELE9BQU87UUFDTixPQUFPLElBQUksQ0FBQyxRQUFRLENBQUMsU0FBUyxFQUFFLENBQUE7SUFDakMsQ0FBQztJQUVELElBQUksQ0FBQyxTQUFpQixFQUFFLE9BQWdCO1FBQ3ZDLElBQUksQ0FBQyxTQUFTLENBQUMsU0FBUyxDQUFDLEdBQUcsT0FBTyxDQUFBO0lBQ3BDLENBQUM7SUFFRCxNQUFNLENBQUMsU0FBaUI7UUFDdkIsT0FBTyxJQUFJLENBQUMsU0FBUyxDQUFDLFNBQVMsQ0FBQyxDQUFBO0lBQ2pDLENBQUM7SUFFRCxNQUFNLENBQUMsU0FBaUI7UUFDdkIsSUFBSSxDQUFDLE1BQU0sQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLFNBQVMsRUFBRSxTQUFTLENBQUMsRUFBRSxDQUFDO1lBQy9DLE9BQU8sS0FBSyxDQUFBO1FBQ2IsQ0FBQztRQUVELE1BQU0sT0FB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