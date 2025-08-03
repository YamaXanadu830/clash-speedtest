export namespace main {
	
	export class ConfigInfo {
	    proxyCount: number;
	    configPath: string;
	    filter: string;
	    block: string;
	
	    static createFrom(source: any = {}) {
	        return new ConfigInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.proxyCount = source["proxyCount"];
	        this.configPath = source["configPath"];
	        this.filter = source["filter"];
	        this.block = source["block"];
	    }
	}
	export class MonitorConfig {
	    duration: number;
	    interval: number;
	    targetURL: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new MonitorConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.duration = source["duration"];
	        this.interval = source["interval"];
	        this.targetURL = source["targetURL"];
	        this.type = source["type"];
	    }
	}
	export class SystemInfo {
	    platform: string;
	    version: string;
	    buildTime: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.platform = source["platform"];
	        this.version = source["version"];
	        this.buildTime = source["buildTime"];
	    }
	}
	export class TestConfig {
	    configPath: string;
	    filterRegex: string;
	    blockRegex: string;
	    serverURL: string;
	    downloadSize: number;
	    uploadSize: number;
	    timeout: number;
	    concurrent: number;
	    maxLatency: number;
	    minDownloadSpeed: number;
	    minUploadSpeed: number;
	    fastMode: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TestConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.configPath = source["configPath"];
	        this.filterRegex = source["filterRegex"];
	        this.blockRegex = source["blockRegex"];
	        this.serverURL = source["serverURL"];
	        this.downloadSize = source["downloadSize"];
	        this.uploadSize = source["uploadSize"];
	        this.timeout = source["timeout"];
	        this.concurrent = source["concurrent"];
	        this.maxLatency = source["maxLatency"];
	        this.minDownloadSpeed = source["minDownloadSpeed"];
	        this.minUploadSpeed = source["minUploadSpeed"];
	        this.fastMode = source["fastMode"];
	    }
	}

}

