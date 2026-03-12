export namespace main {
	
	export class Config {
	    videoDir: string;
	    port: number;
	    allowedExts: string[];
	    isAllowDel: boolean;
	    isDeep: boolean;
	    imageExts: string[];
	    isAllowDelImg: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.videoDir = source["videoDir"];
	        this.port = source["port"];
	        this.allowedExts = source["allowedExts"];
	        this.isAllowDel = source["isAllowDel"];
	        this.isDeep = source["isDeep"];
	        this.imageExts = source["imageExts"];
	        this.isAllowDelImg = source["isAllowDelImg"];
	    }
	}
	export class DeleteImageResult {
	    success: boolean;
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new DeleteImageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.msg = source["msg"];
	    }
	}
	export class DeleteResult {
	    success: boolean;
	    msg: string;
	    currentVideo: string;
	    hasVideo: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DeleteResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.msg = source["msg"];
	        this.currentVideo = source["currentVideo"];
	        this.hasVideo = source["hasVideo"];
	    }
	}
	export class ImageInfo {
	    pic: string;
	    picW: number;
	    picH: number;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pic = source["pic"];
	        this.picW = source["picW"];
	        this.picH = source["picH"];
	        this.path = source["path"];
	    }
	}
	export class ImageListResult {
	    success: boolean;
	    imageList: ImageInfo[];
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageListResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.imageList = this.convertValues(source["imageList"], ImageInfo);
	        this.msg = source["msg"];
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
	export class VideoInfo {
	    videoList: string[];
	    currentVideo: string;
	    hasVideo: boolean;
	    port: number;
	
	    static createFrom(source: any = {}) {
	        return new VideoInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.videoList = source["videoList"];
	        this.currentVideo = source["currentVideo"];
	        this.hasVideo = source["hasVideo"];
	        this.port = source["port"];
	    }
	}

}

