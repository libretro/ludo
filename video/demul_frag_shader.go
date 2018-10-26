package video

var demulFragmentShader = `
#version 330

uniform sampler2D tex;
uniform float mask;
uniform vec4 texColor;

in vec2 fragTexCoord;

out vec4 outputColor;

vec4 demultiply(in vec4 c) {
  return vec4(c.rgb/c.a, c.a);
}

void main() {
  vec4 color = demultiply(texture(tex, fragTexCoord));
  outputColor = texColor * color;
}
` + "\x00"
