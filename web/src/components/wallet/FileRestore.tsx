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
import { RequestConfig, RestUrl } from "../../api/Config";
import { ImportMnemonicRequest } from "../../shared/Mnemonic";
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
    this.componentWillUnmount = this.componentWillUnmount.bind(this);
    this.filesSelectedHandler = this.filesSelectedHandler.bind(this);
    this.loadSecureFileData = this.loadSecureFileData.bind(this);
    this.onMnemonic = this.onMnemonic.bind(this);
    this.waitForRescan = this.waitForRescan.bind(this);
    this.state = {
      loading: false
    }
  }

  componentDidMount(): void {
    this.setState({ fileContents: undefined });
  }

  componentWillUnmount(): void {}

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
    var self = this;
    console.log('onMnemonic');
    const wordCount = wordlist.split(" ").length;
    if (
      wordCount === 12 ||
      wordCount === 13 ||
      wordCount === 24 ||
      wordCount === 25
    ) {
      var request: ImportMnemonicRequest = {
        mnemonic: wordlist
      }
      axios.post<ImportMnemonicRequest>(RestUrl + "wallet/mnemonic", request, RequestConfig).then(function (response) {
        console.log(JSON.stringify(response.data, null, 2));
        self.setState( { loading: true, mnemonic: wordlist }, self.waitForRescan());
      }).catch(function (error) {
        console.log("onMnemonic execute wallet/mnemonic [Post] Error: " + error);
      });
    } else {
      var err: DropzoneError = {
        title: "Incorrect Word Count",
        message: "Count needs to be 12, 13, 24 or 25 words"
      };
      this.setState({ error: err });
      return;
    }
  };

  // Returns a Promise that resolves after "ms" Milliseconds
  private timer(ms: number) {
    return new Promise(res => setTimeout(res, ms));
  }

  private waitForRescan(): | undefined {
    //TODO: Use redux-saga for these types of actions
    console.log('waitForRescan');
    var self = this;
    while (this.state.loading) {
      console.log('waitForRescan loop before sleep');
      setTimeout(() => 
      {
        console.log('waitForRescan loop after sleep');
        axios.get(RestUrl + "wallet/defaultaddress", RequestConfig)
          .then(function (response) {
            console.log('waitForRescan response', response);
            self.setState( { loading: false } );
            self.props.onComplete();
        })
        .catch(function (error) {
          console.log("waitForRescan execute wallet/defaultaddress [Get] Error: " + error);
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
                    <Box width="60px" margin="0">
                      <SecureFileIcon width="60px" height="60px" />
                    </Box>
                    <Box
                      direction="column"
                      width="500px"
                      align="center"
                      margin="0 auto 0 auto"
                    >
                      <H3 margin="0 0 1em 0">
                        Restore using Secure Restore File{" "}
                      </H3>
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
