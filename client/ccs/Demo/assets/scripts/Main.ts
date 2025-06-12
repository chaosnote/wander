import { _decorator, Component, Node } from 'cc';
import { DEBUG } from 'cc/env';
//-----------------------------------------------
import { ILogger, IWebSocketConn } from "kernel";
import kernel from "kernel";
//-----------------------------------------------
import { BaseInitComplete } from './components/Events';
//-----------------------------------------------
const { ccclass, property } = _decorator;

@ccclass('Main')
export class Main extends Component {

    @kernel.logger(DEBUG)
    private logger: ILogger;

    protected onLoad(): void {
        BaseInitComplete.on(this.onBaseInitComplete.bind(this));
    }
    private onBaseInitComplete(): void {
        const msg = "onBaseInitComplete" ;
        this.logger.debug(msg, `
            註冊網路封包行為
        `);
    }
}
