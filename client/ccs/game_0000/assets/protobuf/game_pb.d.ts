// DO NOT EDIT! This is a generated file. Edit the JSDoc in src/*.js instead and run 'npm run build:types'.

/** Action enum. */
export enum Action {
    INIT = 0,
    BET = 1,
    COMPLETE = 2
}

/** Represents a Player. */
export class Player implements IPlayer {

    /**
     * Constructs a new Player.
     * @param [properties] Properties to set
     */
    constructor(properties?: IPlayer);

    /** Player name. */
    public name: string;

    /** Player wallet. */
    public wallet: number;

    /**
     * Creates a new Player instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Player instance
     */
    public static create(properties?: IPlayer): Player;

    /**
     * Encodes the specified Player message. Does not implicitly {@link Player.verify|verify} messages.
     * @param message Player message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IPlayer, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Player message, length delimited. Does not implicitly {@link Player.verify|verify} messages.
     * @param message Player message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IPlayer, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Player message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Player
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Player;

    /**
     * Decodes a Player message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Player
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Player;

    /**
     * Verifies a Player message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Player message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Player
     */
    public static fromObject(object: { [k: string]: any }): Player;

    /**
     * Creates a plain object from a Player message. Also converts values to other types if specified.
     * @param message Player
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Player, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Player to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };

    /**
     * Gets the default type url for Player
     * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns The default type url
     */
    public static getTypeUrl(typeUrlPrefix?: string): string;
}

/** Represents an Init. */
export class Init implements IInit {

    /**
     * Constructs a new Init.
     * @param [properties] Properties to set
     */
    constructor(properties?: IInit);

    /** Init player. */
    public player?: (IPlayer|null);

    /**
     * Creates a new Init instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Init instance
     */
    public static create(properties?: IInit): Init;

    /**
     * Encodes the specified Init message. Does not implicitly {@link Init.verify|verify} messages.
     * @param message Init message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IInit, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Init message, length delimited. Does not implicitly {@link Init.verify|verify} messages.
     * @param message Init message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IInit, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an Init message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Init
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Init;

    /**
     * Decodes an Init message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Init
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Init;

    /**
     * Verifies an Init message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an Init message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Init
     */
    public static fromObject(object: { [k: string]: any }): Init;

    /**
     * Creates a plain object from an Init message. Also converts values to other types if specified.
     * @param message Init
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Init, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Init to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };

    /**
     * Gets the default type url for Init
     * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns The default type url
     */
    public static getTypeUrl(typeUrlPrefix?: string): string;
}

/** Represents a GameMessage. */
export class GameMessage implements IGameMessage {

    /**
     * Constructs a new GameMessage.
     * @param [properties] Properties to set
     */
    constructor(properties?: IGameMessage);

    /** GameMessage error. */
    public error?: (string|null);

    /** GameMessage type. */
    public type: GameMessage.MessageType;

    /** GameMessage action. */
    public action: string;

    /** GameMessage payload. */
    public payload: Uint8Array;

    /** GameMessage timestamp. */
    public timestamp: (number|Long);

    /** GameMessage _error. */
    public _error?: "error";

    /**
     * Creates a new GameMessage instance using the specified properties.
     * @param [properties] Properties to set
     * @returns GameMessage instance
     */
    public static create(properties?: IGameMessage): GameMessage;

    /**
     * Encodes the specified GameMessage message. Does not implicitly {@link GameMessage.verify|verify} messages.
     * @param message GameMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IGameMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified GameMessage message, length delimited. Does not implicitly {@link GameMessage.verify|verify} messages.
     * @param message GameMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IGameMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a GameMessage message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns GameMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): GameMessage;

    /**
     * Decodes a GameMessage message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns GameMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): GameMessage;

    /**
     * Verifies a GameMessage message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a GameMessage message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns GameMessage
     */
    public static fromObject(object: { [k: string]: any }): GameMessage;

    /**
     * Creates a plain object from a GameMessage message. Also converts values to other types if specified.
     * @param message GameMessage
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: GameMessage, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this GameMessage to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };

    /**
     * Gets the default type url for GameMessage
     * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns The default type url
     */
    public static getTypeUrl(typeUrlPrefix?: string): string;
}

export namespace GameMessage {

    /** MessageType enum. */
    enum MessageType {
        REQUEST = 0,
        RESPONSE = 1,
        NOTIFY = 2,
        ALERT = 3
    }
}
