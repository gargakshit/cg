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

const zoomCenter = [0.0, 0.0];
const zoomSize = 4.0;
const maxIters = 500;

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

render();
// requestAnimationFrame(render);
