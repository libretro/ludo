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

uniform sampler2D Texture;
uniform vec4 color;

COMPAT_VARYING vec2 fragTexCoord;

vec4 demultiply(vec4 c) {
  return vec4(c.rgb/c.a, c.a);
}

void main() {
  COMPAT_FRAGCOLOR = color * demultiply(COMPAT_TEXTURE(Texture, fragTexCoord));
}
` + "\x00"
