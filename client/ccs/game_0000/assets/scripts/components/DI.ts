export class DIKey {
    static WSConn : string = "websocket.conn" ;
    static LOGGER : string = "logger" ;
}
export class DI {
    private static store:{[key: string]:any} = {} ;
    public static set(key:string, value:any) {
        if(this.store[key] != undefined ){
            throw new Error(`Key: ${key} already exists.`) ;
        }
        this.store[key] = value ;
    }
    public static must_get<T>(key:string){
        if(this.store[key] == undefined ){
            throw new Error(`Key: ${key} has an empty value.`) ;
        }
        return <T>this.store[key] ;
    }
}
