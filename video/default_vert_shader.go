package video

var vertexShader = `
#if __VERSION__ >= 130
#define COMPAT_VARYING out
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

COMPAT_ATTRIBUTE vec2 vert;
COMPAT_ATTRIBUTE vec2 vertTexCoord;

COMPAT_VARYING vec2 fragTexCoord;

void main() {
  fragTexCoord = vertTexCoord;
  gl_Position = vec4(vert, 0, 1);
}
` + "\x00"
