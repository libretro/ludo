package video

var darkenFragmentShader = `
#version 330

uniform sampler2D tex;
uniform float mask;
uniform vec4 texColor;

in vec2 fragTexCoord;

out vec4 outputColor;

vec4 grayscale(in vec4 c) {
  float average = (c.r + c.g + c.b) / 3.0;
  return vec4(average, average, average, 1.0);
}

vec4 darken(in vec4 c) {
  return vec4(c.r/4, c.g/4, c.b/4, 1.0);
}

void main() {
  vec4 color = texture(tex, fragTexCoord);
  outputColor = texColor * mix(color, darken(grayscale(color)), mask);
}
` + "\x00"
