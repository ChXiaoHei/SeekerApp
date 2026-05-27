export namespace config {
	
	export class Config {
	    hotkey: string;
	    maxHistory: number;
	    dllPath: string;
	    theme: string;
	    windowWidth: number;
	    windowHeight: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hotkey = source["hotkey"];
	        this.maxHistory = source["maxHistory"];
	        this.dllPath = source["dllPath"];
	        this.theme = source["theme"];
	        this.windowWidth = source["windowWidth"];
	        this.windowHeight = source["windowHeight"];
	    }
	}

}

export namespace everything {
	
	export class SearchOptions {
	    matchCase: boolean;
	    matchWholeWord: boolean;
	    matchPath: boolean;
	    useRegex: boolean;
	    maxResults: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matchCase = source["matchCase"];
	        this.matchWholeWord = source["matchWholeWord"];
	        this.matchPath = source["matchPath"];
	        this.useRegex = source["useRegex"];
	        this.maxResults = source["maxResults"];
	    }
	}
	export class SearchResult {
	    fileName: string;
	    fullPath: string;
	    isFolder: boolean;
	    size: number;
	    // Go type: time
	    dateModified: any;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileName = source["fileName"];
	        this.fullPath = source["fullPath"];
	        this.isFolder = source["isFolder"];
	        this.size = source["size"];
	        this.dateModified = this.convertValues(source["dateModified"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

