package video

var darkenFragmentShader = `
uniform sampler2D tex;
uniform float mask;
uniform vec4 texColor;

varying vec2 fragTexCoord;

vec4 grayscale(vec4 c) {
  float average = (c.r + c.g + c.b) / 3.0;
  return vec4(average, average, average, 1.0);
}

vec4 darken(vec4 c) {
  return vec4(c.r/4.0, c.g/4.0, c.b/4.0, 1.0);
}

void main() {
  vec4 color = texture2D(tex, fragTexCoord);
  gl_FragColor = texColor * mix(color, darken(grayscale(color)), mask);
}
` + "\x00"
