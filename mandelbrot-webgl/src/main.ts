import "./style.css";

import vertShaderSrc from "./shaders/vert.glsl?raw";
import fragShaderSrc from "./shaders/frag.glsl?raw";

const canvas = document.querySelector("#view")! as HTMLCanvasElement;
const gl = canvas.getContext("webgl2")!;

const vertShader = gl.createShader(gl.VERTEX_SHADER)!;
gl.shaderSource(vertShader, vertShaderSrc);
gl.compileShader(vertShader);
console.log(gl.getShaderInfoLog(vertShader));

const fragShader = gl.createShader(gl.FRAGMENT_SHADER)!;
gl.shaderSource(fragShader, fragShaderSrc);
gl.compileShader(fragShader);
console.log(gl.getShaderInfoLog(fragShader));

const glProgram = gl.createProgram()!;
gl.attachShader(glProgram, vertShader);
gl.attachShader(glProgram, fragShader);
gl.linkProgram(glProgram);
gl.useProgram(glProgram);
console.log(gl.getProgramInfoLog(glProgram));

const vertexBuf = gl.createBuffer();
gl.bindBuffer(gl.ARRAY_BUFFER, vertexBuf);
gl.bufferData(
  gl.ARRAY_BUFFER,
  new Float32Array([-1, -1, 3, -1, -1, 3]),
  gl.STATIC_DRAW
);

const positionAttrib = gl.getAttribLocation(glProgram, "a_Position");
gl.enableVertexAttribArray(positionAttrib);
gl.bindBuffer(gl.ARRAY_BUFFER, vertexBuf);
gl.vertexAttribPointer(positionAttrib, 2, gl.FLOAT, false, 0, 0);

const zoomCenterUniform = gl.getUniformLocation(glProgram, "u_zoomCenter");
const zoomSizeUniform = gl.getUniformLocation(glProgram, "u_zoomSize");
const maxIterationsUniform = gl.getUniformLocation(
  glProgram,
  "u_maxIterations"
);
const resolutionUniform = gl.getUniformLocation(glProgram, "u_resolution");

const maxIters = 512;
const zoomCenter = [0.0, 0.0];
let zoomSize = 4.0;

gl.viewport(0, 0, canvas.width, canvas.height);

function render() {
  gl.uniform2f(zoomCenterUniform, zoomCenter[0], zoomCenter[1]);
  gl.uniform1f(zoomSizeUniform, zoomSize);
  gl.uniform1i(maxIterationsUniform, maxIters);
  gl.uniform2f(resolutionUniform, canvas.width, canvas.height);

  gl.clearColor(0, 0, 0, 1);
  gl.clear(gl.COLOR_BUFFER_BIT);
  gl.drawArrays(gl.TRIANGLES, 0, 3);

  // requestAnimationFrame(render);
}

canvas.addEventListener("wheel", (e) => {
  const scaleDelta = e.deltaY * 0.001;
  if (scaleDelta < 0) {
    zoomSize *= 0.99 + scaleDelta;
  } else {
    zoomSize *= 1.01 + scaleDelta;
  }

  requestAnimationFrame(render);
});

let shouldPan = false;
const panStart = [0, 0];
const oldZoomCenter = [zoomCenter[0], zoomCenter[1]];

canvas.addEventListener("mousedown", (e) => {
  shouldPan = true;
  panStart[0] = e.clientX;
  panStart[1] = e.clientY;
  oldZoomCenter[0] = zoomCenter[0];
  oldZoomCenter[1] = zoomCenter[1];
});

canvas.addEventListener("mouseup", (_) => {
  shouldPan = false;
});

canvas.addEventListener("mousemove", (e) => {
  if (shouldPan) {
    const x = panStart[0] - e.clientX;
    const y = e.clientY - panStart[1];

    zoomCenter[0] = oldZoomCenter[0] + x * zoomSize * 0.00116;
    zoomCenter[1] = oldZoomCenter[1] + y * zoomSize * 0.00116;

    requestAnimationFrame(render);
  }
});

requestAnimationFrame(render);
