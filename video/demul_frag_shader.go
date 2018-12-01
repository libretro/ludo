package video

var demulFragmentShader = `
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

vec4 demultiply(vec4 c) {
  return vec4(c.rgb/c.a, c.a);
}

void main() {
  vec4 color = demultiply(COMPAT_TEXTURE(tex, fragTexCoord));
  COMPAT_FRAGCOLOR = texColor * color;
}
` + "\x00"
