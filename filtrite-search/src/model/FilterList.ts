export class FilterListData {
    public display_name: string;
    public urls: Array<URLInfo>;
    public filter_file_url: string;
    public stars: number | null;
    public repo_owner: string;
    public repo_name: string;
    public list_url: string;

    public static listFromJSON(obj: unknown): Array<FilterListData> {
        if (typeof obj === "string") {
            obj = JSON.parse(obj);
        }
        return (obj as Array<FilterListData>);
    }
}

export class URLInfo {
    public url: string;
    public title: string;
}

