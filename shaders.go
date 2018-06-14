package drawstate

var vertexShader = `
#version 430

layout(location = 0) out vec2 UV;

void main() {
  gl_Position.x = (gl_VertexID == 1) ? 3 : -1;
  gl_Position.y = (gl_VertexID == 2) ? 3 : -1;
  UV = (gl_Position.xy + 1) / 2;
  gl_Position.zw = vec2(1, 1);
}
` + "\x00"

var fragmentShader = `
#version 430

layout(location = 0) in vec2 UV;
layout(location = 0) out vec4 fragColor;
layout(binding = 0, r32ui) uniform uimage2D state0;

vec4 h2rgb(float i) {
  float h = 6 * fract(i);
  float x = fract(h);
  vec4 color;
  if (h < 1.0) {
    color = vec4(1, x, 0, 1);
  } else if (h < 2.0) {
    color = vec4(1 - x, 1, 0, 1);
  } else if (h < 3.0) {
    color = vec4(0, 1, x, 1);
  } else if (h < 4.0) {
    color = vec4(0.0, 1 - x, 1, 1);
  } else if (h < 5.0) {
    color = vec4(x, 0.0, 1, 1);
  } else {
    color = vec4(1, 0.0, 1 - x, 1);
  }
  return color;
}

// Generates a hue from an index.
// steps: Number of steps to go around the colorwheel once.
// cycles: Number of cycles until the exact same colors are repeated.
float hueWheel(uint i, uint steps, uint cycles) {
  i = i % (cycles * steps);
  return fract(float(i) / float(steps) +
               float(i / steps) / float(cycles * steps));
}

void main() {
  uint s =
      imageLoad(state0, ivec2(fract(UV + 0.5) * vec2(imageSize(state0).xy))).r;
  fragColor = h2rgb(hueWheel(s, 7, 13));
}
` + "\x00"
