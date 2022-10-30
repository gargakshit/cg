#define lerp mix

precision mediump float;

uniform vec2 u_resolution;
uniform float u_time;

float hash(vec2 n) {
  return fract(sin(dot(n, vec2(12.9898, 4.1414))) * 43758.5453);
}

float noise(vec2 x) {
  vec2 p = floor(x);
  vec2 u = fract(p);
  u = u * u * (3.0 - 2.0 * u);

  float res =
      mix(mix(hash(p), hash(p + vec2(1.0, 0.0)), u.x),
          mix(hash(p + vec2(0.0, 1.0)), hash(p + vec2(1.0, 1.0)), u.x), u.y);
  return res * res;
}

float fbm(vec2 p) {
  float f = 0.0;
  f += 0.5 * noise(p + u_time);
  p = 2.32 * p;
  f += 0.25 * noise(p);
  p = 3.03 * p;
  f += 0.125 * noise(p);
  p = 2.61 * p;
  f += 0.0625 * noise(p);
  p = 2.04 * p;
  f += 0.015625 * noise(p + sin(u_time));

  return f / 0.9375;
}

float pattern(vec2 p) { return fbm(p + fbm(p + fbm(p))); }

float colormap_red(float x) {
  if (x < 0.0) {
    return 54.0 / 255.0;
  } else if (x < 20049.0 / 82979.0) {
    return (829.79 * x + 54.51) / 255.0;
  } else {
    return 1.0;
  }
}

float colormap_green(float x) {
  if (x < 20049.0 / 82979.0) {
    return 0.0;
  } else if (x < 327013.0 / 810990.0) {
    return (8546482679670.0 / 10875673217.0 * x -
            2064961390770.0 / 10875673217.0) /
           255.0;
  } else if (x <= 1.0) {
    return (103806720.0 / 483977.0 * x + 19607415.0 / 483977.0) / 255.0;
  } else {
    return 1.0;
  }
}

float colormap_blue(float x) {
  if (x < 0.0) {
    return 54.0 / 255.0;
  } else if (x < 7249.0 / 82979.0) {
    return (829.79 * x + 54.51) / 255.0;
  } else if (x < 20049.0 / 82979.0) {
    return 127.0 / 255.0;
  } else if (x < 327013.0 / 810990.0) {
    return (792.02249341361393720147485376583 * x -
            64.364790735602331034989206222672) /
           255.0;
  } else {
    return 1.0;
  }
}

vec4 colormap(float x) {
  return vec4(colormap_red(x), colormap_green(x), colormap_blue(x), 1.0);
}

void main() {
  vec2 uv = gl_FragCoord.xy / u_resolution;
  float shade = pattern(uv);

  gl_FragColor = vec4(colormap(shade).rgb, shade);
}
