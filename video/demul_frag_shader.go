package video

var demulFragmentShader = `
uniform sampler2D tex;
uniform float mask;
uniform vec4 texColor;

varying vec2 fragTexCoord;

vec4 demultiply(vec4 c) {
  return vec4(c.rgb/c.a, c.a);
}

void main() {
  vec4 color = demultiply(texture2D(tex, fragTexCoord));
  gl_FragColor = texColor * color;
}
` + "\x00"
