import * as React from "react";
import styled from "styled-components";
import backarrow from "../../assets/backarrow.svg";
import importWhite from "../../assets/importWhite.svg";
import logo from "../../assets/logo.svg";
import passphrase from "../../assets/passphrase.svg";
import passphraseWhite from "../../assets/passphraseWhite.svg";
import passwordentry from "../../assets/passwordentry.svg";
import restoreWallet from "../../assets/restoreWhite.svg";
import safe from "../../assets/safe.svg";
import securefile from "../../assets/securefile.svg";
import securefileWhite from "../../assets/securefileWhite.svg";

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

const ImportIconWhite: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={importWhite}
    width={width || "50px"}
    height={height || "30px"}
    margin={margin || "0 10px"}
  />
);

const RestoreIconWhite: React.FunctionComponent<ImageProps & Styleable> = ({
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

const PassphraseIconWhite: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={passphraseWhite}
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

const SecureFileIconWhite: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={securefileWhite}
    width={width || "50px"}
    height={height || "30px"}
    margin={margin || "0 10px"}
  />
);

const PasswordEntry: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={passwordentry}
    width={width || "50px"}
    height={height || "30px"}
    margin={margin || "0 10px"}
  />
);

const SafeImage: React.FunctionComponent<ImageProps & Styleable> = ({
  width,
  height,
  margin,
  onClick,
  style
}) => (
  <SvgIcon
    onClick={onClick}
    style={{ cursor: "pointer", ...style }}
    src={safe}
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
  BackArrow,
  ImportIconWhite,
  PassphraseIcon,
  PassphraseIconWhite,
  PasswordEntry,
  PlainAppLogo,
  RestoreIconWhite,
  SafeImage,
  SecureFileIcon,
  SecureFileIconWhite
};
