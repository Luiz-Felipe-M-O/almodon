const OriginToRoot = {
    "http://localhost:3000": {
        client: "http://localhost:3000/dist/",
        server: "",
    },
    "http://localhost:4545": {
        client: "http://localhost:4545/",
        server: "http://localhost:4545/api/v1/",
    },
    "https://alan-b-lima.github.io": {
        client: "https://alan-b-lima.github.io/almodon/ui/web/dist/",
        server: "",
    },
};
function urls() {
    const origin = location.origin;
    if (!Object.hasOwn(OriginToRoot, origin)) {
        throw new Error("Unknown location " + origin);
    }
    const root = OriginToRoot[origin];
    return root;
}
var Source;
(function (Source) {
    function From(path, origin = Source.client) {
        return new URL(path, origin).href;
    }
    Source.From = From;
    const source = urls();
    Source.client = source.client;
    Source.server = source.server;
})(Source || (Source = {}));
export default Source;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic291cmNlLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2ludGVybmFsL3N1cHBvcnQvc291cmNlLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBLE1BQU0sWUFBWSxHQUF1RDtJQUNyRSx1QkFBdUIsRUFBRTtRQUNyQixNQUFNLEVBQUUsNkJBQTZCO1FBQ3JDLE1BQU0sRUFBRSxFQUFFO0tBQ2I7SUFDRCx1QkFBdUIsRUFBRTtRQUNyQixNQUFNLEVBQUUsd0JBQXdCO1FBQ2hDLE1BQU0sRUFBRSwrQkFBK0I7S0FDMUM7SUFDRCwrQkFBK0IsRUFBRTtRQUM3QixNQUFNLEVBQUUsb0RBQW9EO1FBQzVELE1BQU0sRUFBRSxFQUFFO0tBQ2I7Q0FDSixDQUFBO0FBRUQsU0FBUyxJQUFJO0lBQ1QsTUFBTSxNQUFNLEdBQUcsUUFBUSxDQUFDLE1BQU0sQ0FBQTtJQUM5QixJQUFJLENBQUMsTUFBTSxDQUFDLE1BQU0sQ0FBQyxZQUFZLEVBQUUsTUFBTSxDQUFDLEVBQUUsQ0FBQztRQUN2QyxNQUFNLElBQUksS0FBSyxDQUFDLG1CQUFtQixHQUFHLE1BQU0sQ0FBQyxDQUFBO0lBQ2pELENBQUM7SUFFRCxNQUFNLElBQUksR0FBRyxZQUFZLENBQUMsTUFBTSxDQUFDLENBQUE7SUFDakMsT0FBTyxJQUFJLENBQUE7QUFDZixDQUFDO0FBRUQsSUFBVSxNQUFNLENBU2Y7QUFURCxXQUFVLE1BQU07SUFDWixTQUFnQixJQUFJLENBQUMsSUFBWSxFQUFFLFNBQWlCLE9BQUEsTUFBTTtRQUN0RCxPQUFPLElBQUksR0FBRyxDQUFDLElBQUksRUFBRSxNQUFNLENBQUMsQ0FBQyxJQUFJLENBQUE7SUFDckMsQ0FBQztJQUZlLFdBQUksT0FFbkIsQ0FBQTtJQUVELE1BQU0sTUFBTSxHQUFHLElBQUksRUFBRSxDQUFBO0lBRVIsYUFBTSxHQUFHLE1BQU0sQ0FBQyxNQUFNLENBQUE7SUFDdEIsYUFBTSxHQUFHLE1BQU0sQ0FBQyxNQUFNLENBQUE7QUFDdkMsQ0FBQyxFQVRTLE1BQU0sS0FBTixNQUFNLFFBU2Y7QUFFRCxlQUFlLE1BQU0sQ0FBQSJ9