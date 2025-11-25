var jsxmm;
(function (jsxmm) {
    function Element(tag, properties = {}, ...children) {
        const element = document.createElement(tag);
        replace(element, properties);
        for (let i = 0; i < children.length; i++) {
            const child = children[i];
            if (typeof child === "object" && Object.hasOwn(child, "HTML") && typeof child.HTML === "function") {
                element.append(child.HTML());
            }
            else {
                element.append(child);
            }
        }
        return element;
    }
    jsxmm.Element = Element;
    function Style(element, style) {
        replace(element.style, style);
    }
    jsxmm.Style = Style;
    function replace(base, replacement) {
        for (const key in replacement) {
            if (!(key in base)) {
                console.error(`${key} not present in ${base} element`);
            }
            if (typeof replacement[key] === "object") {
                replace(base[key], replacement[key]);
            }
            else {
                base[key] = replacement[key];
            }
        }
    }
})(jsxmm || (jsxmm = {}));
export default jsxmm;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZWxlbWVudC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9tb2R1bGUvanN4bW0vZWxlbWVudC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQSxJQUFVLEtBQUssQ0F1Q2Q7QUF2Q0QsV0FBVSxLQUFLO0lBR2QsU0FBZ0IsT0FBTyxDQUF3QyxHQUFNLEVBQUUsYUFBNEIsRUFBRSxFQUFFLEdBQUcsUUFBdUM7UUFDaEosTUFBTSxPQUFPLEdBQUcsUUFBUSxDQUFDLGFBQWEsQ0FBQyxHQUFHLENBQUMsQ0FBQTtRQUMzQyxPQUFPLENBQUMsT0FBTyxFQUFFLFVBQVUsQ0FBQyxDQUFBO1FBRTVCLEtBQUssSUFBSSxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsR0FBRyxRQUFRLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFLENBQUM7WUFDMUMsTUFBTSxLQUFLLEdBQUcsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ3pCLElBQUksT0FBTyxLQUFLLEtBQUssUUFBUSxJQUFJLE1BQU0sQ0FBQyxNQUFNLENBQUMsS0FBSyxFQUFFLE1BQU0sQ0FBQyxJQUFJLE9BQVEsS0FBYSxDQUFDLElBQUksS0FBSyxVQUFVLEVBQUUsQ0FBQztnQkFDNUcsT0FBTyxDQUFDLE1BQU0sQ0FBRSxLQUFtQixDQUFDLElBQUksRUFBRSxDQUFDLENBQUE7WUFDNUMsQ0FBQztpQkFBTSxDQUFDO2dCQUNQLE9BQU8sQ0FBQyxNQUFNLENBQUMsS0FBc0IsQ0FBQyxDQUFBO1lBQ3ZDLENBQUM7UUFDRixDQUFDO1FBRUQsT0FBTyxPQUFPLENBQUE7SUFDZixDQUFDO0lBZGUsYUFBTyxVQWN0QixDQUFBO0lBRUQsU0FBZ0IsS0FBSyxDQUFDLE9BQW9CLEVBQUUsS0FBbUM7UUFDOUUsT0FBTyxDQUFDLE9BQU8sQ0FBQyxLQUFLLEVBQUUsS0FBSyxDQUFDLENBQUE7SUFDOUIsQ0FBQztJQUZlLFdBQUssUUFFcEIsQ0FBQTtJQUVELFNBQVMsT0FBTyxDQUFDLElBQThCLEVBQUUsV0FBcUM7UUFDckYsS0FBSyxNQUFNLEdBQUcsSUFBSSxXQUFXLEVBQUUsQ0FBQztZQUMvQixJQUFJLENBQUMsQ0FBQyxHQUFHLElBQUksSUFBSSxDQUFDLEVBQUUsQ0FBQztnQkFDcEIsT0FBTyxDQUFDLEtBQUssQ0FBQyxHQUFHLEdBQUcsbUJBQW1CLElBQUksVUFBVSxDQUFDLENBQUE7WUFFdkQsQ0FBQztZQUVELElBQUksT0FBTyxXQUFXLENBQUMsR0FBRyxDQUFDLEtBQUssUUFBUSxFQUFFLENBQUM7Z0JBQzFDLE9BQU8sQ0FBQyxJQUFJLENBQUMsR0FBRyxDQUFDLEVBQUUsV0FBVyxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUE7WUFDckMsQ0FBQztpQkFBTSxDQUFDO2dCQUNQLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxXQUFXLENBQUMsR0FBRyxDQUFDLENBQUE7WUFDN0IsQ0FBQztRQUNGLENBQUM7SUFDRixDQUFDO0FBR0YsQ0FBQyxFQXZDUyxLQUFLLEtBQUwsS0FBSyxRQXVDZDtBQUVELGVBQWUsS0FBSyxDQUFBIn0=