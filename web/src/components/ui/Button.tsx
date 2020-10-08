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
  min-width: ${(props) => props.width || "218px"};
  min-height: ${(props) => props.minHeight || "2em"};
  font-size: ${(props) => props.fontSize || "1em"};
  background: ${(props) => (props.primary ? props.theme.blue : "white")};
  border-radius: 3px;
  border: 1px solid ${(props) => (props.primary ? "#0073e6" : "#d2d2d2")};
  color: ${(props) => (props.primary ? "white" : props.theme.blue)};
  margin: ${(props) => props.margin || "0 0 0 0"};
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

const ArrowButton:React.FunctionComponent<ArrowButtonProps> = ({ label, onClick, type, disabled, focus }) => (
<StyledButton autoFocus={focus} onClick={onClick} primary direction="row-reverse" align="flex-end" type={type} disabled={disabled} >
    {label} <span style={{float:"right"}}>&#8594;</span>
</StyledButton>
)

export default StyledButton;

export { BackButton, ArrowButton };
