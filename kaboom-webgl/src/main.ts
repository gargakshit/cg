import "../../style.css";

import vertShaderSrc from "./shaders/vert.glsl?raw";
import fragShaderSrc from "./shaders/frag.glsl?raw";

const canvas = document.querySelector("#view")! as HTMLCanvasElement;
const fpsCounter = document.querySelector("#fps")! as HTMLSpanElement;
const gl = canvas.getContext("webgl2")!;

// function adaptCanvasSize(canvas: HTMLCanvasElement) {
//   canvas.width = window.innerWidth;
//   canvas.height = window.innerHeight;
// }

// adaptCanvasSize(canvas);

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

const resolutionUniform = gl.getUniformLocation(glProgram, "u_resolution");
const timeUniform = gl.getUniformLocation(glProgram, "u_time");

let frameTime = 0;
let lastFrameTime = performance.now();

gl.viewport(0, 0, canvas.width, canvas.height);

function render() {
  gl.uniform2f(resolutionUniform, canvas.width, canvas.height);
  gl.uniform1f(timeUniform, frameTime);

  gl.clearColor(0, 0, 0, 1);
  gl.clear(gl.COLOR_BUFFER_BIT);
  gl.drawArrays(gl.TRIANGLES, 0, 3);

  const currentFrameTime = performance.now();
  frameTime += (currentFrameTime - lastFrameTime) / 1000;
  frameTime %= 8.4;

  fpsCounter.innerText = (1000 / (currentFrameTime - lastFrameTime))
    .toFixed(2)
    .toString();

  lastFrameTime = currentFrameTime;

  requestAnimationFrame(render);
}

requestAnimationFrame(render);
