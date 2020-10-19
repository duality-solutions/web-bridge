import React, { Component } from "react";
import { Box } from "../ui/Box";
import { BackButton } from "../ui/Button";
import { SCard } from "../ui/Card";
import { Container } from "../ui/Container";
import { PassphraseIconWhite, SecureFileIconWhite } from "../ui/Images";
import { Text } from "../ui/Text";

export interface WalletRestoreProps {
  restoreUsingMnemonic: () => void;
  restoreWithSecureFile: () => void;
  cancelRestore: () => void;
}

export interface WalletRestoreState {
  useMnemonic?: boolean;
}

export class WalletRestore extends Component<
  WalletRestoreProps,
  WalletRestoreState
> {
  private newWallet?: boolean;
  constructor(props: WalletRestoreProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
  }

  componentDidMount(): void {
    this.setState({ useMnemonic: undefined });
  }

  componentDidUnmount(): void {}

  render() {
    return (
      <>
        <Container height="50vh">
          <Box direction="column" align="center" width="100%">
            <Box display="flex" direction="row" align="center" width="100%">
              <BackButton
                onClick={() => this.props.cancelRestore()}
                margin="70px 0 0 -350px"
              />
              <SCard
                onClick={() => this.props.restoreUsingMnemonic()}
                padding="2em 1em 1em 1em"
                height="140px"
                width="220px"
              >
                <PassphraseIconWhite
                  height="60px"
                  width="60px"
                  style={{ margin: "0 0 1em 0" }}
                />
                <Text
                  align="center"
                  color="white"
                  fontSize="1em"
                  fontWeight="bold"
                >
                  Restore using mnemonic
                </Text>
              </SCard>
              <SCard
                onClick={() => this.props.restoreWithSecureFile()}
                padding="2em 1em 1em 1em"
                height="140px"
                width="220px"
              >
                <SecureFileIconWhite
                  height="60px"
                  width="60px"
                  style={{ margin: " 0 0 1em 0" }}
                />
                <Text
                  align="center"
                  color="white"
                  fontSize="1em"
                  fontWeight="bold"
                >
                  Restore using secure file
                </Text>
              </SCard>
            </Box>
          </Box>
        </Container>
      </>
    );
  }
}
