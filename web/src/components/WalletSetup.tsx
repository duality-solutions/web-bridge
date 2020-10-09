import React, { Component } from "react";
import { Container } from "./ui/Container";
import { SCard } from "./ui/Card";
import { Box } from "./ui/Box";
import { H1, H3, Text } from "./ui/Text";
import { ImportIcon, RestoreIcon } from "./ui/Images";
import { WalletRestore } from "./WalletRestore";
import { WalletFileRestore } from "./WalletFileRestore";
import { WalletMnemonicRestore } from "./WalletMnemonicRestore";

export interface WalletSetupProps {
  onComplete: () => void;
}

export interface WalletSetupState {
  newWallet?: boolean;
  restoreUsingMnemonic: boolean;
  restoreWithPassphrase: boolean;
}

export class WalletSetup extends Component<WalletSetupProps, WalletSetupState> {
  private newWallet?: boolean;
  constructor(props: WalletSetupProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentWillUnmount = this.componentWillUnmount.bind(this);
    this.onNewWallet = this.onNewWallet.bind(this);
  }

  componentDidMount(): void {
    this.setState({
      newWallet: undefined,
      restoreUsingMnemonic: false,
      restoreWithPassphrase: false
    });
    console.log(this.newWallet);
  }

  componentWillUnmount(): void {}

  private onNewWallet(): void {
    this.props.onComplete();
  }

  render() {
    return (
      <>
        <H1 align="center" color="black">
          Manage Wallet
        </H1>
        <p></p>
        {this.state && this.state.newWallet === undefined && (
          <Container height="50vh">
            <p></p>
            <Box direction="column" align="center" width="100%">
              <Box display="flex" direction="row" align="center" width="100%">
                <SCard onClick={() => this.onNewWallet()}>
                  <ImportIcon height="80px" width="80px" />
                  <H3 align="center" color="white">
                    Create Wallet
                  </H3>
                  <Text color="white" align="center">
                    Create a new wallet
                  </Text>
                </SCard>
                <SCard onClick={() => this.setState({ newWallet: false })}>
                  <RestoreIcon height="80px" width="80px" />
                  <H3 align="center" color="white">
                    Restore Wallet
                  </H3>
                  <Text color="white" align="center">
                    You have a backed up mnemonic or file you would like to
                    restore from
                  </Text>
                </SCard>
              </Box>
            </Box>
          </Container>
        )}
        {this.state &&
          this.state.newWallet === false &&
          this.state.restoreUsingMnemonic === false &&
          this.state.restoreWithPassphrase === false && (
            <WalletRestore
              cancelRestore={() => this.setState({ newWallet: undefined })}
              restoreUsingMnemonic={() =>
                this.setState({
                  restoreUsingMnemonic: true,
                  restoreWithPassphrase: false
                })
              }
              restoreWithSecureFile={() =>
                this.setState({
                  restoreUsingMnemonic: false,
                  restoreWithPassphrase: true
                })
              }
            />
          )}
        {this.state &&
          this.state.newWallet === false &&
          this.state.restoreUsingMnemonic === true &&
          this.state.restoreWithPassphrase === false && (
            <WalletMnemonicRestore
              onComplete={() => this.props.onComplete()}
              onCancel={() => this.setState({ restoreUsingMnemonic: false })}
            />
          )}
        {this.state &&
          this.state.newWallet === false &&
          this.state.restoreUsingMnemonic === false &&
          this.state.restoreWithPassphrase === true && (
            <WalletFileRestore
              onComplete={() => this.props.onComplete()}
              onCancel={() => this.setState({ restoreWithPassphrase: false })}
            />
          )}
      </>
    );
  }
}
