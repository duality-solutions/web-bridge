import React, { Component } from "react";
import { Box } from "../ui/Box";
import { ArrowButton, BackButton } from "../ui/Button";
import { Card } from "../ui/Card";
import { Container } from "../ui/Container";
import { Input } from "../ui/Input";
import { SecureFileIcon } from "../ui/Images";
import { H3, Text } from "../ui/Text";
import { getAesEncryptor } from "../../shared/AesEncryption";

export interface WalletRestoreFilePasswordProps {
  fileContents: string | ArrayBuffer;
  cancelPassword: () => void;
  onMnemonic: (mnemonic: string) => void;
}

export interface WalletRestoreFilePasswordState {
  passwordText?: string;
  error?: string;
  mnemonic?: string;
}

export class WalletRestoreFilePassword extends Component<
  WalletRestoreFilePasswordProps,
  WalletRestoreFilePasswordState
> {
  constructor(props: WalletRestoreFilePasswordProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
    this.decryptSecureMnemonicFile = this.decryptSecureMnemonicFile.bind(this);
  }

  componentDidMount(): void {}

  componentDidUnmount(): void {}

  private decryptSecureMnemonicFile(): void {
    let mnemonicFilePassphrase = this.state.passwordText;
    if (mnemonicFilePassphrase) {
      const { decrypt } = getAesEncryptor(mnemonicFilePassphrase);
      try {
        const mnemonic = decrypt(this.props.fileContents.toString());
        this.setState( { mnemonic: mnemonic } );
        this.props.onMnemonic(mnemonic);
      }
      catch(e) {
        console.log('Error:', e);
        this.setState( { error: "Incorrect password" } );
      }
    }
  }

  render() {
    return (
      <>
        <H3 align="center" colored fontWeight="600">
          Restore Account
        </H3>
        <Container height="50vh" margin="10% 0 0 0">
          <Box direction="column" align="center" width="100%">
            <Box
              direction="column"
              width="700px"
              align="start"
              margin="0 auto 0 auto"
            >
              <BackButton
                onClick={() => this.props.cancelPassword()}
                margin="130px 0 0 -100px"
              />
              <Card
                width="100%"
                align="center"
                minHeight="300px"
                padding="2em 4em 2em 2em"
              >
                <Box display="flex" direction="row" margin="0">
                  <Box width="60px" margin="0">
                    <SecureFileIcon width="60px" height="60px" />
                  </Box>
                  <Box margin="1em 0 0 2em">
                    <H3 margin="0 0 1em 0">
                      Restore using Secure Restore File{" "}
                    </H3>
                    <Text fontSize="0.8em">
                      Enter your secure file password
                    </Text>
                    <Input
                      value={
                        this.state && this.state.passwordText
                          ? this.state.passwordText
                          : ""
                      }
                      name="password"
                      onChange={(e) =>
                        this.setState({ passwordText: e.target.value })
                      }
                      placeholder="Password"
                      type="password"
                      margin="1em 0 1em 0"
                      padding="0 1em 0 1em"
                      autoFocus={true}
                    />
                    <Text fontSize="0.8em" margin="0">
                      This is the password used to create the secure file
                    </Text>
                    {this.state && this.state.error ? (
                      <Text align="center" color="#e30429">
                        {this.state.error}
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
              margin="0 auto 0 auto"
            >
              <ArrowButton
                label="Continue"
                type="button"
                onClick={() => this.decryptSecureMnemonicFile()}
              />
            </Box>
          </Box>
        </Container>
      </>
    );
  }
}
