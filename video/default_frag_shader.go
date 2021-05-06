package video

var defaultFragmentShader = `
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

uniform vec2 OutputSize;
uniform vec2 TextureSize;
uniform vec2 InputSize;
uniform sampler2D Texture;
COMPAT_VARYING vec2 fragTexCoord;

void main() {
  vec4 c = COMPAT_TEXTURE(Texture, fragTexCoord);
  COMPAT_FRAGCOLOR = vec4(c.rgb, 1.0);
}
` + "\x00"
