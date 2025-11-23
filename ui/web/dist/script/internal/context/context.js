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
    #current;
    #swapper;
    constructor(placeholder, swapper = new hash()) {
        this.#room = placeholder;
        this.#contexts = {};
        this.#current = undefined;
        this.#swapper = swapper;
    }
    SwapperCurrent() {
        return this.#swapper.Namespace();
    }
    Current() {
        return this.#current;
    }
    Link(namespace, context) {
        this.#contexts[namespace] = context;
    }
    Unlink(namespace) {
        delete this.#contexts[namespace];
    }
    SwapTo(namespace) {
        if (this.Current() === namespace) {
            return true;
        }
        const context = this.#contexts[namespace];
        if (context === undefined) {
            return false;
        }
        const content = context.HTML();
        this.#room.replaceChildren(content);
        this.#swapper.SwapNamespace(namespace);
        this.#current = namespace;
        return true;
    }
}
export class context {
    static parser = new DOMParser();
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
    HTML() {
        if (this.#retry) {
            try_callback(this.onpreload);
            context.load(this.#url).then(([result, ok]) => {
                this.#content.replaceWith(result);
                this.#content = result;
                this.#retry = !ok;
                try_callback(this.onload);
            });
            return this.#content;
        }
        return this.#content;
    }
    static async load(url) {
        const [result, error] = await AsyncTry(fetch, url);
        if (error !== null) {
            throw error;
        }
        if (!result.ok) {
            return [StatusPage(result.status), false];
        }
        const page = await result.text();
        const new_document = context.parser.parseFromString(page, "text/html");
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29udGV4dC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9pbnRlcm5hbC9jb250ZXh0L2NvbnRleHQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7QUFBQSxPQUFPLEVBQUUsUUFBUSxFQUFFLE1BQU0sNEJBQTRCLENBQUE7QUFDckQsT0FBTyxLQUFLLE1BQU0sK0JBQStCLENBQUE7QUFDakQsT0FBTyxFQUFFLFVBQVUsRUFBRSxNQUFNLHdCQUF3QixDQUFBO0FBQ25ELE9BQU8sTUFBTSxNQUFNLHNCQUFzQixDQUFBO0FBT3pDLE1BQU0sT0FBTyxZQUFZO0lBQ3hCLEtBQUssQ0FBYTtJQUVsQixTQUFTLENBQTJCO0lBQ3BDLFFBQVEsQ0FBb0I7SUFDNUIsUUFBUSxDQUFTO0lBRWpCLFlBQVksV0FBd0IsRUFBRSxVQUFtQixJQUFJLElBQUksRUFBRTtRQUNsRSxJQUFJLENBQUMsS0FBSyxHQUFHLFdBQVcsQ0FBQTtRQUV4QixJQUFJLENBQUMsU0FBUyxHQUFHLEVBQUUsQ0FBQTtRQUNuQixJQUFJLENBQUMsUUFBUSxHQUFHLFNBQVMsQ0FBQTtRQUN6QixJQUFJLENBQUMsUUFBUSxHQUFHLE9BQU8sQ0FBQTtJQUN4QixDQUFDO0lBRUQsY0FBYztRQUNiLE9BQU8sSUFBSSxDQUFDLFFBQVEsQ0FBQyxTQUFTLEVBQUUsQ0FBQTtJQUNqQyxDQUFDO0lBRUQsT0FBTztRQUNOLE9BQU8sSUFBSSxDQUFDLFFBQVEsQ0FBQTtJQUNyQixDQUFDO0lBRUQsSUFBSSxDQUFDLFNBQWlCLEVBQUUsT0FBa0I7UUFDekMsSUFBSSxDQUFDLFNBQVMsQ0FBQyxTQUFTLENBQUMsR0FBRyxPQUFPLENBQUE7SUFDcEMsQ0FBQztJQUVELE1BQU0sQ0FBQyxTQUFpQjtRQUN2QixPQUFPLElBQUksQ0FBQyxTQUFTLENBQUMsU0FBUyxDQUFDLENBQUE7SUFDakMsQ0FBQztJQUVELE1BQU0sQ0FBQyxTQUFpQjtRQUN2QixJQUFJLElBQUks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