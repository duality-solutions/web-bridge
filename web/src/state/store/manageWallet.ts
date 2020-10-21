import { connect } from "react-redux";
import { MapPropsToDispatchObj } from "../shared/MapPropsToDispatchObj";
import { ManageWalletActions } from "../actions/manageWallet";
import { WalletSetup, WalletSetupProps, WalletViewDispatch } from "../../components/wallet/Setup";
import { RootState } from '../reducers';

const mapStateToProps = (state: RootState): WalletSetupProps => {
    return {
        complete: true,
        onCancel: null,
        onComplete: null,
    };
};

const mapDispatchToProps: MapPropsToDispatchObj<WalletViewDispatch> = { ...ManageWalletActions };

export default connect(mapStateToProps, mapDispatchToProps)(WalletSetup);
