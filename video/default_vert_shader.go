package video

var vertexShader = `
#if __VERSION__ >= 130
#define COMPAT_VARYING out
#define COMPAT_ATTRIBUTE in
#else
#define COMPAT_VARYING varying
#define COMPAT_ATTRIBUTE attribute
#endif

COMPAT_ATTRIBUTE vec2 vert;
COMPAT_ATTRIBUTE vec2 vertTexCoord;

COMPAT_VARYING vec2 fragTexCoord;

uniform mat4 MVP;

void main() {
  fragTexCoord = vertTexCoord;
  gl_Position = vec4(vert, 0.0, 1.0) * MVP;
}
` + "\x00"
