precision highp float;

uniform vec2 u_resolution;
uniform vec2 u_zoomCenter;
uniform float u_zoomSize;
uniform int u_maxIterations;

vec2 f(vec2 x, vec2 c) {
	return mat2(x, -x.y, x.x) * x + c;
}

void main() {
  vec2 uv = gl_FragCoord.xy / u_resolution;

  vec2 c = u_zoomCenter + (uv * 4.0 - vec2(2.0)) * (u_zoomSize / 4.0);
  vec2 z = vec2(0.0);
  bool escaped = false;

  for (int i = 0; i < 10000; i++) {
    if (i > u_maxIterations) {
      break;
    }

    z = f(z, c);
    if (length(z) > 2.0) {
      escaped = true;
      break;
    }
  }

  gl_FragColor = escaped ? vec4(vec3(1), 1.0) : vec4(vec3(0.0), 1.0);
}
