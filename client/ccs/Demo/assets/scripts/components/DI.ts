export class DIKey {
    static WSConn: string = "websocket.conn";
    static WSObserver: string = "websocket.observer";
    static LOGGER: string = "logger";
}
class di_config {
    is_only: boolean;
    instance: any;
    create_func: CreateFunc;
}
export type CreateFunc = (...args: any) => any;
export class DI {
    private static store: { [key: string]: di_config } = {};
    public static set(key: string, create_func: CreateFunc) {
        if (this.store[key] != undefined) {
            throw new Error(`Key: ${key} already exists.`);
        }
        let config = new di_config();
        config.is_only = false;
        config.create_func = create_func;
        this.store[key] = config;
    }
    public static set_share(key: string, create_func: CreateFunc) {
        if (this.store[key] != undefined) {
            throw new Error(`Key: ${key} already exists.`);
        }
        let config = new di_config();
        config.is_only = true;
        config.create_func = create_func;
        this.store[key] = config;
    }
    public static must_get<T>(key: string, ...args): T {
        if (this.store[key] == undefined) {
            throw new Error(`Key: ${key} has an empty value.`);
        }
        let config = this.store[key];
        if (config.is_only) {
            if (config.instance == undefined) {
                config.instance = config.create_func(...args) ;
            }
            return config.instance ;
        }
        return <T>config.create_func(...args) ;
    }
}
