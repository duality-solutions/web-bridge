import styled from "styled-components";

interface BoxProps {
  align?: string;
  margin?: string;
  width?: string;
  direction?: string;
  display?: string;
  minWidth?: string;
  height?: string;
  background?: string;
  borderRadius?: string;
  padding?: string;
  alignContents?: string;
  minHeight?: string;
  border?: string;
  justifyContent?: string;
}

const StyledBox = styled("div")<BoxProps>`
  display: ${(props) => props.display || "block"};
  justify-content: ${(props) => props.justifyContent || "center"};
  background: ${(props) => (props.background ? props.background : "")};
  text-align: ${(props) => props.align || "start"};
  width: ${(props) => (props.width ? props.width : "500px")};
  margin: ${(props) => (props.margin ? props.margin : "0 0 0 0")};
  box-sizing: border-box;
  min-width: ${(props) => props.minWidth || "10px"};
  height: ${(props) => (props.height ? props.height : "10px")};
  border-radius: ${(props) => (props.borderRadius ? props.borderRadius : "")};
  padding: ${(props) => (props.padding ? props.padding : "")};
  align-items: ${(props) => (props.alignContents ? props.alignContents : "")};
  min-height: ${(props) => (props.minHeight ? props.minHeight : "")};
  border: ${(props) => (props.border ? props.border : "")};
`;

export default StyledBox;

export { StyledBox as Box };
