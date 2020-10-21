import * as React from "react";
import styled from "styled-components";
import { BackArrow } from "./Images";

interface ButtonProps {
  align?: string;
  primary?: boolean;
  theme?: { blue: string };
  direction?: string;
  width?: string;
  minHeight?: string;
  fontSize?: string;
  margin?: string;
  type?: string;
  background?: string;
  color?: string;
}

const StyledButton = styled("button")<ButtonProps>`
  align-self: ${(props) => (props.align ? props.align : "center")};
  justify-content: center;
  min-width: ${(props) => (props.width ? props.width : "218px")};
  min-height: ${(props) => (props.minHeight ? props.minHeight : "2em")};
  font-size: ${(props) => (props.fontSize ? props.fontSize : "1em")};
  background: ${(props) => (props.primary ? "#0073e6" : "white")};
  border-radius: 3px;
  border: 1px solid ${(props) => (props.primary ? "#0073e6" : "#d2d2d2")};
  color: ${(props) => (props.primary ? "white" : "#0073e6")};
  margin: ${(props) => (props.margin ? props.margin : "0 0 0 0")};
  padding: 0.5em 1em;
  cursor: pointer;
`;

const BackButton: React.FunctionComponent<{
  margin?: string;
  onClick?: () => void;
}> = ({ margin, onClick }) => (
  <BackArrow
    onClick={onClick}
    style={{
      float: "left",
      margin,
      position: "fixed",
      cursor: "pointer"
    }}
  />
);

interface ArrowButtonProps {
  label: string;
  onClick?: () => void;
  type: "button" | "submit" | "reset" | undefined;
  disabled?: boolean;
  focus?: boolean;
}

const ArrowButton: React.FunctionComponent<ArrowButtonProps> = ({
  label,
  onClick,
  type,
  disabled,
  focus
}) => (
  <StyledButton
    autoFocus={focus}
    onClick={onClick}
    direction="row-reverse"
    align="flex-end"
    type={type}
    disabled={disabled}
  >
    {label} <span style={{ float: "right" }}>&#8594;</span>
  </StyledButton>
);

const StyledBackArrowButton = styled("div")`
  display: block;
  direction: column;
  width: 100px;
  cursor: pointer;
`;

const BackArrowButton: React.FunctionComponent<{
  onClick: () => void;
  marginTop?: string;
}> = ({ onClick, marginTop }) => (
  <StyledBackArrowButton onClick={() => onClick()}>
    <BackArrow style={{ marginTop: marginTop || "175%" }} />
  </StyledBackArrowButton>
);

const LightButton = styled("button")<ButtonProps>`
  align-self: ${(props) => (props.align ? props.align : "center")};
  justify-content: center;
  min-width: 218px;
  min-height: 2em;
  font-size: 1em;
  background: white;
  border-radius: 3px;
  border: 2px solid ${(props) => props.theme.blue};
  color: ${(props) => props.theme.blue};
  /* margin: 0.5em 1em; */
  margin: 0 0 0 0;
  padding: 0.5em 1em;
  cursor: pointer;
`;

export default StyledButton;

export { ArrowButton, BackArrowButton, BackButton, LightButton };
