import axios from 'axios';
import { RequestConfig }from "./Config";

interface WalletSetupStatusResponse {
    MnemonicBackup: boolean;
    HasAccounts: boolean;
    HasLinks: boolean;
    HasTransactions: boolean;
    WalletEncrypted: boolean;
    UnlockedUntil: number;
};

export type WalletSetupStatus = WalletSetupStatusResponse;

export const GetWalletSetupStatus = async (): Promise<WalletSetupStatus> => {
    return await axios.get<WalletSetupStatus>("/wallet/setup", RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        var errMessage = "GetWalletSetupStatus execute [Get] /wallet/setup error: " + error;
        console.log(errMessage);
        var errResponse: WalletSetupStatus = {
            MnemonicBackup: false,
            HasAccounts: false,
            HasLinks: false,
            HasTransactions: false,
            WalletEncrypted: false,
            UnlockedUntil: 0,
        }
        return errResponse;
    });
}
