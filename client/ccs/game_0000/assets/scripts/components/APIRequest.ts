export interface ResGuestToken {
    Code: string;
    Content: { token:string };
}

export async function ReqGuestToken(addr: string): Promise<string> {
    try {
        const response = await fetch(addr);

        if (!response.ok) {
            // 直接拋出錯誤，async 函數會自動將其包裝為 rejected Promise
            throw new Error(`HTTP error! Status: ${response.status} - ${response.statusText}`);
        }

        const data: ResGuestToken = await response.json();
        if(data.Code != "OK" ){
            throw new Error(`Response error! Code: ${data.Code}`);
        }
        // 直接 return 資料，async 函數會自動將其包裝為 resolved Promise
        return data.Content.token ; 

    } catch (error: any) {
        // 在這裡可以處理錯誤日誌，然後重新拋出，或者根據需求回傳一個特定的錯誤值
        console.error(`Error in GetGuestToken for address ${addr}:`, error.message);
        // 重新拋出錯誤，讓調用者去處理 rejected Promise
        throw error;
    }
}