import { _decorator, Component, log } from 'cc';
import { DEBUG, DEV } from 'cc/env';
//-----------------------------------------------
import { ILogger, IWebSocketConn } from "kernel";
import kernel from "kernel";
//-----------------------------------------------
import { ReqGuestToken } from './APIRequest';
import { BaseInitComplete, SendGamePacket } from './Events';
import { DI, DIKey } from './DI';
//-----------------------------------------------
import pb from "../../protobuf/game_pb.js"; // 留意 tsconfig.json 與 *.js
//-----------------------------------------------
const { ccclass, property } = _decorator;
//-----------------------------------------------
export class BaseConfig {
    WSAddr: string
}
//
// 底層主體附加至畫布(extends and attach to canvas)
//
@ccclass('Base')
export class Base extends Component {
    @property({ tooltip: "後端請求網址" })
    backend_addr: string = "localhost:8080";

    @property({ tooltip: "遊戲伺服器連結網址" })
    websocket_addr: string = "localhost:8081";

    @property({ tooltip: "game_id(len:4)" })
    game_id: string = "0000";

    @property({ tooltip: "token" })
    token: string = "";

    // @property({ tooltip: "是否使用訪客身份" })
    // is_guest: boolean = true;

    @kernel.logger(DEBUG)
    private logger: ILogger;

    async onLoad() {
        // 載入靜態資訊( 文件、圖片 ... )、觸發 loading bar
        //
        //
        SendGamePacket.on(this.onSendGamePack.bind(this));
    }

    async start() {
        const msg = "start";

        if (this.token == "") {
            this.token = await ReqGuestToken(`http://${this.backend_addr}/guest/new`);
        }
        this.logger.debug(msg, `token:\n ${this.token}`);
        //
        let config = new BaseConfig();
        config.WSAddr = `http://${this.websocket_addr}/ws/00/${this.game_id}?token=${this.token}`;
        //
        const ws = kernel.genWebSocketConn(DEBUG);
        ws.open.on(() => { this.logger.debug(msg, "open"); });
        ws.close.on((reson: string) => this.logger.debug(msg, reson));
        ws.error.on(() => this.logger.error(msg, "error"));
        ws.message_text.on((pack: string) => this.logger.debug(msg, pack));
        ws.message_binary.on((pack: Uint8Array) => {

            let message = pb.GameMessage.decode(pack);
            if (message.type == pb.GameMessage.MessageType.RESPONSE) {
                this.logger.debug(msg+".onReceiveGameMessage", `action:\t ${message.action}`);

                let list = DI.must_get<{[key:string]:(payload:Uint8Array)=>void}>(DIKey.WSObserver) ;
                let func = list[message.action] ;
                if(func == undefined){
                    throw new Error(`unset action:\t ${message.action}`) ;
                }
                func(message.payload) ;
                return;
            }
            // 處理非遊戲封包
            //
        });
        DI.set_share(DIKey.WSConn, ()=>ws);
        //
        BaseInitComplete.emit(config);
    }

    private onSendGamePack(action: string, pack: Uint8Array) {
        const msg = "onSendGamePack";
        this.logger.debug(msg, `action: ${action}`);
        let output = pb.GameMessage.encode({
            action: action,
            payload: pack,
            timestamp: new Date().getTime(),
            type: pb.GameMessage.MessageType.REQUEST
        }) ;

        const ws = DI.must_get<IWebSocketConn>(DIKey.WSConn);
        ws.send(output.finish()) ;
    }
}
