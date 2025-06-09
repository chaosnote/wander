import { ISignal } from "kernel" ;
import kernel from "kernel";

/**
 * 遊戲層與共用層事件
 */

/**
 * 事件:底層載入設定完成
 */
export const BaseInitComplete: ISignal = kernel.genSignal();
/**
 * 事件:前端發送遊戲封包
 */
export const SendPacket: ISignal = kernel.genSignal();