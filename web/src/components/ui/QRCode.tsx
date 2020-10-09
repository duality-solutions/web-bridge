import * as React from "react";
import * as ReactDOM from "react-dom";
import qrcode from "qrcode-generator";

export interface QRCodeProps {
  value?: string;
  ecLevel?: "L" | "M" | "Q" | "H";
  enableCORS?: boolean;
  size?: number;
  minimumCellSize?: number;
  minPadding?: number;
  bgColor?: string;
  fgColor?: string;
  qrStyle?: "squares" | "dots";
  style?: object;
}

function utf16to8(str: string): string {
  let out: string = "",
    i: number,
    c: number;
  const len: number = str.length;
  for (i = 0; i < len; i++) {
    c = str.charCodeAt(i);
    if (c >= 0x0001 && c <= 0x007f) {
      out += str.charAt(i);
    } else if (c > 0x07ff) {
      out += String.fromCharCode(0xe0 | ((c >> 12) & 0x0f));
      out += String.fromCharCode(0x80 | ((c >> 6) & 0x3f));
      out += String.fromCharCode(0x80 | ((c >> 0) & 0x3f));
    } else {
      out += String.fromCharCode(0xc0 | ((c >> 6) & 0x1f));
      out += String.fromCharCode(0x80 | ((c >> 0) & 0x3f));
    }
  }
  return out;
}

function drawPositioningPattern(
  row: number,
  col: number,
  length: number,
  props: QRCodeProps,
  ctx: CanvasRenderingContext2D
) {
  const cellSize = Math.max(
    Math.trunc(props.size! / length),
    props.minimumCellSize!
  );

  for (let r = -1; r <= 7; r++) {
    if (!(row + r <= -1 || length <= row + r)) {
      for (let c = -1; c <= 7; c++) {
        if (
          (!(col + c <= -1 || length <= col + c) &&
            0 <= r &&
            r <= 6 &&
            (c === 0 || c === 6)) ||
          (0 <= c && c <= 6 && (r === 0 || r === 6)) ||
          (2 <= r && r <= 4 && 2 <= c && c <= 4)
        ) {
          const w =
            Math.ceil((row + r + 1) * cellSize) -
            Math.floor((row + r) * cellSize);
          const h =
            Math.ceil((col + c + 1) * cellSize) -
            Math.floor((col + c) * cellSize);

          ctx.fillStyle = props.fgColor!;
          ctx.fillRect(
            Math.round((row + r) * cellSize),
            Math.round((col + c) * cellSize),
            w,
            h
          );
        }
      }
    }
  }
}

const defaultProps: QRCodeProps = {
  value: "",
  ecLevel: "M",
  size: 150,
  minPadding: 10,
  bgColor: "#FFFFFF",
  fgColor: "#000000",
  qrStyle: "squares",
  minimumCellSize: 1
};

export const QRCode: React.FunctionComponent<QRCodeProps> = (props) => {
  const { value, ecLevel, size, bgColor, fgColor, qrStyle, minimumCellSize } = {
    ...defaultProps,
    ...props
  };
  const canvasRef = React.useRef(null);

  const qrCode = qrcode(0, ecLevel!);
  qrCode.addData(utf16to8(value!));
  qrCode.make();
  const length = qrCode.getModuleCount();
  const cellSize = Math.max(Math.trunc(size! / length), minimumCellSize!);
  const actualSize = cellSize * length;

  React.useEffect(() => {
    const canvas: HTMLCanvasElement = ReactDOM.findDOMNode(
      canvasRef.current!
    ) as HTMLCanvasElement;
    const ctx: CanvasRenderingContext2D = canvas.getContext("2d")!;

    const scale = window.devicePixelRatio || 1;
    canvas.height = canvas.width = actualSize * scale;
    ctx.scale(scale, scale);
    ctx.fillStyle = bgColor!;
    ctx.fillRect(0, 0, actualSize, actualSize);
    const s = qrStyle === "dots" ? 1 : 0;
    for (let row = 0; row < length; row++) {
      for (let col = 0; col < length; col++) {
        if (qrCode.isDark(row, col)) {
          ctx.fillStyle = fgColor!;
          const w =
            Math.ceil((col + 1) * cellSize) - Math.floor(col * cellSize) - s;
          const h =
            Math.ceil((row + 1) * cellSize) - Math.floor(row * cellSize) - s;
          ctx.fillRect(
            Math.round(col * cellSize),
            Math.round(row * cellSize),
            w,
            h
          );
        }
      }
    }
    drawPositioningPattern(0, 0, length, props, ctx);
    drawPositioningPattern(length - 7, 0, length, props, ctx);
    drawPositioningPattern(0, length - 7, length, props, ctx);
  }, [
    actualSize,
    cellSize,
    length,
    props,
    qrCode,
    value,
    ecLevel,
    size,
    bgColor,
    fgColor,
    qrStyle
  ]);
  return (
    <div
      style={{
        display: "inline-block",
        ...props.style,
        background: props.bgColor
      }}
    >
      <canvas
        height={actualSize}
        width={actualSize}
        style={{
          height: actualSize + "px",
          width: actualSize + "px",
          padding:
            props.minPadding! + Math.max(0, (size! - actualSize) / 2) + "px",
          background: props.bgColor
        }}
        ref={canvasRef}
      />
    </div>
  );
};
