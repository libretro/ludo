package video

var darkenFragmentShader = `
#if __VERSION__ >= 130
#define COMPAT_VARYING in
#define COMPAT_ATTRIBUTE in
#define COMPAT_TEXTURE texture
#define COMPAT_FRAGCOLOR FragColor
out vec4 COMPAT_FRAGCOLOR;
#else
#define COMPAT_VARYING varying
#define COMPAT_ATTRIBUTE attribute
#define COMPAT_TEXTURE texture2D
#define COMPAT_FRAGCOLOR gl_FragColor
#endif

uniform sampler2D tex;
uniform float mask;
uniform vec4 texColor;

COMPAT_VARYING vec2 fragTexCoord;

vec4 grayscale(vec4 c) {
  float average = (c.r + c.g + c.b) / 3.0;
  return vec4(average, average, average, 1.0);
}

vec4 darken(vec4 c) {
  return vec4(c.r/4.0, c.g/4.0, c.b/4.0, 1.0);
}

void main() {
  vec4 color = COMPAT_TEXTURE(tex, fragTexCoord);
  COMPAT_FRAGCOLOR = texColor * mix(color, darken(grayscale(color)), mask);
}
` + "\x00"
