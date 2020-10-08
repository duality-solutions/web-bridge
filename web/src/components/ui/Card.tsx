import styled from "styled-components";

interface CardProps {
  width?: string;
  align?: string;
  padding?: string;
  minHeight?: string;
  minWidth?: string;
  border?: string;
  background?: string;
}

const StyledCard = styled("div")<CardProps>`
  align-items: ${(props) => props.align || "center"};
  min-width: 2em;
  width: ${(props) => props.width || "500px"};
  min-height: ${(props) => props.minHeight || "5em"};
  min-width: ${(props) => props.minWidth || "200px"};
  box-shadow: 0 0 14px 0 rgba(0, 0, 0, 0.05);
  border: ${(props) => props.border || "solid 1px #f1f1f1"};
  border-radius: 4px;
  background-color: ${(props) => props.background || "#ffffff"};
  padding: ${(props) => props.padding || "2em 4em"};
  margin: 1em 0;
  box-sizing: border-box;
`;

const SquareCard = styled("div")<{
  padding?: string;
  height?: string;
  width?: string;
}>`
  display: flex;
  flex-direction: column;
  align-items: center;
  color: white;
  background: #2e77d0;
  box-shadow: 0 0 14px 0 rgba(0, 0, 0, 0.05);
  border: solid 2px #f7f6f6;
  border-radius: 6px;
  // min-height: 10em;
  width: ${(props) => props.width || "200px"};
  height: ${(props) => props.height || "200px"};
  padding: ${(props) => props.padding || "0.25em 1em"};
  margin: 0 0.5em 0 0.5em;
  cursor: pointer;
`;

export default StyledCard;

export { StyledCard as Card, SquareCard as SCard };
