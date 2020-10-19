import React, { Component } from "react";
import { Container } from "../ui/Container";
import { SCard } from "../ui/Card";
import { Box } from "../ui/Box";
import { H1, H3, Text } from "../ui/Text";
import { ImportIconWhite, RestoreIconWhite } from "../ui/Images";
import { MnemonicBackup } from "./MnemonicBackup";
import { MnemonicWarning } from "./MnemonicWarning";
import { WalletRestore } from "./Restore";
import { WalletFileRestore } from "./FileRestore";
import { WalletMnemonicRestore } from "./MnemonicRestore";
import { WalletPassword } from "./WalletPassword";

enum SetupState {
  Init = 1,
  New,
  NewWarned,
  Restore,
  RestoreWithMnemonic,
  RestoreWithSecureFile,
  CreatePassword
}

export interface WalletSetupProps {
  onComplete: () => void;
}

export interface WalletSetupState {
  setupState: SetupState;
}

export class WalletSetup extends Component<WalletSetupProps, WalletSetupState> {
  constructor(props: WalletSetupProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
    this.onNewWallet = this.onNewWallet.bind(this);
  }

  componentDidMount(): void {
    this.setState({
      setupState: SetupState.Init
    });
  }

  componentDidUnmount(): void {}

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
        {this.state && this.state.setupState === SetupState.Init && (
          <Container height="50vh">
            <p></p>
            <Box direction="column" align="center" width="100%">
              <Box display="flex" direction="row" align="center" width="100%">
                <SCard
                  onClick={() => this.setState({ setupState: SetupState.New })}
                >
                  <ImportIconWhite height="80px" width="80px" />
                  <H3 align="start" color="white" minwidth="50px">
                    Create Wallet
                  </H3>
                  <Text color="white" align="center">
                    Create a new wallet
                  </Text>
                </SCard>
                <SCard
                  onClick={() =>
                    this.setState({ setupState: SetupState.Restore })
                  }
                >
                  <RestoreIconWhite height="80px" width="80px" />
                  <H3 align="start" color="white" minwidth="50px">
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
        {this.state && this.state.setupState === SetupState.New && (
          <MnemonicWarning
            onComplete={() =>
              this.setState({ setupState: SetupState.NewWarned })
            }
            onCancel={() => this.setState({ setupState: SetupState.Init })}
          />
        )}
        {this.state && this.state.setupState === SetupState.NewWarned && (
          <MnemonicBackup
            onCancel={() => this.setState({ setupState: SetupState.Init })}
            onComplete={() =>
              this.setState({ setupState: SetupState.CreatePassword })
            }
          />
        )}
        {this.state && this.state.setupState === SetupState.CreatePassword && (
          <WalletPassword
            onComplete={() => this.props.onComplete()}
            uiType={"CREATE"}
            onCancel={() => this.setState({ setupState: SetupState.Init })}
          />
        )}
        {this.state && this.state.setupState === SetupState.Restore && (
          <WalletRestore
            cancelRestore={() => this.setState({ setupState: SetupState.Init })}
            restoreUsingMnemonic={() =>
              this.setState({ setupState: SetupState.RestoreWithMnemonic })
            }
            restoreWithSecureFile={() =>
              this.setState({ setupState: SetupState.RestoreWithSecureFile })
            }
          />
        )}
        {this.state &&
          this.state.setupState === SetupState.RestoreWithMnemonic && (
            <WalletMnemonicRestore
              onComplete={() => this.props.onComplete()}
              onCancel={() => this.setState({ setupState: SetupState.Restore })}
            />
          )}
        {this.state &&
          this.state.setupState === SetupState.RestoreWithSecureFile && (
            <WalletFileRestore
              onComplete={() => this.props.onComplete()}
              onCancel={() => this.setState({ setupState: SetupState.Restore })}
            />
          )}
      </>
    );
  }
}
