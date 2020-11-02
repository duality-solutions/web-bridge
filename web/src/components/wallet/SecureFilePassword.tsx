import React, { Component } from "react";
import { Box } from "../ui/Box";
import { ArrowButton, BackArrowButton } from "../ui/Button";
import { Card } from "../ui/Card";
import { Container } from "../ui/Container";
import { Input } from "../ui/Input";
import { SecureFileIcon } from "../ui/Images";
import { H3, Text } from "../ui/Text";
import { getAesEncryptor } from "../../shared/AesEncryption";
import { ValidationResult } from "../../shared/ValidationResult";
import { SaveFile } from "../../shared/SaveFile";

export interface WalletSecureFilePasswordProps {
  mnemonic: string;
  onCancel: () => void;
  validationResult?: ValidationResult<string>;
}

export interface WalletSecureFilePasswordState {
  password: string;
  confirmPassword: string;
  error?: string;
  fileCipherText?: string;
}

export class WalletSecureFilePassword extends Component<
    WalletSecureFilePasswordProps,
    WalletSecureFilePasswordState
> {
  constructor(props: WalletSecureFilePasswordProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
    this.encryptSecureMnemonicFile = this.encryptSecureMnemonicFile.bind(this);
    // init state
    this.state = { password: "", confirmPassword: "" };
  }

  componentDidMount(): void {}

  componentDidUnmount(): void {}

  private encryptSecureMnemonicFile = async () => {
    if (this.state.password !== this.state.confirmPassword) {
      this.setState({ error: "Passwords do not match" });
      return;
    } else if (!/.{6,}/.test(this.state.password)) {
      this.setState({ error: "Password must be > 6 characters" });
      return;
    }
    if (this.state.password) {
      try {
        const { encrypt } = getAesEncryptor(this.state.password);
        const encryptedMnemonic = encrypt(this.props.mnemonic);
        SaveFile("webbridge-backup.psh.json", encryptedMnemonic);
      } catch (e) {
        console.log("Error:", e);
      }
    }
  }

  render() {
    const { onCancel } = this.props;
    const { error } = this.state;
    
    return (
      <>
        <Container height="50vh" margin="10% 5% 0 0">
          <Box direction="column" align="center" width="100%">
            <Box
              direction="column"
              width="800px"
              align="start"
              margin="0 auto 0 auto"
            >
              <div style={{ display: "flex" }}>
                <BackArrowButton onClick={() => onCancel()} />
                <Card
                  width="100%"
                  align="center"
                  minHeight="225px"
                  padding="2em 4em 2em 2em"
                >
                  <Box display="flex" direction="row" margin="0">
                    <Box width="120px" margin="0">
                      <SecureFileIcon width="60px" height="60px" />
                    </Box>
                    <Box margin="0 0 0 2em">
                      <H3>Secure file</H3>
                      <Text fontSize="14px">
                        Create a secure file password{" "}
                      </Text>
                      <Input
                        value={this.state.password}
                        name="password"
                        onChange={(e: any) =>
                          this.setState({ password: e.target.value })
                        }
                        placeholder="Password"
                        type="password"
                        margin="1em 0 1em 0"
                        padding="0 1em 0 1em"
                        autoFocus={true}
                        error={error ? true : false}
                      />
                      <Text fontSize="14px">Confirm Password</Text>
                      <Input
                        value={this.state.confirmPassword}
                        name="confirmPassword"
                        onChange={(e: any) =>
                          this.setState({ confirmPassword: e.target.value })
                        }
                        placeholder="Password"
                        type="password"
                        margin="1em 0 1em 0"
                        padding="0 1em 0 1em"
                        error={error ? true : false}
                      />
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
              </div>
            </Box>
            <Box
              direction="column"
              width="800px"
              align="right"
              margin="0 auto 0 auto"
            >
              <ArrowButton
                onClick={() => this.encryptSecureMnemonicFile()}
                type="button"
                label="Save Backup"
              />
            </Box>
          </Box>
        </Container>
      </>
    );
  }
}
