import React, { useRef, useEffect, MutableRefObject } from "react";
import { FunctionComponent } from "react";
import { FilePathInfo } from "../../shared/FilePathInfo";
import { Text } from "./Text";
import { Card } from "./Card";
import Button from "./Button";

export interface DropzoneError {
  title: string;
  message: string;
}

interface DropzoneDispatchProps {
  filesSelected: (files: FilePathInfo[]) => void;
  directoriesSelected?: (directories: FilePathInfo[]) => void;
}

interface DropzoneStateProps {
  error?: DropzoneError;
  multiple?: boolean;
  accept?: string;
}

type DropzoneProps = DropzoneDispatchProps & DropzoneStateProps;

export const Dropzone: FunctionComponent<DropzoneProps> = ({
  error,
  filesSelected,
  directoriesSelected,
  multiple,
  accept
}) => {
  const dirFileInputRef: MutableRefObject<HTMLInputElement | null> = useRef(
    null
  );
  useEffect(() => {
    const elem = dirFileInputRef.current;
    if (elem) {
      elem.setAttribute("webkitdirectory", ""); //annoyingly, we don't seem to be able to set this as an attr with JSX :(
    }
  });
  return (
    <Card
      background="white"
      border={error ? "dashed 2px #ea4964" : "dashed 2px #b0b0b0"}
      minHeight="266px"
      padding="75px 0"
      onDragOver={(e) => {
        e.preventDefault();
        e.dataTransfer.dropEffect = "move";
      }}
      onDrop={(e) => {
        e.preventDefault();
        const files = [...e.dataTransfer.files];
        filesSelected(
          files.map((f) => ({ path: f.name, type: f.type, size: f.size }))
        );
      }}
    >
      {error ? (
        <>
          <Text
            fontSize="18px"
            fontWeight="bold"
            color="#ea4964"
            margin="0"
            align="center"
          >
            {error.title}
          </Text>
          <Text align="center" margin="20px 0" fontSize="0.8em">
            {error.message}
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
        multiple={multiple}
        accept={accept || "*/*"}
        onChange={(e) => {
          e.preventDefault();
          if (!e.currentTarget.files) {
            return;
          }
          const files = [...e.currentTarget.files];
          filesSelected(
            files.map((f) => ({ path: f.name, type: f.type, size: f.size }))
          );
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
      {directoriesSelected && (
        <>
          <input
            ref={dirFileInputRef}
            type="file"
            id="dirElem"
            multiple={multiple}
            accept={accept || "*/*"}
            onChange={(e) => {
              e.preventDefault();
              if (!e.currentTarget.files) {
                return;
              }
              const files = [...e.currentTarget.files];
              directoriesSelected(
                files.map((f) => ({ path: f.name, type: "", size: 0 }))
              );
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
  );
};
