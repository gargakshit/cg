#define PI 3.1415926538
#define noise_amplitude 1.0
#define sphere_radius 1.5
#define lerp mix
#define num_iterations 96

precision mediump float;

uniform vec2 u_resolution;
uniform float u_time;

const float fov = PI / 3.0;

float add_xyz(vec3 p) { return p.x + p.y + p.z; }

float hash(float n) {
  float x = sin(n) * 43758.5453;
  return x - floor(x);
}

float noise(vec3 x) {
  vec3 p = floor(x);
  vec3 f = x - p;
  f = f * (f * vec3(3.0) - f * 2.0);
  float n = add_xyz(p * vec3(1.0, 57.0, 113.0));

  return lerp(lerp(lerp(hash(n + 0.0), hash(n + 1.0), f.x),
                   lerp(hash(n + 57.0), hash(n + 58.0), f.x), f.y),
              lerp(lerp(hash(n + 113.0), hash(n + 114.0), f.x),
                   lerp(hash(n + 170.0), hash(n + 171.0), f.x), f.y),
              f.z);
}

vec3 rotate(vec3 v) {
  return vec3(add_xyz(vec3(0.0, 0.8, 0.6) * v),
              add_xyz(vec3(-0.8, 0.36, -0.48) * v),
              add_xyz(vec3(-0.6, -0.48, 0.64) * v));
}

float fbm(vec3 x) {
  vec3 p = rotate(x);

  float f = 0.0;
  f += 0.5 * noise(p);
  p = 2.32 * p;
  f += 0.25 * noise(p);
  p = 3.03 * p;
  f += 0.125 * noise(p);
  p = 2.61 * p;
  f += 0.0625 * noise(p);

  return f / 0.9375;
}

float signed_distance(vec3 p) {
  float rad = sphere_radius * (sin((u_time - 2.0) / 4.0) + 0.5);
  float displacement = -fbm(3.4 * p) * noise_amplitude;

  return length(p) - (rad + displacement);
}

bool trace_sphere(vec3 orig, vec3 dir, out vec3 pos) {
  pos = orig;

  for (int i = 0; i < num_iterations; i++) {
    float d = signed_distance(pos);
    if (d < 0.0) {
      return true;
    }

    pos = pos + dir * max(0.1 * d, 0.01);
  }

  return false;
}

const float eps = 0.1;

vec3 distance_normal_field(vec3 pos) {
  float d = signed_distance(pos);
  float nx = signed_distance(pos + vec3(eps, 0.0, 0.0)) - d;
  float ny = signed_distance(pos + vec3(0.0, eps, 0.0)) - d;
  float nz = signed_distance(pos + vec3(0.0, 0.0, eps)) - d;

  return normalize(vec3(nx, ny, nz));
}

const vec3 yellow = vec3(1.7, 1.3, 1);
const vec3 orange = vec3(1.0, 0.6, 0.0);
const vec3 red = vec3(1.0, 0.0, 0.0);
const vec3 darkGray = vec3(0.2);
const vec3 gray = vec3(0.4);

vec3 palette(float d) {
  float x = clamp(d, 0.0, 1.0);

  if (x < 0.25) {
    return lerp(gray, darkGray, x * 4.0);
  }

  if (x < 0.5) {
    return lerp(darkGray, red, x * 4.0 - 1.0);
  }

  if (x < 0.75) {
    return lerp(red, orange, x * 4.0 - 2.0);
  }

  return lerp(orange, yellow, x * 4.0 - 3.0);
}

void main() {
  vec3 dir = vec3((gl_FragCoord.x + 0.5) - u_resolution.x / 2.0,
                  (gl_FragCoord.y + 0.5) - u_resolution.y / 2.0,
                  -float(u_resolution.y) / (2.0 * tan(fov / 2.0)));

  vec3 hit;
  if (trace_sphere(vec3(0, 0, 3), normalize(dir), hit)) {
    vec3 light_dir = normalize(vec3(10.0) - hit);
    float intensity = max(0.4, add_xyz(light_dir * distance_normal_field(hit)));
    float noise_level = (sphere_radius - length(hit)) / noise_amplitude;

    gl_FragColor = vec4(intensity * palette((-0.2 + noise_level) * 2.0), 1.0);
  } else {
    gl_FragColor = vec4(0.2, 0.7, 0.8, 1.0);
  }
}
