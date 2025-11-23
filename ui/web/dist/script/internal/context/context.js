export class Orquestrator {
    #room;
    #contexts;
    #current;
    constructor(switchable) {
        this.#room = switchable;
        this.#contexts = {};
        this.#current = undefined;
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
        if (this.#current === namespace) {
            return true;
        }
        const context = this.#contexts[namespace];
        if (context === undefined) {
            return false;
        }
        context.Build().then((content) => {
            this.#room.replaceChildren(content);
        }).catch(() => {
            this.#room.replaceChildren("something went wrong");
            this.Unlink(namespace);
        });
        return true;
    }
}
export class context {
    static parser = new DOMParser();
    #url;
    #content;
    constructor(url) {
        this.#content = undefined;
        this.#url = url;
    }
    async Build() {
        if (this.#content !== undefined) {
            return this.#content;
        }
        const result = await fetch(this.#url);
        if (!result.ok) {
            throw new Error(result.statusText);
        }
        const page = await result.text();
        const document = context.parser.parseFromString(page, "text/html");
        const content = document.getElementById("almodon");
        if (content === null) {
            throw new Error("No content");
        }
        this.#content = content;
        return content;
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29udGV4dC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9pbnRlcm5hbC9jb250ZXh0L2NvbnRleHQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBSUEsTUFBTSxPQUFPLFlBQVk7SUFDeEIsS0FBSyxDQUFhO0lBRWxCLFNBQVMsQ0FBeUI7SUFDbEMsUUFBUSxDQUFvQjtJQUU1QixZQUFZLFVBQXVCO1FBQ2xDLElBQUksQ0FBQyxLQUFLLEdBQUcsVUFBVSxDQUFBO1FBRXZCLElBQUksQ0FBQyxTQUFTLEdBQUcsRUFBRSxDQUFBO1FBQ25CLElBQUksQ0FBQyxRQUFRLEdBQUcsU0FBUyxDQUFBO0lBQzFCLENBQUM7SUFFRCxPQUFPO1FBQ04sT0FBTyxJQUFJLENBQUMsUUFBUSxDQUFBO0lBQ3JCLENBQUM7SUFFRCxJQUFJLENBQUMsU0FBaUIsRUFBRSxPQUFnQjtRQUN2QyxJQUFJLENBQUMsU0FBUyxDQUFDLFNBQVMsQ0FBQyxHQUFHLE9BQU8sQ0FBQTtJQUNwQyxDQUFDO0lBRUQsTUFBTSxDQUFDLFNBQWlCO1FBQ3ZCLE9BQU8sSUFBSSxDQUFDLFNBQVMsQ0FBQyxTQUFTLENBQUMsQ0FBQTtJQUNqQyxDQUFDO0lBRUQsTUFBTSxDQUFDLFNBQWlCO1FBQ3ZCLElBQUksSUFBSSxDQUFDLFFBQVEsS0FBSyxTQUFTLEVBQUUsQ0FBQztZQUNqQyxPQUFPLElBQUksQ0FBQTtRQUNaLENBQUM7UUFFRCxNQUFNLE9BQU8sR0FBRyxJQUFJLENBQUMsU0FBUyxDQUFDLFNBQVMsQ0FBQyxDQUFBO1FBQ3pDLElBQUksT0FBTyxLQUFLLFNBQVMsRUFBRSxDQUFDO1lBQzNCLE9BQU8sS0FBSyxDQUFBO1FBQ2IsQ0FBQztRQUVELE9BQU8sQ0FBQyxLQUFLLEVBQUUsQ0FBQyxJQUFJLENBQUMsQ0FBQyxPQUFPLEVBQUUsRUFBRTtZQUNoQyxJQUFJLENBQUMsS0FBSyxDQUFDLGVBQWUsQ0FBQyxPQUFPLENBQUMsQ0FBQTtRQUNwQyxDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsR0FBRyxFQUFFO1lBQ2IsSUFBSSxDQUFDLEtBQUssQ0FBQyxlQUFlLENBQUMsc0JBQXNCLENBQUMsQ0FBQTtZQUNsRCxJQUFJLENBQUMsTUFBTSxDQUFDLFNBQVMsQ0FBQyxDQUFBO1FBQ3ZCLENBQUMsQ0FBQyxDQUFBO1FBRUYsT0FBTyxJQUFJLENBQUE7SUFDWixDQUFDO0NBQ0Q7QUFFRCxNQUFNLE9BQU8sT0FBTztJQUNYLE1BQU0sQ0FBQyxNQUFNLEdBQUcsSUFBSSxTQUFTLEVBQUUsQ0FBQTtJQUV2QyxJQUFJLENBQVE7SUFDWixRQUFRLENBQXlCO0lBRWpDLFlBQVksR0FBVztRQUN0QixJQUFJLENBQUMsUUFBUSxHQUFHLFNBQVMsQ0FBQTtRQUN6QixJQUFJLENBQUMsSUFBSSxHQUFHLEdBQUcsQ0FBQTtJQUNoQixDQUFDO0lBRUQsS0FBSyxDQUFDLEtBQUs7UUFDVixJQUFJLElBQUksQ0FBQyxRQUFRLEtBQUssU0FBUyxFQUFFLENBQUM7WUFDakMsT0FBTyxJQUFJLENBQUMsUUFBUSxDQUFBO1FBQ3JCLENBQUM7UUFFRCxNQUFNLE1BQU0sR0FBRyxNQUFNLEtBQUssQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7UUFDckMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxFQUFFLEVBQUUsQ0FBQztZQUNoQixNQUFNLElBQUksS0FBSyxDQUFDLE1BQU0sQ0FBQyxVQUFVLENBQUMsQ0FBQTtRQUNuQyxDQUFDO1FBRUQsTUFBTSxJQUFJLEdBQUcsTUFBTSxNQUFNLENBQUMsSUFBSSxFQUFFLENBQUE7UUFDaEMsTUFBTSxRQUFRLEdBQUcsT0FBTyxDQUFDLE1BQU0sQ0FBQyxlQUFlLENBQUMsSUFBSSxFQUFFLFdBQVcsQ0FBQyxDQUFBO1FBRWxFLE1BQU0sT0FBTyxHQUFHLFFBQVEsQ0FBQyxjQUFjLENBQUMsU0FBUyxDQUFDLENBQUE7UUFDbEQsSUFBSSxPQUFPLEtBQUssSUFBSSxFQUFFLENBQUM7WUFDdEIsTUFBTSxJQUFJLEtBQUssQ0FBQyxZQUFZLENBQUMsQ0FBQTtRQUM5QixDQUFDO1FBRUQsSUFBSSxDQUFDLFFBQVEsR0FBRyxPQUFPLENBQUE7UUFDdkIsT0FBTyxPQUFPLENBQUE7SUFDZixDQUFDIn0=