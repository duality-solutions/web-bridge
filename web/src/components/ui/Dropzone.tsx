import React, { Component } from "react";
import { FilePathInfo } from "../../shared/FilePathInfo";
import { Text } from "./Text";
import { Card } from "./Card";
import Button from "./Button";

export interface DropzoneError {
  title: string;
  message: string;
}

interface DropzoneProps {
  multiple?: boolean;
  accept?: string;
  filesSelected: (files: FilePathInfo[]) => void;
  directoriesSelected?: (directories: FilePathInfo[]) => void;
  error?: DropzoneError;
}

interface DropzoneState {
  fileContents: string | ArrayBuffer | null;
}

export class Dropzone extends Component<DropzoneProps, DropzoneState> {
  private dirFileInputRef: React.RefObject<HTMLInputElement>;
  constructor(props: DropzoneProps) {
    super(props);
    this.dirFileInputRef = React.createRef();
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentWillUnmount = this.componentWillUnmount.bind(this);
    this.loadFilesDataReader = this.loadFilesDataReader.bind(this);
  }

  componentDidMount(): void {
    if (this.dirFileInputRef.current) {
      this.dirFileInputRef.current.setAttribute("webkitdirectory", ""); //annoyingly, we don't seem to be able to set this as an attr with JSX :(
    }
  }

  componentWillUnmount(): void {}

  private loadFilesDataReader = (files: File[]) => {
    let selectedFiles: FilePathInfo[] = [];
    files.forEach((file) => {
      let reader = new FileReader();
      reader.readAsDataURL(file);
      const fileInfo: FilePathInfo = {
        path: file.name,
        type: file.type,
        size: file.size,
        fileReader: reader
      };
      selectedFiles.push(fileInfo);
    });
    this.props.filesSelected(selectedFiles);
  };

  render() {
    return (
      <>
        <Card
          background="white"
          border={
            this.props.error
              ? "dashed 2px #ea4964"
              : "dashed 2px #b0b0b0"
          }
          minHeight="266px"
          padding="75px 0"
          onDragOver={(e) => {
            e.preventDefault();
            e.dataTransfer.dropEffect = "move";
          }}
          onDrop={(e) => {
            e.preventDefault();
            const files = [...e.dataTransfer.files];
            this.loadFilesDataReader(files);
          }}
        >
          {this.props.error ? (
            <>
              <Text
                fontSize="18px"
                fontWeight="bold"
                color="#ea4964"
                margin="0"
                align="center"
              >
                {this.props.error.title}
              </Text>
              <Text align="center" margin="20px 0" fontSize="0.8em">
                {this.props.error.message}
              </Text>
            </>
          ) : (
            <>
              <Text
                fontSize="18px"
                fontWeight="bold"
                color="#9b9b9b"
                margin="0"
                align="center"
              >
                Drag file here
              </Text>
              <Text align="center" margin="20px 0" fontSize="0.8em">
                or
              </Text>
            </>
          )}
          <input
            type="file"
            id="fileElem"
            multiple={this.props.multiple}
            accept={this.props.accept || "*/*"}
            onChange={(e) => {
              e.preventDefault();
              if (!e.currentTarget.files) {
                return;
              }
              const files = [...e.currentTarget.files];
              this.loadFilesDataReader(files);
            }}
            style={{ display: "none" }}
          />
          <Button color="#0055c4" width="175px" type="button">
            <label
              style={{
                width: "100%",
                height: "100%",
                display: "block",
                cursor: "pointer"
              }}
              className="button"
              htmlFor="fileElem"
            >
              Select file
            </label>
          </Button>
          {this.props.directoriesSelected && (
            <>
              <input
                ref={this.dirFileInputRef}
                type="file"
                id="dirElem"
                multiple={this.props.multiple}
                accept={this.props.accept || "*/*"}
                onChange={(e) => {
                  e.preventDefault();
                  if (!e.currentTarget.files) {
                    return;
                  }
                  const files = [...e.currentTarget.files];
                  this.loadFilesDataReader(files);
                }}
                style={{ display: "none" }}
              />
              <Button color="#0055c4" width="175px" type="button">
                <label
                  style={{
                    width: "100%",
                    height: "100%",
                    display: "block",
                    cursor: "pointer"
                  }}
                  className="button"
                  htmlFor="dirElem"
                >
                  Select directory
                </label>
              </Button>
            </>
          )}
        </Card>
      </>
    );
  }
}
