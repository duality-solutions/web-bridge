import * as React from "react";
import styled from "styled-components";
import backarrow from "../../assets/backarrow.svg";
import importwallet from "../../assets/import.svg";
import logo from "../../assets/logo.svg";
import passphrase from "../../assets/passphrase.svg";
import restoreWallet from "../../assets/restore.svg";
import securefile from "../../assets/securefile.svg";


interface ImageProps {
  src?: string;
  width?: string;
  height?: string;
  margin?: string;
  white?: boolean;
  onClick?: (e: React.MouseEvent) => void;
  cursor?: string;
  float?: string;
}

const SvgIcon = styled("img")<ImageProps>`
  src: ${(props) => (props.src ? props.src : "")};
  width: ${(props) => (props.width ? props.width : "100%")};
  height: ${(props) => (props.height ? props.height : "100px")};
  margin: ${(props) => (props.margin ? props.margin : "0")};
  vertical-align: middle;
  background: ${(props) => (props.white ? "#737373" : "")};
  cursor: ${(props) => (props.cursor ? props.cursor : "")};
  float: ${(props) => (props.float ? props.float : "")};
`;

const PlainAppLogo = () => <SvgIcon src={logo} height="70px" />;

interface Styleable {
  style?: Record<string, any>;
}

const ImportIcon: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={importwallet}
    width={width || "50px"}
    height={height || "30px"}
    margin={margin || "0 10px"}
  />
);

const RestoreIcon: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={restoreWallet}
    width={width || "50px"}
    height={height || "30px"}
    margin={margin || "0 10px"}
  />
);

const PassphraseIcon: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={passphrase}
    width={width || "50px"}
    height={height || "30px"}
    margin={margin || "0 10px"}
  />
);

const SecureFileIcon: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={securefile}
    width={width || "50px"}
    height={height || "30px"}
    margin={margin || "0 10px"}
  />
);

const BackArrow: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={backarrow}
    width={width || "50px"}
    height={height || "30px"}
    margin={margin || "0 10px"}
  />
);

export {
  SvgIcon as AppLogo,
  PlainAppLogo,
  ImportIcon,
  RestoreIcon,
  PassphraseIcon,
  SecureFileIcon,
  BackArrow
};
