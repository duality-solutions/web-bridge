import React, { Component } from "react";
import { Box } from "../ui/Box";
import { ArrowButton, BackButton } from "../ui/Button";
import { Card } from "../ui/Card";
import { Container } from "../ui/Container";
import { SecureFileIcon } from "../ui/Images";
import { MnemonicInput } from "../ui/Input";
import { H1, H3, Text } from "../ui/Text";

export interface WalletMnemonicRestoreProps {
  onComplete: () => void;
  onCancel: () => void;
  error?: string;
}

export interface WalletMnemonicRestoreState {
  stateError?: string;
  mnemonic?: string;
}

export class WalletMnemonicRestore extends Component<
  WalletMnemonicRestoreProps,
  WalletMnemonicRestoreState
> {
  private newWallet?: boolean;
  constructor(props: WalletMnemonicRestoreProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
    this.setMnemonic = this.setMnemonic.bind(this);
  }

  componentDidMount(): void {}

  componentDidUnmount(): void {}

  private setMnemonic(words: string): void {
    this.setState({ mnemonic: words });
  }

  render() {
    return (
      <>
        <H1 align="center" colored fontWeight="600">
          Restore Account
        </H1>
        <Container height="50vh" margin="10% 0 0 0">
          <form /*onSubmit={handleSubmit}*/>
            <Box direction="column" align="center" width="100%">
              <Box
                direction="column"
                width="700px"
                align="start"
                margin="0 auto 0 auto"
              >
                <BackButton
                  onClick={() => this.props.onCancel()}
                  margin="150px 0 0 -100px"
                />
                <Card
                  width="100%"
                  align="center"
                  minHeight="225px"
                  padding="2em 4em 2em 2em"
                >
                  <Box display="flex" direction="row" margin="0">
                    <Box width="60px" margin="0">
                      <SecureFileIcon width="60px" height="60px" />
                    </Box>
                    <Box margin="1em 0 0 2em">
                      <H3 margin="0 0 1em 0">
                        Restore using mnemonic passphrase{" "}
                      </H3>
                      <MnemonicInput
                        placeholder="Enter mnemonic passphrase"
                        name="mnemonic-text"
                        onChange={(e) => this.setMnemonic(e.target.value)}
                        value={
                          this.state && this.state.mnemonic
                            ? this.state.mnemonic
                            : ""
                        }
                      />
                      {this.state && this.state.stateError ? (
                        <Text align="center" color="#e30429">
                          {this.state.stateError}
                        </Text>
                      ) : (
                        <></>
                      )}
                    </Box>
                  </Box>
                </Card>
              </Box>
              <Box
                direction="column"
                width="700px"
                align="right"
                margin="0 auto 10px auto"
              >
                <ArrowButton
                  label="Continue"
                  type="button"
                  onClick={() => this.props.onComplete()}
                />
              </Box>
            </Box>
          </form>
        </Container>
      </>
    );
  }
}
