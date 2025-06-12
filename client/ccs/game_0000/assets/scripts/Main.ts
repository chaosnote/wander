import { _decorator, Component, Node } from 'cc';
import { DEBUG, DEV } from 'cc/env';
//-----------------------------------------------
const { ccclass, property } = _decorator;
//-----------------------------------------------
import { ILogger, ISignal, IWebSocketConn } from 'kernel';
import kernel from 'kernel';
//-----------------------------------------------
import pb from "../protobuf/game_pb.js"; // 留意 tsconfig.json 與 *.js
//-----------------------------------------------
import { BaseInitComplete, SendGamePacket } from './components/Events';
import { DI, DIKey } from './components/DI';
import { BaseConfig } from './components/Base';
//-----------------------------------------------
@ccclass('Main')
export class Main extends Component {

    @kernel.logger(DEBUG)
    private logger: ILogger;

    protected onLoad(): void {
        BaseInitComplete.on(this.onBaseInitComplete.bind(this));
    }

    private onBaseInitComplete(config: BaseConfig, ...args: any[]): void {
        const msg = "onBaseInitComplete";
        this.logger.debug(msg, config);

        // 封包註冊
        DI.set_share(DIKey.WSObserver, () => {
            let observer: { [key: string]: (payload: Uint8Array) => void } = {};
            observer[pb.Action[pb.Action.INIT]] = this.onInit.bind(this);
            observer[pb.Action[pb.Action.BET]] = this.onBetRes.bind(this);
            observer[pb.Action[pb.Action.COMPLETE]] = this.onCompleteRes.bind(this);
            return observer;
        });

        const ws = DI.must_get<IWebSocketConn>(DIKey.WSConn);
        ws.dial(config.WSAddr, "arraybuffer");
    }

    private onInit(payload: Uint8Array) {
        const msg = "onInit";
        let data = pb.Init.decode(payload);
        this.logger.debug(msg, `name:\t ${data.player.name}`, `wallet:\t ${data.player.wallet}`);

        SendGamePacket.emit(pb.Action[pb.Action.BET], null);
    }
    private onBetRes(payload: Uint8Array) {
        const msg = "onBetRes";
        this.logger.debug(msg, `action:\t BET`);

        SendGamePacket.emit(pb.Action[pb.Action.COMPLETE], null);
    }
    private onCompleteRes(payload: Uint8Array) {
        const msg = "onCompleteRes";
        this.logger.debug(msg, `action:\t COMPLETE`);
    }
}
