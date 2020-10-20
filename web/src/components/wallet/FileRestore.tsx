import React, { Component } from "react";
import { Box } from "../ui/Box";
import { BackButton } from "../ui/Button";
import { Card } from "../ui/Card";
import { Container } from "../ui/Container";
import { Dropzone, DropzoneError } from "../ui/Dropzone";
import { SecureFileIcon } from "../ui/Images";
import { LoadingSpinner } from "../ui/LoadingSpinner";
import { H3, Text } from "../ui/Text";
import { FilePathInfo } from "../../shared/FilePathInfo";
import { WalletRestoreFilePassword } from "./RestoreFilePassword";
import { RequestConfig } from "../../api/Config";
import axios from "axios";

export interface WalletFileRestoreProps {
  onComplete: () => void;
  onCancel: () => void;
}

export interface WalletFileRestoreState {
  stateError?: string;
  error?: DropzoneError;
  fileContents?: string | ArrayBuffer;
  mnemonic?: string;
  loading: boolean;
}

export class WalletFileRestore extends Component<
  WalletFileRestoreProps,
  WalletFileRestoreState
> {
  constructor(props: WalletFileRestoreProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
    this.filesSelectedHandler = this.filesSelectedHandler.bind(this);
    this.loadSecureFileData = this.loadSecureFileData.bind(this);
    this.onMnemonic = this.onMnemonic.bind(this);
    this.waitForRescan = this.waitForRescan.bind(this);
    this.state = {
      loading: false
    };
  }

  componentDidMount(): void {
    this.setState({ fileContents: undefined });
  }

  componentDidUnmount(): void {}

  private loadSecureFileData = (file: FilePathInfo, reader?: FileReader) => {
    if (reader) {
      if (reader.result) {
        this.setState({
          fileContents: reader.result,
          error: undefined,
          stateError: undefined
        });
      }
    }
  };

  private onMnemonic = async (wordlist: string) => {
    console.log("onMnemonic");
    const wordCount = wordlist.split(" ").length;
    if (
      wordCount === 12 ||
      wordCount === 13 ||
      wordCount === 24 ||
      wordCount === 25
    ) {
      this.setState({ loading: true });
    } else {
      var err: DropzoneError = {
        title: "Incorrect Word Count",
        message: "Count needs to be 12, 13, 24 or 25 words"
      };
      this.setState({ error: err });
      return;
    }
  };

  private waitForRescan(): undefined {
    //TODO: Use redux-saga for these types of actions
    console.log("waitForRescan");
    var self = this;
    while (this.state.loading) {
      console.log("waitForRescan loop before sleep");
      setTimeout(() => {
        console.log("waitForRescan loop after sleep");
        axios
          .get("/wallet/defaultaddress", RequestConfig)
          .then(function (response) {
            console.log("waitForRescan response", response);
            self.setState({ loading: false });
            self.props.onComplete();
          })
          .catch(function (error) {
            console.log(
              "waitForRescan execute wallet/defaultaddress [Get] Error: " +
                error
            );
          });
      }, 5000);
    }
    return;
  }

  private filesSelectedHandler = (files: FilePathInfo[]) => {
    var dropError: DropzoneError;
    if (files.length !== 1) {
      dropError = {
        title: "More that one file selected",
        message: "Please select only one file"
      };
      this.setState({ error: dropError });
      return;
    }
    const file: FilePathInfo = files[0];
    if (file.size > 131072) { //128KiB
      dropError = {
        title: "File is too large",
        message: "Please select a mnemonic recovery file"
      };
      this.setState({ error: dropError, stateError: dropError.message });
      return;
    }
    if (file.fileReader) {
      file.fileReader.onload = () => {
        this.loadSecureFileData(file, file.fileReader);
      };
    }
  };

  render() {
    return (
      <>
        {this.state && !this.state.fileContents && !this.state.loading && (
          <Container height="50vh" margin="10% 0 0 0">
            <Box direction="column" align="center" width="100%">
              <Box
                direction="column"
                width="700px"
                align="start"
                margin="0 auto 0 auto"
              >
                <BackButton
                  onClick={() => this.props.onCancel()}
                  margin="150px 0 0 -80px"
                />

                <Card
                  width="100%"
                  align="center"
                  minHeight="225px"
                  padding="2em 4em 2em 2em"
                >
                  <Box display="flex" direction="row" margin="0">
                    <SecureFileIcon width="60px" height="60px" />
                    <Box
                      direction="column"
                      width="500px"
                      align="center"
                      margin="0 auto 0 auto"
                    >
                      <Box width="60px" margin="0">
                        <H3 margin="0 0 1em 0">
                          Restore using Secure Restore File
                        </H3>
                      </Box>
                      <Dropzone
                        multiple={false}
                        accept={".psh.json"}
                        filesSelected={this.filesSelectedHandler}
                        error={this.state && this.state.error}
                      ></Dropzone>
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
            </Box>
          </Container>
        )}
        {this.state && this.state.fileContents && !this.state.loading && (
          <WalletRestoreFilePassword
            cancelPassword={() => this.setState({ fileContents: undefined })}
            onMnemonic={(wordlist) => this.onMnemonic(wordlist)}
            fileContents={this.state.fileContents}
          />
        )}
        {this.state && this.state.loading && (
          <>
            <LoadingSpinner
              active={typeof this.state.loading !== "undefined"}
              label="Restoring your wallet using the imported mnemonic file"
              size={50}
              opaque
            />
          </>
        )}
      </>
    );
  }
}
